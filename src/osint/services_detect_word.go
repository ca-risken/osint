package main

import (
	"context"
	"errors"

	"github.com/CyberAgent/mimosa-osint/pkg/model"
	"github.com/CyberAgent/mimosa-osint/proto/osint"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (s *osintService) ListOsintDetectWord(ctx context.Context, req *osint.ListOsintDetectWordRequest) (*osint.ListOsintDetectWordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOsintDetectWord(req.ProjectId, req.RelOsintDataSourceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.ListOsintDetectWordResponse{}, nil
		}
		appLogger.Errorf("Failed to List OsintDetectWord, error: %v", err)
		return nil, err
	}
	data := osint.ListOsintDetectWordResponse{}
	for _, d := range *list {
		data.OsintDetectWord = append(data.OsintDetectWord, convertOsintDetectWord(&d))
	}
	return &data, nil
}

func (s *osintService) GetOsintDetectWord(ctx context.Context, req *osint.GetOsintDetectWordRequest) (*osint.GetOsintDetectWordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetOsintDetectWord(req.ProjectId, req.OsintDetectWordId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &osint.GetOsintDetectWordResponse{}, nil
		}
		appLogger.Errorf("Failed to Get OsintDetectWord, error: %v", err)
		return nil, err
	}

	return &osint.GetOsintDetectWordResponse{OsintDetectWord: convertOsintDetectWord(getData)}, nil
}

func (s *osintService) PutOsintDetectWord(ctx context.Context, req *osint.PutOsintDetectWordRequest) (*osint.PutOsintDetectWordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.OsintDetectWord{
		OsintDetectWordID:    req.OsintDetectWord.OsintDetectWordId,
		RelOsintDataSourceID: req.OsintDetectWord.RelOsintDataSourceId,
		Word:                 req.OsintDetectWord.Word,
		ProjectID:            req.OsintDetectWord.ProjectId,
	}

	registerdData, err := s.repository.UpsertOsintDetectWord(data)
	if err != nil {
		appLogger.Errorf("Failed to Put OsintDetectWord, error: %v", err)
		return nil, err
	}
	return &osint.PutOsintDetectWordResponse{OsintDetectWord: convertOsintDetectWord(registerdData)}, nil
}

func (s *osintService) DeleteOsintDetectWord(ctx context.Context, req *osint.DeleteOsintDetectWordRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteOsintDetectWord(req.ProjectId, req.OsintDetectWordId); err != nil {
		appLogger.Errorf("Failed to Delete OsintDetectWord, error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertOsintDetectWord(data *model.OsintDetectWord) *osint.OsintDetectWord {
	if data == nil {
		return &osint.OsintDetectWord{}
	}
	return &osint.OsintDetectWord{
		OsintDetectWordId:    data.OsintDetectWordID,
		RelOsintDataSourceId: data.RelOsintDataSourceID,
		Word:                 data.Word,
		ProjectId:            data.ProjectID,
		CreatedAt:            data.CreatedAt.Unix(),
		UpdatedAt:            data.CreatedAt.Unix(),
	}
}
