package main

import (
	"github.com/CyberAgent/mimosa-osint-go/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vikyd/zero"
)

func (r *osintRepository) ListRelOsintDetectWord(projectID, relOsintDataSourceID uint32) (*[]model.RelOsintDetectWord, error) {
	query := `select * from rel_osint_detect_word where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(relOsintDataSourceID) {
		query += " and rel_osint_data_source_id = ?"
		params = append(params, relOsintDataSourceID)
	}
	var data []model.RelOsintDetectWord
	if err := r.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) GetRelOsintDetectWord(projectID uint32, relOsintDetectWordID uint32) (*model.RelOsintDetectWord, error) {
	var data model.RelOsintDetectWord
	if err := r.SlaveDB.Where("project_id = ? AND rel_osint_detect_word_id = ?", projectID, relOsintDetectWordID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) UpsertRelOsintDetectWord(data *model.RelOsintDetectWord) (*model.RelOsintDetectWord, error) {
	var savedData model.RelOsintDetectWord
	update := relOsintDetectWordToMap(data)
	if err := r.MasterDB.Where("project_id = ? AND rel_osint_detect_word_id = ?", data.ProjectID, data.RelOsintDetectWordID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *osintRepository) DeleteRelOsintDetectWord(projectID uint32, relOsintDetectWordID uint32) error {
	if err := r.MasterDB.Where("project_id = ? AND rel_osint_detect_word_id = ?", projectID, relOsintDetectWordID).Delete(&model.RelOsintDetectWord{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *osintRepository) ListOsintDetectWord(projectID uint32) (*[]model.OsintDetectWord, error) {
	query := `select * from osint_detect_word where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
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
	update := osintDetectWordToMap(data)
	if err := r.MasterDB.Where("project_id = ? AND osint_detect_word_id = ?", data.ProjectID, data.OsintDetectWordID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
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

func (r *osintRepository) ListOsintDetectWordFromRelOsintDataSourceID(projectID, relOsintDataSourceID uint32) (*[]model.OsintDetectWord, error) {
	query := `select * from osint_detect_word where osint_detect_word_id = any (select osint_detect_word_id from rel_osint_detect_word where project_id = ? and rel_osint_data_source_id = ?);`
	var params []interface{}
	params = append(params, projectID, relOsintDataSourceID)
	var data []model.OsintDetectWord
	if err := r.MasterDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func relOsintDetectWordToMap(relOsintDetectWord *model.RelOsintDetectWord) map[string]interface{} {
	return map[string]interface{}{
		"rel_osint_detect_word_id": relOsintDetectWord.RelOsintDetectWordID,
		"rel_osint_data_source_id": relOsintDetectWord.RelOsintDataSourceID,
		"osint_detect_word_id":     relOsintDetectWord.OsintDetectWordID,
		"project_id":               relOsintDetectWord.ProjectID,
	}
}

func osintDetectWordToMap(osintDetectWord *model.OsintDetectWord) map[string]interface{} {
	return map[string]interface{}{
		"osint_detect_word_id": osintDetectWord.OsintDetectWordID,
		"word":                 osintDetectWord.Word,
		"project_id":           osintDetectWord.ProjectID,
	}
}
