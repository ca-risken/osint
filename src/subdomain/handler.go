package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-osint-go/pkg/common"
	"github.com/CyberAgent/mimosa-osint-go/pkg/message"
	"github.com/CyberAgent/mimosa-osint-go/pkg/model"
	"github.com/CyberAgent/mimosa-osint-go/proto/osint"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/vikyd/zero"
)

type sqsHandler struct {
	findingClient   finding.FindingServiceClient
	alertClient     alert.AlertServiceClient
	osintClient     osint.OsintServiceClient
	subdomainSearch subdomainSearchAPI
}

func newHandler() *sqsHandler {
	h := &sqsHandler{}
	h.subdomainSearch = newSubdomainSearchClient()
	appLogger.Info("Start subdomainsearch Client")
	h.findingClient = newFindingClient()
	appLogger.Info("Start Finding Client")
	h.alertClient = newAlertClient()
	appLogger.Info("Start Alert Client")
	h.osintClient = newOsintClient()
	appLogger.Info("Start Osint Client")
	return h
}

func (s *sqsHandler) HandleMessage(msg *sqs.Message) error {
	msgBody := aws.StringValue(msg.Body)
	appLogger.Infof("got message. message: %v", msgBody)
	// Parse message
	message, err := parseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message. message: %v, error: %v", message, err)
		return err
	}

	detectList, err := getDetectList(message.DetectWord)
	if err != nil {
		appLogger.Errorf("Failed getting detect list, error: %v", err)
		return err
	}

	// Run Harvester
	hosts, err := s.subdomainSearch.run(message.ResourceName, message.RelOsintDataSourceID)
	detectedHosts := detectHost(hosts, message.ResourceName, detectList)
	statusList := getHTTPStatus(detectedHosts)
	findings, err := makeFindings(statusList, message)
	if err != nil {
		appLogger.Errorf("Failed making Findings, error: %v", err)
		return err
	}

	// Put Finding and Tag Finding
	ctx := context.Background()
	if err := s.putFindings(ctx, findings); err != nil {
		appLogger.Errorf("Faild to put findngs. relOsintDataSourceID: %v, error: %v", message.RelOsintDataSourceID, err)
		return err
	}

	// Put RelOsintDataSource
	relOsintDataSource := &osint.RelOsintDataSourceForUpsert{
		RelOsintDataSourceId: message.RelOsintDataSourceID,
		OsintId:              message.OsintID,
		OsintDataSourceId:    message.OsintDataSourceID,
		ProjectId:            message.ProjectID,
		ScanAt:               time.Now().Unix(),
	}

	if err := s.putRelOsintDataSource(relOsintDataSource, true, ""); err != nil {
		appLogger.Errorf("Faild to put rel_osint_data_source. relOsintDataSourceID: %v, error: %v", message.RelOsintDataSourceID, err)
		return err
	}

	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, message.ProjectID); err != nil {
		appLogger.Errorf("Faild to analyze alert. relOsintDataSourceID: %v, error: %v", message.RelOsintDataSourceID, err)
		return err
	}
	return nil

}

func (s *sqsHandler) putFindings(ctx context.Context, findings []*finding.FindingForUpsert) error {
	for _, f := range findings {
		res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
		if err != nil {
			return err
		}
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagOsint)
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagSubdomainSearch)
		s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDomain)
		appLogger.Infof("Success to PutFinding. finding: %v", f)
	}
	return nil
}

func (s *sqsHandler) tagFinding(ctx context.Context, projectID uint32, findingID uint64, tag string) error {

	_, err := s.findingClient.TagFinding(ctx, &finding.TagFindingRequest{
		ProjectId: projectID,
		Tag: &finding.FindingTagForUpsert{
			FindingId: findingID,
			ProjectId: projectID,
			Tag:       tag,
		}})
	if err != nil {
		appLogger.Errorf("Failed to TagFinding. error: %v", err)
		return err
	}
	return nil
}

func (s *sqsHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	appLogger.Info("Success to analyze alert.")
	return nil
}

func (s *sqsHandler) putRelOsintDataSource(relOsintDataSource *osint.RelOsintDataSourceForUpsert, isSuccess bool, errStr string) error {
	ctx := context.Background()

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

func makeFindings(arrStatus *[]httpStatus, message *message.OsintQueueMessage) ([]*finding.FindingForUpsert, error) {
	findings := []*finding.FindingForUpsert{}
	for _, status := range *arrStatus {
		score := getScore(&status)
		description := getDescription(&status)
		data, err := json.Marshal(map[string]httpStatus{"data": status})
		if err != nil {
			return nil, err
		}
		findings = append(findings, &finding.FindingForUpsert{
			Description:      description,
			DataSource:       message.DataSource,
			DataSourceId:     status.HostName,
			ResourceName:     fmt.Sprintf("%v:%v", message.ResourceType, message.ResourceName),
			ProjectId:        message.ProjectID,
			OriginalScore:    score,
			OriginalMaxScore: 10.0,
			Data:             string(data),
		})
	}
	return findings, nil
}

func getScore(status *httpStatus) float32 {
	var score float32 = 1.0
	if !zero.IsZeroVal(status.HTTP) {
		score = score + 1.0
	}
	if !zero.IsZeroVal(status.HTTPS) {
		score = score + 1.0
	}
	if status.HTTP != 401 && status.HTTP != 403 {
		score = score + 1.0
	}
	if status.HTTPS != 401 && status.HTTPS != 403 {
		score = score + 1.0
	}
	return score
}

func getDescription(status *httpStatus) string {
	desc := fmt.Sprintf("%s is accessible from public.", status.HostName)
	if !zero.IsZeroVal(status.HTTP) && !zero.IsZeroVal(status.HTTPS) {
		desc = desc + " (http/https)"
	} else if !zero.IsZeroVal(status.HTTP) {
		desc = desc + " (http)"
	} else {
		desc = desc + " (https)"
	}
	return desc
}

func getStatus(isSuccess bool) osint.Status {
	if isSuccess {
		return osint.Status_OK
	}
	return osint.Status_ERROR
}
