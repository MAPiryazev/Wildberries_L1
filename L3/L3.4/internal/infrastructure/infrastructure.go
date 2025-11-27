package infrastructure

import (
	"fmt"
	"log"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/config"
)

// NewStorage создает новый Storage клиент (MinIO)
func NewStorage(minioConfig *config.MinioConfgig) (Storage, error) {
	return NewMinIOStorage(minioConfig)
}

// NewMessageQueue создает новый MessageQueue клиент (Kafka)
func NewMessageQueue(kafkaConfig *config.KafkaConfig) (MessageQueue, error) {
	return NewKafkaQueue(kafkaConfig)
}

// InitInfrastructure инициализирует все infrastructure компоненты
func InitInfrastructure(envPath string) (Storage, MessageQueue, error) {
	log.Printf("[infra] start init, env=%s", envPath)
	// Загружаем конфигурацию MinIO
	minioConfig, err := config.LoadMinioConfig(envPath)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка загрузки MinIO конфигурации: %w", err)
	}
	log.Printf("[infra] minio config ok (endpoint=%s)", minioConfig.MinioEndpoint)

	// Создаем Storage клиент
	storage, err := NewStorage(minioConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка инициализации Storage: %w", err)
	}
	log.Printf("[infra] storage client ready")

	// Загружаем конфигурацию Kafka
	kafkaConfig, err := config.LoadKafkaConfig(envPath)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка загрузки Kafka конфигурации: %w", err)
	}
	log.Printf("[infra] kafka config ok (brokers=%s)", kafkaConfig.Brokers)

	// Создаем MessageQueue клиент
	messageQueue, err := NewMessageQueue(kafkaConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка инициализации MessageQueue: %w", err)
	}
	log.Printf("[infra] message queue ready")

	return storage, messageQueue, nil
}
