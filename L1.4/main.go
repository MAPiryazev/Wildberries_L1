package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("не указано количество воркеров на запуске")
		return
	}

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Ошибка преобразования типа, введите целое число")
		return
	}
	fmt.Println("Вы ввели ", num)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	producerChan := make(chan int)
	defer close(producerChan)

	for i := 0; i < num; i++ {
		go func(id int) {
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d завершает работу\n", id)
					return
				case val, ok := <-producerChan:
					if ok {
						fmt.Printf("Worker %d, значение %d\n", i, val)
					} else {
						return
					}
				}

			}
		}(i)
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Пишушая горутина заканчивает работу")
			return
		default:
			producerChan <- rand.Int()
			time.Sleep(200 * time.Millisecond)
		}
	}

}
