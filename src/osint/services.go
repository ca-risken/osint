package main

import (
	"github.com/ca-risken/core/proto/project"
)

type osintService struct {
	repository    osintRepoInterface
	sqs           sqsAPI
	projectClient project.ProjectServiceClient
}
