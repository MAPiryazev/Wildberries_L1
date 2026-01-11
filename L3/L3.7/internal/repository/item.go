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

type itemRepo struct {
	db  *dbpg.DB
	log *slog.Logger
}

func NewItemRepository(db *dbpg.DB, cfg *config.Config, log *slog.Logger) ItemRepository {
	return &itemRepo{
		db:  db,
		log: log,
	}
}

func (r *itemRepo) Create(ctx context.Context, item *models.Item) (*models.Item, error) {
	const query = `
		INSERT INTO items (id, name, sku, quantity, reserved_qty, location, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, sku, quantity, reserved_qty, location, created_by, updated_by, created_at, updated_at
	`

	newItem := &models.Item{}
	err := r.db.QueryRowContext(ctx, query,
		item.ID,
		item.Name,
		item.SKU,
		item.Quantity,
		item.ReservedQty,
		item.Location,
		item.CreatedBy,
		item.UpdatedBy,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(
		&newItem.ID,
		&newItem.Name,
		&newItem.SKU,
		&newItem.Quantity,
		&newItem.ReservedQty,
		&newItem.Location,
		&newItem.CreatedBy,
		&newItem.UpdatedBy,
		&newItem.CreatedAt,
		&newItem.UpdatedAt,
	)

	if err != nil {
		r.log.Error("failed to create item", "err", err, "sku", item.SKU)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}

	r.log.Info("item created", "item_id", newItem.ID, "sku", newItem.SKU, "created_by", newItem.CreatedBy)
	return newItem, nil
}

func (r *itemRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	const query = `
		SELECT id, name, sku, quantity, reserved_qty, location, created_by, updated_by, created_at, updated_at
		FROM items
		WHERE id = $1
	`

	item := &models.Item{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.SKU,
		&item.Quantity,
		&item.ReservedQty,
		&item.Location,
		&item.CreatedBy,
		&item.UpdatedBy,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		r.log.Warn("item not found", "id", id)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrItemNotFound, err)
	}

	return item, nil
}

func (r *itemRepo) GetAll(ctx context.Context, limit, offset int) ([]*models.Item, error) {
	const query = `
		SELECT id, name, sku, quantity, reserved_qty, location, created_by, updated_by, created_at, updated_at
		FROM items
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		r.log.Error("failed to query items", "err", err)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		item := &models.Item{}
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.SKU,
			&item.Quantity,
			&item.ReservedQty,
			&item.Location,
			&item.CreatedBy,
			&item.UpdatedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			r.log.Error("failed to scan item", "err", err)
			return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("rows iteration error", "err", err)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}

	r.log.Debug("fetched items", "count", len(items), "limit", limit, "offset", offset)
	return items, nil
}

func (r *itemRepo) Update(ctx context.Context, item *models.Item) (*models.Item, error) {
	const query = `
		UPDATE items
		SET name = $1, sku = $2, quantity = $3, reserved_qty = $4, location = $5, updated_by = $6, updated_at = $7
		WHERE id = $8
		RETURNING id, name, sku, quantity, reserved_qty, location, created_by, updated_by, created_at, updated_at
	`

	updated := &models.Item{}
	err := r.db.QueryRowContext(ctx, query,
		item.Name,
		item.SKU,
		item.Quantity,
		item.ReservedQty,
		item.Location,
		item.UpdatedBy,
		item.UpdatedAt,
		item.ID,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.SKU,
		&updated.Quantity,
		&updated.ReservedQty,
		&updated.Location,
		&updated.CreatedBy,
		&updated.UpdatedBy,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		r.log.Error("failed to update item", "err", err, "id", item.ID)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}

	r.log.Info("item updated", "item_id", updated.ID, "updated_by", updated.UpdatedBy)
	return updated, nil
}

func (r *itemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM items WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Error("failed to delete item", "err", err, "id", id)
		return fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		r.log.Error("failed to get affected rows", "err", err)
		return fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}

	if affected == 0 {
		r.log.Warn("item not found for deletion", "id", id)
		return fmt.Errorf("%s: item with id %s not found", appErrors.ErrItemNotFound, id)
	}

	r.log.Info("item deleted", "item_id", id)
	return nil
}

func (r *itemRepo) GetHistory(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]*models.ItemHistory, error) {
	const query = `
		SELECT id, item_id, changed_by, action, changes, created_at
		FROM item_history
		WHERE item_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, itemID, limit, offset)
	if err != nil {
		r.log.Error("failed to query item history", "err", err, "item_id", itemID)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}
	defer rows.Close()

	var history []*models.ItemHistory
	for rows.Next() {
		h := &models.ItemHistory{}
		if err := rows.Scan(
			&h.ID,
			&h.ItemID,
			&h.ChangedBy,
			&h.Action,
			&h.Changes,
			&h.CreatedAt,
		); err != nil {
			r.log.Error("failed to scan history record", "err", err)
			return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
		}
		history = append(history, h)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("history rows iteration error", "err", err)
		return nil, fmt.Errorf("%s: %w", appErrors.ErrDatabaseQuery, err)
	}

	r.log.Debug("fetched item history", "item_id", itemID, "count", len(history))
	return history, nil
}
