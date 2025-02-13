package lfu

import (
	"container/list"

	"github.com/golrice/e-fis/internal/cache/basic"
)

type LfuCache basic.Cache

// lactual list element
type entry struct {
	key     string
	value   basic.Value
	hitTime int
}

func New(maxBytes int64, onRemove func(key string, value basic.Value)) *LfuCache {
	return &LfuCache{
		Mem: basic.MemInfo{
			MaxBytes:  maxBytes,
			UsedBytes: 0,
		},
		Bl:       list.New(),
		Cache:    make(map[string]*list.Element),
		OnRemove: onRemove,
	}
}

func (c *LfuCache) Get(key string) (value basic.Value, ok bool) {
	if v, ok := c.Cache[key]; ok {
		vv := v.Value.(*entry)
		vv.hitTime += 1
		return vv.value, true
	}
	return
}

func (c *LfuCache) RemoveByStrategy() {
	target := c.Bl.Front()
	if target == nil {
		return
	}
	tarV := target.Value.(*entry)

	// fine the entry which has min hitTime
	for e := target.Next(); e != nil; e = e.Next() {
		ev := e.Value.(*entry)
		if tarV.hitTime > ev.hitTime {
			target = e
			tarV = ev
		}
	}

	// remove target
	c.Bl.Remove(target)
	delete(c.Cache, tarV.key)
	c.Mem.UsedBytes -= int64(len(tarV.key)) + int64(tarV.value.Len())

	if c.OnRemove != nil {
		c.OnRemove(tarV.key, tarV.value)
	}
}

func (c *LfuCache) Add(key string, value basic.Value) {
	// check whether the kv is in cache
	if e, ok := c.Cache[key]; ok {
		// in cache, update
		v := e.Value.(*entry)

		v.value = value

		c.Mem.UsedBytes += int64(value.Len()) - int64(v.value.Len())
	} else {
		// if not in cache, add it in link & update cache
		e := c.Bl.PushBack(&entry{
			key:     key,
			value:   value,
			hitTime: 0,
		})
		c.Cache[key] = e
		c.Mem.UsedBytes += int64(len(key)) + int64(value.Len())
	}

	// check our mem size
	for c.Mem.MaxBytes != 0 && c.Mem.MaxBytes < c.Mem.UsedBytes {
		c.RemoveByStrategy()
	}
}

func (c *LfuCache) Delete(key string) {
	e, ok := c.Cache[key]
	if !ok {
		return
	}

	v := e.Value.(*entry)

	c.Mem.UsedBytes -= int64(len(key)) + int64(v.value.Len())
	delete(c.Cache, key)
	c.Bl.Remove(e)
}

func (c *LfuCache) Update(key string, value basic.Value) (ok bool) {
	e, ok := c.Cache[key]
	if !ok {
		return
	}

	v := e.Value.(*entry)
	c.Mem.UsedBytes += int64(value.Len()) - int64(v.value.Len())
	v.value = value
	v.hitTime += 1
	c.Bl.MoveToFront(e)

	if c.Mem.MaxBytes != 0 && c.Mem.MaxBytes < c.Mem.UsedBytes {
		c.RemoveByStrategy()
	}

	return
}

func (c *LfuCache) Len() int {
	return c.Bl.Len()
}
