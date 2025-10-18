package repository

import (
	"context"

	"calendar/internal/models"
)

// Storage интерфейс который описывает методы repository слоя
type Storage interface {
	CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	GetEvent(ctx context.Context, eventId int) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventId int) error
	EventsForDay(ctx context.Context, userId int, date string) ([]*models.Event, error)
	EventsForWeek(ctx context.Context, userId int, date string) ([]*models.Event, error)
	EventsForMonth(ctx context.Context, userId int, date string) ([]*models.Event, error)
}
