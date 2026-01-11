package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/config"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/db"
	appErrors "github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	defer func() {
		if err := db.Close(database); err != nil {
			log.Printf("[ERROR] close database: %v", err)
		}
	}()

	log.Printf("[INFO] Starting database migrations...")

	if err := db.RunMigrations(database); err != nil {
		return fmt.Errorf("migration failed: %w - %s", err, appErrors.ErrMigrationFailed)
	}

	log.Printf("[INFO] Migrations completed successfully")
	log.Printf("[INFO] Server configuration loaded - Port: %d, Env: %s", cfg.Server.Port, cfg.App.Env)
	log.Printf("[INFO] Database connected - Host: %s:%d, DB: %s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := database.Master.Close(); err != nil {
		log.Printf("[ERROR] Error closing database: %v", err)
	}

	log.Printf("[INFO] Server shutdown completed")

	return nil
}
