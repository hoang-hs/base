package base

import (
	"context"
	"encoding/json"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/hoang-hs/base/log"
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
		log.Fatal("connect grpc error, domain:[%s], err:[%s]", domain, err.Error())
	}
	return conn
}

func (b *Service) grpToIError(ctx context.Context, inputErr error) *Error {
	var ierr Error
	grpcErr, ok := status.FromError(inputErr)
	if !ok {
		return ErrSystemError(ctx, fmt.Sprintf("grpc error convert failed, err:[%s]", inputErr.Error()))
	}

	err := json.Unmarshal([]byte(grpcErr.Message()), &ierr)
	if err != nil {
		return ErrSystemError(ctx, fmt.Sprintf("grpc error unmarshal failed with input [%s]", inputErr.Error()))
	}

	return &ierr
}
