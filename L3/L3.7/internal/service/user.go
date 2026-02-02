package service

import (
	"context"
	"fmt"
	"log/slog"

	appErrors "github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/errors"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/repository"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/validator"
	"github.com/google/uuid"
)

type userService struct {
	userRepo repository.UserRepository
	log      *slog.Logger
}

func NewUserService(userRepo repository.UserRepository, log *slog.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		log:      log,
	}
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Warn("user not found", "id", id)
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Warn("user not found by email", "email", email)
		return nil, err
	}
	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	v := validator.New()
	v.ValidatePagination(limit, offset)
	if !v.IsValid() {
		s.log.Warn("invalid pagination params", "errors", v.ErrorMessage())
		return nil, fmt.Errorf("%s: %s", appErrors.ErrInvalidItemData, v.ErrorMessage())
	}

	users, err := s.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		s.log.Error("failed to list users", "err", err)
		return nil, err
	}

	s.log.Debug("users listed", "count", len(users))
	return users, nil
}

func (s *userService) IsPermitted(role models.Role, action string) bool {
	permissions := map[models.Role][]string{
		models.RoleAdmin: {
			"create_item",
			"read_item",
			"update_item",
			"delete_item",
			"view_history",
			"list_users",
		},
		models.RoleManager: {
			"create_item",
			"read_item",
			"update_item",
			"view_history",
		},
		models.RoleViewer: {
			"read_item",
			"view_history",
		},
		models.RoleAuditor: {
			"read_item",
			"view_history",
		},
	}

	allowed, exists := permissions[role]
	if !exists {
		s.log.Warn("unknown role", "role", role)
		return false
	}

	for _, perm := range allowed {
		if perm == action {
			return true
		}
	}

	s.log.Debug("permission denied", "role", role, "action", action)
	return false
}
