package model

import (
	"time"
)

// Osint entity
type Osint struct {
	OsintID      uint32 `gorm:"primary_key"`
	ProjectID    uint32
	ResourceType string
	ResourceName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// OsintDataSource entity
type OsintDataSource struct {
	OsintDataSourceID uint32 `gorm:"primary_key"`
	Name              string
	Description       string
	MaxScore          float32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// RelOsintDataSource entity
type RelOsintDataSource struct {
	RelOsintDataSourceID uint32 `gorm:"primary_key"`
	OsintID              uint32
	OsintDataSourceID    uint32
	ProjectID            uint32
	Status               string
	StatusDetail         string
	ScanAt               time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// OsintDetectWord entity
type OsintDetectWord struct {
	OsintDetectWordID    uint32 `gorm:"primary_key"`
	RelOsintDataSourceID uint32
	Word                 string
	ProjectID            uint32
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// OsintResource entity
type OsintResource struct {
	OsintResourceID uint32 `gorm:"primary_key"`
	ResourceName    string
	ResourceType    string
	ScanAt          time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// OsintRelatedResource entity
type OsintRelatedResource struct {
	OsintResourceID     uint32
	RelatedResourceName string
	RelatedResourceType string
	ScanAt              time.Time
}
