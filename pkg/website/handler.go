package website

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
	findingClient finding.FindingServiceClient
	alertClient   alert.AlertServiceClient
	osintClient   osint.OsintServiceClient
	websiteClient *WebsiteClient
	logger        logging.Logger
}

func NewSQSHandler(
	fc finding.FindingServiceClient,
	ac alert.AlertServiceClient,
	oc osint.OsintServiceClient,
	wc *WebsiteClient,
	l logging.Logger,

) *SQSHandler {
	return &SQSHandler{
		findingClient: fc,
		alertClient:   ac,
		osintClient:   oc,
		websiteClient: wc,
		logger:        l,
	}
}

func (s *SQSHandler) HandleMessage(ctx context.Context, sqsMsg *types.Message) error {
	msgBody := aws.ToString(sqsMsg.Body)
	s.logger.Infof(ctx, "got message. message: %v", msgBody)
	// Parse message
	msg, err := message.ParseMessageOSINT(msgBody)
	if err != nil {
		s.logger.Errorf(ctx, "Invalid message. message: %v, error: %v", msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	beforeScanAt := time.Now()
	requestID, err := s.logger.GenerateRequestID(fmt.Sprint(msg.ProjectID))
	if err != nil {
		s.logger.Warnf(ctx, "Failed to generate requestID: err=%+v", err)
		requestID = fmt.Sprint(msg.ProjectID)
	}
	s.logger.Infof(ctx, "start Scan, RequestID=%s", requestID)

	// Run website
	cspan, _ := tracer.StartSpanFromContext(ctx, "runwebsite")
	wappalyzerResult, err := s.websiteClient.run(msg.ResourceName)
	cspan.Finish(tracer.WithError(err))
	if err != nil {
		s.logger.Errorf(ctx, "Failed exec wappalyzer, error: %v", err)
		if err.Error() == "signal: killed" {
			err = errors.New("an error occurred while executing wappalyzer. Scan will restart in a little while")
			s.updateStatusToError(ctx, msg, err)
			return err
		}
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Put Finding and Tag Finding
	if err := s.putFindings(ctx, wappalyzerResult, msg); err != nil {
		s.logger.Errorf(ctx, "Faild to put findings. ResourceName: %v, error: %v", msg.ResourceName, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Clear finding score
	if _, err := s.findingClient.ClearScore(ctx, &finding.ClearScoreRequest{
		DataSource: msg.DataSource,
		ProjectId:  msg.ProjectID,
		Tag:        []string{fmt.Sprintf("osint_id:%v", msg.OsintID)},
		BeforeAt:   beforeScanAt.Unix(),
	}); err != nil {
		s.logger.Errorf(ctx, "Failed to clear finding score. ResourceName: %v, error: %v", msg.ResourceName, err)
		s.updateStatusToError(ctx, msg, err)
		return mimosasqs.WrapNonRetryable(err)
	}

	// Update status
	if err := s.updateScanStatusSuccess(ctx, msg); err != nil {
		s.logger.Errorf(ctx, "Faild to update scan status. ResourceName: %v, error: %v", msg.ResourceName, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	s.logger.Infof(ctx, "end Scan, RequestID=%s", requestID)

	if msg.ScanOnly {
		return nil
	}
	// Call AnalyzeAlert
	if err := s.CallAnalyzeAlert(ctx, msg.ProjectID); err != nil {
		s.logger.Notifyf(ctx, logging.ErrorLevel, "Failed to analyzeAlert, project_id=%d, err=%+v", msg.ProjectID, err)
		return mimosasqs.WrapNonRetryable(err)
	}
	return nil
}

func (s *SQSHandler) updateStatusToError(ctx context.Context, msg *message.OsintQueueMessage, err error) {
	if updateErr := s.updateScanStatusError(ctx, msg, err.Error()); updateErr != nil {
		s.logger.Warnf(ctx, "Failed to update scan status error: err=%+v", updateErr)
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
	s.logger.Infof(ctx, "Success to update osint status, response=%+v", resp)
	return nil
}

func (s *SQSHandler) CallAnalyzeAlert(ctx context.Context, projectID uint32) error {
	_, err := s.alertClient.AnalyzeAlert(ctx, &alert.AnalyzeAlertRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	s.logger.Info(ctx, "Success to analyze alert.")
	return nil
}
