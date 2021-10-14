package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/osint/pkg/message"
	"github.com/gassara-kys/envconfig"
)

type sqsConfig struct {
	AWSRegion string `envconfig:"aws_region"   default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	SubdomainQueueURL string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-subdomain"`
}

type sqsAPI interface {
	send(ctx context.Context, msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error)
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
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}
	session := sqs.New(sess, &aws.Config{
		Region:   &conf.AWSRegion,
		Endpoint: &conf.Endpoint,
	})
	xray.AWS(session.Client)

	return &sqsClient{
		svc: session,
		queueURLMap: map[string]string{
			// queueURLMap:
			"osint:subdomain": conf.SubdomainQueueURL,
		},
	}
}

func (s *sqsClient) send(ctx context.Context, msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error) {
	url := s.queueURLMap[msg.DataSource]
	if url == "" {
		return nil, fmt.Errorf("Unknown data_source, value=%s", msg.DataSource)
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		appLogger.Errorf("Failed to parse message error: %v", err)
		return nil, fmt.Errorf("Failed to parse message, err=%+v", err)
	}
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
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
