package repository

import (
	"context"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/models"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	GetEventByID(ctx context.Context, eventID int64) (*models.Event, error)
	GetAllEvents(ctx context.Context) ([]*models.Event, error)

	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)

	CreateBooking(ctx context.Context, booking *models.Booking) (*models.Booking, error)
	GetBooking(ctx context.Context, eventID, userID int64) (*models.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID int64, status string) error

	CountActiveBookings(ctx context.Context, eventID int64) (int64, error)

	FindExpiredBookings(ctx context.Context) ([]*models.Booking, error)
}
