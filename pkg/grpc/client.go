package grpc

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

func NewFindingClient(svcAddr string) (finding.FindingServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection: err=%w", err)
	}
	return finding.NewFindingServiceClient(conn), nil
}

func NewAlertClient(svcAddr string) (alert.AlertServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection: err=%w", err)
	}
	return alert.NewAlertServiceClient(conn), nil
}

func NewOsintClient(svcAddr string) (osint.OsintServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get GRPC connection: err=%w", err)
	}
	return osint.NewOsintServiceClient(conn), nil
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	// gRPCクライアントの呼び出し回数が非常に多くトレーシング情報の送信がエラーになるため、トレースは無効にしておく
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
