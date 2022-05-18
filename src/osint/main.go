package main

import (
	"context"
	"fmt"
	"net"

	"github.com/ca-risken/common/pkg/profiler"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/osint/proto/osint"
	"github.com/gassara-kys/envconfig"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const (
	nameSpace   = "osint"
	serviceName = "osint"
)

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConfig struct {
	// osintConfig
	Port            string   `default:"18081"`
	LogLevel        string   `default:"debug" split_words:"true"`
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// db
	DBMasterHost     string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	DBMasterUser     string `split_words:"true" default:"hoge"`
	DBMasterPassword string `split_words:"true" default:"moge"`
	DBSlaveHost      string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	DBSlaveUser      string `split_words:"true" default:"hoge"`
	DBSlavePassword  string `split_words:"true" default:"moge"`

	DBSchema        string `required:"true"    default:"mimosa"`
	DBPort          int    `required:"true"    default:"3306"`
	DBLogMode       bool   `split_words:"true" default:"false"`
	DBMaxConnection int    `split_words:"true" default:"10"`

	// sqs
	AWSRegion string `envconfig:"aws_region"   default:"ap-northeast-1"`
	Endpoint  string `envconfig:"sqs_endpoint" default:"http://queue.middleware.svc.cluster.local:9324"`

	SubdomainQueueURL string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-subdomain"`
	WebsiteQueueURL   string `split_words:"true" required:"true" default:"http://queue.middleware.svc.cluster.local:9324/queue/osint-website"`

	// grpc
	CoreAddr string `required:"true" split_words:"true" default:"core.core.svc.cluster.local:8080"`
}

func main() {
	ctx := context.Background()
	var conf AppConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(ctx, err)
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

	service := &osintService{}
	dbConfig := &DBConfig{
		MasterHost:     conf.DBMasterHost,
		MasterUser:     conf.DBMasterUser,
		MasterPassword: conf.DBMasterPassword,
		SlaveHost:      conf.DBSlaveHost,
		SlaveUser:      conf.DBSlaveUser,
		SlavePassword:  conf.DBSlavePassword,
		Schema:         conf.DBSchema,
		Port:           conf.DBPort,
		LogMode:        conf.DBLogMode,
		MaxConnection:  conf.DBMaxConnection,
	}
	service.repository = newOsintRepository(dbConfig)
	sqsConfig := &SQSConfig{
		AWSRegion:         conf.AWSRegion,
		Endpoint:          conf.Endpoint,
		SubdomainQueueURL: conf.SubdomainQueueURL,
		WebsiteQueueURL:   conf.WebsiteQueueURL,
	}
	service.sqs = newSQSClient(sqsConfig)
	service.projectClient = newProjectClient(conf.CoreAddr)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Errorf(ctx, "Failed to Opening Port. error: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mimosarpc.LoggingUnaryServerInterceptor(appLogger),
				grpctrace.UnaryServerInterceptor())))
	osint.RegisterOsintServiceServer(server, service)

	reflection.Register(server) // enable reflection API
	appLogger.Infof(ctx, "Starting gRPC server, port: %v", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Errorf(ctx, "Failed to gRPC server, error: %v", err)
	}
}
