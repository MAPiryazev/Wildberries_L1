package models

import "time"

// RabbitMQMessage то что будет в redis о сообщении
type RabbitMQMessage struct {
	ID         string
	Payload    string
	SendAt     time.Time
	RetryCount int
}

// RedisMessage то что будет в redis о сообщении
type RedisMessage struct {
	ID        string
	Status    string
	UpdatedAt time.Time
}
