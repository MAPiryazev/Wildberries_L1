package repository

import (
	"context"

	"calendar/internal/models"
)

// Storage интерфейс который описывает методы repository слоя
type Storage interface {
	CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	GetEvent(ctx context.Context, eventID int) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int) error
	EventsForDay(ctx context.Context, userID int, date string) ([]*models.Event, error)
	EventsForWeek(ctx context.Context, userID int, date string) ([]*models.Event, error)
	EventsForMonth(ctx context.Context, userID int, date string) ([]*models.Event, error)
}
