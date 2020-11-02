package main

import (
	"github.com/CyberAgent/mimosa-osint/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vikyd/zero"
)

func (r *osintRepository) ListOsintDetectWord(projectID, relOsintDataSourceID uint32) (*[]model.OsintDetectWord, error) {
	query := `select * from osint_detect_word where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(relOsintDataSourceID) {
		query += " and rel_osint_data_source_id = ?"
		params = append(params, relOsintDataSourceID)
	}
	var data []model.OsintDetectWord
	if err := r.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) GetOsintDetectWord(projectID uint32, osintDetectWordID uint32) (*model.OsintDetectWord, error) {
	var data model.OsintDetectWord
	if err := r.SlaveDB.Where("project_id = ? AND osint_detect_word_id = ?", projectID, osintDetectWordID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) UpsertOsintDetectWord(data *model.OsintDetectWord) (*model.OsintDetectWord, error) {
	var savedData model.OsintDetectWord
	if err := r.MasterDB.Where("project_id = ? AND osint_detect_word_id = ?", data.ProjectID, data.OsintDetectWordID).Assign(data).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *osintRepository) DeleteOsintDetectWord(projectID uint32, osintDetectWordID uint32) error {
	if err := r.MasterDB.Where("project_id = ? AND osint_detect_word_id = ?", projectID, osintDetectWordID).Delete(&model.OsintDetectWord{}).Error; err != nil {
		return err
	}
	return nil
}
