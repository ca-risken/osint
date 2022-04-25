package main

import (
	"context"
	"time"

	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/osint/proto/osint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newFindingClient(svcAddr string) finding.FindingServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		panic(err)
	}
	return finding.NewFindingServiceClient(conn)
}

func newAlertClient(svcAddr string) alert.AlertServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		panic(err)
	}
	return alert.NewAlertServiceClient(conn)
}

func newOsintClient(svcAddr string) osint.OsintServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		panic(err)
	}
	return osint.NewOsintServiceClient(conn)
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
