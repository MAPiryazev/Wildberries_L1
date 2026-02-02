package service

import (
	"context"
	"fmt"
	"time"

	"calendar/internal/models"
	"calendar/internal/repository"
)

// Service - методы для сервиса
type Service interface {
	CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID int) error
	EventsForDay(ctx context.Context, userID int, dateStr string) ([]*models.Event, error)
	EventsForWeek(ctx context.Context, userID int, dateStr string) ([]*models.Event, error)
	EventsForMonth(ctx context.Context, userID int, dateStr string) ([]*models.Event, error)
}

// DefaultService стандартная реализация сервиса
type DefaultService struct {
	repo repository.Storage
}

// NewDefaultService создаёт новый DefaultService с указанным репозиторием.
func NewDefaultService(repo repository.Storage) (*DefaultService, error) {
	if repo == nil {
		return nil, fmt.Errorf("repo service не должен быть nil")
	}
	return &DefaultService{
		repo: repo,
	}, nil
}

// CreateEvent создаёт новое событие и проверяет, что дата не из прошлого.
func (d *DefaultService) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	if isPastDateHelper(event.Date) {
		return nil, fmt.Errorf("дата события не может быть в прошлом")
	}
	return d.repo.CreateEvent(ctx, event)
}

// UpdateEvent обновляет существующее событие после проверки, что оно существует.
func (d *DefaultService) UpdateEvent(ctx context.Context, event *models.Event) error {
	tempEvent, err := d.repo.GetEvent(ctx, event.ID)
	if tempEvent == nil || err != nil {
		return fmt.Errorf("изменяемое событие не найдено или произошла ошибка при поиске")
	}
	return d.repo.UpdateEvent(ctx, event)
}

// DeleteEvent удаляет событие по ID.
func (d *DefaultService) DeleteEvent(ctx context.Context, eventID int) error {
	return d.repo.DeleteEvent(ctx, eventID)
}

// EventsForDay возвращает список событий пользователя за указанный день.
func (d *DefaultService) EventsForDay(ctx context.Context, userID int, dateStr string) ([]*models.Event, error) {
	return d.repo.EventsForDay(ctx, userID, dateStr)
}

// EventsForWeek возвращает список событий пользователя за указанную неделю.
func (d *DefaultService) EventsForWeek(ctx context.Context, userID int, dateStr string) ([]*models.Event, error) {
	return d.repo.EventsForWeek(ctx, userID, dateStr)
}

// EventsForMonth возвращает список событий пользователя за указанный месяц.
func (d *DefaultService) EventsForMonth(ctx context.Context, userID int, dateStr string) ([]*models.Event, error) {
	return d.repo.EventsForMonth(ctx, userID, dateStr)
}

// isPastDateHelper проверяет, находится ли дата в прошлом (с учётом дня без времени).
func isPastDateHelper(t time.Time) bool {
	now := time.Now().Truncate(24 * time.Hour)
	date := t.Truncate(24 * time.Hour)
	return date.Before(now)
}
