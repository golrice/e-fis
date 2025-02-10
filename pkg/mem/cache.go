package pkg

import "container/list"

// build a simple cache, which is might not safe for concurrency still.
type Cache struct {
	// record the memory capatity and memory used
	mem memInfo
	// bidirectional linked list
	bl *list.List
	// interate with user
	cache map[string]*list.Element

	// callback
	onRemove func(key string, value Value)
}

// memory capatity and memory used
type memInfo struct {
	maxBytes  int64
	usedBytes int64
}

// lactual list element
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onRemove func(key string, value Value)) *Cache {
	return &Cache{
		mem: memInfo{
			maxBytes:  maxBytes,
			usedBytes: 0,
		},
		bl:       list.New(),
		cache:    make(map[string]*list.Element),
		onRemove: onRemove,
	}
}

// we get the kv and change position according to the mem strategy
func (c *Cache) Get(key string) (value Value, ok bool) {
	if v, ok := c.cache[key]; ok {
		c.bl.MoveToFront(v)
		vv := v.Value.(*entry)
		return vv.value, true
	}
	return
}

func (c *Cache) RemoveByStrategy() {
	item := c.bl.Back()

	// nil if empty
	if item == nil {
		return
	}

	v := item.Value.(*entry)
	// we need to remove the item from list and flush mem and cache
	c.bl.Remove(item)
	delete(c.cache, v.key)
	c.mem.usedBytes -= int64(len(v.key)) + int64(v.value.Len())

	if c.onRemove != nil {
		c.onRemove(v.key, v.value)
	}
}

func (c *Cache) Add(key string, value Value) {
	// check whether the kv is in cache
	if e, ok := c.cache[key]; ok {
		// in cache, update
		v := e.Value.(*entry)

		c.bl.MoveToFront(e)
		v.value = value

		c.mem.usedBytes += int64(value.Len()) - int64(v.value.Len())
	} else {
		// if not in cache, add it in link & update cache
		e := c.bl.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = e
		c.mem.usedBytes += int64(len(key)) + int64(value.Len())
	}

	// check our mem size
	for c.mem.maxBytes != 0 && c.mem.maxBytes < c.mem.usedBytes {
		c.RemoveByStrategy()
	}
}

func (c *Cache) Len() int {
	return c.bl.Len()
}
