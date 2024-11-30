package conc

import (
	"context"
	"errors"
	"time"
)

var (
	ErrClosed = errors.New("closed")
)

type Batch[T any] struct {
	in        chan *future[T]
	errorFunc func(err error)
}

func NewBatch[T any]() *Batch[T] {
	return &Batch[T]{
		in: make(chan *future[T]),
	}
}

func NewBatchf[T any](errorFunc func(err error)) *Batch[T] {
	return &Batch[T]{
		in:        make(chan *future[T]),
		errorFunc: errorFunc,
	}
}

func (b *Batch[T]) CollectAndExec(ctx context.Context, maxBatchSize int, every time.Duration, f func(context.Context, []T) error) {
	if b == nil {
		return
	}
	var buf []*future[T]

	defer func() {
		for _, v := range buf {
			v.outErr = ErrClosed
			close(v.done)
		}
	}()

	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if len(buf) > 0 {
				go b.batchExecute(ctx, buf, f)
				buf = nil
			}
		case param := <-b.in:
			buf = append(buf, param)
			if len(buf) >= maxBatchSize {
				ticker.Reset(every)
				go b.batchExecute(ctx, buf, f)
				buf = nil
			}
		}
	}
}

func (b *Batch[T]) batchExecute(ctx context.Context, buf []*future[T], f func(context.Context, []T) error) {
	if b == nil {
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	params := make([]T, 0, len(buf))
	for _, v := range buf {
		params = append(params, v.in)
	}

	err := f(ctx, params)

	for _, v := range buf {
		v.outErr = err
		close(v.done)
	}
	if err != nil && b.errorFunc != nil {
		b.errorFunc(err)
	}
}

func (b *Batch[T]) Execute(ctx context.Context, t T) error {
	if b == nil {
		return nil
	}
	f := future[T]{
		in:   t,
		done: make(chan struct{}),
	}
	b.in <- &f
	return f.wait(ctx)
}

type future[T any] struct {
	in     T
	outErr error
	done   chan struct{}
}

func (f *future[T]) wait(ctx context.Context) error {
	select {
	case <-f.done:
		return f.outErr
	case <-ctx.Done():
		return ctx.Err()
	}
}
