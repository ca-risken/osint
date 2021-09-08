package main

import (
	"context"

	"github.com/ca-risken/osint/pkg/model"
)

func (r *osintRepository) ListOsintDataSource(ctx context.Context, projectID uint32, name string) (*[]model.OsintDataSource, error) {
	var data []model.OsintDataSource
	paramName := "%" + name + "%"
	if err := r.SlaveDB.WithContext(ctx).Where("name like ?", paramName).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) GetOsintDataSource(ctx context.Context, projectID uint32, osintDataSourceID uint32) (*model.OsintDataSource, error) {
	var data model.OsintDataSource
	if err := r.SlaveDB.WithContext(ctx).Where("osint_data_source_id = ?", osintDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) UpsertOsintDataSource(ctx context.Context, input *model.OsintDataSource) (*model.OsintDataSource, error) {
	var data model.OsintDataSource
	if err := r.MasterDB.WithContext(ctx).Where("osint_data_source_id = ?", input.OsintDataSourceID).Assign(&model.OsintDataSource{Name: input.Name, Description: input.Description, MaxScore: input.MaxScore}).FirstOrCreate(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) DeleteOsintDataSource(ctx context.Context, projectID uint32, osintDataSourceID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("osint_data_source_id =  ?", osintDataSourceID).Delete(&model.OsintDataSource{}).Error; err != nil {
		return err
	}
	return nil
}
