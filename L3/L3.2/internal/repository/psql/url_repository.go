package psql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"shortener/internal/models"
)

// Repo реализация UrlRepository на Postgres (таблица shortcuts)
type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, original string, shortCode string, clientID *uuid.UUID, expiresAt *time.Time) (int64, error) {
	const q = `
        insert into shortcuts (original_url, short_code, client_id, expires_at)
        values ($1, $2, $3, $4)
        returning id;
    `
	var id int64
	err := r.pool.QueryRow(ctx, q, original, shortCode, clientID, expiresAt).Scan(&id)
	return id, err
}

func (r *Repo) GetByShortCode(ctx context.Context, shortCode string) (*models.ShortURL, error) {
	const q = `
        select id, original_url, short_code, client_id, expires_at, created_at
        from shortcuts
        where short_code = $1
        limit 1;
    `
	row := r.pool.QueryRow(ctx, q, shortCode)
	var s models.ShortURL
	var clientID *uuid.UUID
	var expiresAt *time.Time
	if err := row.Scan(&s.ID, &s.Original, &s.ShortCode, &clientID, &expiresAt, &s.CreatedAt); err != nil {
		return nil, err
	}
	s.ClientID = clientID
	s.ExpiresAt = expiresAt
	return &s, nil
}

func (r *Repo) ExistsByOriginalURL(ctx context.Context, original string) (bool, error) {
	const q = `select exists(select 1 from shortcuts where original_url = $1);`
	var exists bool
	if err := r.pool.QueryRow(ctx, q, original).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
