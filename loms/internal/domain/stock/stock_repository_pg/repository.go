package stock_repository_pg

import (
	"context"
	"errors"
	"fmt"
	"route256/loms/internal/domain/model"
	"route256/loms/internal/domain/stock/stock_repository_pg/query"
	"route256/loms/internal/infra/sre"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

const (
	ConstrainGteTotalReserved = "stocks_totalCount_gte_reserved"
	ConstrainTotalUint32      = "stocks_totalCount_uint32"
	ConstrainReservedUint32   = "stocks_reserved_uint32"
)

type StockRepository struct {
	master  *pgxpool.Pool
	replica *pgxpool.Pool
}

func NewOrderRepository(master *pgxpool.Pool, replica *pgxpool.Pool) *StockRepository {
	return &StockRepository{
		master:  master,
		replica: replica,
	}
}

// CancelReserved implements app.StockRepository.
func (r *StockRepository) CancelReserved(ctx context.Context, items []model.OrderItem) error {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "stock_repository.CancelReserved")
	defer span.End()

	return r.modifyReserved(ctx, items, true)
}

// GetBySkuId implements app.StockRepository.
func (r *StockRepository) GetBySkuId(ctx context.Context, sku int64) (uint32, error) {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "stock_repository.GetBySkuId")
	defer span.End()

	repository := query.New(r.replica)

	startTime := time.Now()
	stocks, err := repository.GetStocksBySkuId(ctx, sku)
	sre.TrackDbRequest("stock_get_by_sku", "select", err, startTime)
	if err != nil {
		return 0, fmt.Errorf("failed to get stocks by sku id: %w", err)
	}

	return uint32(stocks.TotalCount - stocks.Reserved), nil
}

// RemoveReserved implements app.StockRepository.
func (r *StockRepository) RemoveReserved(ctx context.Context, items []model.OrderItem) error {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "stock_repository.RemoveReserved")
	defer span.End()

	return pgx.BeginTxFunc(ctx, r.master, pgx.TxOptions{}, func(tx pgx.Tx) error {
		repository := query.New(tx)

		var requests = make([]query.RemoveReserveParams, 0, len(items))

		for _, item := range items {
			requests = append(requests, query.RemoveReserveParams{
				Sku:   item.Sku,
				Count: int64(item.Count),
			})
		}

		startTime := time.Now()
		result := repository.RemoveReserve(ctx, requests)
		var errResult error = nil
		result.Exec(func(i int, err error) {
			if errResult == nil && err != nil {
				errResult = handleRemoveReservedErr(err, i, items, repository)
			}
		})
		sre.TrackDbRequest("stock_remove_reserved", "update", errResult, startTime)

		return errResult
	})
}

// Reserve implements app.StockRepository.
func (r *StockRepository) Reserve(ctx context.Context, items []model.OrderItem) error {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "stock_repository.Reserve")
	defer span.End()

	return r.modifyReserved(ctx, items, false)
}

func (r *StockRepository) modifyReserved(ctx context.Context, items []model.OrderItem, cancel bool) error {
	var operationName = "stock_reserve_cancel"
	if !cancel {
		operationName = "stock_reserve"
	}

	return pgx.BeginTxFunc(ctx, r.master, pgx.TxOptions{}, func(tx pgx.Tx) error {
		repository := query.New(tx)

		var requests = make([]query.ReserveParams, 0, len(items))

		for _, item := range items {
			var modifier int64 = 1
			if cancel {
				modifier = -1
			}

			requests = append(requests, query.ReserveParams{
				Sku:   item.Sku,
				Count: int64(item.Count) * modifier,
			})
		}

		startTime := time.Now()
		result := repository.Reserve(ctx, requests)
		var errResult error = nil

		result.QueryRow(func(i int, updated int64, err error) {
			if errResult == nil && err != nil {
				errResult = handleReservedError(err, i, cancel, items, repository)
			}
		})
		sre.TrackDbRequest(operationName, "update", errResult, startTime)

		return errResult
	})
}

func handleRemoveReservedErr(err error, i int, items []model.OrderItem, queries *query.Queries) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return &model.ErrStockNotFound{Sku: items[i].Sku}
	}

	err = checkForKnownConstrainErr(err, &items[i], queries)
	return fmt.Errorf("failed to remove reserved in batch %w", err)
}

func handleReservedError(err error, i int, cancel bool, items []model.OrderItem, queries *query.Queries) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return &model.ErrStockNotFound{Sku: items[i].Sku}
	}

	err = checkForKnownConstrainErr(err, &items[i], queries)
	if cancel {
		return fmt.Errorf("failed to cancel reserved in batch %w", err)
	}

	return fmt.Errorf("failed to reserve in batch %w", err)
}

func checkForKnownConstrainErr(err error, item *model.OrderItem, queries *query.Queries) error {
	var pgErr *pgconn.PgError = nil
	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName == ConstrainGteTotalReserved ||
			pgErr.ConstraintName == ConstrainTotalUint32 ||
			pgErr.ConstraintName == ConstrainReservedUint32 {

			stock, err := queries.GetStocksBySkuId(context.Background(), item.Sku)
			if err != nil {
				return fmt.Errorf("%w, %w", &model.ErrStockOutOfBounds{
					Sku:        item.Sku,
					Reserved:   0,
					TotalCount: 0,
					Change:     item.Count,
				}, err)
			}

			return &model.ErrStockOutOfBounds{
				Sku:        item.Sku,
				Reserved:   uint32(stock.Reserved),
				TotalCount: uint32(stock.TotalCount),
				Change:     item.Count,
			}
		}
	}

	return err
}
