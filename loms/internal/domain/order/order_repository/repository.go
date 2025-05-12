package order_repository

import (
	"context"
	"route256/loms/internal/domain/model"
	"sort"
	"sync"
)

type OrderRepository struct {
	mtx       sync.RWMutex
	orders    map[int64]*model.OrderModel
	idCounter int64
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		mtx:    sync.RWMutex{},
		orders: make(map[int64]*model.OrderModel),
	}
}

// Create implements order_service.OrderRepository.
func (o *OrderRepository) Create(_ context.Context, createOrder *model.CreateOrderModel) (int64, error) {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	o.idCounter = o.idCounter + 1
	order := &model.OrderModel{
		Id:     o.idCounter,
		UserId: createOrder.UserId,
		Items:  createOrder.Items,
		Status: model.OrderStatusNew,
	}
	o.orders[o.idCounter] = order

	return o.idCounter, nil
}

// GetById implements order_service.OrderRepository.
func (o *OrderRepository) GetById(_ context.Context, orderId int64) (*model.OrderModel, error) {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	item, ok := o.orders[orderId]
	if !ok {
		return nil, &model.ErrOrderNotFound{OrderId: orderId}
	}

	sort.Slice(item.Items, func(i, j int) bool {
		return item.Items[i].Sku < item.Items[j].Sku
	})

	orderCopy := *item
	return &orderCopy, nil
}

// UpdateStatus implements order_service.OrderRepository.
func (o *OrderRepository) UpdateStatus(_ context.Context, orderId int64, status string, expertStatus string) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	order, ok := o.orders[orderId]
	if !ok {
		return &model.ErrOrderNotFound{OrderId: orderId}
	}

	if order.Status != expertStatus {
		return &model.ErrOrderStatusMismatch{
			OrderId:       orderId,
			CurrentStatus: order.Status,
			ExpectedState: expertStatus,
		}
	}

	order.Status = status
	return nil
}
