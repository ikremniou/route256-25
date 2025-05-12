package order_repository_test

import (
	"context"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/domain/order/order_repository"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderRepository_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		args           []*model.CreateOrderModel
		userId         int64
		resultingOrder []*model.OrderModel
		wantId         []int64
		wantErr        []bool
	}{
		{
			name: "should create 2 separate orders under 1 userId successfully",
			args: []*model.CreateOrderModel{
				{
					UserId: 5,
					Items: []model.OrderItem{
						{Sku: 1, Count: 2}, {Sku: 2, Count: 3}, {Sku: 1, Count: 2},
					},
				},
				{
					UserId: 5,
					Items: []model.OrderItem{
						{Sku: 1, Count: 2}, {Sku: 2, Count: 3}, {Sku: 1, Count: 2},
					},
				},
			},
			resultingOrder: []*model.OrderModel{
				{
					Id: 1, UserId: 5, Status: model.OrderStatusNew,
					Items: []model.OrderItem{
						{Sku: 1, Count: 4},
						{Sku: 2, Count: 3},
					},
				},
			},
			wantId: []int64{1, 2}, wantErr: []bool{false, false},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := order_repository.NewOrderRepository()

			for index, createItem := range tt.args {
				gotId, err := repo.Create(context.Background(), createItem)

				if (err != nil) != tt.wantErr[index] {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				if gotId != tt.wantId[index] {
					t.Errorf("Create() got = %d, want %d", gotId, tt.wantId)
				}
			}
		})
	}
}

func TestOrderRepository_GetById(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		orderId       int64
		prepareOrders []*model.CreateOrderModel
		want          *model.OrderModel
		wantErr       bool
	}{
		{
			name:    "should get order successfully",
			orderId: 2,
			prepareOrders: []*model.CreateOrderModel{
				{
					UserId: 5,
					Items:  []model.OrderItem{{Sku: 1, Count: 2}},
				},
				{
					UserId: 5,
					Items:  []model.OrderItem{{Sku: 1, Count: 3}},
				},
			},
			want: &model.OrderModel{
				Id:     2,
				UserId: 5,
				Status: model.OrderStatusNew,
				Items:  []model.OrderItem{{Sku: 1, Count: 3}},
			},
			wantErr: false,
		},
		{
			name:          "should return error if no order found",
			orderId:       999,
			prepareOrders: []*model.CreateOrderModel{},
			want:          nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := order_repository.NewOrderRepository()

			for _, order := range tt.prepareOrders {
				repo.Create(context.Background(), order)
			}

			got, err := repo.GetById(context.Background(), tt.orderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetById() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.Equal(t, tt.want, got, "GetById() = %v, want %v", got, tt.want)
		})
	}
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	t.Parallel()

	type updates struct {
		OrderId      int64
		Status       string
		ExpectStatus string
	}

	tests := []struct {
		name          string
		updates       []updates
		prepareOrders []*model.CreateOrderModel
		want          []*model.OrderModel
		wantErr       []bool
	}{
		{
			name: "should update orders status successfully",
			prepareOrders: []*model.CreateOrderModel{
				{
					UserId: 5,
					Items:  []model.OrderItem{{Sku: 1, Count: 2}},
				},
				{
					UserId: 5,
					Items:  []model.OrderItem{{Sku: 1, Count: 3}},
				},
			},
			updates: []updates{
				{
					OrderId:      1,
					Status:       model.OrderStatusAwaitingPayment,
					ExpectStatus: model.OrderStatusNew,
				},
				{
					OrderId:      2,
					Status:       model.OrderStatusFailed,
					ExpectStatus: model.OrderStatusNew,
				},
			},
			want: []*model.OrderModel{
				{
					Id:     1,
					UserId: 5,
					Status: model.OrderStatusAwaitingPayment,
					Items:  []model.OrderItem{{Sku: 1, Count: 2}},
				},
				{
					Id:     2,
					UserId: 5,
					Status: model.OrderStatusFailed,
					Items:  []model.OrderItem{{Sku: 1, Count: 3}},
				},
			},
			wantErr: []bool{false, false},
		},
		{
			name:          "should return error if no order found when updating status",
			updates:       []updates{{OrderId: 999, Status: model.OrderStatusFailed}},
			prepareOrders: []*model.CreateOrderModel{},
			want:          []*model.OrderModel{},
			wantErr:       []bool{true},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := order_repository.NewOrderRepository()

			for _, order := range tt.prepareOrders {
				repo.Create(context.Background(), order)
			}

			for index, update := range tt.updates {
				err := repo.UpdateStatus(context.Background(), update.OrderId, update.Status, update.ExpectStatus)

				if (err != nil) != tt.wantErr[index] {
					t.Errorf("UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			for index, wantOrder := range tt.want {
				got, _ := repo.GetById(context.Background(), wantOrder.Id)

				require.Equal(t, tt.want[index], got, "UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
