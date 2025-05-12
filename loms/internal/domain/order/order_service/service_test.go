package order_service_test

import (
	"context"
	"errors"
	"testing"

	"route256/loms/internal/domain/model"
	"route256/loms/internal/domain/order/order_service"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestOrderService_CreateOrder(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	type deps struct {
		orderRepo *OrderRepositoryMock
		stockRepo *StockRepositoryMock
	}
	tests := []struct {
		name    string
		deps    deps
		order   *model.CreateOrderModel
		want    int64
		wantErr bool
	}{
		{
			name: "should create order and reserve stocks successfully",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					CreateMock.Return(123, nil).
					UpdateStatusMock.Return(nil),
				stockRepo: NewStockRepositoryMock(mc).ReserveMock.Return(nil),
			},
			order: &model.CreateOrderModel{
				UserId: 1, Items: []model.OrderItem{{Sku: 1, Count: 1}},
			},
			want: 123, wantErr: false,
		},
		{
			name: "should return error if items is empty",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc),
				stockRepo: NewStockRepositoryMock(mc),
			},
			order: &model.CreateOrderModel{
				UserId: 1, Items: []model.OrderItem{},
			},
			want: 0, wantErr: true,
		},
		{
			name: "should return error if create order failed",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).CreateMock.Return(0, errors.New("error")),
				stockRepo: NewStockRepositoryMock(mc),
			},
			order: &model.CreateOrderModel{
				UserId: 1, Items: []model.OrderItem{{Sku: 1, Count: 1}},
			},
			want: 0, wantErr: true,
		},
		{
			name: "should update status to failed if reserve stock failed",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					CreateMock.Return(321, nil).
					UpdateStatusMock.When(minimock.AnyContext, 321, model.OrderStatusReserving, model.OrderStatusNew).
					Then(nil).UpdateStatusMock.When(minimock.AnyContext, 321, model.OrderStatusFailed, model.OrderStatusReserving).
					Then(nil),
				stockRepo: NewStockRepositoryMock(mc).ReserveMock.Return(errors.New("error")),
			},
			order: &model.CreateOrderModel{
				UserId: 1, Items: []model.OrderItem{{Sku: 1, Count: 1}},
			},
			want: 0, wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := order_service.NewOrderService(tt.deps.orderRepo, tt.deps.stockRepo)

			got, err := s.CreateOrder(context.Background(), tt.order)

			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("OrderService.CreateOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_OrderInfo(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	type deps struct {
		orderRepo *OrderRepositoryMock
	}
	tests := []struct {
		name    string
		deps    deps
		orderId int64
		want    *model.OrderModel
		wantErr bool
	}{
		{
			name: "should return order info successfully",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(&model.OrderModel{
					Id: 1, UserId: 1, Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{{Sku: 1, Count: 1}}}, nil),
			},
			orderId: 1,
			want: &model.OrderModel{
				Id: 1, UserId: 1, Status: model.OrderStatusAwaitingPayment,
				Items: []model.OrderItem{{Sku: 1, Count: 1}},
			},
			wantErr: false,
		},
		{
			name: "should return error if orderId is invalid",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc),
			},
			orderId: 0,
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error if GetById fails",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(nil, errors.New("error")),
			},
			orderId: 1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := order_service.NewOrderService(tt.deps.orderRepo, nil)
			got, err := s.OrderInfo(context.Background(), tt.orderId)

			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.OrderInfo() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestOrderService_PayOrder(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	type deps struct {
		orderRepo *OrderRepositoryMock
		stockRepo *StockRepositoryMock
	}
	tests := []struct {
		name    string
		deps    deps
		orderId int64
		wantErr bool
	}{
		{
			name: "should pay order successfully",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(&model.OrderModel{
					Id: 1, UserId: 1, Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{{Sku: 1, Count: 1}}}, nil).
					UpdateStatusMock.Return(nil),
				stockRepo: NewStockRepositoryMock(mc).
					RemoveReservedMock.Return(nil),
			},
			orderId: 1,
			wantErr: false,
		},
		{
			name: "should return error if orderId is invalid",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc),
				stockRepo: NewStockRepositoryMock(mc),
			},
			orderId: 0,
			wantErr: true,
		},
		{
			name: "should return error if GetById fails",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(nil, errors.New("error")),
				stockRepo: NewStockRepositoryMock(mc),
			},
			orderId: 1,
			wantErr: true,
		},
		{
			name: "should return error if and rollback status if remove reserved fails",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(&model.OrderModel{
					Id: 1, UserId: 1, Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{{Sku: 1, Count: 1}}}, nil).
					UpdateStatusMock.When(minimock.AnyContext, 1, model.OrderStatusPaying, model.OrderStatusAwaitingPayment).
					Then(nil).UpdateStatusMock.When(minimock.AnyContext, 1, model.OrderStatusAwaitingPayment, model.OrderStatusPaying).
					Then(nil),
				stockRepo: NewStockRepositoryMock(mc).
					RemoveReservedMock.Return(errors.New("error")),
			},
			orderId: 1,
			wantErr: true,
		},
		{
			name: "should return error if UpdateStatus fails",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(&model.OrderModel{
					Id: 1, UserId: 1, Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{{Sku: 1, Count: 1}}}, nil).
					UpdateStatusMock.When(minimock.AnyContext, 1, model.OrderStatusPaying, model.OrderStatusAwaitingPayment).
					Then(nil).UpdateStatusMock.When(minimock.AnyContext, 1, model.OrderStatusPayed, model.OrderStatusPaying).
					Then(errors.New("error")),
				stockRepo: NewStockRepositoryMock(mc).
					RemoveReservedMock.Return(nil),
			},
			orderId: 1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := order_service.NewOrderService(tt.deps.orderRepo, tt.deps.stockRepo)
			err := s.PayOrder(context.Background(), tt.orderId)

			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.PayOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderService_CancelOrder(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	type deps struct {
		orderRepo *OrderRepositoryMock
		stockRepo *StockRepositoryMock
	}
	tests := []struct {
		name    string
		deps    deps
		orderId int64
		wantErr bool
	}{
		{
			name: "should cancel order successfully",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(&model.OrderModel{
					Id: 1, UserId: 1, Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{{Sku: 1, Count: 1}},
				}, nil).
					UpdateStatusMock.Return(nil),
				stockRepo: NewStockRepositoryMock(mc).
					CancelReservedMock.Return(nil),
			},
			orderId: 1,
			wantErr: false,
		},
		{
			name: "should return error if orderId is invalid",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc),
				stockRepo: NewStockRepositoryMock(mc),
			},
			orderId: 0,
			wantErr: true,
		},
		{
			name: "should return error if GetById fails",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(nil, errors.New("error")),
				stockRepo: NewStockRepositoryMock(mc),
			},
			orderId: 1,
			wantErr: true,
		},
		{
			name: "should not return error if order is already cancelled",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(&model.OrderModel{
					Id: 1, UserId: 1, Status: model.OrderStatusCancelled,
					Items: []model.OrderItem{{Sku: 1, Count: 1}}}, nil),
				stockRepo: NewStockRepositoryMock(mc),
			},
			orderId: 1,
			wantErr: false,
		},
		{
			name: "should return error and rollback status if remove reserved fails",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(&model.OrderModel{
					Id: 1, UserId: 1, Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{{Sku: 1, Count: 1}}}, nil).
					UpdateStatusMock.When(minimock.AnyContext, 1, model.OrderStatusCancelling, model.OrderStatusAwaitingPayment).
					Then(nil).UpdateStatusMock.When(minimock.AnyContext, 1, model.OrderStatusAwaitingPayment, model.OrderStatusCancelling).
					Then(nil),
				stockRepo: NewStockRepositoryMock(mc).
					CancelReservedMock.Return(errors.New("error")),
			},
			orderId: 1,
			wantErr: true,
		},
		{
			name: "should return error if update status fails",
			deps: deps{
				orderRepo: NewOrderRepositoryMock(mc).
					GetByIdMock.Return(&model.OrderModel{
					Id: 1, UserId: 1, Status: model.OrderStatusAwaitingPayment,
					Items: []model.OrderItem{{Sku: 1, Count: 1}}}, nil).
					UpdateStatusMock.When(minimock.AnyContext, 1, model.OrderStatusCancelling, model.OrderStatusAwaitingPayment).
					Then(nil).UpdateStatusMock.When(minimock.AnyContext, 1, model.OrderStatusCancelled, model.OrderStatusCancelling).
					Then(errors.New("error")),
				stockRepo: NewStockRepositoryMock(mc).
					CancelReservedMock.Return(nil),
			},
			orderId: 1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := order_service.NewOrderService(tt.deps.orderRepo, tt.deps.stockRepo)
			err := s.CancelOrder(context.Background(), tt.orderId)

			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.CancelOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
