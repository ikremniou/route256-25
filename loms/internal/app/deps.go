package app

import (
	"context"
	"fmt"
	"route256/loms/internal/app/controllers"
	"route256/loms/internal/domain/notifier"
	"route256/loms/internal/domain/order/order_repository"
	"route256/loms/internal/domain/order/order_repository_pg"
	"route256/loms/internal/domain/order/order_service"
	"route256/loms/internal/domain/stock/stock_repository"
	"route256/loms/internal/domain/stock/stock_repository_pg"
	"route256/loms/internal/domain/stock/stock_service"
	"route256/loms/internal/infra/logger"
	"route256/loms/internal/infra/loms_config"
	orders_v1 "route256/loms/pkg/api/orders/v1"
	stocks_v1 "route256/loms/pkg/api/stocks/v1"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
)

type OrderRepository interface {
	order_service.OrderRepository
}

type StockRepository interface {
	stock_service.StockRepository
	order_service.StockRepository
}

type NotifierProducer interface {
	Close() error
}

type Deps struct {
	notifier NotifierProducer
}

func InitializeDeps(ctx context.Context, grpcServer *grpc.Server, config *loms_config.Config) *Deps {
	var orderRepository OrderRepository
	var stockRepository StockRepository
	var notifyProducer NotifierProducer

	ctx, span := otel.GetTracerProvider().Tracer("initialize").Start(ctx, "initialize.deps")
	defer span.End()

	if config.Server.IsInMemory {
		orderRepository = order_repository.NewOrderRepository()
		stockRepository = stock_repository.NewStockRepository()
	} else {
		master, replica, err := connectToDatabases(ctx, config)
		if err != nil {
			logger.Fatal("Failed to connect to databases", "error", err)
		}

		orderRepository = order_repository_pg.NewOrderRepository(master, replica, config)
		stockRepository = stock_repository_pg.NewOrderRepository(master, replica)
		outboxRepository := notifier.NewOutboxRepository(master)

		notifyProducer = notifier.NewNotifierProducer(ctx, config, outboxRepository)
	}

	var orderService = order_service.NewOrderService(orderRepository, stockRepository)
	var stocksService = stock_service.NewStockService(stockRepository)

	orders_v1.RegisterOrdersServiceServer(grpcServer, controllers.NewOrderController(orderService))
	stocks_v1.RegisterStocksServiceServer(grpcServer, controllers.NewStocksController(stocksService))

	return &Deps{
		notifier: notifyProducer,
	}
}

func connectToDatabases(ctx context.Context, config *loms_config.Config) (*pgxpool.Pool, *pgxpool.Pool, error) {
	const addressTemplate = "postgresql://%s:%s@%s:%s/%s?sslmode=disable"

	masterCfg, err := pgxpool.ParseConfig(fmt.Sprintf(addressTemplate,
		config.MasterDb.User, config.MasterDb.Password, config.MasterDb.Host,
		config.MasterDb.Port, config.MasterDb.DbName))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse master db config: %w", err)
	}

	replicaCfg, err := pgxpool.ParseConfig(fmt.Sprintf(addressTemplate,
		config.ReplicaDb.User, config.ReplicaDb.Password, config.ReplicaDb.Host,
		config.ReplicaDb.Port, config.ReplicaDb.DbName))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse replica db config: %w", err)
	}

	logger.Info("Connecting to master db", "db", config.MasterDb.DbName)
	masterPool, err := pgxpool.NewWithConfig(ctx, masterCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to master db: %w", err)
	}

	logger.Info("Connecting to replica db", "db", config.ReplicaDb.DbName)
	replicaPool, err := pgxpool.NewWithConfig(ctx, replicaCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to replica db: %w", err)
	}

	_, span := otel.GetTracerProvider().Tracer("initialize").Start(ctx, "db.connect.master")
	if err := waitForDatabase(ctx, masterPool); err != nil {
		return nil, nil, err
	}
	span.End()

	_, span = otel.GetTracerProvider().Tracer("initialize").Start(ctx, "db.connect.replica")
	if err := waitForDatabase(ctx, replicaPool); err != nil {
		return nil, nil, err
	}
	span.End()

	return masterPool, replicaPool, nil
}

func waitForDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	times := 0
	for {
		err := pool.Ping(ctx)
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			return fmt.Errorf("context cancelled while waiting for db: %w", ctx.Err())
		}
		logger.Info("Waiting for db to be ready", "times", times)
		time.Sleep(time.Second)
		times += 1
	}

	return nil
}

func setupSre(ctx context.Context, config *loms_config.Config) error {
	tracingResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("loms"),
		),
	)

	if err != nil {
		return fmt.Errorf("failed to create jaeger resource: %w", err)
	}

	hostname := fmt.Sprintf("http://%s:%s", config.Jaeger.Host, config.Jaeger.Port)
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(hostname))
	if err != nil {
		return fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(tracingResource),
	)

	otel.SetTracerProvider(traceProvider)

	return nil
}
