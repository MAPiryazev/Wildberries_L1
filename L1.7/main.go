package main

import (
	"fmt"
	"sync"
)

type SyncMap struct {
	mapa map[int]int
	mu   sync.RWMutex
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		mapa: make(map[int]int),
		mu:   sync.RWMutex{},
	}
}

func main() {
	syncMap := NewSyncMap()
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(iteration int) {
			defer wg.Done()
			syncMap.mu.Lock()
			defer syncMap.mu.Unlock()
			_, ok := syncMap.mapa[iteration]
			if !ok {
				syncMap.mapa[iteration] = 1
			} else {
				syncMap.mapa[iteration]++
			}

		}(i)
	}

	wg.Wait()
	fmt.Println(syncMap.mapa)

}
