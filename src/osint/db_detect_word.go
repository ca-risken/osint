package main

import (
	"context"

	"github.com/CyberAgent/mimosa-osint/pkg/model"
	"github.com/vikyd/zero"
)

func (r *osintRepository) ListOsintDetectWord(ctx context.Context, projectID, relOsintDataSourceID uint32) (*[]model.OsintDetectWord, error) {
	query := `select * from osint_detect_word where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(relOsintDataSourceID) {
		query += " and rel_osint_data_source_id = ?"
		params = append(params, relOsintDataSourceID)
	}
	var data []model.OsintDetectWord
	if err := r.SlaveDB.WithContext(ctx).Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) GetOsintDetectWord(ctx context.Context, projectID uint32, osintDetectWordID uint32) (*model.OsintDetectWord, error) {
	var data model.OsintDetectWord
	if err := r.SlaveDB.WithContext(ctx).Where("project_id = ? AND osint_detect_word_id = ?", projectID, osintDetectWordID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) UpsertOsintDetectWord(ctx context.Context, data *model.OsintDetectWord) (*model.OsintDetectWord, error) {
	var savedData model.OsintDetectWord
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND osint_detect_word_id = ?", data.ProjectID, data.OsintDetectWordID).Assign(data).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *osintRepository) DeleteOsintDetectWord(ctx context.Context, projectID uint32, osintDetectWordID uint32) error {
	if err := r.MasterDB.WithContext(ctx).Where("project_id = ? AND osint_detect_word_id = ?", projectID, osintDetectWordID).Delete(&model.OsintDetectWord{}).Error; err != nil {
		return err
	}
	return nil
}
