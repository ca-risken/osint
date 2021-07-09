package message

// OsintQueueMessage is the message for SQS queue
type OsintQueueMessage struct {
	DataSource           string `json:"data_source"`
	RelOsintDataSourceID uint32 `json:"rel_osint_data_source_id"`
	OsintID              uint32 `json:"osint_id"`
	OsintDataSourceID    uint32 `json:"osint_data_source_id"`
	ProjectID            uint32 `json:"project_id"`
	ResourceName         string `json:"resource_name"`
	ResourceType         string `json:"resoorce_type"`
	DetectWord           string `json:"detect_word"`
	ScanOnly             bool   `json:"scan_only,string"`
}
