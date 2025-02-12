package cache

import (
	"sync"

	"github.com/golrice/e-fis/internal/cache/basic"
	"github.com/golrice/e-fis/internal/cache/fifo"
	"github.com/golrice/e-fis/internal/cache/lfu"
	"github.com/golrice/e-fis/internal/cache/lru"
)

type cache struct {
	mu       sync.Mutex
	bc       basic.BasicCache
	capacity int64
}

func NewCache(capacity int64, bc string) *cache {
	var bcache basic.BasicCache
	switch bc {
	case "lru":
		bcache = lru.New(capacity, nil)
	case "fifo":
		bcache = fifo.New(capacity, nil)
	case "lfu":
		bcache = lfu.New(capacity, nil)
	}
	return &cache{
		mu:       sync.Mutex{},
		bc:       bcache,
		capacity: capacity,
	}
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.bc == nil {
		c.bc = lru.New(c.capacity, nil)
	}

	c.bc.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.bc == nil {
		return
	}

	if v, ok := c.bc.Get(key); ok {
		bv := v.(ByteView)
		return NewByteView(bv.ByteSlice()), ok
	}

	return
}
