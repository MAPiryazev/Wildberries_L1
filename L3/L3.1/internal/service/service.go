package service

import (
	"context"
	"fmt"
	"time"

	"L3.1/internal/infrastructure"
	"L3.1/internal/models"
	"github.com/google/uuid"
)

// Notifier интерфей с необходимыми методами
type Notifier interface {
	Send(ctx context.Context, msg *models.RabbitMQMessage)
}

// NotificationService логика работы
type NotificationService struct {
	cache            infrastructure.CacheClient
	queue            infrastructure.QueueMQClient
	notifier         Notifier
	queueNameDelayed string
	queueNameReady   string

	retryBaseSeonds int
	maxRetries      int
}

func NewNotificationService(
	cache infrastructure.CacheClient,
	queue infrastructure.QueueMQClient,
	notifier Notifier,
	delayedQueue string,
	readyQueue string,
) (*NotificationService, error) {
	if cache == nil || queue == nil || notifier == nil {
		return nil, fmt.Errorf("не все поля переданы в конструктор")
	}
	return &NotificationService{
		cache:           cache,
		queue:           queue,
		notifier:        notifier,
		queueNameReady:  readyQueue,
		retryBaseSeonds: 2,
		maxRetries:      5,
	}, nil
}

func (ns *NotificationService) CreateNotification(ctx context.Context, payload string, sendAt time.Time) (string, error) {
	id := uuid.New().String()

	rabbitMsg := &models.RabbitMQMessage{
		ID:         id,
		Payload:    payload,
		SendAt:     sendAt,
		RetryCount: 0,
	}
	redisMsg := &models.RedisMessage{
		Status:    "scheduled",
		UpdatedAt: time.Now(),
	}

	err := ns.cache.Set(ctx, id, redisMsg)
	if err != nil {
		return "", err
	}

	delay := time.Until(sendAt)
	if delay <= 0 {
		err := ns.queue.Publish(ctx, ns.queueNameReady, rabbitMsg)
		if err != nil {
			return id, fmt.Errorf("немедленная отправка не удалась: %w", err)
		}
		return id, nil
	}

	err = ns.queue.Publish(ctx, ns.queueNameDelayed, rabbitMsg)
	if err != nil {
		return id, fmt.Errorf("отложенная отправка сообщения не удалась: %w", err)
	}
	return id, nil

}
