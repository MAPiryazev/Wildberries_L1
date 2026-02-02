package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/config"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/handlers"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/repository"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/service"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/templates"
	"github.com/robfig/cron/v3"
)

func main() {
	//ctx := context.Background()

	// Определяем путь к .env файлу
	envPath := filepath.Join("../../environment/.env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Пробуем альтернативный путь
		envPath = "environment/.env"
	}

	// Загружаем конфигурацию
	apiCfg := config.LoadAPIConfig(envPath)
	masterDSN, dbOpts, err := config.LoadDBPSQL(envPath)
	if err != nil {
		log.Fatalf("не удалось загрузить конфиг БД: %v", err)
	}

	// Инициализируем репозиторий
	repo, err := repository.NewPostgresRepository(masterDSN, dbOpts)
	if err != nil {
		log.Fatalf("ошибка инициализации репозитория: %v", err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			log.Printf("ошибка закрытия репозитория: %v", err)
		}
	}()

	// Инициализируем сервис
	svc := service.NewService(repo)

	// Инициализируем хендлеры
	handler := handlers.NewHandler(svc)

	// Загружаем шаблоны
	templatesPath := filepath.Join("../../internal/templates/*.html")
	if _, err := os.Stat("../../internal/templates"); os.IsNotExist(err) {
		// Пробуем альтернативный путь
		templatesPath = "internal/templates/*.html"
	}
	templates.LoadTemplates(templatesPath)

	// Настраиваем роутинг
	router := handlers.SetupRouter(handler)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    ":" + apiCfg.Port,
		Handler: router,
	}

	// Запускаем фоновый процесс для очистки просроченных бронирований
	c := cron.New()
	_, err = c.AddFunc("@every 1m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cancelled, err := svc.CancelExpiredBookings(ctx)
		if err != nil {
			log.Printf("ошибка при отмене просроченных бронирований: %v", err)
		} else if len(cancelled) > 0 {
			log.Printf("отменено просроченных бронирований: %d", len(cancelled))
		}
	})
	if err != nil {
		log.Fatalf("ошибка настройки cron: %v", err)
	}
	c.Start()
	defer c.Stop()

	log.Printf("Фоновый процесс очистки бронирований запущен (каждую минуту)")

	// Запускаем HTTP сервер в горутине
	go func() {
		log.Printf("Сервер запущен на порту %s", apiCfg.Port)
		log.Printf("Пользовательский интерфейс: http://localhost:%s/ui/user/events", apiCfg.Port)
		log.Printf("Админ-панель: http://localhost:%s/ui/admin/events", apiCfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("ошибка при остановке сервера: %v", err)
	}

	log.Println("Сервер остановлен")
}
