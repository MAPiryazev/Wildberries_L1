package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const ARRAY_LENGTH = 100_000

func main() {
	array := make([]int, 0, ARRAY_LENGTH)
	for i := 0; i < ARRAY_LENGTH; i++ {
		array = append(array, rand.Intn(100000))
	}

	chan1 := make(chan int)
	chan2 := make(chan int)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(chan1)
		for _, val := range array {
			chan1 <- val
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(chan2)
		for val := range chan1 {
			chan2 <- val * 2
		}
	}()

	for val := range chan2 {
		fmt.Println(val)
	}
	wg.Wait()
}
