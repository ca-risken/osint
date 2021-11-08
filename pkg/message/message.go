package message

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// GuardDutyDataSource is the specific data_source label for subdomain
	SubdomainDataSource = "osint:subdomain"
	// AccessAnalyzerDataSource is the specific data_source label for website
	WebsiteDataSource = "osint:website"
)

// OsintQueueMessage is the message for SQS queue
type OsintQueueMessage struct {
	DataSource           string `json:"data_source"`
	RelOsintDataSourceID uint32 `json:"rel_osint_data_source_id"`
	OsintID              uint32 `json:"osint_id"`
	OsintDataSourceID    uint32 `json:"osint_data_source_id"`
	ProjectID            uint32 `json:"project_id"`
	ResourceName         string `json:"resource_name"`
	ResourceType         string `json:"resource_type"`
	DetectWord           string `json:"detect_word"`
	ScanOnly             bool   `json:"scan_only,string"`
}

// Validate is the validation to OsintQueueMessage
func (o *OsintQueueMessage) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.OsintID, validation.Required),
		validation.Field(&o.OsintDataSourceID, validation.Required),
		validation.Field(&o.RelOsintDataSourceID, validation.Required),
		validation.Field(&o.DataSource, validation.Required, validation.In(
			SubdomainDataSource,
			WebsiteDataSource,
		)),
		validation.Field(&o.ProjectID, validation.Required),
		validation.Field(&o.ResourceName, validation.Required),
		validation.Field(&o.ResourceType, validation.Required),
	)
}

// ParseMessage parse message & validation
func ParseMessage(msg string) (*OsintQueueMessage, error) {
	message := &OsintQueueMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return nil, err
	}
	if err := message.Validate(); err != nil {
		return nil, err
	}
	return message, nil
}
