package errros

type Kind int

const (
	KindUnknown Kind = iota
	KindValidation
	KindNotFound
	KindConflict
	KindUnauthorized
	KindForbidden
	KindInternal
)

var (
	// Validation errors
	ErrItemNotFound      = "item not found"
	ErrItemAlreadyExists = "item already exists"
	ErrInvalidItemData   = "invalid item data"
	ErrInvalidItemID     = "invalid item id"
	ErrInvalidQuantity   = "quantity must be positive"

	// Authorization errors
	ErrUnauthorized = "unauthorized"
	ErrForbidden    = "forbidden"
	ErrInvalidRole  = "invalid role"

	// Database errors
	ErrDatabaseConnection  = "database connection error"
	ErrDatabaseQuery       = "database query error"
	ErrDatabaseTransaction = "database transaction error"

	// Internal errors
	ErrInternal        = "internal server error"
	ErrMigrationFailed = "migration failed"
)
