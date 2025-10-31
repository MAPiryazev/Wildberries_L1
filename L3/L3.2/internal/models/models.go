package models

import (
	"time"

	"github.com/google/uuid"
)

// ShortURL доменная модель хранения короткой ссылки (psql)
type ShortURL struct {
	ID        int64
	Original  string
	ShortCode string
	ClientID  *uuid.UUID
	ExpiresAt *time.Time
	CreatedAt time.Time
}

// ClickEvent событие перехода (ClickHouse)
type ClickEvent struct {
	ShortCode string
	ClientID  uuid.UUID
	UserAgent string
	IP        string
	At        time.Time
}

// AggPoint агрегированная точка для графиков/сводок
type AggPoint struct {
	Key   string
	Count int64
}
