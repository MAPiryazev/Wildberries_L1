package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleViewer  Role = "viewer"
	RoleAuditor Role = "auditor"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Role      Role      `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

type Item struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	SKU         string    `db:"sku"`
	Quantity    int       `db:"quantity"`
	ReservedQty int       `db:"reserved_qty"`
	Location    string    `db:"location"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	CreatedBy   uuid.UUID `db:"created_by"`
	UpdatedBy   uuid.UUID `db:"updated_by"`
}

type ItemHistory struct {
	ID        uuid.UUID `db:"id"`
	ItemID    uuid.UUID `db:"item_id"`
	ChangedBy uuid.UUID `db:"changed_by"`
	Action    string    `db:"action"`
	Changes   string    `db:"changes"`
	CreatedAt time.Time `db:"created_at"`
}

type AuditLog struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Action    string    `db:"action"`
	Entity    string    `db:"entity"`
	EntityID  uuid.UUID `db:"entity_id"`
	OldValue  string    `db:"old_value"`
	NewValue  string    `db:"new_value"`
	IPAddress string    `db:"ip_address"`
	CreatedAt time.Time `db:"created_at"`
}

type CreateItemRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	SKU      string `json:"sku" validate:"required,min=1,max=50"`
	Quantity int    `json:"quantity" validate:"required,min=0"`
	Location string `json:"location" validate:"required,min=1,max=255"`
}

type UpdateItemRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	SKU      *string `json:"sku,omitempty" validate:"omitempty,min=1,max=50"`
	Quantity *int    `json:"quantity,omitempty" validate:"omitempty,min=0"`
	Location *string `json:"location,omitempty" validate:"omitempty,min=1,max=255"`
}

type ItemHistoryResponse struct {
	ID        uuid.UUID `json:"id"`
	ItemID    uuid.UUID `json:"item_id"`
	ChangedBy uuid.UUID `json:"changed_by"`
	Action    string    `json:"action"`
	Changes   string    `json:"changes"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}
