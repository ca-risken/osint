package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/datasource-api/proto/osint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newFindingClient(ctx context.Context, svcAddr string) (finding.FindingServiceClient, error) {
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection: err=%w", err)
	}
	return finding.NewFindingServiceClient(conn), nil
}

func newAlertClient(ctx context.Context, svcAddr string) (alert.AlertServiceClient, error) {
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection: err=%w", err)
	}
	return alert.NewAlertServiceClient(conn), nil
}

func newOsintClient(ctx context.Context, svcAddr string) (osint.OsintServiceClient, error) {
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection: err=%w", err)
	}
	return osint.NewOsintServiceClient(conn), nil
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
