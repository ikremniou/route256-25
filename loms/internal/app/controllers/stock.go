package controllers

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
	stocks_v1 "route256/loms/pkg/api/stocks/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StockService interface {
	StocksInfo(ctx context.Context, skuId int64) (uint32, error)
}

type StockController struct {
	stocks_v1.UnimplementedStocksServiceServer
	service StockService
}

func NewStocksController(service StockService) *StockController {
	return &StockController{
		service: service,
	}
}

func (c *StockController) StocksInfo(ctx context.Context, stockInfo *stocks_v1.StocksInfoRequest) (*stocks_v1.StocksInfoResponse, error) {
	count, err := c.service.StocksInfo(ctx, stockInfo.Sku)
	if err != nil {
		var targetErr *model.ErrStockNotFound
		if errors.As(err, &targetErr) {
			return nil, status.Errorf(codes.NotFound, "stocksInfo: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "stocksInfo: %v", err)
	}

	return &stocks_v1.StocksInfoResponse{
		Count: count,
	}, nil
}
