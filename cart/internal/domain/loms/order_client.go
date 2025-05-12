package loms

import (
	"context"
	"fmt"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/cart_config"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/sre"
	orders_v1 "route256/cart/internal/pb/orders/v1"
	"time"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type OrderClient struct {
	client orders_v1.OrdersServiceClient
}

func NewOrderClient(config *cart_config.Config) *OrderClient {
	var address = fmt.Sprintf("%s:%s", config.Loms.Host, config.Loms.Port)
	grpcClient, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()), sre.WithTracingDial(config))
	if err != nil {
		logger.Fatal("Failed to create orders grpc client", "err", err)
	}

	realClient := orders_v1.NewOrdersServiceClient(grpcClient)
	return &OrderClient{client: realClient}
}

func (c *OrderClient) CreateOrder(ctx context.Context, userId int64, items []model.CartItemModel) (int64, error) {
	ctx, span := otel.Tracer("client").Start(ctx, "order_client.CreateOrder")
	defer span.End()

	mappedItems := make([]*orders_v1.OrderItem, 0, len(items))

	for index := 0; index < len(items); index++ {
		mappedItems = append(mappedItems, &orders_v1.OrderItem{
			Sku:   items[index].SkuId,
			Count: items[index].Count,
		})
	}

	startTime := time.Now()
	response, err := c.client.CreateOrder(ctx, &orders_v1.CreateOrderRequest{
		User:  userId,
		Items: mappedItems,
	})
	sre.TrackExternalRequest("loms_create_order", err, startTime)
	if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
		return 0, fmt.Errorf("%w: %w", model.ErrCreateOrderPreconditionFailed, err)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to create order, user %d: %w", userId, err)
	}

	return response.OrderId, nil
}
