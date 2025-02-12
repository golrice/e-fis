package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 2^32
type Hash func([]byte) uint32

// a map only store the node basic metadata, we use node_name(string) to access the real node
type DHTMap struct {
	// virtual node : real node
	replicas int
	// virtual nodes vec
	nodes []int
	// virtual id -> real name
	origins map[int]string
	// hash function
	hash Hash
}

func New(replicas int, h Hash) *DHTMap {
	m := &DHTMap{
		replicas: replicas,
		nodes:    make([]int, 0),
		origins:  map[int]string{},
		hash:     h,
	}

	if h == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// add a node in map
func (m *DHTMap) Add(realNodes ...string) {
	for _, realNode := range realNodes {
		// we create some virtual node into map
		for i := 0; i < m.replicas; i += 1 {
			virtualNodeID := int(m.hash([]byte(strconv.Itoa(i) + realNode)))
			m.nodes = append(m.nodes, virtualNodeID)
			m.origins[virtualNodeID] = realNode
		}
	}
	sort.Ints(m.nodes)
}

// only get the real node name, not the value
func (m *DHTMap) Get(key string) string {
	if key == "" {
		return ""
	}

	// calculate the hash value of the key
	hash := m.hash([]byte(key))

	// find the next node of this key, we can use binary search
	idx := sort.SearchInts(m.nodes, int(hash))

	// get the real node
	return m.origins[m.nodes[idx%len(m.nodes)]]
}
