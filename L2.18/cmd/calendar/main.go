package main

import (
	"log"
	"net/http"

	"calendar/internal/app"
	"calendar/internal/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	APIConfig := config.LoadAPIConfig("../../environment/.env")

	server, err := app.NewServer("../../environment/.env")
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Запускаем сервер на порту ", APIConfig.Port)
	err = http.ListenAndServe(":"+APIConfig.Port, server)
	if err != nil {
		log.Panicln("запуск сервера не удался", err)
	}

}
