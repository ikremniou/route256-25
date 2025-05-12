package mw

import (
	"context"
	"route256/loms/internal/infra/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Panic(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.Error("Panic in handler", "method", info.FullMethod, "error", e)
			err = status.Errorf(codes.Internal, "panic: %v", e)
		}
	}()
	return handler(ctx, req)
}
