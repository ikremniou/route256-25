package repository_test

import (
	"context"
	"route256/cart/internal/domain/cart/repository"
	"route256/cart/internal/domain/model"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// I will not add the t.Parallel in the inner tests for
// the sake of the correct result validation
func TestCartRepository_ConcurrentWrites(t *testing.T) {
	t.Parallel()
	const (
		numGoroutines = 100
		userId        = 1
	)
	repo := repository.NewCartRepository()

	t.Run("Should perform concurrent CreateItem", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := range numGoroutines {
			go func(index int) {
				defer wg.Done()
				item := &model.CartItemModel{
					UserId: userId,
					SkuId:  int64(index),
					Count:  1,
				}
				created, err := repo.CreateItem(context.Background(), item)
				require.True(t, created)
				require.NoError(t, err)
			}(i)
		}

		wg.Wait()

		items := repo.GetByUserId(context.Background(), userId)
		require.Len(t, items, numGoroutines)
	})

	t.Run("Should perform concurrent updates", func(t *testing.T) {
		const (
			singleSkuId    = 500
			updatesPerItem = 10
			expectedCount  = updatesPerItem
		)

		var wg sync.WaitGroup
		wg.Add(updatesPerItem)

		for range updatesPerItem {
			go func() {
				defer wg.Done()
				item := &model.CartItemModel{
					UserId: userId,
					SkuId:  singleSkuId,
					Count:  1,
				}
				_, err := repo.CreateItem(context.Background(), item)
				require.NoError(t, err)
			}()
		}

		wg.Wait()

		items := repo.GetByUserId(context.Background(), userId)
		for _, item := range items {
			if item.SkuId == singleSkuId {
				require.Equal(t, uint32(expectedCount), item.Count)
			}
		}
	})

	t.Run("Should perform concurrent deletions", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(numGoroutines / 2)

		for i := range numGoroutines / 2 {
			go func(index int) {
				defer wg.Done()
				repo.DeleteBySku(context.Background(), userId, int64(index))
			}(i)
		}

		wg.Wait()

		items := repo.GetByUserId(context.Background(), userId)
		require.Len(t, items, numGoroutines/2+1)
	})
}

func TestCartRepository_ConcurrentReads(t *testing.T) {
	t.Parallel()

	repo := repository.NewCartRepository()

	t.Run("Should handle multiple concurrent readers", func(t *testing.T) {
		const (
			numItems   = 100
			userId     = 1
			numReaders = 200
		)

		for i := range numItems {
			item := &model.CartItemModel{
				UserId: userId,
				SkuId:  int64(i),
				Count:  uint32(i + 1),
			}
			_, err := repo.CreateItem(context.Background(), item)
			require.NoError(t, err)
		}

		var wg sync.WaitGroup
		wg.Add(numReaders)

		for range numReaders {
			go func() {
				defer wg.Done()
				items := repo.GetByUserId(context.Background(), userId)
				require.Len(t, items, numItems)
			}()
		}

		wg.Wait()
	})

	t.Run("Should handle concurrent sorted reads", func(t *testing.T) {
		const (
			userId     = 2
			numItems   = 50
			numReaders = 100
		)

		for i := numItems - 1; i >= 0; i-- {
			item := &model.CartItemModel{
				UserId: userId,
				SkuId:  int64(i),
				Count:  uint32(i + 1),
			}
			_, err := repo.CreateItem(context.Background(), item)
			require.NoError(t, err)
		}

		var wg sync.WaitGroup
		wg.Add(numReaders)

		for range numReaders {
			go func() {
				defer wg.Done()
				sortedItems := repo.GetAllOrderBySku(context.Background(), userId)
				require.Len(t, sortedItems, numItems)

				for j := 1; j < len(sortedItems); j++ {
					require.LessOrEqual(t, sortedItems[j-1].SkuId, sortedItems[j].SkuId)
				}
			}()
		}

		wg.Wait()
	})
}

func TestCartRepository_ConcurrentReadWrites(t *testing.T) {
	t.Parallel()

	repo := repository.NewCartRepository()

	t.Run("Should handle mixed read/write operations", func(t *testing.T) {
		const (
			numUsers   = 5
			numItems   = 20
			numWorkers = 100
			operations = 1000
		)

		for u := 1; u <= numUsers; u++ {
			for i := range numItems {
				item := &model.CartItemModel{
					UserId: int64(u),
					SkuId:  int64(i),
					Count:  1,
				}
				_, err := repo.CreateItem(context.Background(), item)
				require.NoError(t, err)
			}
		}

		var wg sync.WaitGroup
		wg.Add(numWorkers)

		for w := range numWorkers {
			go func(workerID int) {
				defer wg.Done()

				for i := range operations / numWorkers {
					// Determine operation type: 0=read, 1=write, 2=delete
					opType := (workerID + i) % 3
					userID := int64((workerID+i)%numUsers + 1)

					switch opType {
					case 0: // Read
						items := repo.GetByUserId(context.Background(), userID)
						if items != nil {
							// Just accessing the items is enough for testing concurrency
							_ = len(items)
						}
					case 1: // Write
						skuID := int64((workerID + i) % numItems)
						item := &model.CartItemModel{
							UserId: userID,
							SkuId:  skuID,
							Count:  1,
						}

						// Just accessing the items is enough for testing concurrency
						_, _ = repo.CreateItem(context.Background(), item)
					case 2: // Delete
						if i%2 == 0 {
							skuID := int64((workerID + i) % numItems)
							repo.DeleteBySku(context.Background(), userID, skuID)
						} else {
							if i%(operations/10) == 0 {
								repo.DeleteAll(context.Background(), userID)
							}
						}
					}
				}
			}(w)
		}

		wg.Wait()
	})

	t.Run("Should handle concurrent reads during high-frequency writes", func(t *testing.T) {
		const (
			userId    = 100
			writeOps  = 50
			readOps   = 200
			itemCount = 5
		)

		var wg sync.WaitGroup
		wg.Add(writeOps + readOps)

		for i := range writeOps {
			go func(index int) {
				defer wg.Done()
				item := &model.CartItemModel{
					UserId: userId,
					SkuId:  int64(index % itemCount),
					Count:  1,
				}
				repo.CreateItem(context.Background(), item)
			}(i)
		}

		for range readOps {
			go func() {
				defer wg.Done()
				_ = repo.GetByUserId(context.Background(), userId)
				_ = repo.GetAllOrderBySku(context.Background(), userId)
			}()
		}

		wg.Wait()

		items := repo.GetByUserId(context.Background(), userId)
		require.LessOrEqual(t, len(items), itemCount, "Should have at most itemCount different items")
	})
}
