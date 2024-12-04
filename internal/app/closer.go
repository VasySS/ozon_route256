package app

import (
	"context"
	"errors"
	"sync"
)

type closeFunc func(ctx context.Context) error

type Closer struct {
	mu    sync.Mutex
	funcs []closeFunc
}

func NewCloser() *Closer {
	return &Closer{
		funcs: make([]closeFunc, 0),
	}
}

func (c *Closer) AddWithCtx(f closeFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Add(f func()) {
	c.AddWithCtx(func(context.Context) error {
		f()
		return nil
	})
}

func (c *Closer) AddWithError(f func() error) {
	c.AddWithCtx(func(context.Context) error {
		return f()
	})
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var combinedErr error
	done := make(chan struct{})

	go func() {
		defer close(done)

		for i := len(c.funcs) - 1; i >= 0; i-- {
			if err := c.funcs[i](ctx); err != nil {
				errors.Join(combinedErr, err)
			}
		}
	}()

	select {
	case <-ctx.Done():
		return errors.New("вышло время ожидания, досрочное завершение работы...")
	case <-done:
		return combinedErr
	}
}
