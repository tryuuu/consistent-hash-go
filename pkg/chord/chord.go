package chord

import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"

	"github.com/tryuuu/consistent-hash-go/pkg/distributor"
)

const HashSpaceBits = 32

var _ distributor.HashDistributor = (*ChordRing)(nil)

type ChordNode struct {
	id          uint32
	name        string
	successor   uint32
	predecessor uint32
	fingerTable []uint32
}

type ChordRing struct {
	mu        sync.RWMutex
	hash      func([]byte) uint32
	nodes     map[uint32]*ChordNode
	sortedIDs []uint32
}

func New() *ChordRing {
	return &ChordRing{
		hash: func(b []byte) uint32 {
			h := fnv.New32a()
			h.Write(b)
			return h.Sum32()
		},
		nodes: make(map[uint32]*ChordNode),
	}
}

func (cr *ChordRing) Add(nodeName string) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if nodeName == "" {
		return
	}

	id := cr.hash([]byte(nodeName))

	if _, exists := cr.nodes[id]; exists {
		return
	}

	n := &ChordNode{
		id:          id,
		name:        nodeName,
		fingerTable: make([]uint32, HashSpaceBits),
	}
	cr.nodes[id] = n
	cr.rebuild()
}

func (cr *ChordRing) Remove(nodeName string) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if len(cr.nodes) == 0 {
		return
	}

	id := cr.hash([]byte(nodeName))
	if _, exists := cr.nodes[id]; !exists {
		return
	}
	delete(cr.nodes, id)
	cr.rebuild()
}

func (cr *ChordRing) Get(key string) string {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	if len(cr.sortedIDs) == 0 {
		return ""
	}

	keyID := cr.hash([]byte(key))

	startID := cr.sortedIDs[0]
	startNode := cr.nodes[startID]
	target := cr.findSuccessorViaChord(startNode, keyID)
	if target == nil {
		return ""
	}
	return target.name
}

func (cr *ChordRing) rebuild() {
	cr.sortedIDs = cr.sortedIDs[:0]
	for id := range cr.nodes {
		cr.sortedIDs = append(cr.sortedIDs, id)
	}
	sort.Slice(cr.sortedIDs, func(i, j int) bool { return cr.sortedIDs[i] < cr.sortedIDs[j] })

	if len(cr.sortedIDs) == 0 {
		return
	}

	for i, id := range cr.sortedIDs {
		n := cr.nodes[id]
		succID := cr.sortedIDs[(i+1)%len(cr.sortedIDs)]
		predID := cr.sortedIDs[(i-1+len(cr.sortedIDs))%len(cr.sortedIDs)]
		n.successor = succID
		n.predecessor = predID
	}

	for _, id := range cr.sortedIDs {
		n := cr.nodes[id]
		for i := 0; i < HashSpaceBits; i++ {
			start := id + (1 << uint(i))
			succ := cr.findSuccessorOnRing(start)
			if succ != nil {
				n.fingerTable[i] = succ.id
			} else {
				n.fingerTable[i] = id
			}
		}
	}
}

func (cr *ChordRing) findSuccessorOnRing(id uint32) *ChordNode {
	if len(cr.sortedIDs) == 0 {
		return nil
	}
	idx := sort.Search(len(cr.sortedIDs), func(i int) bool { return cr.sortedIDs[i] >= id })
	if idx == len(cr.sortedIDs) {
		idx = 0
	}
	return cr.nodes[cr.sortedIDs[idx]]
}

func between(start, x, end uint32) bool {
	if start < end {
		return x > start && x <= end
	}
	return x > start || x <= end
}

func (cr *ChordRing) findSuccessorViaChord(node *ChordNode, id uint32) *ChordNode {
	if node == nil {
		return nil
	}
	if len(cr.sortedIDs) == 0 {
		return nil
	}

	maxSteps := len(cr.sortedIDs) + 2
	current := node

	for step := 0; step < maxSteps; step++ {
		succ := cr.nodes[current.successor]
		if succ == nil {
			return cr.findSuccessorOnRing(id)
		}

		if between(current.id, id, succ.id) {
			return succ
		}

		nextID := current.id
		for i := HashSpaceBits - 1; i >= 0; i-- {
			fid := current.fingerTable[i]
			if fid == 0 {
				continue
			}
			if between(current.id, fid, id) {
				nextID = fid
				break
			}
		}

		if nextID == current.id {
			current = succ
		} else {
			n, ok := cr.nodes[nextID]
			if !ok {
				return cr.findSuccessorOnRing(id)
			}
			current = n
		}
	}

	return cr.findSuccessorOnRing(id)
}

func (cr *ChordRing) String() string {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	if len(cr.sortedIDs) == 0 {
		return "ChordRing{empty}"
	}

	result := "ChordRing{\n"
	for _, id := range cr.sortedIDs {
		n := cr.nodes[id]
		result += fmt.Sprintf("  %d (%s) succ=%d pred=%d\n", n.id, n.name, n.successor, n.predecessor)
	}
	result += "}"
	return result
}
