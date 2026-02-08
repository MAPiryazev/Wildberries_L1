package main

import (
	"log"

	"server-calendar/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run("./cfg/config.yaml"); err != nil {
		log.Fatal(err)
	}
}
