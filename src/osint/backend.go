package main

import (
	"github.com/kelseyhightower/envconfig"
)

type osintConfig struct {
	Port     string `default:"18081"`
	LogLevel string `split_words:"true" default:"debug"`

	DB  osintRepoInterface
	SQS *sqsClient
}

func newOsintConfig() (*osintConfig, error) {
	config := &osintConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	config.DB = newOsintRepository()
	config.SQS = newSQSClient()
	return config, nil
}
