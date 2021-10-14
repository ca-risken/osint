package main

import (
	"github.com/gassara-kys/envconfig"
)

type osintConfig struct {
	Port     string `default:"18081"`
	LogLevel string `default:"debug" split_words:"true"`
	EnvName  string `default:"local" split_words:"true"`

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
