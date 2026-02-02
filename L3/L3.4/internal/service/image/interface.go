package image

import (
	"context"
	"io"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/models"
)

// Buckets описывает имена MinIO бакетов
type Buckets struct {
	Original  string
	Processed string
}

// Config параметры service слоя
type Config struct {
	Buckets           Buckets
	TopicImageTasks   string
	DefaultProcessing models.ProcessingType
}

// Service описывает бизнес-логику работы с изображениями
type Service interface {
	// Upload принимает исходный файл, сохраняет его и ставит задачу в очередь
	Upload(ctx context.Context, file io.Reader, size int64, contentType, filename string, processingType models.ProcessingType) (*models.ImageInfo, error)

	// Get возвращает информацию о состоянии обработки изображения
	Get(ctx context.Context, imageID string) (*models.ImageInfo, error)

	// List возвращает список изображений с их статусами
	List(ctx context.Context, limit int) ([]*models.ImageInfo, error)

	// Delete удаляет исходный/обработанный файлы и метаданные
	Delete(ctx context.Context, imageID string) error

	// MarkProcessing переводит изображение в статус "processing" (worker подтверждает старт)
	MarkProcessing(ctx context.Context, imageID string) (*models.ImageInfo, error)

	// MarkReady обновляет статус на "ready" и сохраняет путь до обработанного файла
	MarkReady(ctx context.Context, imageID string, processedObjectKey string) (*models.ImageInfo, error)

	// MarkFailed обновляет статус на "failed" и фиксирует ошибку обработки
	MarkFailed(ctx context.Context, imageID string, processingErr error) (*models.ImageInfo, error)
}
