package main

import (
	"context"

	"github.com/ca-risken/go-sqs-poller/worker/v5"
)

type SQSConfig struct {
	AWSRegion string
	Endpoint  string

	SubdomainQueueName string
	SubdomainQueueURL  string
	MaxNumberOfMessage int32
	WaitTimeSecond     int32
}

func newSQSConsumer(ctx context.Context, conf *SQSConfig) *worker.Worker {

	client, err := worker.CreateSqsClient(ctx, conf.AWSRegion, conf.Endpoint)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create a new client, %v", err)
	}

	return &worker.Worker{
		Config: &worker.Config{
			QueueName:          conf.SubdomainQueueName,
			QueueURL:           conf.SubdomainQueueURL,
			MaxNumberOfMessage: conf.MaxNumberOfMessage,
			WaitTimeSecond:     conf.WaitTimeSecond,
		},
		Log:       appLogger,
		SqsClient: client,
	}
}
