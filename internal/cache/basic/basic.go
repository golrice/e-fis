package basic

import "container/list"

// build a simple cache, which is might not safe for concurrency still.
type Cache struct {
	// record the memory capatity and memory used
	Mem MemInfo
	// bidirectional linked list
	Bl *list.List
	// interate with user
	Cache map[string]*list.Element

	// callback
	OnRemove func(key string, value Value)
}

// memory capatity and memory used
type MemInfo struct {
	MaxBytes  int64
	UsedBytes int64
}

type Value interface {
	Len() int
}

type BasicCache interface {
	Get(key string) (value Value, ok bool)
	RemoveByStrategy()
	Add(key string, value Value)
	Len() int
}
