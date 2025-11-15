package hash

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type SimpleHash struct {
	mu    sync.RWMutex
	hash  func([]byte) uint32
	nodes []string
}

func New() *SimpleHash {
	return &SimpleHash{
		hash: func(b []byte) uint32 {
			h := fnv.New32a()
			h.Write(b)
			return h.Sum32()
		},
		nodes: make([]string, 0),
	}
}

func (sh *SimpleHash) Add(node string) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	// Check if node already exists
	for _, n := range sh.nodes {
		if n == node {
			return
		}
	}
	sh.nodes = append(sh.nodes, node)
}

func (sh *SimpleHash) Remove(node string) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	for i, n := range sh.nodes {
		if n == node {
			sh.nodes = append(sh.nodes[:i], sh.nodes[i+1:]...)
			return
		}
	}
}

func (sh *SimpleHash) Get(key string) string {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	if len(sh.nodes) == 0 {
		return ""
	}

	h := sh.hash([]byte(key))

	idx := int(h) % len(sh.nodes)
	return sh.nodes[idx]
}

func (sh *SimpleHash) String() string {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	if len(sh.nodes) == 0 {
		return "SimpleHash{empty}"
	}

	result := "SimpleHash{\n"
	for i, node := range sh.nodes {
		result += fmt.Sprintf("  [%d] %s\n", i, node)
	}
	result += "}"
	return result
}
