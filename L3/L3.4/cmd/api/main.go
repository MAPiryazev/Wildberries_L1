package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/config"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/httpapi"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/infrastructure"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/service/image"
)

func main() {
	envPath := flag.String("env", "../../environment/.env", "path to .env file")
	addr := flag.String("addr", ":8080", "http listen address")
	flag.Parse()

	log.Printf("[bootstrap] starting API, env=%s addr=%s", *envPath, *addr)

	minioCfg, err := config.LoadMinioConfig(*envPath)
	if err != nil {
		log.Fatalf("load minio config: %v", err)
	}
	log.Printf("[bootstrap] minio config loaded, endpoint=%s buckets=[%s,%s]",
		minioCfg.MinioEndpoint, minioCfg.MinioBucketOriginal, minioCfg.MinioBucketProcessed)

	kafkaCfg, err := config.LoadKafkaConfig(*envPath)
	if err != nil {
		log.Fatalf("load kafka config: %v", err)
	}
	log.Printf("[bootstrap] kafka config loaded, brokers=%s topic=%s group=%s",
		kafkaCfg.Brokers, kafkaCfg.TopicImageTasks, kafkaCfg.GroupID)

	storage, queue, err := infrastructure.InitInfrastructure(*envPath)
	if err != nil {
		log.Fatalf("init infrastructure: %v", err)
	}
	defer queue.Close()

	imageService, err := image.NewService(storage, queue, image.Config{
		Buckets: image.Buckets{
			Original:  minioCfg.MinioBucketOriginal,
			Processed: minioCfg.MinioBucketProcessed,
		},
		TopicImageTasks:   kafkaCfg.TopicImageTasks,
		DefaultProcessing: models.ProcessingResize,
	})
	if err != nil {
		log.Fatalf("init image service: %v", err)
	}
	log.Printf("[bootstrap] image service initialized, default processing=%s", models.ProcessingResize)

	handler := httpapi.NewHandler(imageService, httpapi.HandlerConfig{MaxListItems: 100})
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:    *addr,
		Handler: router,
	}

	go func() {
		log.Printf("[api] HTTP server listening on %s", *addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("[api] shutting down...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}
}
