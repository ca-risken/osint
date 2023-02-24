package main

import (
	"context"
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/common/pkg/profiler"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/datasource-api/pkg/message"
	"github.com/ca-risken/osint/pkg/grpc"
	"github.com/ca-risken/osint/pkg/sqs"
	"github.com/ca-risken/osint/pkg/website"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "osint"
	serviceName = "website"
	settingURL  = "https://docs.security-hub.jp/osint/datasource/"
)

var appLogger = logging.NewLogger()

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConfig struct {
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// sqs
	Debug string `default:"false"`

	AWSRegion   string `envconfig:"aws_region" default:"ap-northeast-1"`
	SQSEndpoint string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	OSINTWebsiteQueueName string `split_words:"true" default:"osint-website"`
	OSINTWebsiteQueueURL  string `split_words:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-website"`
	MaxNumberOfMessage    int32  `split_words:"true" default:"5"`
	WaitTimeSecond        int32  `split_words:"true" default:"20"`

	// grpc
	CoreAddr             string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
	DataSourceAPISvcAddr string `required:"true" split_words:"true" default:"datasource-api.core.svc.cluster.local:8081"`

	// wappalyzer
	WappalyzerPath string `required:"true" split_words:"true" default:"/opt/wappalyzer/src/drivers/npm/cli.js"`
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

	fc, err := grpc.NewFindingClient(conf.CoreAddr)
	if err != nil {
		appLogger.Fatalf(ctx, "failed to create finding client, err=%+v", err)
	}
	ac, err := grpc.NewAlertClient(conf.CoreAddr)
	if err != nil {
		appLogger.Fatalf(ctx, "failed to create alert client, err=%+v", err)
	}
	oc, err := grpc.NewOsintClient(conf.DataSourceAPISvcAddr)
	if err != nil {
		appLogger.Fatalf(ctx, "failed to create osint client, err=%+v", err)
	}
	wc := website.NewWappalyzerClient(conf.WappalyzerPath)
	handler := website.NewSQSHandler(fc, ac, oc, wc, appLogger)

	f, err := mimosasqs.NewFinalizer(message.WebsiteDataSource, settingURL, conf.CoreAddr, &mimosasqs.DataSourceRecommnend{
		ScanFailureRisk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", message.WebsiteDataSource),
		ScanFailureRecommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the network is reachable to the target host.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/osint/datasource/
		- For Website type, make sure the URL format(e.g. http://example.com ) is registerd.
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	})
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create Finalizer, err=%+v", err)
	}

	sqsConf := &sqs.SQSConfig{
		Debug:              conf.Debug,
		AWSRegion:          conf.AWSRegion,
		SQSEndpoint:        conf.SQSEndpoint,
		QueueName:          conf.OSINTWebsiteQueueName,
		QueueURL:           conf.OSINTWebsiteQueueURL,
		MaxNumberOfMessage: conf.MaxNumberOfMessage,
		WaitTimeSecond:     conf.WaitTimeSecond,
	}
	consumer, err := sqs.NewSQSConsumer(ctx, sqsConf, appLogger)
	if err != nil {
		appLogger.Fatalf(ctx, "Failed to create SQS consumer, err=%+v", err)
	}
	appLogger.Info(ctx, "Start the website SQS consumer server...")
	consumer.Start(ctx,
		mimosasqs.InitializeHandler(
			mimosasqs.RetryableErrorHandler(
				mimosasqs.TracingHandler(getFullServiceName(),
					mimosasqs.StatusLoggingHandler(appLogger,
						f.FinalizeHandler(handler))))))
}