package main

import (
	"fmt"
)

func main() {
	input := []string{"cat", "cat", "dog", "cat", "tree"}

	set := make(map[string]struct{})

	for _, val := range input {
		set[val] = struct{}{}
	}
	fmt.Println(set)
}
