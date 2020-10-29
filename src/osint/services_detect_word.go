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
	appLogger.Infof("updata: %v", req.RelOsintDetectWord)
	data := &model.RelOsintDetectWord{
		RelOsintDetectWordID: req.RelOsintDetectWord.RelOsintDetectWordId,
		RelOsintDataSourceID: req.RelOsintDetectWord.RelOsintDataSourceId,
		OsintDetectWordID:    req.RelOsintDetectWord.OsintDetectWordId,
		ProjectID:            req.RelOsintDetectWord.ProjectId,
	}
	appLogger.Infof("updata model: %v", data)
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

func (s *osintService) ListOsintDetectWord(ctx context.Context, req *osint.ListOsintDetectWordRequest) (*osint.ListOsintDetectWordResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := s.repository.ListOsintDetectWord(req.ProjectId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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
		if gorm.IsRecordNotFoundError(err) {
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
		OsintDetectWordID: req.OsintDetectWord.OsintDetectWordId,
		Word:              req.OsintDetectWord.Word,
		ProjectID:         req.OsintDetectWord.ProjectId,
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

func convertRelOsintDetectWord(data *model.RelOsintDetectWord) *osint.RelOsintDetectWord {
	if data == nil {
		return &osint.RelOsintDetectWord{}
	}
	return &osint.RelOsintDetectWord{
		RelOsintDetectWordId: data.RelOsintDetectWordID,
		RelOsintDataSourceId: data.RelOsintDataSourceID,
		OsintDetectWordId:    data.OsintDetectWordID,
		ProjectId:            data.ProjectID,
		CreatedAt:            data.CreatedAt.Unix(),
		UpdatedAt:            data.CreatedAt.Unix(),
	}
}

func convertOsintDetectWord(data *model.OsintDetectWord) *osint.OsintDetectWord {
	if data == nil {
		return &osint.OsintDetectWord{}
	}
	return &osint.OsintDetectWord{
		OsintDetectWordId: data.OsintDetectWordID,
		Word:              data.Word,
		ProjectId:         data.ProjectID,
		CreatedAt:         data.CreatedAt.Unix(),
		UpdatedAt:         data.CreatedAt.Unix(),
	}
}
