package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"shortener/internal/models"
	"shortener/internal/repository"
)

// Service объединяет операции над ссылками и их аналитикой
type Service interface {
	// Ссылки
	CreateShort(ctx context.Context, original string, clientID *uuid.UUID, expiresAt *time.Time) (*models.ShortURL, error)
	Resolve(ctx context.Context, shortCode string) (*models.ShortURL, error)

	// Аналитика
	RecordClick(ctx context.Context, ev models.ClickEvent) error
	Count(ctx context.Context, shortCode string) (int64, error)
	Daily(ctx context.Context, shortCode string) ([]models.AggPoint, error)
	ByUserAgent(ctx context.Context, shortCode string) ([]models.AggPoint, error)
}

type service struct {
	urls         repository.UrlRepository
	analytics    repository.AnalyticsRepository
	genShortCode func() (string, error)
}

func NewService(urls repository.UrlRepository, analytics repository.AnalyticsRepository, genShortCode func() (string, error)) Service {
	return &service{urls: urls, analytics: analytics, genShortCode: genShortCode}
}

func (s *service) CreateShort(ctx context.Context, original string, clientID *uuid.UUID, expiresAt *time.Time) (*models.ShortURL, error) {
	// Генерируем короткий код и пробуем вставить. В редких коллизиях можно попробовать ещё раз.
	const maxAttempts = 3
	var id int64
	var shortCode string
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		shortCode, err = s.genShortCode()
		if err != nil {
			return nil, err
		}
		id, err = s.urls.Create(ctx, original, shortCode, clientID, expiresAt)
		if err == nil {
			break
		}
		// возможна коллизия по unique(short_code) — попробуем другой код
		if attempt == maxAttempts-1 {
			return nil, err
		}
	}
	return &models.ShortURL{
		ID:        id,
		Original:  original,
		ShortCode: shortCode,
		ClientID:  clientID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

func (s *service) Resolve(ctx context.Context, shortCode string) (*models.ShortURL, error) {
	return s.urls.GetByShortCode(ctx, shortCode)
}

func (s *service) RecordClick(ctx context.Context, ev models.ClickEvent) error {
	if ev.At.IsZero() {
		ev.At = time.Now()
	}
	return s.analytics.InsertClick(ctx, ev)
}

func (s *service) Count(ctx context.Context, shortCode string) (int64, error) {
	return s.analytics.CountByShortCode(ctx, shortCode)
}

func (s *service) Daily(ctx context.Context, shortCode string) ([]models.AggPoint, error) {
	return s.analytics.AggregateDaily(ctx, shortCode)
}

func (s *service) ByUserAgent(ctx context.Context, shortCode string) ([]models.AggPoint, error) {
	return s.analytics.AggregateByUserAgent(ctx, shortCode)
}
