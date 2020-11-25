package main

import (
	"encoding/json"
	"fmt"

	"github.com/CyberAgent/mimosa-osint/pkg/message"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/kelseyhightower/envconfig"
)

type sqsConfig struct {
	AWSRegion string `envconfig:"aws_region" default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://localhost:9324"`

	PrivateExposeQueueURL string `split_words:"true" required:"true"`
	SubdomainQueueURL     string `split_words:"true" required:"true"`
}

type sqsAPI interface {
	send(msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error)
}

type sqsClient struct {
	svc         *sqs.SQS
	queueURLMap map[string]string
}

func newSQSClient() *sqsClient {
	var conf sqsConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}
	session := sqs.New(session.New(), &aws.Config{
		Region:   &conf.AWSRegion,
		Endpoint: &conf.Endpoint,
	})

	return &sqsClient{
		svc: session,
		queueURLMap: map[string]string{
			// queueURLMap:
			"osint:private-expose": conf.PrivateExposeQueueURL,
			"osint:subdomain":      conf.SubdomainQueueURL,
		},
	}
}

func (s *sqsClient) send(msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error) {
	url := s.queueURLMap[msg.DataSource]
	if url == "" {
		return nil, fmt.Errorf("Unknown data_source, value=%s", msg.DataSource)
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		appLogger.Errorf("Failed to parse message error: %v", err)
		return nil, fmt.Errorf("Failed to parse message, err=%+v", err)
	}
	resp, err := s.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody:  aws.String(string(buf)),
		QueueUrl:     &url,
		DelaySeconds: aws.Int64(1),
	})
	if err != nil {
		appLogger.Errorf("Failed to send message, error: %v ", err)
		return nil, err
	}
	return resp, nil
}
