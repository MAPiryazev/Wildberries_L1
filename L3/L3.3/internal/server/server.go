package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"L3.3/internal/config"
	"L3.3/internal/handlers"
	"L3.3/internal/middleware"
	"L3.3/internal/repository"
	"L3.3/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Server struct {
	httpServer *http.Server
	repo       repository.Repository
}

func New() (*Server, error) {
	apiCfg := config.LoadAPIConfig("../../environment/.env")

	masterDSN, dbOpts, err := config.LoadDBPSQL("../../environment/.env")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить конфиг БД: %w", err)
	}

	repo, err := repository.NewPostgresRepository(masterDSN, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации репозитория: %w", err)
	}

	svc := service.NewCommentService(repo, apiCfg)
	h := handlers.NewCommentHandler(svc)

	r := ginext.New("commenttree")
	r.Use(middleware.LoggingMiddleware())

	// API
	api := r.Group("/api")
	h.RegisterRoutes(api)

	// Статика
	webDir := filepath.Join("../../web")
	r.GET("/", func(c *ginext.Context) {
		c.File(filepath.Join(webDir, "index.html"))
	})
	r.StaticFile("/styles.css", filepath.Join(webDir, "styles.css"))
	r.StaticFile("/app.js", filepath.Join(webDir, "app.js"))

	srv := &http.Server{
		Addr:    ":" + apiCfg.Port,
		Handler: r,
	}

	return &Server{httpServer: srv, repo: repo}, nil
}

func (s *Server) Run(ctx context.Context) error {
	// graceful shutdown
	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// ловим сигналы OS
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		return s.shutdown()
	case sig := <-sigCh:
		fmt.Printf("получен сигнал: %v, останавливаем сервер...\n", sig)
		return s.shutdown()
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}

func (s *Server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка graceful shutdown: %w", err)
	}
	if closer, ok := s.repo.(*repository.PostgresRepository); ok {
		_ = closer.Close()
	}
	return nil
}
