package main

import (
	"context"

	"github.com/CyberAgent/mimosa-osint/pkg/model"
	"github.com/CyberAgent/mimosa-osint/proto/osint"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

func (s *osintService) ListOsintDataSource(ctx context.Context, req *osint.ListOsintDataSourceRequest) (*osint.ListOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOsintDataSource(req.ProjectId, req.Name)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &osint.ListOsintDataSourceResponse{}, nil
		}
		appLogger.Errorf("Failed to List OsintDataSource, error: %v", err)
		return nil, err
	}
	data := osint.ListOsintDataSourceResponse{}
	for _, d := range *list {
		data.OsintDataSource = append(data.OsintDataSource, convertOsintDataSource(&d))
	}
	return &data, nil
}

func (s *osintService) GetOsintDataSource(ctx context.Context, req *osint.GetOsintDataSourceRequest) (*osint.GetOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetOsintDataSource(req.ProjectId, req.OsintDataSourceId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &osint.GetOsintDataSourceResponse{}, nil
		}
		appLogger.Errorf("Failed to Get OsintDataSource, error: %v", err)
		return nil, err
	}

	return &osint.GetOsintDataSourceResponse{OsintDataSource: convertOsintDataSource(getData)}, nil
}

func (s *osintService) PutOsintDataSource(ctx context.Context, req *osint.PutOsintDataSourceRequest) (*osint.PutOsintDataSourceResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.OsintDataSource{
		OsintDataSourceID: req.OsintDataSource.OsintDataSourceId,
		Name:              req.OsintDataSource.Name,
		Description:       req.OsintDataSource.Description,
		MaxScore:          req.OsintDataSource.MaxScore,
	}

	registerdData, err := s.repository.UpsertOsintDataSource(data)
	if err != nil {
		appLogger.Errorf("Failed to Put OsintDataSource, error: %v", err)
		return nil, err
	}
	return &osint.PutOsintDataSourceResponse{OsintDataSource: convertOsintDataSource(registerdData)}, nil
}

func (s *osintService) DeleteOsintDataSource(ctx context.Context, req *osint.DeleteOsintDataSourceRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOsintDataSource(req.ProjectId, req.OsintDataSourceId); err != nil {
		appLogger.Errorf("Failed to Delete OsintDataSource, error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertOsintDataSource(data *model.OsintDataSource) *osint.OsintDataSource {
	if data == nil {
		return &osint.OsintDataSource{}
	}
	return &osint.OsintDataSource{
		OsintDataSourceId: data.OsintDataSourceID,
		Name:              data.Name,
		Description:       data.Description,
		MaxScore:          data.MaxScore,
		CreatedAt:         data.CreatedAt.Unix(),
		UpdatedAt:         data.CreatedAt.Unix(),
	}
}
