package main

import (
	"fmt"

	"github.com/tryuuu/consistent-hash-go/consistenthash"
	"github.com/tryuuu/consistent-hash-go/hash"
)

func main() {
	fmt.Println("=== Consistent Hash ===")
	ring := consistenthash.New(3)
	ring.Add("NodeA")
	ring.Add("NodeB")
	ring.Add("NodeC")

	// debug
	fmt.Println(ring)
	fmt.Println()

	keys := []string{"apple", "banana", "cherry", "grape", "melon"}
	for _, k := range keys {
		fmt.Printf("%s → %s\n", k, ring.Get(k))
	}

	fmt.Println("\n=== Simple Hash ===")
	simpleHash := hash.New()
	simpleHash.Add("NodeA")
	simpleHash.Add("NodeB")
	simpleHash.Add("NodeC")

	// debug
	fmt.Println(simpleHash)
	fmt.Println()

	for _, k := range keys {
		fmt.Printf("%s → %s\n", k, simpleHash.Get(k))
	}
}
