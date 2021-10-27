package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	mimosasqs "github.com/ca-risken/common/pkg/sqs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/osint/pkg/message"
	"github.com/ca-risken/osint/proto/osint"
)

type sqsHandler struct {
	findingClient finding.FindingServiceClient
	alertClient   alert.AlertServiceClient
	osintClient   osint.OsintServiceClient
}

func newHandler() *sqsHandler {
	return &sqsHandler{
		findingClient: newFindingClient(),
		alertClient:   newAlertClient(),
		osintClient:   newOsintClient(),
	}
}

func (s *sqsHandler) HandleMessage(ctx context.Context, sqsMsg *sqs.Message) error {
	msgBody := aws.StringValue(sqsMsg.Body)
	appLogger.Infof("got message. message: %v", msgBody)
	// Parse message
	msg, err := parseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message. message: %v, error: %v", msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	requestID, err := logging.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		appLogger.Warnf("Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof("start Scan, RequestID=%s", requestID)

	wappalyzerClient, err := newWappalyzerClient()
	if err != nil {
		appLogger.Errorf("Error occured when configure: %v, error: %v", msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	appLogger.Info("Start wappalyzer Client")

	// Run wappalyzer
	_, segment := xray.BeginSubsegment(ctx, "runwappalyzer")
	wappalyzerResult, err := wappalyzerClient.run(msg.ResourceName)
	segment.Close(err)
	if err != nil {
		appLogger.Errorf("Failed exec wappalyzer, error: %v", err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}

	// Clear finding score
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{msg.ResourceName},
	}); err != nil {
		appLogger.Errorf("Failed to clear finding score. ResourceName: %v, error: %v", msg.ResourceName, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}

	// Put Finding and Tag Finding
	if err := s.putFindings(ctx, wappalyzerResult, msg); err != nil {
		appLogger.Errorf("Faild to put findings. ResourceName: %v, error: %v", msg.ResourceName, err)
		return s.handleErrorWithUpdateStatus(ctx, msg, err)
	}

	// Update status
	if err := s.updateScanStatusSuccess(ctx, msg); err != nil {
		appLogger.Errorf("Faild to update scan status. ResourceName: %v, error: %v", msg.ResourceName, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	appLogger.Infof("end Scan, RequestID=%s", requestID)

	if msg.ScanOnly {
		return nil
	}
	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, msg.ProjectID); err != nil {
		appLogger.Errorf("Faild to analyze alert. ResourceName: %v, error: %v", msg.ResourceName, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil
}

func (s *sqsHandler) handleErrorWithUpdateStatus(ctx context.Context, msg *message.OsintQueueMessage, err error) error {
	if updateErr := s.updateScanStatusError(ctx, msg, err.Error()); updateErr != nil {
		appLogger.Warnf("Failed to update scan status error: err=%+v", updateErr)
	}
	return mimosasqs.WrapNonRetryable(err)
}

func (s *sqsHandler) updateScanStatusError(ctx context.Context, msg *message.OsintQueueMessage, statusDetail string) error {
	if len(statusDetail) > 200 {
		statusDetail = statusDetail[:200] + " ..." // cut long text
	}
	req := &osint.PutRelOsintDataSourceRequest{
		ProjectId: msg.ProjectID,
		RelOsintDataSource: &osint.RelOsintDataSourceForUpsert{
			RelOsintDataSourceId: msg.RelOsintDataSourceID,
			OsintId:              msg.OsintID,
			OsintDataSourceId:    msg.OsintDataSourceID,
			ProjectId:            msg.ProjectID,
			Status:               osint.Status_ERROR,
			StatusDetail:         statusDetail,
			ScanAt:               time.Now().Unix(),
		}}

	return s.putRelOsintDataSource(ctx, req)
}

func (s *sqsHandler) updateScanStatusSuccess(ctx context.Context, msg *message.OsintQueueMessage) error {
	req := &osint.PutRelOsintDataSourceRequest{
		ProjectId: msg.ProjectID,
		RelOsintDataSource: &osint.RelOsintDataSourceForUpsert{
			RelOsintDataSourceId: msg.RelOsintDataSourceID,
			OsintId:              msg.OsintID,
			OsintDataSourceId:    msg.OsintDataSourceID,
			ProjectId:            msg.ProjectID,
			Status:               osint.Status_OK,
			StatusDetail:         "",
			ScanAt:               time.Now().Unix(),
		}}
	return s.putRelOsintDataSource(ctx, req)
}

func (s *sqsHandler) putRelOsintDataSource(ctx context.Context, status *osint.PutRelOsintDataSourceRequest) error {
	resp, err := s.osintClient.PutRelOsintDataSource(ctx, status)
	if err != nil {
		return err
	}
	appLogger.Infof("Success to update AWS status, response=%+v", resp)
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

func getStatus(isSuccess bool) osint.Status {
	if isSuccess {
		return osint.Status_OK
	}
	return osint.Status_ERROR
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
