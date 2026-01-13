package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	orchannel "github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.1/utils/or-channel"
)

func chanArrayGen(n, seconds int) []chan interface{} {
	if n <= 0 || seconds <= 0 {
		panic("время и количество каналов обязательно должны быть > 0")
	}

	arr := make([]chan interface{}, n)
	for val := range arr {
		arr[val] = make(chan interface{})
	}

	go func() {
		defer func() {
			for _, val := range arr {
				close(val)
			}
		}()

		<-time.After(time.Duration(seconds) * time.Second)
		num := rand.IntN(n)
		arr[num] <- struct{}{}
	}()

	return arr
}

func inputGen(ctx context.Context) chan int {
	if ctx == nil {
		panic("контекст не передан")
	}

	out := make(chan int)

	go func() {
		defer close(out)

		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				v := rand.IntN(10_000)

				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}()

	return out
}

type workerPool struct {
	amountOfWorkers int
	inputChannel    chan int
	sum             int

	stopChannel <-chan interface{}
	once        sync.Once

	wg sync.WaitGroup
	mu sync.RWMutex
}

func newWorkerPool(amount int, inputCh chan int, stopCh <-chan interface{}) *workerPool {
	if amount <= 0 {
		amount = 1
		log.Println("количество воркеров объявлено как 1")
	}

	return &workerPool{
		amountOfWorkers: amount,
		inputChannel:    inputCh,
		stopChannel:     stopCh,
	}
}

func (wp *workerPool) start() {
	for i := 0; i < wp.amountOfWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()

			for {
				select {
				case <-wp.stopChannel:
					return
				case v, ok := <-wp.inputChannel:
					if !ok {
						return
					}
					wp.mu.Lock()
					wp.sum += v
					wp.mu.Unlock()
				}
			}
		}()
	}
}

func (wp *workerPool) wait() {
	wp.wg.Wait()
}

func (wp *workerPool) result() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.sum
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := inputGen(ctx)
	stopChannelsArray := chanArrayGen(10, 5)

	wp := newWorkerPool(10, input, orchannel.OrChannel(stopChannelsArray...))
	wp.start()
	wp.wait()
	fmt.Println(wp.result())
}
