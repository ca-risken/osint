package main

import (
	"github.com/ca-risken/core/proto/project"
	"github.com/gassara-kys/envconfig"
)

type osintConfig struct {
	Port     string `default:"18081"`
	LogLevel string `default:"debug" split_words:"true"`
	EnvName  string `default:"local" split_words:"true"`

	DB            osintRepoInterface
	SQS           *sqsClient
	projectClient project.ProjectServiceClient
}

func newOsintConfig() (*osintConfig, error) {
	config := &osintConfig{}
	if err := envconfig.Process("", config); err != nil {
		return nil, err
	}
	config.DB = newOsintRepository()
	config.SQS = newSQSClient()
	config.projectClient = newProjectClient()
	return config, nil
}
