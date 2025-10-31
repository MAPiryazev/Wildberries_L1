package repository

import (
	"context"
	"shortener/internal/models"
	"time"

	"github.com/google/uuid"
)

// UrlRepository отвечает за операции с короткими ссылками в Postgres
type UrlRepository interface {
	Create(ctx context.Context, original string, shortCode string, clientID *uuid.UUID, expiresAt *time.Time) (int64, error)
	GetByShortCode(ctx context.Context, shortCode string) (*models.ShortURL, error)
	ExistsByOriginalURL(ctx context.Context, original string) (bool, error)
}

// AnalyticsRepository отвечает за запись и агрегацию кликов в ClickHouse
type AnalyticsRepository interface {
	InsertClick(ctx context.Context, e models.ClickEvent) error
	CountByShortCode(ctx context.Context, shortCode string) (int64, error)
	AggregateDaily(ctx context.Context, shortCode string) ([]models.AggPoint, error)
	AggregateByUserAgent(ctx context.Context, shortCode string) ([]models.AggPoint, error)
}
