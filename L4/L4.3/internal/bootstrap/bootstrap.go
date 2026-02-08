package bootstrap

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"server-calendar/cfg"
	"server-calendar/internal/handler"
	applog "server-calendar/internal/log"
	"server-calendar/internal/server"
	"server-calendar/internal/service"
	"server-calendar/internal/storage"
	"server-calendar/internal/worker"
)

func Run(path string) error {
	config, err := cfg.NewConfig(path)
	if err != nil {
		return err
	}

	applog.SetLogrus(config.Log)

	stg := storage.NewStorage()

	reminderWorker := worker.NewReminderWorker(100)
	reminderWorker.Start()
	defer reminderWorker.Stop()

	svc := service.NewCalendarService(stg, reminderWorker)

	asyncLogger := applog.NewAsyncLogger(100)
	asyncLogger.Start()
	defer asyncLogger.Stop()

	archiveWorker := worker.NewArchiveWorker(5*time.Minute, stg.ArchiveOldEvents)
	archiveWorker.Start()
	defer archiveWorker.Stop()

	router := handler.NewRouter(svc)
	handlerWithLogs := applog.LoggerMiddleware(asyncLogger)(router)

	srv := server.New(
		handlerWithLogs,
		server.Port(config.Port),
		server.ReadTimeout(5*time.Second),
		server.WriteTimeout(10*time.Second),
		server.ShutdownTimeout(5*time.Second),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logrus.Infof("received signal: %v", sig)
		_ = srv.Shutdown()
	case err := <-srv.Notify():
		logrus.Errorf("server exited with error: %v", err)
	}

	return nil
}
