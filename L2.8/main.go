package main

import (
	"fmt"
	"log"
	"myntp/myntp"
)

func main() {
	time, err := myntp.GetNTPTime()
	if err != nil {
		log.Fatal("Ошибка получения времени: ", err)
	}
	fmt.Println(time)
}
