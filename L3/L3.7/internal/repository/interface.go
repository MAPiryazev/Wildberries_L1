package repository

import (
	"context"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
	"github.com/google/uuid"
)

type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) (*models.Item, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Item, error)
	Update(ctx context.Context, item *models.Item) (*models.Item, error)
	Delete(ctx context.Context, id uuid.UUID) error

	GetHistory(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]*models.ItemHistory, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.User, error)
}
