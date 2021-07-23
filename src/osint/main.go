package main

import (
	"fmt"
	"net"

	mimosaxray "github.com/CyberAgent/mimosa-common/pkg/xray"
	"github.com/CyberAgent/mimosa-osint/proto/osint"
	"github.com/aws/aws-xray-sdk-go/xray"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := newOsintConfig()
	if err != nil {
		panic(err)
	}
	mimosaxray.InitXRay(xray.Config{})

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Errorf("Failed to Opening Port. error: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				xray.UnaryServerInterceptor(),
				mimosaxray.AnnotateEnvTracingUnaryServerInterceptor(conf.EnvName))))
	osintServer := newOsintService(conf.DB, conf.SQS)
	osint.RegisterOsintServiceServer(server, osintServer)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server, port: %v", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Errorf("Failed to gRPC server, error: %v", err)
	}
}
