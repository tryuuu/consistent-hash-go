package consistenthash

import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
)

type HashRing struct {
	mu   sync.RWMutex
	hash func([]byte) uint32
	keys []uint32
	ring map[uint32]string
}

func New() *HashRing {
	return &HashRing{
		hash: func(b []byte) uint32 {
			h := fnv.New32a()
			h.Write(b)
			return h.Sum32()
		},
		ring: make(map[uint32]string),
	}
}

func (hr *HashRing) Add(node string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	h := hr.hash([]byte(node))
	hr.ring[h] = node
	hr.keys = append(hr.keys, h)
	sort.Slice(hr.keys, func(i, j int) bool { return hr.keys[i] < hr.keys[j] })
}

func (hr *HashRing) Get(key string) string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if len(hr.keys) == 0 {
		return ""
	}
	h := hr.hash([]byte(key))
	// debug
	fmt.Printf("key=%s, hash=%d\n", key, h)

	idx := sort.Search(len(hr.keys), func(i int) bool { return hr.keys[i] >= h })
	if idx == len(hr.keys) {
		idx = 0
	}
	return hr.ring[hr.keys[idx]]
}

func (hr *HashRing) String() string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if len(hr.keys) == 0 {
		return "HashRing{empty}"
	}
	result := "HashRing{\n"
	for _, key := range hr.keys {
		result += fmt.Sprintf("  %d -> %s\n", key, hr.ring[key])
	}
	result += "}"
	return result
}
