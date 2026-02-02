package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}

	//Выход по условию
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			if i == 9 {
				fmt.Println("Выход по условию")
				return
			}
		}
	}()

	//Через канал уведомления
	notificationChan := make(chan struct{})
	defer close(notificationChan)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-notificationChan:
				fmt.Println("Выход по сигналу из канала")
				return
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
	time.Sleep(3 * time.Second)
	notificationChan <- struct{}{}

	//Через контекст
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Выход по завершении контекста")
				return
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	//Через runtime.Goexit()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Выход по runtime.Goexit")
		runtime.Goexit()
	}()

	wg.Wait()
}
