package main

import (
	"context"
	"fmt"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/osint/pkg/model"
	"gorm.io/gorm"
)

type osintRepoInterface interface {
	ListOsint(context.Context, uint32) (*[]model.Osint, error)
	GetOsint(context.Context, uint32, uint32) (*model.Osint, error)
	UpsertOsint(context.Context, *model.Osint) (*model.Osint, error)
	DeleteOsint(context.Context, uint32, uint32) error
	ListOsintDataSource(context.Context, uint32, string) (*[]model.OsintDataSource, error)
	GetOsintDataSource(context.Context, uint32, uint32) (*model.OsintDataSource, error)
	UpsertOsintDataSource(context.Context, *model.OsintDataSource) (*model.OsintDataSource, error)
	DeleteOsintDataSource(context.Context, uint32, uint32) error
	ListRelOsintDataSource(context.Context, uint32, uint32, uint32) (*[]model.RelOsintDataSource, error)
	GetRelOsintDataSource(context.Context, uint32, uint32) (*model.RelOsintDataSource, error)
	UpsertRelOsintDataSource(context.Context, *model.RelOsintDataSource) (*model.RelOsintDataSource, error)
	DeleteRelOsintDataSource(context.Context, uint32, uint32) error
	ListOsintDetectWord(context.Context, uint32, uint32) (*[]model.OsintDetectWord, error)
	GetOsintDetectWord(context.Context, uint32, uint32) (*model.OsintDetectWord, error)
	UpsertOsintDetectWord(context.Context, *model.OsintDetectWord) (*model.OsintDetectWord, error)
	DeleteOsintDetectWord(context.Context, uint32, uint32) error
	// For Invoke
	ListAllRelOsintDataSource(context.Context, uint32) (*[]model.RelOsintDataSource, error)
}

type osintRepository struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

func newOsintRepository(config *DBConfig) osintRepoInterface {
	repo := osintRepository{}
	repo.MasterDB = initDB(config, true)
	repo.SlaveDB = initDB(config, false)
	return &repo
}

type DBConfig struct {
	MasterHost     string
	MasterUser     string
	MasterPassword string
	SlaveHost      string
	SlaveUser      string
	SlavePassword  string

	Schema        string
	Port          int
	LogMode       bool
	MaxConnection int
}

func initDB(conf *DBConfig, isMaster bool) *gorm.DB {
	var user, pass, host string
	if isMaster {
		user = conf.MasterUser
		pass = conf.MasterPassword
		host = conf.MasterHost
	} else {
		user = conf.SlaveUser
		pass = conf.SlavePassword
		host = conf.SlaveHost
	}

	dsn := fmt.Sprintf("%s:%s@tcp([%s]:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
		user, pass, host, conf.Port, conf.Schema)
	db, err := mimosasql.Open(dsn, conf.LogMode, conf.MaxConnection)
	if err != nil {
		fmt.Printf("Failed to open DB. isMaster: %v", isMaster)
		panic(err)
	}
	return db
}
