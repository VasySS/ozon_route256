package inmem

import "time"

type lruTracker[K comparable] struct {
	keys    []K
	indices map[K]int
	ttl     time.Duration
}

func newLRUTracker[K comparable](capacity int) *lruTracker[K] {
	return &lruTracker[K]{
		keys:    make([]K, 0, capacity),
		indices: make(map[K]int, capacity),
	}
}

func (t *lruTracker[K]) RecordAccess(key K, now time.Time) {
	// если ключ уже присутствует в кэше, то его надо переместить в конец
	if idx, ok := t.indices[key]; ok {
		t.keys = append(t.keys[:idx], t.keys[idx+1:]...)
	}

	t.keys = append(t.keys, key)
	t.indices[key] = len(t.keys) - 1
}

func (t *lruTracker[K]) Evict() K {
	evictedKey := t.keys[0]
	// удаление самого старого элемента из кэша
	t.keys = t.keys[1:]
	delete(t.indices, evictedKey)

	return evictedKey
}
