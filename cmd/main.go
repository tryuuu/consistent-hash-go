package main

import (
	"fmt"

	"github.com/tryuuu/consistent-hash-go/pkg/consistenthash"
	"github.com/tryuuu/consistent-hash-go/pkg/distributor"
	"github.com/tryuuu/consistent-hash-go/pkg/hash"
)

func testHashDistributor(name string, dist distributor.HashDistributor) {
	fmt.Printf("=== %s ===\n", name)
	dist.Add("NodeA")
	dist.Add("NodeB")
	dist.Add("NodeC")

	keys := []string{"apple", "banana", "cherry", "grape", "melon"}
	for _, k := range keys {
		fmt.Printf("%s → %s\n", k, dist.Get(k))
	}
	fmt.Println()
}

func main() {
	testHashDistributor("Consistent Hash", consistenthash.New(3))
	testHashDistributor("Simple Hash", hash.New())
}
