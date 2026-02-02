package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	A := make([]int, 0)
	B := make([]int, 0)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rand.Intn(500); i++ {
		A = append(A, rand.Intn(400))
	}
	rand.Seed(time.Now().UnixNano() + 1)
	for i := 0; i < rand.Intn(500); i++ {
		B = append(B, rand.Intn(400))
	}

	set := make(map[int]struct{})
	for _, val := range A {
		set[val] = struct{}{}
	}
	fmt.Printf("Пересечение = {")
	for _, val := range B {
		_, ok := set[val]
		if ok {
			fmt.Printf(" %d", val)
		}
	}
	fmt.Printf(" }")

}
