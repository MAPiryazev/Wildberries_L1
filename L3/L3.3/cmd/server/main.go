package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"L3.3/internal/server"
)

func main() {
	srv, err := server.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init server: %v\n", err)
		os.Exit(1)
	}

	// Контекст, который отменится по сигналу
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := srv.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
