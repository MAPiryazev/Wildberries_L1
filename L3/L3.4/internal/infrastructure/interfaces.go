package infrastructure

import (
	"context"
	"io"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/models"
)

type Storage interface {
	Upload(bucket string, objectName string, reader io.Reader, size int64, contentType string) error
	Download(bucket string, objectName string) (io.ReadCloser, error)
	Delete(bucket string, objectName string) error
	Exists(bucket string, objectName string) (bool, error)

	// Методы для работы с метаданными изображений
	SaveMetadata(imageID string, info *models.ImageInfo) error
	GetMetadata(imageID string) (*models.ImageInfo, error)
	UpdateMetadata(imageID string, info *models.ImageInfo) error
	DeleteMetadata(imageID string) error
	ListMetadata(ctx context.Context, limit int) ([]*models.ImageInfo, error)

	// Метод для получения URL изображения (для веб-интерфейса)
	GetURL(bucket string, objectName string) string
}

type MessageQueue interface {
	Publish(topic string, message []byte) error
	Subscribe(topic, groupID string, handler func(msg []byte) error) error
	Close() error
}
