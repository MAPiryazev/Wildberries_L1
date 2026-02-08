package service

import (
	"context"
	"errors"
	"server-calendar/internal/storage/entity"
	"server-calendar/internal/worker"
	"time"
)

type CalendarService struct {
	repo     repository
	reminder *worker.ReminderWorker
}

type repository interface {
	CreateEvent(event entity.Event) error
	UpdateEvent(event entity.Event) error
	DeleteEvent(event entity.Event) error
	GetEventsByDateRange(userID entity.UserID, from, to time.Time) ([]entity.Event, error)
}

func NewCalendarService(repo repository, rw *worker.ReminderWorker) *CalendarService {
	return &CalendarService{
		repo:     repo,
		reminder: rw,
	}
}

func (s *CalendarService) CreateEvent(_ context.Context, e entity.Event) error {
	if e.Date.Before(time.Now()) {
		return errors.New("cannot create event in the past")
	}

	if err := s.repo.CreateEvent(e); err != nil {
		return err
	}

	if e.RemindAt != nil {
		s.reminder.Add(worker.Reminder{
			EventID:  int(e.EventID),
			UserID:   int(e.UserID),
			Title:    e.Title,
			RemindAt: *e.RemindAt,
		})
	}

	return nil
}

func (s *CalendarService) UpdateEvent(_ context.Context, e entity.Event) error {
	if e.Date.Before(time.Now()) {
		return errors.New("cannot update event in the past")
	}
	return s.repo.UpdateEvent(e)
}

func (s *CalendarService) DeleteEvent(_ context.Context, id entity.Event) error {
	return s.repo.DeleteEvent(id)
}

func (s *CalendarService) EventsForDay(_ context.Context, userID entity.UserID) ([]entity.Event, error) {
	from := time.Now()
	to := from.
		AddDate(0, 0, 1).
		Truncate(24 * time.Hour)

	return s.repo.GetEventsByDateRange(userID, from, to)
}

func (s *CalendarService) EventsForWeek(_ context.Context, userID entity.UserID) ([]entity.Event, error) {
	from := time.Now()
	to := from.
		AddDate(0, 0, 7).
		Truncate(24 * time.Hour)
	return s.repo.GetEventsByDateRange(userID, from, to)
}

func (s *CalendarService) EventsForMonth(_ context.Context, userID entity.UserID) ([]entity.Event, error) {
	from := time.Now()
	to := from.
		AddDate(0, 1, 0).
		Truncate(24 * time.Hour)
	return s.repo.GetEventsByDateRange(userID, from, to)
}
