package route_err_group_test

import (
	"context"
	"errors"
	"fmt"
	"route256/cart/pkg/route_err_group"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestRouteErrGroup_Success_ShouldCorrectlyHandleRps(t *testing.T) {
	type someStruct struct {
		count int64
	}

	var ctx = context.Background()
	var errGroup = route_err_group.NewRouteErrorGroup[someStruct](ctx, route_err_group.Options{
		Rps:        5,
		BufferSize: 15,
	})
	var counter atomic.Int64

	for range 15 {
		errGroup.Run(func(ctx context.Context) (someStruct, error) {
			var count = counter.Add(1)
			return someStruct{count}, nil
		})
	}

	res, err := errGroup.Await()
	require.NoError(t, err, "Error is unexpected")

	require.Equal(t, len(res), 15)
	require.Contains(t, res, someStruct{count: 11})
}

func TestRouteErrGroup_Success_ShouldHandleBufferOverflow(t *testing.T) {
	type someStruct struct {
		count int64
	}

	var ctx = context.Background()
	var errGroup = route_err_group.NewRouteErrorGroup[someStruct](ctx, route_err_group.Options{
		Rps:        5,
		BufferSize: 15,
	})
	var counter atomic.Int64

	for range 20 {
		errGroup.Run(func(ctx context.Context) (someStruct, error) {
			var count = counter.Add(1)
			return someStruct{count}, nil
		})
	}

	_, err := errGroup.Await()
	require.Error(t, err, "Error should be thrown when static buffer is overflowed")
}

func TestRouteErrGroup_Success_ShouldCorrectlyHandleDynamicBuffer(t *testing.T) {
	t.Parallel()
	var ctx = context.Background()
	var errGroup = route_err_group.NewRouteErrorGroup[int64](ctx, route_err_group.Options{
		Rps: 5,
	})
	var concurrencyCounter atomic.Int64

	for range 50 {
		errGroup.Run(func(ctx context.Context) (int64, error) {
			var newValue = concurrencyCounter.Add(1)
			if newValue > 5 {
				return 0, fmt.Errorf("concurrency limit exceeded: %d", newValue)
			}

			defer concurrencyCounter.Add(-1)

			return 10, nil
		})
	}

	res, err := errGroup.Await()
	require.NoError(t, err, "Error is unexpected")
	require.Equal(t, len(res), 50)
}

func TestRouteErrGroup_Success_ShouldWorkWithEmptyOptions(t *testing.T) {
	t.Parallel()
	var ctx = context.Background()
	var errGroup = route_err_group.NewRouteErrorGroup[int64](ctx, route_err_group.Options{})
	var unsafeCounter = 0

	for range 50 {
		errGroup.Run(func(ctx context.Context) (int64, error) {
			unsafeCounter += 1
			return 10, nil
		})
	}

	result, err := errGroup.Await()
	require.NoError(t, err, "Error is unexpected")
	require.Equal(t, len(result), 50)
	require.Equal(t, unsafeCounter, 50)
}

func TestRouteErrGroup_Error_AndCancelOtherOperations(t *testing.T) {
	t.Parallel()
	var ctx = context.Background()
	var errGroup = route_err_group.NewRouteErrorGroup[int64](ctx, route_err_group.Options{Rps: 1})
	var counter atomic.Int64
	const numberOfOpsBeforeCancel = 5

	for range 10 {
		errGroup.Run(func(ctx context.Context) (int64, error) {
			var newNumber = counter.Add(1)
			if newNumber >= numberOfOpsBeforeCancel {
				return 0, errors.New("test error")
			}

			return 10, nil
		})
	}

	_, err := errGroup.Await()

	require.Error(t, err, "Error is required")
	require.GreaterOrEqual(t, counter.Load(), int64(numberOfOpsBeforeCancel))
}
