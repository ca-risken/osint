package main

import (
	"context"
	"fmt"

	mimosasql "github.com/ca-risken/common/pkg/database/sql"
	"github.com/ca-risken/osint/pkg/model"
	"github.com/gassara-kys/envconfig"
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
	ListAllRelOsintDataSource(context.Context) (*[]model.RelOsintDataSource, error)
}

type osintRepository struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

func newOsintRepository() osintRepoInterface {
	repo := osintRepository{}
	repo.MasterDB = initDB(true)
	repo.SlaveDB = initDB(false)
	return &repo
}

type dbConfig struct {
	MasterHost     string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	MasterUser     string `split_words:"true" default:"hoge"`
	MasterPassword string `split_words:"true" default:"moge"`
	SlaveHost      string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	SlaveUser      string `split_words:"true" default:"hoge"`
	SlavePassword  string `split_words:"true" default:"moge"`

	Schema  string `required:"true"    default:"mimosa"`
	Port    int    `required:"true"    default:"3306"`
	LogMode bool   `split_words:"true" default:"false"`
}

func initDB(isMaster bool) *gorm.DB {
	conf := &dbConfig{}
	if err := envconfig.Process("DB", conf); err != nil {
		panic(err)
	}

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
	db, err := mimosasql.Open(dsn, conf.LogMode)
	if err != nil {
		fmt.Printf("Failed to open DB. isMaster: %v", isMaster)
		panic(err)
	}
	return db
}
