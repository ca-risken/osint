package main

import (
	"fmt"

	"github.com/CyberAgent/mimosa-osint-go/pkg/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
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
	ListRelOsintDetectWord(uint32, uint32) (*[]model.RelOsintDetectWord, error)
	GetRelOsintDetectWord(uint32, uint32) (*model.RelOsintDetectWord, error)
	UpsertRelOsintDetectWord(*model.RelOsintDetectWord) (*model.RelOsintDetectWord, error)
	DeleteRelOsintDetectWord(uint32, uint32) error
	ListOsintDetectWord(uint32) (*[]model.OsintDetectWord, error)
	GetOsintDetectWord(uint32, uint32) (*model.OsintDetectWord, error)
	UpsertOsintDetectWord(*model.OsintDetectWord) (*model.OsintDetectWord, error)
	DeleteOsintDetectWord(uint32, uint32) error
	// For Invoke
	ListOsintDetectWordFromRelOsintDataSourceID(uint32, uint32) (*[]model.OsintDetectWord, error)
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

	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp([%s]:%d)/%s?charset=utf8mb4&interpolateParams=true&parseTime=true&loc=Local",
			user, pass, host, conf.Port, conf.Schema))
	if err != nil {
		fmt.Printf("Failed to open DB. isMaster: %v", isMaster)
		panic(err)
	}
	db.LogMode(conf.LogMode)
	db.SingularTable(true) // if set this to true, `User`'s default table name will be `user`
	return db
}
