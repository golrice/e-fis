package lru

import (
	"reflect"
	"testing"
)

// just for testing
type String string

func (d String) Len() int {
	return len(d)
}

func TestLru_Get(t *testing.T) {
	cache := New(int64(0), nil)

	cache.Add("key1", String("1234"))
	if v, ok := cache.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := cache.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestLru_RemoveByStrategy(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"

	cap := len(k1 + k2 + v1 + v2)

	cache := New(int64(cap), nil)
	cache.Add(k1, String(v1))
	cache.Add(k2, String(v2))
	// flush
	cache.Add(k3, String(v3))

	if _, ok := cache.Get("key1"); ok || cache.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestLru_OnRemove(t *testing.T) {
	keys := make([]string, 0)

	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"

	cap := len(k1 + k2 + v1 + v2)

	cache := New(int64(cap), callback)
	cache.Add(k1, String(v1))
	cache.Add(k2, String(v2))
	cache.Add(k3, String(v3))

	expect := []string{k1}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s, but we get %s", expect, keys)
	}
}
