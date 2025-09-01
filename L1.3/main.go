package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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

	producerChan := make(chan int)
	defer close(producerChan)

	for i := 0; i < num; i++ {
		go func() {
			for num := range producerChan {
				fmt.Println(num)
			}
		}()
	}

	for true {
		producerChan <- rand.Int()
	}

}
