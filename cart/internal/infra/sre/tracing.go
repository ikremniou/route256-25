package sre

import (
	"route256/cart/internal/infra/cart_config"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func WithTracingDial(_ *cart_config.Config) grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler())
}
