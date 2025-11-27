package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/config"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOStorage реализация Storage интерфейса для MinIO
type MinIOStorage struct {
	client   *minio.Client
	config   *config.MinioConfgig
	endpoint string
}

// NewMinIOStorage создает новый клиент MinIO Storage
func NewMinIOStorage(cfg *config.MinioConfgig) (*MinIOStorage, error) {
	// Инициализация MinIO клиента
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации MinIO клиента: %w", err)
	}
	log.Printf("[minio] client connected endpoint=%s ssl=%v", cfg.MinioEndpoint, cfg.MinioUseSSL)

	storage := &MinIOStorage{
		client:   minioClient,
		config:   cfg,
		endpoint: cfg.MinioEndpoint,
	}

	// Создаем необходимые buckets если их нет
	ctx := context.Background()
	buckets := []string{
		cfg.MinioBucketOriginal,
		cfg.MinioBucketProcessed,
		cfg.MinioBucketMetadata,
	}

	for _, bucketName := range buckets {
		exists, err := minioClient.BucketExists(ctx, bucketName)
		if err != nil {
			return nil, fmt.Errorf("ошибка проверки существования bucket %s: %w", bucketName, err)
		}
		if !exists {
			err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			if err != nil {
				return nil, fmt.Errorf("ошибка создания bucket %s: %w", bucketName, err)
			}
			log.Printf("[minio] bucket %s created", bucketName)
		} else {
			log.Printf("[minio] bucket %s exists", bucketName)
		}
	}

	return storage, nil
}

// Upload загружает файл в MinIO bucket
func (m *MinIOStorage) Upload(bucket string, objectName string, reader io.Reader, size int64, contentType string) error {
	ctx := context.Background()

	_, err := m.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("ошибка загрузки объекта %s в bucket %s: %w", objectName, bucket, err)
	}

	return nil
}

// Download скачивает файл из MinIO bucket
func (m *MinIOStorage) Download(bucket string, objectName string) (io.ReadCloser, error) {
	ctx := context.Background()

	object, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка скачивания объекта %s из bucket %s: %w", objectName, bucket, err)
	}

	return object, nil
}

// Delete удаляет файл из MinIO bucket
func (m *MinIOStorage) Delete(bucket string, objectName string) error {
	ctx := context.Background()

	err := m.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("ошибка удаления объекта %s из bucket %s: %w", objectName, bucket, err)
	}

	return nil
}

// Exists проверяет существование файла в MinIO bucket
func (m *MinIOStorage) Exists(bucket string, objectName string) (bool, error) {
	ctx := context.Background()

	_, err := m.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("ошибка проверки существования объекта %s в bucket %s: %w", objectName, bucket, err)
	}

	return true, nil
}

// SaveMetadata сохраняет метаинформацию об изображении в MinIO
func (m *MinIOStorage) SaveMetadata(imageID string, info *models.ImageInfo) error {
	metadataJSON, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("ошибка сериализации метаданных: %w", err)
	}

	metadataPath := fmt.Sprintf("%s.json", imageID)
	reader := bytes.NewReader(metadataJSON)

	return m.Upload(m.config.MinioBucketMetadata, metadataPath, reader, int64(len(metadataJSON)), "application/json")
}

// GetMetadata получает метаинформацию об изображении из MinIO
func (m *MinIOStorage) GetMetadata(imageID string) (*models.ImageInfo, error) {
	metadataPath := fmt.Sprintf("%s.json", imageID)

	reader, err := m.Download(m.config.MinioBucketMetadata, metadataPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения метаданных: %w", err)
	}
	defer reader.Close()

	var info models.ImageInfo
	if err := json.NewDecoder(reader).Decode(&info); err != nil {
		return nil, fmt.Errorf("ошибка десериализации метаданных: %w", err)
	}

	return &info, nil
}

// UpdateMetadata обновляет метаинформацию об изображении в MinIO
func (m *MinIOStorage) UpdateMetadata(imageID string, info *models.ImageInfo) error {
	// В MinIO объекты неизменяемы, поэтому просто перезаписываем
	return m.SaveMetadata(imageID, info)
}

// DeleteMetadata удаляет метаинформацию об изображении
func (m *MinIOStorage) DeleteMetadata(imageID string) error {
	metadataPath := fmt.Sprintf("%s.json", imageID)
	return m.Delete(m.config.MinioBucketMetadata, metadataPath)
}

// ListMetadata перечисляет сохраненные метаданные
func (m *MinIOStorage) ListMetadata(ctx context.Context, limit int) ([]*models.ImageInfo, error) {
	if limit <= 0 {
		limit = 100
	}

	result := make([]*models.ImageInfo, 0, limit)

	for object := range m.client.ListObjects(ctx, m.config.MinioBucketMetadata, minio.ListObjectsOptions{Recursive: true}) {
		if object.Err != nil {
			return nil, fmt.Errorf("list metadata objects: %w", object.Err)
		}

		reader, err := m.client.GetObject(ctx, m.config.MinioBucketMetadata, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, fmt.Errorf("get metadata object %s: %w", object.Key, err)
		}

		var info models.ImageInfo
		if err := json.NewDecoder(reader).Decode(&info); err != nil {
			reader.Close()
			return nil, fmt.Errorf("decode metadata %s: %w", object.Key, err)
		}
		reader.Close()

		result = append(result, &info)
		if len(result) >= limit {
			break
		}
	}

	return result, nil
}

// GetURL возвращает URL для доступа к изображению
func (m *MinIOStorage) GetURL(bucket string, objectName string) string {
	protocol := "http"
	if m.config.MinioUseSSL {
		protocol = "https"
	}

	// Генерируем presigned URL с временем жизни 24 часа
	ctx := context.Background()
	url, err := m.client.PresignedGetObject(ctx, bucket, objectName, 24*time.Hour, nil)
	if err != nil {
		// Если не удалось создать presigned URL, возвращаем обычный URL
		return fmt.Sprintf("%s://%s/%s/%s", protocol, m.endpoint, bucket, objectName)
	}

	return url.String()
}
