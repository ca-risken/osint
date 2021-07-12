package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/CyberAgent/mimosa-osint/pkg/message"
	"github.com/CyberAgent/mimosa-osint/pkg/model"
	"github.com/CyberAgent/mimosa-osint/proto/osint"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *osintService) ListRelOsintDataSource(ctx context.Context, req *osint.ListRelOsintDataSourceRequest) (*osint.ListRelOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListRelOsintDataSource(req.ProjectId, req.OsintId, req.OsintDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.ListRelOsintDataSourceResponse{}, nil
		}
		appLogger.Errorf("Failed to List RelOsintDataSource. error: %v", err)
		return nil, err
	}
	data := osint.ListRelOsintDataSourceResponse{}
	for _, d := range *list {
		data.RelOsintDataSource = append(data.RelOsintDataSource, convertRelOsintDataSource(&d))
	}
	return &data, nil
}

func (s *osintService) GetRelOsintDataSource(ctx context.Context, req *osint.GetRelOsintDataSourceRequest) (*osint.GetRelOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetRelOsintDataSource(req.ProjectId, req.RelOsintDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.GetRelOsintDataSourceResponse{}, nil
		}
		appLogger.Errorf("Failed to Get RelOsintDataSource. error: %v", err)
		return nil, err
	}

	return &osint.GetRelOsintDataSourceResponse{RelOsintDataSource: convertRelOsintDataSource(getData)}, nil
}

func (s *osintService) PutRelOsintDataSource(ctx context.Context, req *osint.PutRelOsintDataSourceRequest) (*osint.PutRelOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.RelOsintDataSource{
		RelOsintDataSourceID: req.RelOsintDataSource.RelOsintDataSourceId,
		ProjectID:            req.ProjectId,
		OsintDataSourceID:    req.RelOsintDataSource.OsintDataSourceId,
		OsintID:              req.RelOsintDataSource.OsintId,
		Status:               req.RelOsintDataSource.Status.String(),
		StatusDetail:         req.RelOsintDataSource.StatusDetail,
		ScanAt:               time.Unix(req.RelOsintDataSource.ScanAt, 0),
	}

	registerdData, err := s.repository.UpsertRelOsintDataSource(data)
	if err != nil {
		appLogger.Errorf("Failed to Put RelOsintDataSource. error: %v", err)
		return nil, err
	}
	return &osint.PutRelOsintDataSourceResponse{RelOsintDataSource: convertRelOsintDataSource(registerdData)}, nil
}

func (s *osintService) DeleteRelOsintDataSource(ctx context.Context, req *osint.DeleteRelOsintDataSourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteRelOsintDataSource(req.ProjectId, req.RelOsintDataSourceId); err != nil {
		appLogger.Errorf("Failed to Delete RelOsintDataSource. error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertRelOsintDataSource(data *model.RelOsintDataSource) *osint.RelOsintDataSource {
	if data == nil {
		return &osint.RelOsintDataSource{}
	}
	return &osint.RelOsintDataSource{
		RelOsintDataSourceId: data.RelOsintDataSourceID,
		OsintDataSourceId:    data.OsintDataSourceID,
		OsintId:              data.OsintID,
		ProjectId:            data.ProjectID,
		CreatedAt:            data.CreatedAt.Unix(),
		UpdatedAt:            data.CreatedAt.Unix(),
		Status:               getStatus(data.Status),
		StatusDetail:         data.StatusDetail,
		ScanAt:               data.ScanAt.Unix(),
	}
}

func (s *osintService) InvokeScan(ctx context.Context, req *osint.InvokeScanRequest) (*osint.InvokeScanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	relOsintDataSourceData, err := s.repository.GetRelOsintDataSource(req.ProjectId, req.RelOsintDataSourceId)
	if err != nil {
		return nil, err
	}
	osintDataSourceData, err := s.repository.GetOsintDataSource(relOsintDataSourceData.ProjectID, relOsintDataSourceData.OsintDataSourceID)
	if err != nil {
		return nil, err
	}
	osintData, err := s.repository.GetOsint(relOsintDataSourceData.ProjectID, relOsintDataSourceData.OsintID)
	if err != nil {
		return nil, err
	}
	detectWord, err := s.repository.ListOsintDetectWord(relOsintDataSourceData.ProjectID, relOsintDataSourceData.RelOsintDataSourceID)
	if err != nil {
		return nil, err
	}
	jsonDetectWord, err := json.Marshal(map[string][]model.OsintDetectWord{"DetectWord": *detectWord})
	if err != nil {
		return nil, err
	}
	msg := &message.OsintQueueMessage{
		DataSource:           osintDataSourceData.Name,
		RelOsintDataSourceID: req.RelOsintDataSourceId,
		OsintID:              relOsintDataSourceData.OsintID,
		OsintDataSourceID:    relOsintDataSourceData.OsintDataSourceID,
		ProjectID:            req.ProjectId,
		ResourceType:         osintData.ResourceType,
		ResourceName:         osintData.ResourceName,
		DetectWord:           string(jsonDetectWord),
		ScanOnly:             req.ScanOnly,
	}
	resp, err := s.sqs.send(msg)
	if err != nil {
		appLogger.Errorf("Invoked scan. Error: %v", err)
		return nil, err
	}
	if _, err = s.repository.UpsertRelOsintDataSource(&model.RelOsintDataSource{
		RelOsintDataSourceID: relOsintDataSourceData.RelOsintDataSourceID,
		OsintID:              relOsintDataSourceData.OsintID,
		OsintDataSourceID:    relOsintDataSourceData.OsintDataSourceID,
		ProjectID:            relOsintDataSourceData.ProjectID,
		Status:               osint.Status_IN_PROGRESS.String(),
		StatusDetail:         fmt.Sprintf("Start scan at %+v", time.Now().Format(time.RFC3339)),
		ScanAt:               relOsintDataSourceData.ScanAt,
	}); err != nil {
		appLogger.Errorf("Failed to update scan status: %+v", err)
		return nil, err
	}
	appLogger.Infof("Invoked scan. MessageId: %v", *resp.MessageId)
	return &osint.InvokeScanResponse{Message: "Invoke Scan."}, nil
}

func (s *osintService) InvokeScanAll(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {

	list, err := s.repository.ListAllRelOsintDataSource()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &empty.Empty{}, nil
		}
		appLogger.Errorf("Failed to List AllRelOsintDataSource. error: %v", err)
		return nil, err
	}

	for _, relOsintDataSource := range *list {
		if _, err := s.InvokeScan(ctx, &osint.InvokeScanRequest{
			ProjectId:            relOsintDataSource.ProjectID,
			RelOsintDataSourceId: relOsintDataSource.RelOsintDataSourceID,
			ScanOnly:             true,
		}); err != nil {
			// errorが出ても続行
			appLogger.Errorf("InvokeScanAll error: err=%+v", err)
		}
	}

	return &empty.Empty{}, nil
}

func getStatus(s string) osint.Status {
	statusKey := strings.ToUpper(s)
	if _, ok := osint.Status_value[statusKey]; !ok {
		return osint.Status_UNKNOWN
	}
	switch statusKey {
	case osint.Status_OK.String():
		return osint.Status_OK
	case osint.Status_CONFIGURED.String():
		return osint.Status_CONFIGURED
	case osint.Status_IN_PROGRESS.String():
		return osint.Status_IN_PROGRESS
	case osint.Status_ERROR.String():
		return osint.Status_ERROR
	default:
		return osint.Status_UNKNOWN
	}
}
