package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/go-sqs-poller/worker/v5"
)

type SQSConfig struct {
	AWSRegion string
	Endpoint  string

	QueueName          string
	QueueURL           string
	MaxNumberOfMessage int32
	WaitTimeSecond     int32
}

func newSQSConsumer(ctx context.Context, conf *SQSConfig) (*worker.Worker, error) {

	client, err := worker.CreateSqsClient(ctx, conf.AWSRegion, conf.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new SQS client, err=%w", err)
	}

	return &worker.Worker{
		Config: &worker.Config{
			QueueName:          conf.QueueName,
			QueueURL:           conf.QueueURL,
			MaxNumberOfMessage: conf.MaxNumberOfMessage,
			WaitTimeSecond:     conf.WaitTimeSecond,
		},
		Log:       appLogger,
		SqsClient: client,
	}, nil
}
