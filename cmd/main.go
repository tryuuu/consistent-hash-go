package main

import (
	"fmt"

	"github.com/tryuuu/consistent-hash-go/consistenthash"
)

func main() {
	ring := consistenthash.New()
	ring.Add("NodeA")
	ring.Add("NodeB")
	ring.Add("NodeC")

	fmt.Println(ring)
	fmt.Println()

	keys := []string{"apple", "banana", "cherry", "grape", "melon"}
	for _, k := range keys {
		fmt.Printf("%s → %s\n", k, ring.Get(k))
	}
}
