package notifier

import (
	"context"
	"fmt"
	"route256/loms/internal/domain/notifier/query"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepository struct {
	master *pgxpool.Pool
}

type OutboxEntity struct {
	Key     string
	Payload []byte
	Topic   string
}

func NewOutboxRepository(master *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{
		master: master,
	}
}

func (r *OutboxRepository) ProcessPendingMessagesFn(
	ctx context.Context,
	batch int,
	predicate func([]OutboxEntity) error,
) error {
	err := pgx.BeginTxFunc(ctx, r.master, pgx.TxOptions{}, func(tx pgx.Tx) error {
		repository := query.New(tx)

		pendingRows, err := repository.GetPending(ctx, int32(batch))
		if err != nil {
			return fmt.Errorf("failed to get pending messages: %w", err)
		}

		entities := make([]OutboxEntity, 0, len(pendingRows))
		for _, pendingRow := range pendingRows {
			entities = append(entities, OutboxEntity{
				Key:     pendingRow.Key,
				Payload: pendingRow.Payload,
				Topic:   pendingRow.Topic,
			})
		}

		err = predicate(entities)

		return handleOutboxBatchStatusUpdate(ctx, repository, pendingRows, err)
	})

	if err != nil {
		return fmt.Errorf("failed process outbox messages in transaction: %w", err)
	}

	return nil
}

func handleOutboxBatchStatusUpdate(
	ctx context.Context,
	repository *query.Queries,
	pendingRows []query.Outbox,
	predicateErr error,
) error {
	status := "sent"
	if predicateErr != nil {
		status = "failed"
	}

	statusBatch := make([]query.UpdateRowStatusParams, 0, len(pendingRows))
	for _, pendingRow := range pendingRows {
		statusBatch = append(statusBatch, query.UpdateRowStatusParams{
			Status: status,
			ID:     pendingRow.ID,
		})
	}

	batchQuery := repository.UpdateRowStatus(ctx, statusBatch)

	var batchErr error = nil
	batchQuery.Exec(func(i int, err error) {
		if batchErr == nil && err != nil {
			batchErr = fmt.Errorf("failed to update pending message status: %w", err)
		}
	})

	return batchErr
}
