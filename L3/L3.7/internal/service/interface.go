package service

import (
	"context"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
	"github.com/google/uuid"
)

type ItemService interface {
	CreateItem(ctx context.Context, userID uuid.UUID, req *models.CreateItemRequest) (*models.Item, error)
	GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error)
	ListItems(ctx context.Context, limit, offset int) ([]*models.Item, error)
	UpdateItem(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *models.UpdateItemRequest) (*models.Item, error)
	DeleteItem(ctx context.Context, userID uuid.UUID, id uuid.UUID) error

	GetItemHistory(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]*models.ItemHistory, error)
}

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error)

	IsPermitted(role models.Role, action string) bool
}
