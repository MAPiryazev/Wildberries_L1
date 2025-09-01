package main

import (
	"fmt"
	"time"
)

func sleep(Duration time.Duration) {
	<-time.After(Duration) //канал
}

func main() {
	fmt.Println("Начало ", time.Now())

	sleep(2 * time.Second)

	fmt.Println("После сна ", time.Now())
}
