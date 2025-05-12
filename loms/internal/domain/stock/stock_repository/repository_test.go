package stock_repository_test

import (
	"context"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/domain/stock/stock_repository"
	"testing"

	"github.com/stretchr/testify/require"
)

const stocks = `[
	{"sku": 1, "total_count": 10, "reserved": 6},
	{"sku": 2, "total_count": 20, "reserved": 0},
	{"sku": 3, "total_count": 30, "reserved": 10}
]`

func TestStockRepository_RemoveReserved(t *testing.T) {
	t.Parallel()
	type args struct {
		items []model.OrderItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should remove reserved stock successfully",
			args: args{items: []model.OrderItem{
				{Sku: 1, Count: 1},
				{Sku: 3, Count: 2},
			}},
			wantErr: false,
		},
		{
			name: "should return error if stock is out of bounds",
			args: args{items: []model.OrderItem{
				{Sku: 2, Count: 1},
			}},
			wantErr: true,
		},
		{
			name: "should fail if stock not found",
			args: args{items: []model.OrderItem{
				{Sku: 5, Count: 1},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := stock_repository.NewStockRepositoryForTest(stocks)

			if err := o.RemoveReserved(context.Background(), tt.args.items); (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.RemoveReserved() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStockRepository_Reserve(t *testing.T) {
	t.Parallel()
	type args struct {
		items []model.OrderItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should reserve stock successfully",
			args: args{items: []model.OrderItem{
				{Sku: 1, Count: 2},
				{Sku: 3, Count: 5},
			}},
			wantErr: false,
		},
		{
			name: "should return error if stock is out of bounds",
			args: args{items: []model.OrderItem{
				{Sku: 1, Count: 20},
			}},
			wantErr: true,
		},
		{
			name: "should fail if any stock not found",
			args: args{items: []model.OrderItem{
				{Sku: 99, Count: 10},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := stock_repository.NewStockRepositoryForTest(stocks)

			if err := o.Reserve(context.Background(), tt.args.items); (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.Reserve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStockRepository_GetBySkuId(t *testing.T) {
	t.Parallel()
	type args struct {
		sku int64
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{
			name:    "should return available stock",
			args:    args{sku: 1},
			want:    4,
			wantErr: false,
		},
		{
			name:    "should return error if stock not found",
			args:    args{sku: 99},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := stock_repository.NewStockRepositoryForTest(stocks)

			got, err := o.GetBySkuId(context.Background(), tt.args.sku)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.GetBySkuId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OrderRepository.GetBySkuId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStockRepository_CancelReserved(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arg     model.OrderItem
		want    uint32
		wantErr bool
	}{
		{
			name:    "should return available stock",
			arg:     model.OrderItem{Sku: 1, Count: 4},
			want:    8,
			wantErr: false,
		},
		{
			name:    "should return error if stock out of bounds and don't change stocks",
			arg:     model.OrderItem{Sku: 1, Count: 40000000},
			want:    4,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := stock_repository.NewStockRepositoryForTest(stocks)

			err := o.CancelReserved(context.Background(), []model.OrderItem{tt.arg})
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.CancelReserved() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			count, err := o.GetBySkuId(context.Background(), tt.arg.Sku)
			require.NoError(t, err)

			if count != tt.want {
				t.Errorf("OrderRepository.CancelReserved() = %v, want %v", count, tt.want)
			}
		})
	}
}
