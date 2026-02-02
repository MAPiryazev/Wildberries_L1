package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/config"
	appErrors "github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/errors"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"
)

type userRepo struct {
	db  *dbpg.DB
	log *slog.Logger
}

func NewUserRepository(db *dbpg.DB, cfg *config.Config, log *slog.Logger) UserRepository {
	return &userRepo{
		db:  db,
		log: log,
	}
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	const query = `
		SELECT id, email, name, role, created_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		r.log.Warn("user not found", "id", id)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrItemNotFound, err)
	}

	return user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `
		SELECT id, email, name, role, created_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		r.log.Warn("user not found by email", "email", email)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrItemNotFound, err)
	}

	return user, nil
}

func (r *userRepo) GetAll(ctx context.Context, limit, offset int) ([]*models.User, error) {
	const query = `
		SELECT id, email, name, role, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		r.log.Error("failed to query users", "err", err)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Role,
			&user.CreatedAt,
		); err != nil {
			r.log.Error("failed to scan user", "err", err)
			return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("rows iteration error", "err", err)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}

	r.log.Debug("fetched users", "count", len(users), "limit", limit, "offset", offset)
	return users, nil
}
