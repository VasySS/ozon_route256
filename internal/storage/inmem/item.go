package inmem

import "time"

type CacheItem[V any] struct {
	val V
	exp time.Time
}

func (c *CacheItem[V]) IsExpired(now time.Time) bool {
	return now.After(c.exp)
}

func (c *CacheItem[V]) Value() V {
	return c.val
}
