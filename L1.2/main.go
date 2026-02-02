package main

import (
	"fmt"
	"sync"
)

func main() {
	arr := [5]int{2, 4, 6, 8, 10}
	wg := sync.WaitGroup{}
	for _, val := range arr {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			fmt.Println(val * val)
		}(val)
	}
	wg.Wait()
}
