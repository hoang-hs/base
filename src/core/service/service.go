package service

import (
	"context"
	"encoding/json"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/hoang-hs/base/src/common"
	log2 "github.com/hoang-hs/base/src/common/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Service struct {
}

func NewBaseService() *Service {
	return &Service{}
}

func (b *Service) connectGrpc(domain string) *grpc.ClientConn {
	conn, err := grpc.Dial(domain,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				otelgrpc.UnaryClientInterceptor(),
			),
		),
	)
	if err != nil {
		log2.Fatal("connect grpc error", log2.String("domain", domain), log2.Err(err))
	}
	return conn
}

func (b *Service) grpToIError(ctx context.Context, inputErr error) *common.Error {
	var ierr common.Error
	grpcErr, ok := status.FromError(inputErr)
	if !ok {
		return common.ErrSystemError(ctx, fmt.Sprintf("grpc error convert failed, err:[%s]", inputErr.Error()))
	}

	err := json.Unmarshal([]byte(grpcErr.Message()), &ierr)
	if err != nil {
		return common.ErrSystemError(ctx, fmt.Sprintf("grpc error unmarshal failed with input [%s]", inputErr.Error()))
	}

	return &ierr
}
