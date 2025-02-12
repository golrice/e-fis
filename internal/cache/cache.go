package cache

import (
	"sync"

	"github.com/golrice/e-fis/internal/cache/lru"
)

type cache struct {
	mu       sync.Mutex
	lru      *lru.Cache
	capacity int64
}

func NewCache(capacity int64) *cache {
	return &cache{
		mu:       sync.Mutex{},
		lru:      lru.New(capacity, nil),
		capacity: capacity,
	}
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.capacity, nil)
	}

	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		bv := v.(ByteView)
		return NewByteView(bv.ByteSlice()), ok
	}

	return
}
