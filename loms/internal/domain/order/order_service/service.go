package order_service

import (
	"context"
	"fmt"
	"route256/loms/internal/domain/model"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
)

//go:generate minimock -i OrderRepository,StockRepository -p order_service_test,order_service_test
type OrderRepository interface {
	Create(ctx context.Context, model *model.CreateOrderModel) (int64, error)
	GetById(ctx context.Context, orderId int64) (*model.OrderModel, error)
	UpdateStatus(ctx context.Context, orderId int64, status string, expectStatus string) error
}

type StockRepository interface {
	Reserve(ctx context.Context, items []model.OrderItem) error
	RemoveReserved(ctx context.Context, items []model.OrderItem) error
	CancelReserved(ctx context.Context, items []model.OrderItem) error
}

type NotifierProducer interface {
	Publish(ctx context.Context, state model.OrderStateMessage) error
}

type OrderService struct {
	orderRepository OrderRepository
	stockRepository StockRepository
	validator       *validator.Validate
}

func NewOrderService(
	orderRepository OrderRepository,
	stockRepository StockRepository,
) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		stockRepository: stockRepository,
		validator:       validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *model.CreateOrderModel) (int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "order_service.CreateOrder")
	defer span.End()

	if err := s.validator.Struct(order); err != nil {
		return 0, fmt.Errorf("createOrder: failed to validate CreateOrderModel, %w", err)
	}

	orderId, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("createOrder: failed to create order, %w", err)
	}

	err = s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusReserving, model.OrderStatusNew)
	if err != nil {
		return 0, fmt.Errorf("createOrder: failed to set reserving status, %w", err)
	}

	err = s.stockRepository.Reserve(ctx, order.Items)
	if err != nil {
		if errS := s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusFailed, model.OrderStatusReserving); errS != nil {
			return 0, fmt.Errorf("createOrder: status update & reserve fail, %w, %w", err, errS)
		}

		return 0, &model.ErrReservedStockFailed{OrderId: orderId}
	}

	err = s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusAwaitingPayment, model.OrderStatusReserving)
	if err != nil {
		return 0, fmt.Errorf("createOrder: failed to set awaiting payment status, %w", err)
	}

	return orderId, nil
}

func (s *OrderService) OrderInfo(ctx context.Context, orderId int64) (*model.OrderModel, error) {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "order_service.OrderInfo")
	defer span.End()

	if orderId < 1 {
		return nil, fmt.Errorf("orderInfo: %w", &model.ErrInvalidOrderId{OrderId: orderId})
	}

	result, err := s.orderRepository.GetById(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("orderInfo: failed to get order by id, %w", err)
	}

	return result, nil
}

func (s *OrderService) PayOrder(ctx context.Context, orderId int64) error {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "order_service.PayOrder")
	defer span.End()

	if orderId < 1 {
		return fmt.Errorf("payOrder: %w", &model.ErrInvalidOrderId{OrderId: orderId})
	}

	order, err := s.orderRepository.GetById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("payOrder: failed to get order by id, %w", err)
	}

	if order.Status == model.OrderStatusPayed {
		return nil
	}

	err = s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusPaying, model.OrderStatusAwaitingPayment)
	if err != nil {
		return fmt.Errorf("payOrder: failed to update order status, %w", err)
	}

	err = s.stockRepository.RemoveReserved(ctx, order.Items)
	if err != nil {
		if errS := s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusAwaitingPayment, model.OrderStatusPaying); errS != nil {
			return fmt.Errorf("payOrder: status update & remove reserved fail, %w, %w", err, errS)
		}
		return fmt.Errorf("payOrder: failed to remove reserved stock, %w", err)
	}

	err = s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusPayed, model.OrderStatusPaying)
	if err != nil {
		return fmt.Errorf("payOrder: failed to update order status, %w", err)
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderId int64) error {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "order_service.CancelOrder")
	defer span.End()

	if orderId < 1 {
		return fmt.Errorf("cancelOrder: %w ", &model.ErrInvalidOrderId{OrderId: orderId})
	}

	order, err := s.orderRepository.GetById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("cancelOrder: failed to get order by id, %w", err)
	}

	if order.Status == model.OrderStatusCancelled {
		return nil
	}

	err = s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusCancelling, model.OrderStatusAwaitingPayment)
	if err != nil {
		return fmt.Errorf("cancelOrder: failed to update order status, %w", err)
	}

	err = s.stockRepository.CancelReserved(ctx, order.Items)
	if err != nil {
		if errS := s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusAwaitingPayment, model.OrderStatusCancelling); errS != nil {
			return fmt.Errorf("cancelOrder: status update & cancel reserved fail, %w, %w", err, errS)
		}

		return fmt.Errorf("cancelOrder: failed to cancel reserved stock, %w", err)
	}

	err = s.orderRepository.UpdateStatus(ctx, orderId, model.OrderStatusCancelled, model.OrderStatusCancelling)
	if err != nil {
		return fmt.Errorf("cancelOrder: failed to update order status, %w", err)
	}

	return nil
}
