package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tryuuu/consistent-hash-go/pkg/consistenthash"
	"github.com/tryuuu/consistent-hash-go/pkg/distributor"
	"github.com/tryuuu/consistent-hash-go/pkg/hash"
)

func generateTestKeys(count int) []string {
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}
	return keys
}

func recordMapping(keys []string, dist distributor.HashDistributor) map[string]string {
	mapping := make(map[string]string, len(keys))
	for _, key := range keys {
		mapping[key] = dist.Get(key)
	}
	return mapping
}

func measureRemapping(keys []string, beforeMapping map[string]string, dist distributor.HashDistributor) int {
	remapped := 0
	for _, key := range keys {
		afterNode := dist.Get(key)
		if beforeMapping[key] != afterNode {
			remapped++
		}
	}
	return remapped
}

func benchmarkHashDistributor(
	name string,
	dist distributor.HashDistributor,
	testKeys []string,
	numKeys, numInitialNodes, numNodesToAdd int,
) {
	fmt.Printf("--- %s ---\n", name)

	// Initialize nodes
	initialNodes := make([]string, numInitialNodes)
	for i := 0; i < numInitialNodes; i++ {
		nodeName := fmt.Sprintf("Node%d", i)
		initialNodes[i] = nodeName
		dist.Add(nodeName)
	}

	// Record initial mapping
	beforeMapping := recordMapping(testKeys, dist)

	// Add new nodes
	for i := 0; i < numNodesToAdd; i++ {
		nodeName := fmt.Sprintf("Node%d", numInitialNodes+i)
		dist.Add(nodeName)
	}

	// Measure remapping after adding nodes
	remappedAfterAdd := measureRemapping(testKeys, beforeMapping, dist)
	afterAddMapping := recordMapping(testKeys, dist)

	// Remove a node
	dist.Remove(initialNodes[0])

	// Measure remapping after removing node
	remappedAfterRemove := measureRemapping(testKeys, afterAddMapping, dist)

	// Print results
	fmt.Printf("After adding %d nodes:\n", numNodesToAdd)
	fmt.Printf("  Remapped keys: %d (%.2f%%)\n", remappedAfterAdd, float64(remappedAfterAdd)/float64(numKeys)*100)
	fmt.Printf("After removing 1 node:\n")
	fmt.Printf("  Remapped keys: %d (%.2f%%)\n", remappedAfterRemove, float64(remappedAfterRemove)/float64(numKeys)*100)
	fmt.Println()
}

func runBenchmark() {
	const (
		numKeys         = 10000
		numInitialNodes = 5
		numNodesToAdd   = 2
	)

	testKeys := generateTestKeys(numKeys)

	benchmarkHashDistributor(
		"Consistent Hash",
		consistenthash.New(150),
		testKeys,
		numKeys,
		numInitialNodes,
		numNodesToAdd,
	)

	benchmarkHashDistributor(
		"Simple Hash (Normal Hash Function)",
		hash.New(),
		testKeys,
		numKeys,
		numInitialNodes,
		numNodesToAdd,
	)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	runBenchmark()
}
