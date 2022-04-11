package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/common/pkg/profiler"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/osint/pkg/message"
	"github.com/gassara-kys/envconfig"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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

	// The number of host to be inspected at a time in goroutine
	InspectConcurrency int64 `split_words:"true" default:"50"`
}

func main() {
	var conf AppConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(conf.ProfileTypes)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(conf.ProfileExporter)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      conf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	defer pc.Stop()

	traceOpts := []tracer.StartOption{
		tracer.WithEnv(conf.EnvName),
		tracer.WithService(getFullServiceName()),
		tracer.WithDebugMode(conf.TraceDebug),
	}
	tracer.Start(traceOpts...)
	defer tracer.Stop()

	handler := &SQSHandler{}
	handler.harvesterConfig = newHarvesterConfig(conf.ResultPath, conf.HarvesterPath)
	appLogger.Info("Load Harvester Config")
	handler.inspectConcurrency = conf.InspectConcurrency
	appLogger.Info("Load Concurrency Config")
	handler.findingClient = newFindingClient(conf.FindingSvcAddr)
	appLogger.Info("Start Finding Client")
	handler.alertClient = newAlertClient(conf.AlertSvcAddr)
	appLogger.Info("Start Alert Client")
	handler.osintClient = newOsintClient(conf.OsintSvcAddr)
	appLogger.Info("Start Osint Client")
	f, err := mimosasqs.NewFinalizer(message.SubdomainDataSource, settingURL, conf.FindingSvcAddr, &mimosasqs.DataSourceRecommnend{
		ScanFailureRisk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", message.SubdomainDataSource),
		ScanFailureRecommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/osint/datasource/
		- For Domain type, make sure the FQDN format is registered.
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	})
	if err != nil {
		appLogger.Fatalf("Failed to create Finalizer, err=%+v", err)
	}

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
					messageTracingHandler(getFullServiceName(),
						f.FinalizeHandler(handler))))))
}

func messageTracingHandler(serviceName string, next mimosasqs.Handler) mimosasqs.Handler {
	return mimosasqs.HandlerFunc(func(ctx context.Context, msg *sqs.Message) error {
		span, tctx := tracer.StartSpanFromContext(ctx, serviceName)
		// TODO inherit trace from message
		defer span.Finish()

		err := next.HandleMessage(tctx, msg)
		return err
	})
}
