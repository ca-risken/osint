package main

import (
	"github.com/CyberAgent/mimosa-osint/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vikyd/zero"
)

func (r *osintRepository) ListRelOsintDataSource(projectID, osintID, osintDataSourceID uint32) (*[]model.RelOsintDataSource, error) {
	query := `select * from rel_osint_data_source where project_id = ?`
	var params []interface{}
	params = append(params, projectID)
	if !zero.IsZeroVal(osintID) {
		query += " and osint_id = ?"
		params = append(params, osintID)
	}
	if !zero.IsZeroVal(osintDataSourceID) {
		query += " and osint_data_source_id = ?"
		params = append(params, osintDataSourceID)
	}
	var data []model.RelOsintDataSource
	if err := r.SlaveDB.Raw(query, params...).Scan(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) GetRelOsintDataSource(projectID uint32, relOsintDataSourceID uint32) (*model.RelOsintDataSource, error) {
	var data model.RelOsintDataSource
	if err := r.SlaveDB.Where("project_id = ? AND rel_osint_data_source_id = ?", projectID, relOsintDataSourceID).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *osintRepository) UpsertRelOsintDataSource(data *model.RelOsintDataSource) (*model.RelOsintDataSource, error) {
	var savedData model.RelOsintDataSource
	update := relOsintDataSourceToMap(data)
	if err := r.MasterDB.Where("project_id = ? AND rel_osint_data_source_id = ?", data.ProjectID, data.RelOsintDataSourceID).Assign(update).FirstOrCreate(&savedData).Error; err != nil {
		return nil, err
	}
	return &savedData, nil
}

func (r *osintRepository) DeleteRelOsintDataSource(projectID uint32, relOsintDataSourceID uint32) error {
	if err := r.MasterDB.Where("project_id = ? AND rel_osint_data_source_id = ?", projectID, relOsintDataSourceID).Delete(&model.RelOsintDataSource{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *osintRepository) ListAllRelOsintDataSource() (*[]model.RelOsintDataSource, error) {
	var data []model.RelOsintDataSource
	if err := r.SlaveDB.Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func relOsintDataSourceToMap(relOsintDataSource *model.RelOsintDataSource) map[string]interface{} {
	return map[string]interface{}{
		"rel_osint_data_source_id": relOsintDataSource.RelOsintDataSourceID,
		"osint_id":                 relOsintDataSource.OsintID,
		"osint_data_source_id":     relOsintDataSource.OsintDataSourceID,
		"project_id":               relOsintDataSource.ProjectID,
		"status":                   relOsintDataSource.Status,
		"status_detail":            relOsintDataSource.StatusDetail,
		"scan_at":                  relOsintDataSource.ScanAt,
	}
}
