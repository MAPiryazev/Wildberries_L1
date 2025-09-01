package main

import (
	"fmt"
	"sync"
)

type concurrentCounter struct {
	counter int
	mu      sync.Mutex
}

func NewConcurrentCounter() *concurrentCounter {
	return &concurrentCounter{counter: 0, mu: sync.Mutex{}}
}

func main() {
	wg := sync.WaitGroup{}
	cC := NewConcurrentCounter()

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cC.mu.Lock()
			defer cC.mu.Unlock()
			cC.counter++
		}(i)
	}
	wg.Wait()
	fmt.Println(cC.counter)
}
