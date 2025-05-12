package repository_test

import (
	"context"
	"route256/cart/internal/domain/cart/repository"
	"route256/cart/internal/domain/model"
	"testing"
)

func BenchmarkCreateItem(b *testing.B) {
	c := repository.NewCartRepository()
	item := &model.CartItemModel{UserId: 1, SkuId: 1, Count: 1}

	for i := 0; i < b.N; i++ {
		_, err := c.CreateItem(context.Background(), item)
		if err != nil && err != model.ErrTotalCountExceeded {
			b.Errorf("CreateItem() error = %v", err)
		}
	}
}

func BenchmarkGetAllOrderBySku(b *testing.B) {
	const NumberOfItems = 1000

	c := repository.NewCartRepository()
	userId := int64(1)
	for index := range NumberOfItems {
		item := &model.CartItemModel{UserId: userId, SkuId: int64(index), Count: 1}
		_, err := c.CreateItem(context.Background(), item)
		if err != nil {
			b.Errorf("CreateItem() error = %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.GetAllOrderBySku(context.Background(), userId)
	}
}
