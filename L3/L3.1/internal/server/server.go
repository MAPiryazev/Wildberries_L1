package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"L3.1/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, h *handlers.DefaultHandler) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/notify", func(r chi.Router) {
		r.Post("/", h.CreateNotificationHandler)
		r.Get("/status", h.GetStatusHandler)
		r.Delete("/cancel", h.CancelNotificationHandler)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	return &Server{httpServer: srv}
}

func (s *Server) Start() {
	go func() {
		log.Printf("Сервер запущен на порту %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер остановлен корректно")
}
