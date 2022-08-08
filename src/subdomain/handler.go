package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ca-risken/common/pkg/logging"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/pkg/model"
	"github.com/ca-risken/datasource-api/proto/osint"
	"golang.org/x/sync/semaphore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type SQSHandler struct {
	findingClient      finding.FindingServiceClient
	alertClient        alert.AlertServiceClient
	osintClient        osint.OsintServiceClient
	harvesterConfig    HarvesterConfig
	inspectConcurrency int64
}

func (s *SQSHandler) HandleMessage(ctx context.Context, sqsMsg *types.Message) error {
	msgBody := aws.ToString(sqsMsg.Body)
	appLogger.Infof(ctx, "got message. message: %v", msgBody)
	// Parse message
	msg, err := message.ParseMessageOSINT(msgBody)
	if err != nil {
		appLogger.Errorf(ctx, "Invalid message. message: %v, error: %v", msgBody, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	requestID, err := appLogger.GenerateRequestID(fmt.Sprintf("%v-%v", msg.ProjectID, msg.RelOsintDataSourceID))
	if err != nil {
		appLogger.Warnf(ctx, "Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof(ctx, "start Scan, RequestID=%s", requestID)
	detectList, err := getDetectList(msg.DetectWord)
	if err != nil {
		appLogger.Errorf(ctx, "Failed getting detect list, error: %v", err)
		return err
	}

	// Run Harvester
	cspan, cctx := tracer.StartSpanFromContext(ctx, "runHarvester")
	appLogger.Infof(cctx, "start harvester, RequestID=%s", requestID)
	hosts, err := s.harvesterConfig.run(cctx, msg.ResourceName, msg.RelOsintDataSourceID)
	cspan.Finish(tracer.WithError(err))
	if err != nil {
		appLogger.Errorf(cctx, "Failed exec theHarvester, error: %v", err)
		strError := "An error occured while executing osint tool. Ask the system administrator."
		if err.Error() == "signal: killed" {
			strError = "An error occured while executing osint tool. Scan will restart in a little while."
		}
		updateErr := s.putRelOsintDataSource(ctx, msg, false, strError)
		if updateErr != nil {
			appLogger.Warnf(ctx, "Failed to update scan status error: err=%+v", updateErr)
		}
		return err
	}
	appLogger.Infof(cctx, "end harvester, RequestID=%s", requestID)

	wg := sync.WaitGroup{}
	mutex := &sync.Mutex{}
	osintResults := []osintResult{}
	sem := semaphore.NewWeighted(s.inspectConcurrency)
	appLogger.Infof(ctx, "start scan hosts, RequestID=%s", requestID)
	for _, h := range *hosts {
		wg.Add(1)
		if err := sem.Acquire(ctx, 1); err != nil {
			appLogger.Errorf(ctx, "failed to acquire semaphore: %v", err)
			wg.Done()
			return mimosasqs.WrapNonRetryable(err)
		}

		go func(h host) {
			defer func() {
				sem.Release(1)
				wg.Done()
			}()
			privateExpose := searchPrivateExpose(h, detectList)
			takeover := searchTakeover(h.HostName)
			certificateExpiration := privateExpose.checkCertificateExpiration()

			mutex.Lock()
			osintResults = append(osintResults, osintResult{Host: h, PrivateExpose: privateExpose, Takeover: takeover, CertificateExpiration: certificateExpiration})
			mutex.Unlock()
		}(h)
	}
	wg.Wait()
	appLogger.Infof(ctx, "end scan hosts, RequestID=%s", requestID)

	findings, err := makeFindings(ctx, &osintResults, msg)
	if err != nil {
		appLogger.Errorf(ctx, "Failed making Findings, error: %v", err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Clear finding score
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{msg.ResourceName},
	}); err != nil {
		appLogger.Errorf(ctx, "Failed to clear finding score. ResourceName: %v, error: %v", msg.ResourceName, err)
		updateErr := s.putRelOsintDataSource(ctx, msg, false, err.Error())
		if updateErr != nil {
			appLogger.Warnf(ctx, "Failed to update scan status error: err=%+v", updateErr)
		}
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put Finding and Tag Finding
	if err := s.putFindings(ctx, findings, msg.ProjectID, msg.ResourceName); err != nil {
		appLogger.Errorf(ctx, "Failed to put findings. relOsintDataSourceID: %v, error: %v", msg.RelOsintDataSourceID, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put RelOsintDataSource
	if err := s.putRelOsintDataSource(ctx, msg, true, ""); err != nil {
		appLogger.Errorf(ctx, "Failed to put rel_osint_data_source. relOsintDataSourceID: %v, error: %v", msg.RelOsintDataSourceID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	appLogger.Infof(ctx, "end Scan, RequestID=%s", requestID)

	if msg.ScanOnly {
		return nil
	}
	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, msg.ProjectID); err != nil {
		appLogger.Notifyf(ctx, logging.ErrorLevel, "Failed to analyzeAlert, project_id=%d, err=%+v", msg.ProjectID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil

}

func (s *SQSHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	appLogger.Info(ctx, "Success to analyze alert.")
	return nil
}

func (s *SQSHandler) putRelOsintDataSource(ctx context.Context, message *message.OsintQueueMessage, isSuccess bool, errStr string) error {

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
