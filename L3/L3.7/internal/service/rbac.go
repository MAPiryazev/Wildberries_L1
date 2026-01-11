package service

import (
	"context"
	"fmt"

	appErrors "github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/errors"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
	"github.com/google/uuid"
)

type contextKey string

const (
	userContextKey contextKey = "user"
	roleContextKey contextKey = "role"
)

type AuthContext struct {
	UserID uuid.UUID
	Role   models.Role
	Email  string
}

func SetUserContext(ctx context.Context, auth *AuthContext) context.Context {
	ctx = context.WithValue(ctx, userContextKey, auth.UserID)
	ctx = context.WithValue(ctx, roleContextKey, auth.Role)
	return ctx
}

func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	val := ctx.Value(userContextKey)
	if id, ok := val.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

func GetRoleFromContext(ctx context.Context) models.Role {
	val := ctx.Value(roleContextKey)
	if role, ok := val.(models.Role); ok {
		return role
	}
	return ""
}

func RequirePermission(role models.Role, action string) error {
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
		return fmt.Errorf("%s: invalid role %s", appErrors.ErrForbidden, role)
	}

	for _, perm := range allowed {
		if perm == action {
			return nil
		}
	}

	return fmt.Errorf("%s: action %s not allowed for role %s", appErrors.ErrForbidden, action, role)
}
