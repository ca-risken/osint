package main

import (
	"context"
	"errors"

	"github.com/ca-risken/osint/pkg/model"
	"github.com/ca-risken/osint/proto/osint"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *osintService) ListOsint(ctx context.Context, req *osint.ListOsintRequest) (*osint.ListOsintResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOsint(ctx, req.ProjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.ListOsintResponse{}, nil
		}
		appLogger.Errorf("Failed to List Osint. error: %v", err)
		return nil, err
	}
	data := osint.ListOsintResponse{}
	for _, d := range *list {
		data.Osint = append(data.Osint, convertOsint(&d))
	}
	return &data, nil
}

func (s *osintService) GetOsint(ctx context.Context, req *osint.GetOsintRequest) (*osint.GetOsintResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetOsint(ctx, req.ProjectId, req.OsintId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.GetOsintResponse{}, nil
		}
		appLogger.Errorf("Failed to Get Osint. error: %v", err)
		return nil, err
	}

	return &osint.GetOsintResponse{Osint: convertOsint(getData)}, nil
}

func (s *osintService) PutOsint(ctx context.Context, req *osint.PutOsintRequest) (*osint.PutOsintResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.Osint{
		OsintID:      req.Osint.OsintId,
		ResourceType: req.Osint.ResourceType,
		ResourceName: req.Osint.ResourceName,
		ProjectID:    req.Osint.ProjectId,
	}

	registerdData, err := s.repository.UpsertOsint(ctx, data)
	if err != nil {
		appLogger.Errorf("Failed to Put Osint. error: %v", err)
		return nil, err
	}
	return &osint.PutOsintResponse{Osint: convertOsint(registerdData)}, nil
}

func (s *osintService) DeleteOsint(ctx context.Context, req *osint.DeleteOsintRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOsint(ctx, req.ProjectId, req.OsintId); err != nil {
		appLogger.Errorf("Failed to Delete Osint. error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertOsint(data *model.Osint) *osint.Osint {
	if data == nil {
		return &osint.Osint{}
	}
	return &osint.Osint{
		OsintId:      data.OsintID,
		ResourceType: data.ResourceType,
		ResourceName: data.ResourceName,
		ProjectId:    data.ProjectID,
		CreatedAt:    data.CreatedAt.Unix(),
		UpdatedAt:    data.CreatedAt.Unix(),
	}
}
