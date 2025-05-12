package stock_service

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
)

type StockRepository interface {
	GetBySkuId(ctx context.Context, sku int64) (uint32, error)
}

type StockService struct {
	repository StockRepository
}

func NewStockService(repository StockRepository) *StockService {
	return &StockService{repository: repository}
}

// StocksInfo implements controllers.StockService.
func (s *StockService) StocksInfo(ctx context.Context, skuId int64) (uint32, error) {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "stock_service.StocksInfo")
	defer span.End()

	if skuId < 1 {
		return 0, errors.New("SkuId must be greater than 0")
	}

	count, err := s.repository.GetBySkuId(ctx, skuId)
	if err != nil {
		return 0, fmt.Errorf("failed to get stocks by sku, %w", err)
	}

	return count, nil
}
