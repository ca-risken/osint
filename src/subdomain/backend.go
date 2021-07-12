package main

import (
	"github.com/kelseyhightower/envconfig"
)

type backendConfig struct {
	//	Port     string `default:"19001"`
	LogLevel     string `split_words:"true" default:"debug"`
	HarvesterDir string `split_words:"true" default:"/theHarvester"`

	//	DB  osintRepoInterface
	//	SQS *sqsClient
}

func newBackendConfig() (*backendConfig, error) {
	config := &backendConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	return config, nil
}
