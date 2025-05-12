package controllers

import (
	"context"
	"errors"
	"route256/loms/internal/domain/model"
	orders_v1 "route256/loms/pkg/api/orders/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *model.CreateOrderModel) (int64, error)
	OrderInfo(ctx context.Context, orderId int64) (*model.OrderModel, error)
	PayOrder(ctx context.Context, orderId int64) error
	CancelOrder(ctx context.Context, orderId int64) error
}

type OrderController struct {
	orders_v1.UnimplementedOrdersServiceServer
	service OrderService
}

func NewOrderController(service OrderService) *OrderController {
	return &OrderController{
		service: service,
	}
}

func (c *OrderController) CreateOrder(ctx context.Context,
	createRequest *orders_v1.CreateOrderRequest) (*orders_v1.CreateOrderResponse, error) {
	orderItems := make([]model.OrderItem, 0, len(createRequest.Items))
	for _, item := range createRequest.Items {
		orderItems = append(orderItems, model.OrderItem{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	order := &model.CreateOrderModel{
		UserId: createRequest.User,
		Items:  orderItems,
	}

	orderId, err := c.service.CreateOrder(ctx, order)
	if err != nil {
		var targetErr *model.ErrReservedStockFailed
		var statusMismatchErr *model.ErrOrderStatusMismatch
		if errors.As(err, &targetErr) || errors.As(err, &statusMismatchErr) {
			return nil, status.Errorf(codes.FailedPrecondition, "createOrder: %v", err)
		}
		return nil, status.Errorf(codes.Unknown, "createOrder: %v", err)
	}

	return &orders_v1.CreateOrderResponse{
		OrderId: orderId,
	}, nil
}

func (c *OrderController) OrderInfo(ctx context.Context,
	infoRequest *orders_v1.OrderInfoRequest) (*orders_v1.OrderInfoResponse, error) {
	orderInfo, err := c.service.OrderInfo(ctx, infoRequest.OrderId)
	if err != nil {
		var notFoundErr *model.ErrOrderNotFound
		if errors.As(err, &notFoundErr) {
			return nil, status.Errorf(codes.NotFound, "orderInfo: %v", notFoundErr)
		}

		return nil, status.Errorf(codes.Unknown, "orderInfo: %v", err)
	}

	var orderItems = make([]*orders_v1.OrderItem, 0, len(orderInfo.Items))
	for _, item := range orderInfo.Items {
		orderItems = append(orderItems, &orders_v1.OrderItem{
			Sku:   item.Sku,
			Count: item.Count,
		})
	}

	return &orders_v1.OrderInfoResponse{
		User:   orderInfo.UserId,
		Status: orderInfo.Status,
		Items:  orderItems,
	}, nil
}

func (c *OrderController) PayOrder(ctx context.Context,
	orderInfo *orders_v1.PayOrderRequest) (*orders_v1.PayOrderResponse, error) {
	err := c.service.PayOrder(ctx, orderInfo.OrderId)
	if err != nil {
		var notFoundErr *model.ErrOrderNotFound
		if errors.As(err, &notFoundErr) {
			return nil, status.Errorf(codes.NotFound, "payOrder: %v", notFoundErr)
		}

		var invalidStatusErr *model.ErrInvalidOrderStatus
		var statusMismatchErr *model.ErrOrderStatusMismatch
		if errors.As(err, &invalidStatusErr) || errors.As(err, &statusMismatchErr) {
			return nil, status.Errorf(codes.FailedPrecondition, "payOrder: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "payOrder: %v", err)
	}

	return &orders_v1.PayOrderResponse{}, nil
}

func (c *OrderController) CancelOrder(ctx context.Context,
	orderInfo *orders_v1.CancelOrderRequest) (*orders_v1.CancelOrderResponse, error) {
	err := c.service.CancelOrder(ctx, orderInfo.OrderId)
	if err != nil {
		var notFoundErr *model.ErrOrderNotFound
		if errors.As(err, &notFoundErr) {
			return nil, status.Errorf(codes.NotFound, "cancelOrder: %v", err)
		}

		var invalidStatusErr *model.ErrInvalidOrderStatus
		var statusMismatchErr *model.ErrOrderStatusMismatch
		if errors.As(err, &invalidStatusErr) || errors.As(err, &statusMismatchErr) {
			return nil, status.Errorf(codes.FailedPrecondition, "cancel order: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "cancelOrder: %v", err)
	}

	return &orders_v1.CancelOrderResponse{}, nil
}
