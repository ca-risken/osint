package main

import (
	"github.com/CyberAgent/mimosa-osint/pkg/model"
	_ "github.com/go-sql-driver/mysql"
)

func (r *osintRepository) ListOsintDataSource(projectID uint32, name string) (*[]model.OsintDataSource, error) {
	var data []model.OsintDataSource
	paramName := "%" + name + "%"
	if err := r.SlaveDB.Where("name like ?", paramName).Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) GetOsintDataSource(projectID uint32, osintDataSourceID uint32) (*model.OsintDataSource, error) {
	var data model.OsintDataSource
	if err := r.SlaveDB.Where("osint_data_source_id = ?", osintDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) UpsertOsintDataSource(input *model.OsintDataSource) (*model.OsintDataSource, error) {
	var data model.OsintDataSource
	if err := r.MasterDB.Where("osint_data_source_id = ?", input.OsintDataSourceID).Assign(&model.OsintDataSource{Name: input.Name, Description: input.Description, MaxScore: input.MaxScore}).FirstOrCreate(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) DeleteOsintDataSource(projectID uint32, osintDataSourceID uint32) error {
	if err := r.MasterDB.Where("osint_data_source_id =  ?", osintDataSourceID).Delete(&model.OsintDataSource{}).Error; err != nil {
		return err
	}
	return nil
}
