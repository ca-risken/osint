package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ca-risken/common/pkg/logging"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/datasource-api/proto/osint"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type SQSHandler struct {
	findingClient  finding.FindingServiceClient
	alertClient    alert.AlertServiceClient
	osintClient    osint.OsintServiceClient
	wappalyzerPath string
}

func (s *SQSHandler) HandleMessage(ctx context.Context, sqsMsg *types.Message) error {
	msgBody := aws.ToString(sqsMsg.Body)
	appLogger.Infof(ctx, "got message. message: %v", msgBody)
	// Parse message
	msg, err := message.ParseMessageOSINT(msgBody)
	if err != nil {
		appLogger.Errorf(ctx, "Invalid message. message: %v, error: %v", msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	requestID, err := appLogger.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		appLogger.Warnf(ctx, "Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	appLogger.Infof(ctx, "start Scan, RequestID=%s", requestID)

	websiteClient, err := newWappalyzerClient(s.wappalyzerPath)
	if err != nil {
		appLogger.Errorf(ctx, "Error occured when configure: %v, error: %v", msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	appLogger.Info(ctx, "Start website Client")

	// Run website
	cspan, _ := tracer.StartSpanFromContext(ctx, "runwebsite")
	wappalyzerResult, err := websiteClient.run(msg.ResourceName)
	cspan.Finish(tracer.WithError(err))
	if err != nil {
		appLogger.Errorf(ctx, "Failed exec wappalyzer, error: %v", err)
		if err.Error() == "signal: killed" {
			err = errors.New("an error occurred while executing wappalyzer. Scan will restart in a little while")
			s.updateStatusToError(ctx, msg, err)
			return err
		}
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Clear finding score
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{fmt.Sprintf("osint_id:%v", msg.OsintID)},
	}); err != nil {
		appLogger.Errorf(ctx, "Failed to clear finding score. ResourceName: %v, error: %v", msg.ResourceName, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put Finding and Tag Finding
	if err := s.putFindings(ctx, wappalyzerResult, msg); err != nil {
		appLogger.Errorf(ctx, "Faild to put findings. ResourceName: %v, error: %v", msg.ResourceName, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Update status
	if err := s.updateScanStatusSuccess(ctx, msg); err != nil {
		appLogger.Errorf(ctx, "Faild to update scan status. ResourceName: %v, error: %v", msg.ResourceName, err)
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

func (s *SQSHandler) updateStatusToError(ctx context.Context, msg *message.OsintQueueMessage, err error) {
	if updateErr := s.updateScanStatusError(ctx, msg, err.Error()); updateErr != nil {
		appLogger.Warnf(ctx, "Failed to update scan status error: err=%+v", updateErr)
	}
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
	appLogger.Infof(ctx, "Success to update osint status, response=%+v", resp)
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
