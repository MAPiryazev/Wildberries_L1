package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/repository"
)

type DefaultService struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *DefaultService {
	return &DefaultService{repo: repo}
}

// CreateEvent создаёт новое событие (только админ) с проверкой бизнес-логики
func (s *DefaultService) CreateEvent(ctx context.Context, event *models.Event, isAdmin bool) (*models.Event, error) {
	if !isAdmin {
		return nil, fmt.Errorf("только администратор может создавать событие")
	}
	if event == nil {
		return nil, fmt.Errorf("event не может быть nil")
	}
	if event.Capacity <= 0 {
		return nil, fmt.Errorf("количество мест на событии должно быть больше 0")
	}
	if event.StartTime.Before(time.Now()) {
		return nil, fmt.Errorf("время начала события не может быть в прошлом")
	}
	if event.Title == "" {
		return nil, fmt.Errorf("название события не может быть пустым")
	}

	return s.repo.CreateEvent(ctx, event)
}

// GetEventByID возвращает событие по ID с проверкой существования
func (s *DefaultService) GetEventByID(ctx context.Context, eventID int64) (*models.Event, error) {
	if eventID <= 0 {
		return nil, fmt.Errorf("невалидный id события")
	}

	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("событие с id %d не найдено", eventID)
	}
	return event, nil
}

// GetAllEvents возвращает список всех событий
func (s *DefaultService) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repo.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var upcoming []*models.Event
	for _, e := range events {
		if e.StartTime.After(now) {
			upcoming = append(upcoming, e)
		}
	}

	return upcoming, nil
}

// CountFreePlaces возвращает количество свободных мест для события
func (s *DefaultService) CountFreePlaces(ctx context.Context, eventID int64) (int64, error) {
	event, err := s.GetEventByID(ctx, eventID)
	if err != nil {
		return 0, err
	}

	activeBookings, err := s.repo.CountActiveBookings(ctx, eventID)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчёта активных бронирований: %w", err)
	}

	freePlaces := event.Capacity - activeBookings
	if freePlaces < 0 {
		freePlaces = 0
	}

	return freePlaces, nil
}

func (s *DefaultService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user не может быть nil")
	}
	if user.Name == "" {
		return nil, fmt.Errorf("имя пользователя не может быть пустым")
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *DefaultService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("невалидный id пользователя")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("пользователь с id %d не найден", userID)
	}

	return user, nil
}

func (s *DefaultService) IsBookedByUserID(ctx context.Context, eventID, userID int64) (bool, *models.Booking, error) {
	if eventID <= 0 || userID <= 0 {
		return false, nil, fmt.Errorf("невалидный id события или пользователя")
	}

	booking, err := s.repo.GetBooking(ctx, eventID, userID)
	if err != nil {
		return false, nil, fmt.Errorf("ошибка получения брони: %w", err)
	}
	if booking == nil {
		return false, nil, nil // бронь не найдена
	}

	if booking.Status == "booked" || booking.Status == "confirmed" {
		return true, booking, nil
	}

	return false, booking, nil
}

func (s *DefaultService) CreateBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error) {
	if booking == nil {
		return nil, fmt.Errorf("booking не может быть nil")
	}
	if booking.UserID <= 0 {
		return nil, fmt.Errorf("невалидный id пользователя")
	}
	if booking.EventID <= 0 {
		return nil, fmt.Errorf("невалидный id события")
	}

	_, err := s.GetUserByID(ctx, booking.UserID)
	if err != nil {
		return nil, err
	}

	event, err := s.GetEventByID(ctx, booking.EventID)
	if err != nil {
		return nil, err
	}

	if event.StartTime.Before(time.Now()) {
		return nil, fmt.Errorf("нельзя бронировать события в прошлом")
	}

	freePlaces, err := s.CountFreePlaces(ctx, booking.EventID)
	if err != nil {
		return nil, err
	}
	if freePlaces <= 0 {
		return nil, fmt.Errorf("нет свободных мест на событие %d", booking.EventID)
	}

	if booking.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("время истечения брони должно быть в будущем")
	}

	booking.Status = "booked"
	return s.repo.CreateBooking(ctx, booking)
}

func (s *DefaultService) GetBooking(ctx context.Context, eventID, userID int64) (*models.Booking, error) {
	if eventID <= 0 || userID <= 0 {
		return nil, fmt.Errorf("невалидный id события или пользователя")
	}

	booking, err := s.repo.GetBooking(ctx, eventID, userID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, fmt.Errorf("бронь не найдена для события %d и пользователя %d", eventID, userID)
	}

	return booking, nil
}

func (s *DefaultService) UpdateBookingStatus(ctx context.Context, bookingID int64, status string, isAdmin bool) error {
	if bookingID <= 0 {
		return fmt.Errorf("невалидный id брони")
	}
	if !validBookingStatuses[status] {
		return fmt.Errorf("недопустимый статус брони: %s", status)
	}

	booking, err := s.repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		return fmt.Errorf("бронь с id %d не найдена", bookingID)
	}

	// Только админ может менять на cancelled или confirmed, обычный пользователь может подтверждать только свою бронь
	if !isAdmin && status != "confirmed" {
		return fmt.Errorf("только администратор может менять статус на '%s'", status)
	}

	return s.repo.UpdateBookingStatus(ctx, bookingID, status)
}

func (s *DefaultService) GetBookingsByEventID(ctx context.Context, eventID int64) ([]*models.Booking, error) {
	if eventID <= 0 {
		return nil, fmt.Errorf("невалидный id события")
	}

	bookings, err := s.repo.GetBookingsByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *DefaultService) CancelExpiredBookings(ctx context.Context) ([]*models.Booking, error) {
	return s.repo.CancelExpiredBookings(ctx)
}
