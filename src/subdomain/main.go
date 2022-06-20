package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/profiler"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "osint"
	serviceName = "subdomain"
	settingURL  = "https://docs.security-hub.jp/osint/datasource/"
)

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConfig struct {
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// harvester
	ResultPath    string `required:"true" split_words:"true" default:"/results"`
	HarvesterPath string `required:"true" split_words:"true" default:"/theHarvester"`

	// grpc
	CoreAddr             string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
	DataSourceAPISvcAddr string `required:"true" split_words:"true" default:"datasource-api.core.svc.cluster.local:8081"`

	// sqs
	AWSRegion string `envconfig:"aws_region"   default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	OSINTSubdomainQueueName string `split_words:"true" default:"osint-subdomain"`
	OSINTSubdomainQueueURL  string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-subdomain"`
	MaxNumberOfMessage      int32  `split_words:"true" default:"3"`
	WaitTimeSecond          int32  `split_words:"true" default:"20"`

	// The number of host to be inspected at a time in goroutine
	InspectConcurrency int64 `split_words:"true" default:"50"`
}

func main() {
	ctx := context.Background()
	var conf AppConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(conf.ProfileTypes)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(conf.ProfileExporter)
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      conf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		appLogger.Fatal(ctx, err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: conf.EnvName,
		Debug:       conf.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	handler := &SQSHandler{}
	handler.harvesterConfig = newHarvesterConfig(conf.ResultPath, conf.HarvesterPath)
	appLogger.Info(ctx, "Load Harvester Config")
	handler.inspectConcurrency = conf.InspectConcurrency
	appLogger.Info(ctx, "Load Concurrency Config")
	handler.findingClient = newFindingClient(conf.CoreAddr)
	appLogger.Info(ctx, "Start Finding Client")
	handler.alertClient = newAlertClient(conf.CoreAddr)
	appLogger.Info(ctx, "Start Alert Client")
	handler.osintClient = newOsintClient(conf.DataSourceAPISvcAddr)
	appLogger.Info(ctx, "Start Osint Client")
	f, err := mimosasqs.NewFinalizer(message.SubdomainDataSource, settingURL, conf.CoreAddr, &mimosasqs.DataSourceRecommnend{
		ScanFailureRisk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", message.SubdomainDataSource),
		ScanFailureRecommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/osint/datasource/
		- For Domain type, make sure the FQDN format is registered.
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	})
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create Finalizer, err=%+v", err)
	}

	sqsConf := &SQSConfig{
		AWSRegion:          conf.AWSRegion,
		Endpoint:           conf.Endpoint,
		QueueName:          conf.OSINTSubdomainQueueName,
		QueueURL:           conf.OSINTSubdomainQueueURL,
		MaxNumberOfMessage: conf.MaxNumberOfMessage,
		WaitTimeSecond:     conf.WaitTimeSecond,
	}
	consumer := newSQSConsumer(ctx, sqsConf)
	appLogger.Info(ctx, "Start the subdomain SQS consumer server...")
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.TracingHandler(getFullServiceName(),
					mimosasqs.StatusLoggingHandler(appLogger,
						f.FinalizeHandler(handler))))))
}
