package consistenthash

import (
	"sort"
	"testing"
)

func TestHashRingAddAndGet(t *testing.T) {
	ring := New(3)
	nodes := []string{"NodeA", "NodeB", "NodeC"}
	for _, node := range nodes {
		ring.Add(node)
	}

	wantKeys := ring.replicas * len(nodes)
	if got := len(ring.keys); got != wantKeys {
		t.Fatalf("ring keys count = %d, want %d", got, wantKeys)
	}

	if !sort.SliceIsSorted(ring.keys, func(i, j int) bool { return ring.keys[i] < ring.keys[j] }) {
		t.Fatalf("ring keys should be sorted")
	}

	nodeSet := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		nodeSet[node] = struct{}{}
	}

	testKeys := []string{"apple", "banana", "cherry", "grape", "melon", "orange", "peach"}
	for _, key := range testKeys {
		got := ring.Get(key)
		if _, ok := nodeSet[got]; !ok {
			t.Fatalf("ring.Get(%q) returned %q which is not part of the ring", key, got)
		}
		want := expectedNode(ring, key)
		if got != want {
			t.Fatalf("ring.Get(%q) = %q, want %q", key, got, want)
		}
	}
}

func expectedNode(hr *HashRing, key string) string {
	hash := hr.hash([]byte(key))
	idx := sort.Search(len(hr.keys), func(i int) bool { return hr.keys[i] >= hash })
	if idx == len(hr.keys) {
		idx = 0
	}
	return hr.ring[hr.keys[idx]]
}
