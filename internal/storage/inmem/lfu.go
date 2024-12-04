package inmem

import (
	"container/heap"
	"time"
)

type lfuTracker[K comparable] struct {
	freqMap  map[K]int
	heap     *lfuHeap[K]
	capacity int
}

func newLFUTracker[K comparable](capacity int) *lfuTracker[K] {
	return &lfuTracker[K]{
		freqMap:  make(map[K]int),
		heap:     &lfuHeap[K]{},
		capacity: capacity,
	}
}

func (t *lfuTracker[K]) RecordAccess(key K, now time.Time) {
	t.freqMap[key]++
	heap.Push(t.heap, &lfuItem[K]{key: key, freq: t.freqMap[key], lastAccess: now})
}

func (t *lfuTracker[K]) Evict() K {
	evicted := heap.Pop(t.heap).(*lfuItem[K])
	delete(t.freqMap, evicted.key)

	return evicted.key
}

type lfuItem[K comparable] struct {
	key        K
	freq       int
	lastAccess time.Time
}

type lfuHeap[K comparable] []*lfuItem[K]

func (h lfuHeap[K]) Len() int {
	return len(h)
}

func (h lfuHeap[K]) Less(i, j int) bool {
	return h[i].freq < h[j].freq || (h[i].freq == h[j].freq && h[i].lastAccess.Before(h[j].lastAccess))
}

func (h lfuHeap[K]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *lfuHeap[K]) Push(x any) {
	*h = append(*h, x.(*lfuItem[K]))
}

func (h *lfuHeap[K]) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]

	return item
}
