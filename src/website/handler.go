package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/osint/pkg/message"
	"github.com/ca-risken/osint/proto/osint"
)

type SQSHandler struct {
	findingClient  finding.FindingServiceClient
	alertClient    alert.AlertServiceClient
	osintClient    osint.OsintServiceClient
	wappalyzerPath string
}

func (s *SQSHandler) HandleMessage(ctx context.Context, sqsMsg *sqs.Message) error {
	msgBody := aws.StringValue(sqsMsg.Body)
	appLogger.Infof("got message. message: %v", msgBody)
	// Parse message
	msg, err := message.ParseMessage(msgBody)
	if err != nil {
		appLogger.Errorf("Invalid message. message: %v, error: %v", msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	requestID, err := appLogger.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		appLogger.Warnf("Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof("start Scan, RequestID=%s", requestID)

	websiteClient, err := newWappalyzerClient(s.wappalyzerPath)
	if err != nil {
		appLogger.Errorf("Error occured when configure: %v, error: %v", msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	appLogger.Info("Start website Client")

	// Run website
	_, segment := xray.BeginSubsegment(ctx, "runwebsite")
	wappalyzerResult, err := websiteClient.run(msg.ResourceName)
	segment.Close(err)
	if err != nil {
		appLogger.Errorf("Failed exec wappalyzer, error: %v", err)
		if err.Error() == "signal: killed" {
			err = errors.New("An error occured while executing wappalyzer. Scan will restart in a little while.")
			_ = s.handleErrorWithUpdateStatus(ctx, msg, err)
			return err
		}
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
		appLogger.Notifyf(logging.ErrorLevel, "Failed to analyzeAlert, project_id=%d, err=%+v", msg.ProjectID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil
}

func (s *SQSHandler) handleErrorWithUpdateStatus(ctx context.Context, msg *message.OsintQueueMessage, err error) error {
	if updateErr := s.updateScanStatusError(ctx, msg, err.Error()); updateErr != nil {
		appLogger.Warnf("Failed to update scan status error: err=%+v", updateErr)
	}
	return mimosasqs.WrapNonRetryable(err)
}

func (s *SQSHandler) updateScanStatusError(ctx context.Context, msg *message.OsintQueueMessage, statusDetail string) error {
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

func (s *SQSHandler) updateScanStatusSuccess(ctx context.Context, msg *message.OsintQueueMessage) error {
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

func (s *SQSHandler) putRelOsintDataSource(ctx context.Context, status *osint.PutRelOsintDataSourceRequest) error {
	resp, err := s.osintClient.PutRelOsintDataSource(ctx, status)
	if err != nil {
		return err
	}
	appLogger.Infof("Success to update osint status, response=%+v", resp)
	return nil
}

func (s *SQSHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	appLogger.Info("Success to analyze alert.")
	return nil
}
