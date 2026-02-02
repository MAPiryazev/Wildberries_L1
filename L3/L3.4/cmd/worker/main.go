package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/config"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/infrastructure"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/service/image"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/worker"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	envPath := flag.String("env", "../../environment/.env", "path to .env file")
	groupID := flag.String("group", "image-worker", "kafka consumer group id")
	flag.Parse()

	minioCfg, err := config.LoadMinioConfig(*envPath)
	if err != nil {
		log.Fatalf("load minio config: %v", err)
	}

	kafkaCfg, err := config.LoadKafkaConfig(*envPath)
	if err != nil {
		log.Fatalf("load kafka config: %v", err)
	}

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

	processor := worker.NewProcessor(storage, queue, imageService)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := processor.Run(ctx, kafkaCfg.TopicImageTasks, *groupID); err != nil {
			log.Fatalf("worker stopped: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("[worker] shutting down")
}
