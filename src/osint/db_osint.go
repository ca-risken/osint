package main

import (
	"github.com/CyberAgent/mimosa-osint/pkg/model"
)

func (r *osintRepository) ListOsint(projectID uint32) (*[]model.Osint, error) {
	query := `select * from osint where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	var data []model.Osint
	if err := r.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) GetOsint(projectID uint32, osintID uint32) (*model.Osint, error) {
	var data model.Osint
	if err := r.SlaveDB.Where("project_id = ? AND osint_id = ?", projectID, osintID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) UpsertOsint(data *model.Osint) (*model.Osint, error) {
	var savedData model.Osint
	if err := r.MasterDB.Where("project_id = ? AND osint_id = ?", data.ProjectID, data.OsintID).Assign(data).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *osintRepository) DeleteOsint(projectID uint32, osintID uint32) error {
	if err := r.MasterDB.Where("project_id = ? AND osint_id = ?", projectID, osintID).Delete(&model.Osint{}).Error; err != nil {
		return err
	}
	return nil
}
