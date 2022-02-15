package main

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
	"github.com/gassara-kys/envconfig"
)

type AppConfig struct {
	EnvName string `default:"local" split_words:"true"`

	// sqs
	Debug string `default:"false"`

	AWSRegion   string `envconfig:"aws_region" default:"ap-northeast-1"`
	SQSEndpoint string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	WebsiteQueueName   string `split_words:"true" default:"osint-website"`
	WebsiteQueueURL    string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-website"`
	MaxNumberOfMessage int64  `split_words:"true" default:"5"`
	WaitTimeSecond     int64  `split_words:"true" default:"20"`

	// grpc
	FindingSvcAddr string `required:"true" split_words:"true" default:"finding.core.svc.cluster.local:8001"`
	AlertSvcAddr   string `required:"true" split_words:"true" default:"alert.core.svc.cluster.local:8004"`
	OsintSvcAddr   string `required:"true" split_words:"true" default:"osint.osint.svc.cluster.local:18081"`

	// wappalyzer
	WappalyzerPath string `required:"true" split_words:"true" default:"/opt/wappalyzer/src/drivers/npm/cli.js"`
}

func main() {
	var conf AppConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	err = mimosaxray.InitXRay(xray.Config{})
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	handler := &SQSHandler{}
	handler.findingClient = newFindingClient(conf.FindingSvcAddr)
	handler.alertClient = newAlertClient(conf.AlertSvcAddr)
	handler.osintClient = newOsintClient(conf.OsintSvcAddr)
	handler.wappalyzerPath = conf.WappalyzerPath

	sqsConf := &SQSConfig{
		Debug:              conf.Debug,
		AWSRegion:          conf.AWSRegion,
		SQSEndpoint:        conf.SQSEndpoint,
		WebsiteQueueName:   conf.WebsiteQueueName,
		WebsiteQueueURL:    conf.WebsiteQueueURL,
		MaxNumberOfMessage: conf.MaxNumberOfMessage,
		WaitTimeSecond:     conf.WaitTimeSecond,
	}
	consumer := newSQSConsumer(sqsConf)
	appLogger.Info("Start the website SQS consumer server...")
	ctx := context.Background()
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.StatusLoggingHandler(appLogger,
					mimosaxray.MessageTracingHandler(conf.EnvName, "osint.website", handler)))))
}
