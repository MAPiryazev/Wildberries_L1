package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/config"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/db"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/handlers"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/middleware"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/repository"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load("environment/.env")
	if err != nil {
		log.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	if err := db.RunMigrations(database); err != nil {
		log.Error("failed to run migrations", "err", err)
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(database); err != nil {
			log.Error("failed to close database", "err", err)
		}
	}()

	itemRepo := repository.NewItemRepository(database, cfg, log)
	userRepo := repository.NewUserRepository(database, cfg, log)

	itemService := service.NewItemService(itemRepo, userRepo, log)
	userService := service.NewUserService(userRepo, log)

	jwtSecret := cfg.JWT.Secret
	if jwtSecret == "" {
		log.Error("jwt secret is empty")
		os.Exit(1)
	}

	handler := handlers.New(itemService, userService, jwtSecret, log)

	router := chi.NewRouter()
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())

	// Front
	router.Get("/", serveIndex)
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web"))))

	// Public auth
	router.Post("/auth/login", handler.Login)

	// Protected API
	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.JWTAuth(jwtSecret, log))
	handler.RegisterRoutes(apiRouter)
	router.Mount("/api", apiRouter)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("shutdown error", "err", err)
	}

	log.Info("server stopped")
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./web/index.html")
}
