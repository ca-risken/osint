package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gassara-kys/go-sqs-poller/worker/v4"
	"github.com/vikyd/zero"
)

type SQSConfig struct {
	AWSRegion string
	Endpoint  string

	SubdomainQueueName string
	SubdomainQueueURL  string
	MaxNumberOfMessage int64
	WaitTimeSecond     int64
}

func newSQSConsumer(conf *SQSConfig) *worker.Worker {
	var sqsClient *sqs.SQS
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		appLogger.Fatalf("Failed to create a new session, %v", err)
	}
	if !zero.IsZeroVal(&conf.Endpoint) {
		sqsClient = sqs.New(sess, &aws.Config{
			Region:   &conf.AWSRegion,
			Endpoint: &conf.Endpoint,
		})
	} else {
		sqsClient = sqs.New(sess, &aws.Config{
			Region: &conf.AWSRegion,
		})
	}
	return &worker.Worker{
		Config: &worker.Config{
			QueueName:          conf.SubdomainQueueName,
			QueueURL:           conf.SubdomainQueueURL,
			MaxNumberOfMessage: conf.MaxNumberOfMessage,
			WaitTimeSecond:     conf.WaitTimeSecond,
		},
		Log:       appLogger,
		SqsClient: sqsClient,
	}
}
