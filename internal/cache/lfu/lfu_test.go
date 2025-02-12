package lfu

import (
	"reflect"
	"testing"

	"github.com/golrice/e-fis/internal/cache/basic"
)

// just for testing
type String string

func (d String) Len() int {
	return len(d)
}

// TestLFU_Add 测试添加元素
func TestLFU_Add(t *testing.T) {
	cache := New(int64(0), nil)

	cache.Add("key1", String("1234"))
	if v, ok := cache.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed, got %v, %v", v, ok)
	}

	cache.Add("key2", String("5678"))
	if v, ok := cache.Get("key2"); !ok || string(v.(String)) != "5678" {
		t.Fatalf("cache hit key2=5678 failed, got %v, %v", v, ok)
	}
}

// TestLFU_Get 测试获取元素
func TestLFU_Get(t *testing.T) {
	cache := New(int64(0), nil)

	cache.Add("key1", String("1234"))
	cache.Get("key1") // increase hitTime

	cache.Add("key2", String("5678"))
	cache.Get("key2") // increase hitTime

	if v, ok := cache.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed, got %v, %v", v, ok)
	}

	if v, ok := cache.Get("key2"); !ok || string(v.(String)) != "5678" {
		t.Fatalf("cache hit key2=5678 failed, got %v, %v", v, ok)
	}
}

// TestLFU_RemoveByStrategy 测试根据策略移除元素
func TestLFU_RemoveByStrategy(t *testing.T) {
	k1, k2, k3 := "k1", "k2", "k3"
	v1, v2, v3 := "v1", "v2", "v3"

	cache := New(int64(len(k1+v1+k2+v2+k3+v3)), nil)
	cache.Add(k1, String(v1))
	cache.Get(k1) // increase hitTime for k1

	cache.Add(k2, String(v2))
	cache.Get(k2) // increase hitTime for k2
	cache.Get(k2) // increase hitTime for k2

	cache.Add(k3, String(v3))

	// Add another element, which should trigger the removal of the least frequently used item (k1)
	cache.Add("k4", String("v4"))

	_, ok1 := cache.Get(k1)
	_, ok2 := cache.Get(k2)
	_, ok3 := cache.Get(k3)
	if ok1 && ok2 && ok3 {
		t.Fatalf("RemoveByStrategy failed, key1/2/3 should not be removed")
	}
}

// TestLFU_OnRemove 测试回调函数
func TestLFU_OnRemove(t *testing.T) {
	keys := make([]string, 0)

	callback := func(key string, value basic.Value) {
		keys = append(keys, key)
	}

	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "1234", "5678", "v3"

	cache := New(int64(len(k1+k2+k3+v1+v2+v3)), callback)

	cache.Add(k1, String(v1))
	cache.Get(k1) // increase hitTime for k1

	cache.Add(k2, String(v2))
	cache.Get(k2) // increase hitTime for k2
	cache.Get(k2) // increase hitTime for k2

	cache.Add(k3, String(v3))
	// Add another element, which should trigger the removal of the least frequently used item (k3)
	cache.Add("k4", String("v4"))

	expect := []string{k3}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("OnRemove callback failed, expect keys equals to %s, but we got %s", expect, keys)
	}
}

// TestLFU_MemoryManagement 测试内存管理
func TestLFU_MemoryManagement(t *testing.T) {
	cap := len("k1") + len("k2") + len("v1") + len("v2")

	cache := New(int64(cap), nil)
	cache.Add("k1", String("v1"))
	cache.Add("k2", String("v2"))

	// Now add another element, which should trigger the removal of the least frequently used item (k1)
	cache.Add("k3", String("v3"))

	if v, ok := cache.Get("k1"); ok {
		t.Fatalf("MemoryManagement failed, key1 should be removed, got %v, %v", v, ok)
	}

	if v, ok := cache.Get("k2"); !ok || string(v.(String)) != "v2" {
		t.Fatalf("MemoryManagement failed, key2 should not be removed, got %v, %v", v, ok)
	}

	if v, ok := cache.Get("k3"); !ok || string(v.(String)) != "v3" {
		t.Fatalf("MemoryManagement failed, key3 should not be removed, got %v, %v", v, ok)
	}

	if cache.Mem.UsedBytes != int64(cap) {
		t.Fatalf("MemoryManagement failed, expect usedBytes equals to %d, but we got %d", cap, cache.Mem.UsedBytes)
	}
}
