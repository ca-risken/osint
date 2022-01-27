package main

import (
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/osint/proto/osint"
)

type osintService struct {
	repository    osintRepoInterface
	sqs           sqsAPI
	projectClient project.ProjectServiceClient
}

func newOsintService(config *osintConfig) osint.OsintServiceServer {
	return &osintService{
		repository:    config.DB,
		sqs:           config.SQS,
		projectClient: config.projectClient,
	}
}
