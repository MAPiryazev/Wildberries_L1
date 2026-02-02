package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/middleware"
	"calendar/internal/repository/postgres"

	"calendar/internal/service"
)

// NewServer содержит логику для создания сервера чтобы вынести ее из main
func NewServer(envPath string) (http.Handler, error) {
	DBConfig, err := config.LoadDBPSQLConfig(envPath)
	if err != nil {
		return nil, err
	}

	repo, err := postgres.NewPostgresStorage(DBConfig)
	if err != nil {
		return nil, err
	}

	svc, err := service.NewDefaultService(repo)
	if err != nil {
		return nil, err
	}

	handler, err := handlers.NewDefaultHandler(svc)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(middleware.LoggingMiddleware)

	router.Post("/create_event", handler.CreateEvent)
	router.Post("/update_event", handler.UpdateEvent)
	router.Post("/delete_event", handler.DeleteEvent)
	router.Get("/events_for_day", handler.EventsForDay)
	router.Get("/events_for_week", handler.EventsForWeek)
	router.Get("/events_for_month", handler.EventsForMonth)

	return router, nil
}
