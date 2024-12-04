// Пакет inmem реализует интерфейс storage.Cache и хранит данные в оперативной памяти.
// На выбор доступны алгоритмы LRU и LFU, для инвалидации данных используется TTL.
package inmem

import (
	"context"
	"errors"
	"sync"
	"time"
	"workshop-1/internal/storage"
)

var _ storage.Cache[struct{}, struct{}] = (*Cache[struct{}, struct{}])(nil)

type accessTracker[K comparable] interface {
	RecordAccess(key K, now time.Time)
	Evict() K
}

type CacheStrategy int

const (
	LRU CacheStrategy = iota
	LFU
)

type Cache[K comparable, V any] struct {
	mp            map[K]CacheItem[V]
	mx            sync.RWMutex
	ttl           time.Duration
	strategy      CacheStrategy
	capacity      int
	accessTracker accessTracker[K]
}

func New[K comparable, V any](ttl time.Duration, strategy CacheStrategy, capacity int) *Cache[K, V] {
	var tracker accessTracker[K]
	switch strategy {
	case LRU:
		tracker = newLRUTracker[K](capacity)
	case LFU:
		tracker = newLFUTracker[K](capacity)
	}

	return &Cache[K, V]{
		mp:            make(map[K]CacheItem[V]),
		ttl:           ttl,
		strategy:      strategy,
		capacity:      capacity,
		accessTracker: tracker,
	}
}

func (c *Cache[K, V]) Set(ctx context.Context, key K, value V, now time.Time) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if len(c.mp) >= c.capacity {
		evictedKey := c.accessTracker.Evict()
		delete(c.mp, evictedKey)
	}

	item := CacheItem[V]{
		val: value,
		exp: now.Add(c.ttl),
	}

	c.mp[key] = item
	c.accessTracker.RecordAccess(key, now)

	return nil
}

func (c *Cache[K, V]) Get(ctx context.Context, key K, now time.Time) (V, error) {
	c.mx.RLock()
	item, ok := c.mp[key]
	c.mx.RUnlock()

	if !ok || item.IsExpired(now) {
		var zero V
		return zero, errors.New("ошибка получения данных из кэша")
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	c.accessTracker.RecordAccess(key, now)

	c.mp[key] = item

	return item.Value(), nil
}
