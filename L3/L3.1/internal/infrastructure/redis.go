package infrastructure

import (
	"context"
	"encoding/json"
	"time"

	"L3.1/internal/models"
	"github.com/redis/go-redis/v9"
)

// RedisClient клиент для взаимодействия с redis
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient конструктор для структуры выше
func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisClient{client: rdb}
}

// Set добавляет сообщение в redis
func (rc *RedisClient) Set(ctx context.Context, key string, message *models.RedisMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return rc.client.Set(ctx, key, data, 0).Err()
}

// SetWithTTL добавляет сообщение в redis с времененем жизни
func (rc *RedisClient) SetWithTTL(ctx context.Context, key string, message *models.RedisMessage, ttlSeconds int) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	// после истечения TTL redis удалит сообщение автоматически
	return rc.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

// Get получает сообщение из redis
func (rc *RedisClient) Get(ctx context.Context, key string) (*models.RedisMessage, error) {
	var msg models.RedisMessage
	data, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		return &msg, err
	}
	err = json.Unmarshal([]byte(data), &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Exists проверяет есть ли сообщение с указанным id в redis
func (rc *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, err
}

// Delete удаляет сообщение из redis
func (rc *RedisClient) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}
