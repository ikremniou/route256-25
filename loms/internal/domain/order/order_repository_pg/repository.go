package order_repository_pg

import (
	"context"
	"errors"
	"fmt"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/domain/order/order_repository_pg/query"
	"route256/loms/internal/infra/loms_config"
	"route256/loms/internal/infra/sre"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

type OrderRepository struct {
	master  *pgxpool.Pool
	replica *pgxpool.Pool

	statusTopic string
}

func NewOrderRepository(master *pgxpool.Pool, replica *pgxpool.Pool, cfg *loms_config.Config) *OrderRepository {
	return &OrderRepository{
		master:      master,
		replica:     replica,
		statusTopic: cfg.Kafka.OrderTopic,
	}
}

// Create implements app.OrderRepository.
func (r *OrderRepository) Create(ctx context.Context, createOrder *model.CreateOrderModel) (int64, error) {
	ctx, span := otel.GetTracerProvider().Tracer("repo").Start(ctx, "order_repository.CreateOrder")
	defer span.End()

	var orderIdResult int64 = 0
	err := pgx.BeginTxFunc(ctx, r.master, pgx.TxOptions{}, func(tx pgx.Tx) error {
		var repository = query.New(tx)

		startTime := time.Now()
		orderId, err := repository.CreateOrder(ctx, query.CreateOrderParams{
			UserID: createOrder.UserId,
			Topic:  r.statusTopic,
			Status: model.OrderStatusNew,
		})
		sre.TrackDbRequest("order_create", "insert", err, startTime)
		if err != nil {
			return fmt.Errorf("failed to create db order: %w", err)
		}

		var itemsToInsert = make([]query.CreateOrderItemsParams, 0, len(createOrder.Items))
		for _, item := range createOrder.Items {
			itemsToInsert = append(itemsToInsert, query.CreateOrderItemsParams{
				OrderID:  orderId,
				Sku:      item.Sku,
				Quantity: int64(item.Count),
			})
		}

		_, err = repository.CreateOrderItems(ctx, itemsToInsert)
		if err != nil {
			return fmt.Errorf("failed to create db order items: %w", err)
		}

		orderIdResult = orderId
		return nil
	})

	return orderIdResult, err
}

// GetById implements app.OrderRepository.
func (r *OrderRepository) GetById(ctx context.Context, orderId int64) (*model.OrderModel, error) {
	ctx, span := otel.GetTracerProvider().Tracer("repo").Start(ctx, "order_repository.GetById")
	defer span.End()

	var repository = query.New(r.replica)

	startTime := time.Now()
	order, err := repository.GetByID(ctx, orderId)
	sre.TrackDbRequest("order_get_by_id", "select", err, startTime)

	if err != nil {
		return nil, fmt.Errorf("failed to get db order: %w", err)
	}

	if len(order) == 0 {
		return nil, &model.ErrOrderNotFound{OrderId: orderId}
	}

	var items = make([]model.OrderItem, len(order))
	for i, item := range order {
		items[i] = model.OrderItem{
			Sku:   item.Sku,
			Count: uint32(item.Quantity),
		}
	}

	return &model.OrderModel{
		Id:     order[0].ID,
		UserId: order[0].UserID,
		Status: order[0].Status,
		Items:  items,
	}, nil
}

// UpdateStatus implements app.OrderRepository.
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderId int64, status string, expectStatus string) error {
	ctx, span := otel.GetTracerProvider().Tracer("repo").Start(ctx, "order_repository.UpdateStatus")
	defer span.End()

	var repository = query.New(r.master)

	startTime := time.Now()
	_, err := repository.UpdateStatus(ctx, query.UpdateStatusParams{
		ID:        orderId,
		Status:    status,
		Topic:     r.statusTopic,
		OldStatus: expectStatus,
	})
	sre.TrackDbRequest("order_update_status", "update", err, startTime)
	if err != nil {
		order, getErr := repository.GetByID(ctx, orderId)
		if getErr != nil {
			return fmt.Errorf("failed to compose error info: %w, %w", err, getErr)
		}

		if len(order) == 0 {
			return &model.ErrOrderNotFound{OrderId: orderId}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return &model.ErrOrderStatusMismatch{
				OrderId:       orderId,
				CurrentStatus: order[0].Status,
				ExpectedState: expectStatus,
			}
		}

		return err
	}

	return nil
}
