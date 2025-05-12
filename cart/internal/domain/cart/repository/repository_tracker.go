package repository

import (
	"context"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/sre"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func StartCollectingRepositoryStats(ctx context.Context, repo *CartRepository) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping repository stats collection")
			return
		case <-ticker.C:
			total := repo.getTotalItemCount()
			sre.InMemoryCartItems.With(prometheus.Labels{}).Set(float64(total))
		}
	}
}
