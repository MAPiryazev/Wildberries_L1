package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Слишком мало аргументов на вызове программы")
		return
	}

	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Неправильный формат ввода")
		return
	}

	done := time.After(time.Duration(N) * time.Second)
	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for val := range ch {
			fmt.Println("Получено значение", val)
		}
		fmt.Println("Читающая горутина завершилась")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Println("Пишущая горутина завершилась")
				close(ch)
				return
			default:
				ch <- rand.Int()
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	wg.Wait()
}
