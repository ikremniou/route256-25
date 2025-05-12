package order_repository_pg_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"route256/loms/internal/domain/model"
	"route256/loms/internal/domain/order/order_repository_pg"
	"route256/loms/internal/infra/loms_config"
)

type OrderRepositorySuite struct {
	suite.Suite
	container      testcontainers.Container
	connectionPool *pgxpool.Pool
	repository     *order_repository_pg.OrderRepository
	ctx            context.Context
}

func (s *OrderRepositorySuite) SetupSuite() {
	s.ctx = context.Background()
	const (
		user     = "postgres"
		password = "postgres"
		db       = "test_db"
		dbWait   = 30
	)
	// create a docker container with a postgres database using test containers
	req := testcontainers.ContainerRequest{
		Image:        "gitlab-registry.ozon.dev/go/classroom-16/students/base/postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRESQL_PASSWORD": user,
			"POSTGRESQL_USERNAME": password,
			"POSTGRESQL_DATABASE": db,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(time.Second * dbWait),
	}

	container, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err, "Failed to start container")
	s.container = container

	mappedPort, err := container.MappedPort(s.ctx, "5432")
	s.Require().NoError(err, "Failed to get mapped port")

	host, err := container.Host(s.ctx)
	s.Require().NoError(err, "Failed to get container host")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, mappedPort.Port(), db)

	// run the migrations to populate the database using goose
	migrationConnection, err := sql.Open("pgx", dbURL)
	s.Require().NoError(err, "Failed to open DB connection for migrations")
	defer migrationConnection.Close()

	for range dbWait {
		err = migrationConnection.Ping()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	migrationsDir := filepath.Join("..", "..", "..", "..", "migrations")
	err = goose.Up(migrationConnection, migrationsDir)
	s.Require().NoError(err, "Failed to run migrations")

	// create a new connection pool backed on the docker container
	config, err := pgxpool.ParseConfig(dbURL)
	s.Require().NoError(err, "Failed to parse connection config")

	connectionPool, err := pgxpool.NewWithConfig(s.ctx, config)
	s.Require().NoError(err, "Failed to create connection pool")

	err = connectionPool.Ping(s.ctx)
	s.Require().NoError(err, "Failed to ping database")
	s.connectionPool = connectionPool

	s.repository = order_repository_pg.NewOrderRepository(connectionPool, connectionPool,
		&loms_config.Config{Kafka: loms_config.KafkaConfig{OrderTopic: "test_topic"}})
}

func (s *OrderRepositorySuite) TearDownSuite() {
	if s.connectionPool != nil {
		s.connectionPool.Close()
	}

	if s.container != nil {
		s.Require().NoError(s.container.Terminate(s.ctx), "Failed to terminate container")
	}
}

func (s *OrderRepositorySuite) TestOrderRepository_CreateOrder_Success() {
	orderId := s.createOrder()

	require.True(s.T(), orderId > 0, "Invalid order ID")
}

func (s *OrderRepositorySuite) TestOrderRepository_GetOrder_Success() {
	orderId := s.createOrder()

	order, err := s.repository.GetById(s.ctx, orderId)
	require.NoError(s.T(), err, "Failed to get order")
	require.NotNil(s.T(), order, "Order not found")
}

func (s *OrderRepositorySuite) TestOrderRepository_GetOrder_NotFound() {
	order, err := s.repository.GetById(s.ctx, 99999999)
	require.Error(s.T(), err, "Failed to get order")

	var notFount *model.ErrOrderNotFound
	require.True(s.T(), errors.As(err, &notFount), "Invalid error")
	require.Nil(s.T(), order, "Order found")
}

func (s *OrderRepositorySuite) TestOrderRepository_UpdateStatus_Success() {
	orderId := s.createOrder()
	err := s.repository.UpdateStatus(s.ctx, orderId, model.OrderStatusAwaitingPayment, model.OrderStatusNew)

	require.NoError(s.T(), err, "Failed to update status")
}

func (s *OrderRepositorySuite) TestOrderRepository_UpdateStatus_Mismatch() {
	orderId := s.createOrder()
	err := s.repository.UpdateStatus(s.ctx, orderId, model.OrderStatusAwaitingPayment, model.OrderStatusCancelled)

	var statusMismatch *model.ErrOrderStatusMismatch
	require.Error(s.T(), err, "Failed to update status")
	require.True(s.T(), errors.As(err, &statusMismatch), "Invalid error type")
}

func TestOrderRepository(t *testing.T) {
	t.Skip("Skipping this test as CI failing with docker")
	suite.Run(t, new(OrderRepositorySuite))
}

func (s *OrderRepositorySuite) createOrder() int64 {
	orderId, err := s.repository.Create(s.ctx, &model.CreateOrderModel{
		UserId: 1,
		Items: []model.OrderItem{
			{
				Sku:   100,
				Count: 1,
			},
			{
				Sku:   200,
				Count: 2,
			},
		},
	})

	require.NoError(s.T(), err, "Failed to create order")
	return orderId
}
