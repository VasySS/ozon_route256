package storage

import (
	"context"
	"errors"
	"time"
)

type Cache[K comparable, V any] interface {
	Set(ctx context.Context, key K, value V, now time.Time) error
	Get(ctx context.Context, key K, now time.Time) (V, error)
}

var _ Cache[struct{}, struct{}] = (*NoCache[struct{}, struct{}])(nil)

type NoCache[K comparable, V any] struct{}

var ErrNoCache = errors.New("no cache")

func (NoCache[K, V]) Set(ctx context.Context, key K, value V, now time.Time) error {
	return ErrNoCache
}

func (NoCache[K, V]) Get(ctx context.Context, key K, now time.Time) (V, error) {
	var zero V
	return zero, ErrNoCache
}
