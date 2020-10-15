package message

// OsintQueueMessage is the message for SQS queue
type OsintQueueMessage struct {
	DataSource           string `json:"data_source"`
	RelOsintDataSourceID uint32 `json:"jira_setting_id"`
	ProjectID            uint32 `json:"project_id"`
	ResourceName         string `json:"identity_field"`
	ResourceType         string `json:"identity_value"`
}
