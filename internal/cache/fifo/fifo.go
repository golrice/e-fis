package fifo

import (
	"container/list"

	"github.com/golrice/e-fis/internal/cache/basic"
)

type FifoCache basic.Cache

// lactual list element
type entry struct {
	key   string
	value basic.Value
}

func New(maxBytes int64, onRemove func(key string, value basic.Value)) *FifoCache {
	return &FifoCache{
		Mem: basic.MemInfo{
			MaxBytes:  maxBytes,
			UsedBytes: 0,
		},
		Bl:       list.New(),
		Cache:    make(map[string]*list.Element),
		OnRemove: onRemove,
	}
}

func (c *FifoCache) Get(key string) (value basic.Value, ok bool) {
	if v, ok := c.Cache[key]; ok {
		vv := v.Value.(*entry)
		return vv.value, true
	}
	return
}

func (c *FifoCache) RemoveByStrategy() {
	target := c.Bl.Front()
	if target == nil {
		return
	}
	tarV := target.Value.(*entry)

	// remove target
	c.Bl.Remove(target)
	delete(c.Cache, tarV.key)
	c.Mem.UsedBytes -= int64(len(tarV.key)) + int64(tarV.value.Len())

	if c.OnRemove != nil {
		c.OnRemove(tarV.key, tarV.value)
	}
}

func (c *FifoCache) Add(key string, value basic.Value) {
	// check whether the kv is in cache
	if e, ok := c.Cache[key]; ok {
		// in cache, update
		v := e.Value.(*entry)

		v.value = value

		c.Mem.UsedBytes += int64(value.Len()) - int64(v.value.Len())
	} else {
		// if not in cache, add it in link & update cache
		e := c.Bl.PushBack(&entry{
			key:   key,
			value: value,
		})
		c.Cache[key] = e
		c.Mem.UsedBytes += int64(len(key)) + int64(value.Len())
	}

	// check our mem size
	for c.Mem.MaxBytes != 0 && c.Mem.MaxBytes < c.Mem.UsedBytes {
		c.RemoveByStrategy()
	}
}

func (c *FifoCache) Len() int {
	return c.Bl.Len()
}
