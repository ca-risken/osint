package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/osint/pkg/message"
)

type SQSConfig struct {
	AWSRegion string
	Endpoint  string

	SubdomainQueueURL string
	WebsiteQueueURL   string
}

type sqsAPI interface {
	send(ctx context.Context, msg *message.OsintQueueMessage) (*sqs.SendMessageOutput, error)
}

type sqsClient struct {
	svc         *sqs.SQS
	queueURLMap map[string]string
}

func newSQSClient(conf *SQSConfig) *sqsClient {
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

	return &sqsClient{
		svc: session,
		queueURLMap: map[string]string{
			// queueURLMap:
			"osint:subdomain": conf.SubdomainQueueURL,
			"osint:website":   conf.WebsiteQueueURL,
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
		appLogger.Errorf(ctx, "Failed to parse message error: %v", err)
		return nil, fmt.Errorf("Failed to parse message, err=%+v", err)
	}
	resp, err := s.svc.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageBody:  aws.String(string(buf)),
		QueueUrl:     &url,
		DelaySeconds: aws.Int64(1),
	})
	if err != nil {
		appLogger.Errorf(ctx, "Failed to send message, error: %v ", err)
		return nil, err
	}
	return resp, nil
}
