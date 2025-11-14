package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tryuuu/consistent-hash-go/consistenthash"
	"github.com/tryuuu/consistent-hash-go/hash"
)

func generateTestKeys(count int) []string {
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}
	return keys
}

// measures how many keys are remapped after a node operation
func measureRemapping(
	keys []string,
	beforeMapping map[string]string,
	getNode func(string) string,
) int {
	remapped := 0
	for _, key := range keys {
		afterNode := getNode(key)
		if beforeMapping[key] != afterNode {
			remapped++
		}
	}
	return remapped
}

func runBenchmark() {
	numKeys := 10000
	numInitialNodes := 5
	numNodesToAdd := 2

	testKeys := generateTestKeys(numKeys)

	// Consistent Hash
	fmt.Println("--- Consistent Hash ---")
	chRing := consistenthash.New(150) // 150 replicas for better distribution
	initialNodes := make([]string, numInitialNodes)
	for i := 0; i < numInitialNodes; i++ {
		nodeName := fmt.Sprintf("Node%d", i)
		initialNodes[i] = nodeName
		chRing.Add(nodeName)
	}

	// Record initial mapping
	chBeforeMapping := make(map[string]string)
	for _, key := range testKeys {
		chBeforeMapping[key] = chRing.Get(key)
	}

	// Add nodes
	nodesToAdd := make([]string, numNodesToAdd)
	for i := 0; i < numNodesToAdd; i++ {
		nodeName := fmt.Sprintf("Node%d", numInitialNodes+i)
		nodesToAdd[i] = nodeName
		chRing.Add(nodeName)
	}

	// Measure remapping after adding nodes
	chRemappedAfterAdd := measureRemapping(testKeys, chBeforeMapping, chRing.Get)
	chAfterAddMapping := make(map[string]string)
	for _, key := range testKeys {
		chAfterAddMapping[key] = chRing.Get(key)
	}

	// Remove nodes
	nodeToRemove := initialNodes[0]
	chRing.Remove(nodeToRemove)

	// Measure remapping after removing node
	chRemappedAfterRemove := measureRemapping(testKeys, chAfterAddMapping, chRing.Get)

	fmt.Printf("After adding %d nodes:\n", numNodesToAdd)
	fmt.Printf("  Remapped keys: %d (%.2f%%)\n", chRemappedAfterAdd, float64(chRemappedAfterAdd)/float64(numKeys)*100)
	fmt.Printf("After removing 1 node:\n")
	fmt.Printf("  Remapped keys: %d (%.2f%%)\n", chRemappedAfterRemove, float64(chRemappedAfterRemove)/float64(numKeys)*100)
	fmt.Println()

	// Simple Hash
	fmt.Println("--- Simple Hash (Normal Hash Function) ---")
	shHash := hash.New()
	for i := 0; i < numInitialNodes; i++ {
		shHash.Add(initialNodes[i])
	}

	// Record initial mapping
	shBeforeMapping := make(map[string]string)
	for _, key := range testKeys {
		shBeforeMapping[key] = shHash.Get(key)
	}

	// Add nodes
	for i := 0; i < numNodesToAdd; i++ {
		shHash.Add(nodesToAdd[i])
	}

	// Measure remapping after adding nodes
	shRemappedAfterAdd := measureRemapping(testKeys, shBeforeMapping, shHash.Get)
	shAfterAddMapping := make(map[string]string)
	for _, key := range testKeys {
		shAfterAddMapping[key] = shHash.Get(key)
	}

	// Remove node
	shHash.Remove(nodeToRemove)

	// Measure remapping after removing node
	shRemappedAfterRemove := measureRemapping(testKeys, shAfterAddMapping, shHash.Get)

	fmt.Printf("After adding %d nodes:\n", numNodesToAdd)
	fmt.Printf("  Remapped keys: %d (%.2f%%)\n", shRemappedAfterAdd, float64(shRemappedAfterAdd)/float64(numKeys)*100)
	fmt.Printf("After removing 1 node:\n")
	fmt.Printf("  Remapped keys: %d (%.2f%%)\n", shRemappedAfterRemove, float64(shRemappedAfterRemove)/float64(numKeys)*100)
	fmt.Println()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	runBenchmark()
}
