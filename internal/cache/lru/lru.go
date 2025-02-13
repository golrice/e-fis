package lru

import (
	"container/list"

	"github.com/golrice/e-fis/internal/cache/basic"
)

type LruCache basic.Cache

// lactual list element
type entry struct {
	key   string
	value basic.Value
}

func New(maxBytes int64, onRemove func(key string, value basic.Value)) *LruCache {
	return &LruCache{
		Mem: basic.MemInfo{
			MaxBytes:  maxBytes,
			UsedBytes: 0,
		},
		Bl:       list.New(),
		Cache:    make(map[string]*list.Element),
		OnRemove: onRemove,
	}
}

// we get the kv and change position according to the mem strategy
func (c *LruCache) Get(key string) (value basic.Value, ok bool) {
	if v, ok := c.Cache[key]; ok {
		c.Bl.MoveToFront(v)
		vv := v.Value.(*entry)
		return vv.value, true
	}
	return
}

func (c *LruCache) RemoveByStrategy() {
	item := c.Bl.Back()

	// nil if empty
	if item == nil {
		return
	}

	v := item.Value.(*entry)
	// we need to remove the item from list and flush mem and cache
	c.Bl.Remove(item)
	delete(c.Cache, v.key)
	c.Mem.UsedBytes -= int64(len(v.key)) + int64(v.value.Len())

	if c.OnRemove != nil {
		c.OnRemove(v.key, v.value)
	}
}

func (c *LruCache) Add(key string, value basic.Value) {
	// check whether the kv is in cache
	if e, ok := c.Cache[key]; ok {
		// in cache, update
		v := e.Value.(*entry)

		c.Bl.MoveToFront(e)
		v.value = value

		c.Mem.UsedBytes += int64(value.Len()) - int64(v.value.Len())
	} else {
		// if not in cache, add it in link & update cache
		e := c.Bl.PushFront(&entry{
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

func (c *LruCache) Delete(key string) {
	// check whether the kv is in cache
	e, ok := c.Cache[key]
	if !ok {
		return
	}

	v := e.Value.(*entry)

	c.Mem.UsedBytes -= int64(len(key)) + int64(v.value.Len())
	delete(c.Cache, key)
	c.Bl.Remove(e)
}

func (c *LruCache) Update(key string, value basic.Value) (ok bool) {
	// check whether the kv is in cache
	e, ok := c.Cache[key]
	if !ok {
		return
	}

	v := e.Value.(*entry)
	c.Mem.UsedBytes += int64(value.Len()) - int64(v.value.Len())
	v.value = value
	c.Bl.MoveToFront(e)

	if c.Mem.MaxBytes != 0 && c.Mem.MaxBytes < c.Mem.UsedBytes {
		c.RemoveByStrategy()
	}

	return
}

func (c *LruCache) Len() int {
	return c.Bl.Len()
}
