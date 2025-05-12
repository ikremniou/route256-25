package loms

import (
	"context"
	"fmt"
	"route256/cart/internal/infra/cart_config"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/sre"
	stocks_v1 "route256/cart/internal/pb/stocks/v1"
	"time"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StocksClient struct {
	client stocks_v1.StocksServiceClient
}

func NewStocksClient(config *cart_config.Config) *StocksClient {
	var address = fmt.Sprintf("%s:%s", config.Loms.Host, config.Loms.Port)
	grpcClient, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()), sre.WithTracingDial(config))
	if err != nil {
		logger.Fatal("Failed to create stocks grpc client", "err", err)
	}

	realClient := stocks_v1.NewStocksServiceClient(grpcClient)
	return &StocksClient{client: realClient}
}

func (c *StocksClient) StockInfo(ctx context.Context, skuId int64) (uint32, error) {
	ctx, span := otel.Tracer("client").Start(ctx, "stock_client.StockInfo")
	defer span.End()

	startTime := time.Now()
	response, err := c.client.StocksInfo(ctx, &stocks_v1.StocksInfoRequest{
		Sku: skuId,
	})
	sre.TrackExternalRequest("loms_stocks_info", err, startTime)
	if err != nil {
		return 0, fmt.Errorf("failed to get stock info for sku %d: %w", skuId, err)
	}

	return response.Count, nil
}
