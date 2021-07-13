package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-osint/pkg/model"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type osintRepoInterface interface {
	ListOsint(uint32) (*[]model.Osint, error)
	GetOsint(uint32, uint32) (*model.Osint, error)
	UpsertOsint(*model.Osint) (*model.Osint, error)
	DeleteOsint(uint32, uint32) error
	ListOsintDataSource(uint32, string) (*[]model.OsintDataSource, error)
	GetOsintDataSource(uint32, uint32) (*model.OsintDataSource, error)
	UpsertOsintDataSource(*model.OsintDataSource) (*model.OsintDataSource, error)
	DeleteOsintDataSource(uint32, uint32) error
	ListRelOsintDataSource(uint32, uint32, uint32) (*[]model.RelOsintDataSource, error)
	GetRelOsintDataSource(uint32, uint32) (*model.RelOsintDataSource, error)
	UpsertRelOsintDataSource(*model.RelOsintDataSource) (*model.RelOsintDataSource, error)
	DeleteRelOsintDataSource(uint32, uint32) error
	ListOsintDetectWord(uint32, uint32) (*[]model.OsintDetectWord, error)
	GetOsintDetectWord(uint32, uint32) (*model.OsintDetectWord, error)
	UpsertOsintDetectWord(*model.OsintDetectWord) (*model.OsintDetectWord, error)
	DeleteOsintDetectWord(uint32, uint32) error
	// For Invoke
	ListAllRelOsintDataSource() (*[]model.RelOsintDataSource, error)
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
	MasterHost     string `split_words:"true" required:"true"`
	MasterUser     string `split_words:"true" required:"true"`
	MasterPassword string `split_words:"true" required:"true"`
	SlaveHost      string `split_words:"true"`
	SlaveUser      string `split_words:"true"`
	SlavePassword  string `split_words:"true"`

	Schema  string `required:"true"`
	Port    int    `required:"true"`
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		fmt.Printf("Failed to open DB. isMaster: %v", isMaster)
		panic(err)
	}
	if conf.LogMode {
		db.Logger.LogMode(glogger.Info)
	}
	return db
}
