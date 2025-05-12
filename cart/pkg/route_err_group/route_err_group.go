package route_err_group

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrBufferExceeded = fmt.Errorf("err group buffer size exceeded")

type Options struct {
	Rps        int
	BufferSize int
}

type RouteErrGroup[T any] struct {
	wg       sync.WaitGroup
	sem      chan struct{}
	counter  atomic.Int64
	groupMtx sync.Mutex

	ctx        context.Context
	cancelFunc context.CancelCauseFunc

	err error

	isBuffered bool
	results    []T
}

func NewRouteErrorGroup[T any](ctx context.Context, options Options) *RouteErrGroup[T] {
	if options.Rps <= 0 {
		options.Rps = 1
	}

	if options.BufferSize < 0 {
		options.BufferSize = 0
	}

	ctx, cancel := context.WithCancelCause(ctx)
	return &RouteErrGroup[T]{
		wg:         sync.WaitGroup{},
		sem:        make(chan struct{}, options.Rps),
		results:    make([]T, options.BufferSize),
		ctx:        ctx,
		cancelFunc: cancel,
		err:        nil,
		isBuffered: options.BufferSize > 0,
	}
}

func (eg *RouteErrGroup[T]) Run(function func(ctx context.Context) (T, error)) {
	eg.wg.Add(1)
	var position = eg.counter.Add(1) - 1

	go func() {
		defer eg.wg.Done()

		select {
		case eg.sem <- struct{}{}:
		case <-eg.ctx.Done():
			return
		}
		defer func() { <-eg.sem }()

		result, err := function(eg.ctx)
		if err != nil {
			eg.handleError(err)
			return
		}

		eg.handleResult(result, position)
	}()
}

func (eg *RouteErrGroup[T]) Await() ([]T, error) {
	eg.wg.Wait()

	eg.cancelFunc(eg.err)
	return eg.results, eg.err
}

func (eg *RouteErrGroup[T]) handleResult(result T, position int64) {
	if eg.isBuffered {
		if position >= int64(cap(eg.results)) {
			eg.handleError(fmt.Errorf("%w, position: %d, cap: %d", ErrBufferExceeded, position, cap(eg.results)))

			return
		}

		eg.results[position] = result

		return
	}

	eg.groupMtx.Lock()
	defer eg.groupMtx.Unlock()
	eg.results = append(eg.results, result)
}

func (eg *RouteErrGroup[T]) handleError(err error) {
	eg.groupMtx.Lock()
	defer eg.groupMtx.Unlock()

	if eg.err == nil {
		eg.err = err
		eg.cancelFunc(err)
	}
}
