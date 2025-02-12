package cache

import (
	"strconv"
	"sync"
	"testing"
)

func TestCache_New(t *testing.T) {
	cache1 := NewCache(0)
	cache2 := NewCache(10)
	cache3 := NewCache(-10)

	if cache1 == nil || cache2 == nil || cache3 == nil {
		t.Fatal("fail to init cache")
	}
}

func TestCache_Add(t *testing.T) {
	cache := NewCache(0)

	// use for loop to store 100 numbers
	for i := 0; i < 100; i += 1 {
		cache.add(strconv.Itoa(i), NewByteView([]byte(strconv.Itoa(i+1))))
	}
}

func TestCache_Get(t *testing.T) {
	cache := NewCache(0)

	var group sync.WaitGroup
	error_info := make(chan [2]int)

	for i := 0; i < 10; i += 1 {
		group.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			cache.add(strconv.Itoa(i), NewByteView([]byte(strconv.Itoa(i+1))))
		}(&group, i)
	}

	group.Wait()

	for i := 0; i < 10; i += 1 {
		group.Add(1)
		go func(e chan<- [2]int, wg *sync.WaitGroup, i int) {
			defer wg.Done()
			if v, ok := cache.get(strconv.Itoa(i)); v.String() != strconv.Itoa(i+1) || !ok {
				e <- [2]int{i, i + 1}
			}
		}(error_info, &group, i)
	}

	group.Wait()

	select {
	case k := <-error_info:
		{
			t.Fatalf("we want pair <%d, %d>", k[0], k[1])
		}
	default:
	}
}

// func TestCache_Remove(t *testing.T) {

// }
