package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	appErrors "github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/errors"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/repository"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/validator"
	"github.com/google/uuid"
)

type itemService struct {
	itemRepo repository.ItemRepository
	userRepo repository.UserRepository
	log      *slog.Logger
}

func NewItemService(itemRepo repository.ItemRepository, userRepo repository.UserRepository, log *slog.Logger) ItemService {
	return &itemService{
		itemRepo: itemRepo,
		userRepo: userRepo,
		log:      log,
	}
}

func (s *itemService) CreateItem(ctx context.Context, userID uuid.UUID, req *models.CreateItemRequest) (*models.Item, error) {
	v := validator.New()
	v.ValidateCreateItemRequest(req)
	if !v.IsValid() {
		s.log.Warn("invalid create item request", "errors", v.ErrorMessage())
		return nil, fmt.Errorf("%s: %s", appErrors.ErrInvalidItemData, v.ErrorMessage())
	}

	now := time.Now()
	item := &models.Item{
		ID:        uuid.New(),
		Name:      req.Name,
		SKU:       req.SKU,
		Quantity:  req.Quantity,
		Location:  req.Location,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	created, err := s.itemRepo.Create(ctx, item)
	if err != nil {
		s.log.Error("failed to create item in repository", "err", err, "user_id", userID)
		return nil, err
	}

	s.log.Info("item created successfully", "item_id", created.ID, "sku", created.SKU, "user_id", userID)
	return created, nil
}

func (s *itemService) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Warn("item not found", "id", id)
		return nil, err
	}
	return item, nil
}

func (s *itemService) ListItems(ctx context.Context, limit, offset int) ([]*models.Item, error) {
	v := validator.New()
	v.ValidatePagination(limit, offset)
	if !v.IsValid() {
		s.log.Warn("invalid pagination params", "errors", v.ErrorMessage())
		return nil, fmt.Errorf("%s: %s", appErrors.ErrInvalidItemData, v.ErrorMessage())
	}

	items, err := s.itemRepo.GetAll(ctx, limit, offset)
	if err != nil {
		s.log.Error("failed to list items", "err", err)
		return nil, err
	}

	s.log.Debug("items listed", "count", len(items), "limit", limit, "offset", offset)
	return items, nil
}

func (s *itemService) UpdateItem(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *models.UpdateItemRequest) (*models.Item, error) {
	v := validator.New()
	v.ValidateUpdateItemRequest(req)
	if !v.IsValid() {
		s.log.Warn("invalid update item request", "errors", v.ErrorMessage())
		return nil, fmt.Errorf("%s: %s", appErrors.ErrInvalidItemData, v.ErrorMessage())
	}

	current, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Warn("item not found for update", "id", id)
		return nil, err
	}

	if req.Name != nil {
		current.Name = *req.Name
	}
	if req.SKU != nil {
		current.SKU = *req.SKU
	}
	if req.Quantity != nil {
		current.Quantity = *req.Quantity
	}
	if req.Location != nil {
		current.Location = *req.Location
	}

	current.UpdatedAt = time.Now()
	current.UpdatedBy = userID

	updated, err := s.itemRepo.Update(ctx, current)
	if err != nil {
		s.log.Error("failed to update item in repository", "err", err, "item_id", id, "user_id", userID)
		return nil, err
	}

	s.log.Info("item updated successfully", "item_id", updated.ID, "updated_by", userID)
	return updated, nil
}

func (s *itemService) DeleteItem(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	_, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Warn("item not found for deletion", "id", id)
		return err
	}

	if err := s.itemRepo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete item", "err", err, "item_id", id, "user_id", userID)
		return err
	}

	s.log.Info("item deleted successfully", "item_id", id, "deleted_by", userID)
	return nil
}

func (s *itemService) GetItemHistory(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]*models.ItemHistory, error) {
	v := validator.New()
	v.ValidatePagination(limit, offset)
	if !v.IsValid() {
		s.log.Warn("invalid pagination params for history", "errors", v.ErrorMessage())
		return nil, fmt.Errorf("%s: %s", appErrors.ErrInvalidItemData, v.ErrorMessage())
	}

	_, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		s.log.Warn("item not found for history", "id", itemID)
		return nil, err
	}

	history, err := s.itemRepo.GetHistory(ctx, itemID, limit, offset)
	if err != nil {
		s.log.Error("failed to fetch item history", "err", err, "item_id", itemID)
		return nil, err
	}

	s.log.Debug("item history fetched", "item_id", itemID, "count", len(history))
	return history, nil
}
