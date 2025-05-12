package stock_repository

import (
	"context"
	"embed"
	"encoding/json"
	"sync"

	"route256/loms/internal/domain/model"
	"route256/loms/internal/infra/logger"
)

//go:embed stock-data.json
var stockData embed.FS

type StockRepository struct {
	mtx    sync.RWMutex
	Stocks map[int64]*model.StockModel
}

func NewStockRepository() *StockRepository {
	file, err := stockData.ReadFile("stock-data.json")
	if err != nil {
		logger.Fatal("Failed to read stock data file", "error", err)
	}

	var stocks []model.StockModel
	if err := json.Unmarshal(file, &stocks); err != nil {
		logger.Fatal("Failed to unmarshal stock data", "error", err)
	}

	stockMap := make(map[int64]*model.StockModel)
	for _, stock := range stocks {
		stockMap[stock.Sku] = &stock
	}

	return &StockRepository{Stocks: stockMap}
}

// RemoveReserved implements order_service.StockRepository.
func (o *StockRepository) RemoveReserved(_ context.Context, items []model.OrderItem) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	if err := o.validateStocksDecreaseCapacity(items); err != nil {
		return err
	}

	for _, item := range items {
		if stock, ok := o.Stocks[item.Sku]; ok {
			stock.TotalCount = stock.TotalCount - item.Count
			stock.Reserved = stock.Reserved - item.Count
		}
	}

	return nil
}

// CancelReserve implements order_service.StockRepository.
func (o *StockRepository) CancelReserved(_ context.Context, items []model.OrderItem) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	if err := o.validateStocksDecreaseCapacity(items); err != nil {
		return err
	}

	for _, item := range items {
		if stock, ok := o.Stocks[item.Sku]; ok {
			stock.Reserved = stock.Reserved - item.Count
		}
	}

	return nil
}

// Reserve implements order_service.StockRepository.
func (o *StockRepository) Reserve(_ context.Context, items []model.OrderItem) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	for _, item := range items {
		stock, ok := o.Stocks[item.Sku]
		if !ok {
			return &model.ErrStockNotFound{Sku: item.Sku}
		} else if uint64(stock.Reserved)+uint64(item.Count) > uint64(stock.TotalCount) {
			return &model.ErrStockOutOfBounds{
				Sku:        item.Sku,
				Reserved:   stock.Reserved,
				TotalCount: stock.TotalCount,
				Change:     item.Count,
			}
		}
	}

	for _, item := range items {
		if stock, ok := o.Stocks[item.Sku]; ok {
			stock.Reserved = stock.Reserved + item.Count
		}
	}

	return nil
}

// GetBySkuId implements stock_service.StockRepository.
func (o *StockRepository) GetBySkuId(_ context.Context, sku int64) (uint32, error) {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	if stock, ok := o.Stocks[sku]; ok {
		return stock.TotalCount - stock.Reserved, nil
	}

	return 0, &model.ErrStockNotFound{Sku: sku}
}

func (o *StockRepository) validateStocksDecreaseCapacity(items []model.OrderItem) error {
	for _, item := range items {
		stock, ok := o.Stocks[item.Sku]
		if !ok {
			return &model.ErrStockNotFound{Sku: item.Sku}
		} else if stock.Reserved < item.Count {
			return &model.ErrStockOutOfBounds{
				Sku:        item.Sku,
				Reserved:   stock.Reserved,
				TotalCount: stock.TotalCount,
				Change:     item.Count,
			}
		}
	}

	return nil
}
