package cache

import (
	"fmt"
	"testing"
)

func TestNode_New(t *testing.T) {
	node := NewNode("test", 2<<10, func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	if node == nil {
		t.Fatal("fail to init node")
	}
}

func TestGraph_Get(t *testing.T) {
	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	loadCounts := make(map[string]int, len(db))
	node := NewNode("score", 2<<10, func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}

			loadCounts[key] += 1

			return []byte(v), nil
		}

		return nil, fmt.Errorf("no key: %s", key)
	})

	if node == nil {
		t.Fatal("fail to init node")
	}

	// we try to get all item in db
	for k, v := range db {
		// we get this item for the first time, we will call the callback function
		bv, err := node.Get(k)

		if err != nil {
			t.Fatal("callback function cause error: ", err)
		}

		if bv.String() != v {
			t.Fatal("we get wrong value")
		}

		// we get item for second time, and will not call the function
		if _, err := node.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatal("fail to get for second time")
		}
	}

	// we try to get item which is not exists
	if _, err := node.Get("Not Exists"); err == nil {
		t.Fatal("it does not cause error when get a Not Exists item")
	}
}
