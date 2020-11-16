package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/h2ik/go-sqs-poller/v3/worker"
	"github.com/kelseyhightower/envconfig"
)

type sqsConfig struct {
	AWSRegion string `envconfig:"aws_region" default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://localhost:9324"`

	PrivateExposeQueueName string `split_words:"true" default:"osint-privateexpose"`
	PrivateExposeQueueURL  string `split_words:"true" default:"http://localhost:9324/queue/osint-privateexpose"`
	MaxNumberOfMessage     int64  `split_words:"true" default:"10"`
	WaitTimeSecond         int64  `split_words:"true" default:"20"`
}

func newSQSConsumer() *worker.Worker {
	var conf sqsConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Errorf("Failed to start sqs consumer. error:%v ", err)
	}
	sqsClient := sqs.New(session.New(), &aws.Config{
		Region:   &conf.AWSRegion,
		Endpoint: &conf.Endpoint,
	})
	return &worker.Worker{
		Config: &worker.Config{
			QueueName:          conf.PrivateExposeQueueName,
			QueueURL:           conf.PrivateExposeQueueURL,
			MaxNumberOfMessage: conf.MaxNumberOfMessage,
			WaitTimeSecond:     conf.WaitTimeSecond,
		},
		Log:       appLogger,
		SqsClient: sqsClient,
	}
}
