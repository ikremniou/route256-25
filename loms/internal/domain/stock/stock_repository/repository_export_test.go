package stock_repository

import (
	"encoding/json"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/infra/logger"
)

func NewStockRepositoryForTest(stocksAsString string) *StockRepository {
	var stocks []model.StockModel
	if err := json.Unmarshal([]byte(stocksAsString), &stocks); err != nil {
		logger.Fatal("Failed to unmarshal stock data", "error", err)
	}

	stockMap := make(map[int64]*model.StockModel)
	for _, stock := range stocks {
		stockMap[stock.Sku] = &stock
	}

	return &StockRepository{Stocks: stockMap}
}
