package main

import (
	"github.com/CyberAgent/mimosa-osint-go/proto/osint"
)

type osintService struct {
	repository osintRepoInterface
	sqs        sqsAPI
}

func newOsintService(db osintRepoInterface, s sqsAPI) osint.OsintServiceServer {
	return &osintService{
		repository: db,
		sqs:        s}
}
