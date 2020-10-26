package main

import (
	"context"

	"github.com/CyberAgent/mimosa-osint-go/pkg/model"
	"github.com/CyberAgent/mimosa-osint-go/proto/osint"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

func (s *osintService) ListRelOsintDetectWord(ctx context.Context, req *osint.ListRelOsintDetectWordRequest) (*osint.ListRelOsintDetectWordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListRelOsintDetectWord(req.ProjectId, req.RelOsintDataSourceId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &osint.ListRelOsintDetectWordResponse{}, nil
		}
		appLogger.Errorf("Failed to List RelOsintDetectWord, error: %v", err)
		return nil, err
	}
	data := osint.ListRelOsintDetectWordResponse{}
	for _, d := range *list {
		data.RelOsintDetectWord = append(data.RelOsintDetectWord, convertRelOsintDetectWord(&d))
	}
	return &data, nil
}

func (s *osintService) GetRelOsintDetectWord(ctx context.Context, req *osint.GetRelOsintDetectWordRequest) (*osint.GetRelOsintDetectWordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	getData, err := s.repository.GetRelOsintDetectWord(req.ProjectId, req.RelOsintDetectWordId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &osint.GetRelOsintDetectWordResponse{}, nil
		}
		appLogger.Errorf("Failed to Get RelOsintDetectWord, error: %v", err)
		return nil, err
	}

	return &osint.GetRelOsintDetectWordResponse{RelOsintDetectWord: convertRelOsintDetectWord(getData)}, nil
}

func (s *osintService) PutRelOsintDetectWord(ctx context.Context, req *osint.PutRelOsintDetectWordRequest) (*osint.PutRelOsintDetectWordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := &model.RelOsintDetectWord{
		RelOsintDetectWordID: req.RelOsintDetectWord.RelOsintDetectWordId,
		RelOsintDataSourceID: req.RelOsintDetectWord.RelOsintDataSourceId,
		OsintDetectWordID:    req.RelOsintDetectWord.OsintDetectWordId,
	}

	registerdData, err := s.repository.UpsertRelOsintDetectWord(data)
	if err != nil {
		appLogger.Errorf("Failed to Put RelOsintDetectWord, error: %v", err)
		return nil, err
	}
	return &osint.PutRelOsintDetectWordResponse{RelOsintDetectWord: convertRelOsintDetectWord(registerdData)}, nil
}

func (s *osintService) DeleteRelOsintDetectWord(ctx context.Context, req *osint.DeleteRelOsintDetectWordRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DeleteRelOsintDetectWord(req.ProjectId, req.RelOsintDetectWordId); err != nil {
		appLogger.Errorf("Failed to Delete RelOsintDetectWord, error: %v", err)
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertRelOsintDetectWord(data *model.RelOsintDetectWord) *osint.RelOsintDetectWord {
	if data == nil {
		return &osint.RelOsintDetectWord{}
	}
	return &osint.RelOsintDetectWord{
		RelOsintDetectWordId: data.RelOsintDetectWordID,
		RelOsintDataSourceId: data.RelOsintDataSourceID,
		OsintDetectWordId:    data.OsintDetectWordID,
		CreatedAt:            data.CreatedAt.Unix(),
		UpdatedAt:            data.CreatedAt.Unix(),
	}
}
