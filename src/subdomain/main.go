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

	// harvester
	ResultPath    string `required:"true" split_words:"true" default:"/results"`
	HarvesterPath string `required:"true" split_words:"true" default:"/theHarvester"`

	// grpc
	FindingSvcAddr string `required:"true" split_words:"true" default:"finding.core.svc.cluster.local:8001"`
	AlertSvcAddr   string `required:"true" split_words:"true" default:"alert.core.svc.cluster.local:8004"`
	OsintSvcAddr   string `required:"true" split_words:"true" default:"osint.osint.svc.cluster.local:18081"`

	// sqs
	AWSRegion string `envconfig:"aws_region"   default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	SubdomainQueueName string `split_words:"true" default:"osint-subdomain"`
	SubdomainQueueURL  string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-subdomain"`
	MaxNumberOfMessage int64  `split_words:"true" default:"3"`
	WaitTimeSecond     int64  `split_words:"true" default:"20"`
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
	handler.harvesterConfig = newHarvesterConfig(conf.ResultPath, conf.HarvesterPath)
	appLogger.Info("Load Harvester Config")
	handler.findingClient = newFindingClient(conf.FindingSvcAddr)
	appLogger.Info("Start Finding Client")
	handler.alertClient = newAlertClient(conf.AlertSvcAddr)
	appLogger.Info("Start Alert Client")
	handler.osintClient = newOsintClient(conf.OsintSvcAddr)
	appLogger.Info("Start Osint Client")

	sqsConf := &SQSConfig{
		AWSRegion:          conf.AWSRegion,
		Endpoint:           conf.Endpoint,
		SubdomainQueueName: conf.SubdomainQueueName,
		SubdomainQueueURL:  conf.SubdomainQueueURL,
		MaxNumberOfMessage: conf.MaxNumberOfMessage,
		WaitTimeSecond:     conf.WaitTimeSecond,
	}
	consumer := newSQSConsumer(sqsConf)
	appLogger.Info("Start the subdomain SQS consumer server...")
	ctx := context.Background()
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.StatusLoggingHandler(appLogger,
					mimosaxray.MessageTracingHandler(conf.EnvName, "osint.subDomain", handler)))))
}
