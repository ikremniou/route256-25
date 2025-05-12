package mw

import (
	"context"
	"route256/loms/internal/infra/logger"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)
	duration := time.Since(start)
	st := status.Convert(err)
	code := st.Code()

	logger.Debug("Processed requests", "method", info.FullMethod, "duration", duration, "code", code.String())

	return resp, err
}
