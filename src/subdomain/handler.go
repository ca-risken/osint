package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CyberAgent/mimosa-common/pkg/logging"
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-osint/pkg/message"
	"github.com/CyberAgent/mimosa-osint/pkg/model"
	"github.com/CyberAgent/mimosa-osint/proto/osint"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsHandler struct {
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	osintClient     osint.OsintServiceClient
	harvesterConfig harvesterConfig
}

func newHandler() *sqsHandler {
	h := &sqsHandler{}
	h.harvesterConfig = newHarvesterConfig()
	appLogger.Info("Load Harvester Config")
	h.findingClient = newFindingClient()
	appLogger.Info("Start Finding Client")
	h.alertClient = newAlertClient()
	appLogger.Info("Start Alert Client")
	h.osintClient = newOsintClient()
	appLogger.Info("Start Osint Client")
	return h
}

func (s *sqsHandler) HandleMessage(sqsMsg *sqs.Message) error {
	msgBody := aws.StringValue(sqsMsg.Body)
	appLogger.Infof("got message. message: %v", msgBody)
	// Parse message
	msg, err := parseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message. message: %v, error: %v", msgBody, err)
		return err
	}

	requestID, err := logging.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		appLogger.Warnf("Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof("start Scan, RequestID=%s", requestID)
	detectList, err := getDetectList(msg.DetectWord)
	if err != nil {
		appLogger.Errorf("Failed getting detect list, error: %v", err)
		return err
	}

	// Run Harvester
	hosts, err := s.harvesterConfig.run(msg.ResourceName, msg.RelOsintDataSourceID)
	if err != nil {
		appLogger.Errorf("Failed exec theHarvester, error: %v", err)
		strError := "An error occured while executing osint tool. Ask the system administrator."
		if err.Error() == "signal: killed" {
			strError = "An error occured while executing osint tool. Scan will restart in a little while."
		}
		_ = s.putRelOsintDataSource(msg, false, strError)
		return err
	}

	//hosts, err := tmpRun()
	osintResults, err := inspectDomain(hosts, detectList)
	if err != nil {
		appLogger.Errorf("Failed get osintResults, error: %v", err)
		_ = s.putRelOsintDataSource(msg, false, "An error occured while investing resource. Ask the system administrator.")
		return err
	}
	findings, err := makeFindings(osintResults.OsintResults, msg)
	if err != nil {
		appLogger.Errorf("Failed making Findings, error: %v", err)
		return err
	}

	// Put Finding and Tag Finding
	ctx := context.Background()
	if err := s.putFindings(ctx, findings); err != nil {
		appLogger.Errorf("Faild to put findngs. relOsintDataSourceID: %v, error: %v", msg.RelOsintDataSourceID, err)
		return err
	}

	// Put RelOsintDataSource
	if err := s.putRelOsintDataSource(msg, true, ""); err != nil {
		appLogger.Errorf("Faild to put rel_osint_data_source. relOsintDataSourceID: %v, error: %v", msg.RelOsintDataSourceID, err)
		return err
	}
	appLogger.Infof("end Scan, RequestID=%s", requestID)

	if msg.ScanOnly {
		return nil
	}
	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, msg.ProjectID); err != nil {
		appLogger.Errorf("Faild to analyze alert. relOsintDataSourceID: %v, error: %v", msg.RelOsintDataSourceID, err)
		return err
	}
	return nil

}

func inspectDomain(hosts *[]host, detectList *[]string) (*osintResults, error) {
	arr := []osintResult{}
	for _, h := range *hosts {
		privateExpose := searchPrivateExpose(h, detectList)
		takeover := searchTakeover(h.HostName)
		certificateExpiration := privateExpose.checkCertificateExpiration()
		osintResult := osintResult{Host: h, PrivateExpose: privateExpose, Takeover: takeover, CertificateExpiration: certificateExpiration}
		arr = append(arr, osintResult)
	}
	return &osintResults{OsintResults: &arr}, nil
}

func (s *sqsHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	appLogger.Info("Success to analyze alert.")
	return nil
}

func (s *sqsHandler) putRelOsintDataSource(message *message.OsintQueueMessage, isSuccess bool, errStr string) error {
	ctx := context.Background()

	relOsintDataSource := &osint.RelOsintDataSourceForUpsert{
		RelOsintDataSourceId: message.RelOsintDataSourceID,
		OsintId:              message.OsintID,
		OsintDataSourceId:    message.OsintDataSourceID,
		ProjectId:            message.ProjectID,
		ScanAt:               time.Now().Unix(),
	}

	relOsintDataSource.Status = getStatus(isSuccess)
	if isSuccess {
		relOsintDataSource.StatusDetail = ""
	} else {
		errDetail := errStr
		relOsintDataSource.StatusDetail = errDetail
	}
	_, err := s.osintClient.PutRelOsintDataSource(ctx, &osint.PutRelOsintDataSourceRequest{ProjectId: relOsintDataSource.ProjectId, RelOsintDataSource: relOsintDataSource})
	if err != nil {
		return err
	}

	return nil
}

func parseMessage(msg string) (*message.OsintQueueMessage, error) {
	message := &message.OsintQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	//	if err := message.Validate(); err != nil {
	//		return nil, err
	//	}
	return message, nil
}

func getDetectList(detectWord string) (*[]string, error) {
	var mapDetect map[string][]model.OsintDetectWord
	if err := json.Unmarshal([]byte(detectWord), &mapDetect); err != nil {
		return nil, err
	}
	ret := []string{}
	for _, detectWord := range mapDetect["DetectWord"] {
		ret = append(ret, detectWord.Word)
	}
	return &ret, nil
}

func getStatus(isSuccess bool) osint.Status {
	if isSuccess {
		return osint.Status_OK
	}
	return osint.Status_ERROR
}
