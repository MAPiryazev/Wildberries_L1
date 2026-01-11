package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
)

type ValidationError struct {
	Field   string
	Message string
}

type Validator struct {
	errors []ValidationError
}

func New() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() []ValidationError {
	return v.errors
}

func (v *Validator) ErrorMessage() string {
	if v.IsValid() {
		return ""
	}
	var sb strings.Builder
	for i, err := range v.errors {
		sb.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Message))
		if i < len(v.errors)-1 {
			sb.WriteString("; ")
		}
	}
	return sb.String()
}

func (v *Validator) addError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (v *Validator) ValidateCreateItemRequest(req *models.CreateItemRequest) {
	if req == nil {
		v.addError("request", "cannot be nil")
		return
	}

	v.validateItemName(req.Name)
	v.validateItemSKU(req.SKU)
	v.validateQuantity(req.Quantity)
	v.validateLocation(req.Location)
}

func (v *Validator) ValidateUpdateItemRequest(req *models.UpdateItemRequest) {
	if req == nil {
		v.addError("request", "cannot be nil")
		return
	}

	if req.Name != nil {
		v.validateItemName(*req.Name)
	}
	if req.SKU != nil {
		v.validateItemSKU(*req.SKU)
	}
	if req.Quantity != nil {
		v.validateQuantity(*req.Quantity)
	}
	if req.Location != nil {
		v.validateLocation(*req.Location)
	}
}

func (v *Validator) validateItemName(name string) {
	if len(strings.TrimSpace(name)) == 0 {
		v.addError("name", "cannot be empty")
		return
	}
	if len(name) < 1 || len(name) > 255 {
		v.addError("name", "must be between 1 and 255 characters")
	}
}

func (v *Validator) validateItemSKU(sku string) {
	if len(strings.TrimSpace(sku)) == 0 {
		v.addError("sku", "cannot be empty")
		return
	}
	if len(sku) < 1 || len(sku) > 50 {
		v.addError("sku", "must be between 1 and 50 characters")
		return
	}
	if !isSKUFormat(sku) {
		v.addError("sku", "must contain only alphanumeric characters and hyphens")
	}
}

func (v *Validator) validateQuantity(qty int) {
	if qty < 0 {
		v.addError("quantity", "must be non-negative")
	}
}

func (v *Validator) validateLocation(location string) {
	if len(strings.TrimSpace(location)) == 0 {
		v.addError("location", "cannot be empty")
		return
	}
	if len(location) < 1 || len(location) > 255 {
		v.addError("location", "must be between 1 and 255 characters")
	}
}

func isSKUFormat(sku string) bool {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
	return pattern.MatchString(sku)
}

func (v *Validator) ValidateUserRole(role models.Role) {
	validRoles := map[models.Role]bool{
		models.RoleAdmin:   true,
		models.RoleManager: true,
		models.RoleViewer:  true,
		models.RoleAuditor: true,
	}

	if !validRoles[role] {
		v.addError("role", fmt.Sprintf("%s is not a valid role", role))
	}
}

func (v *Validator) ValidatePagination(limit, offset int) {
	if limit <= 0 {
		v.addError("limit", "must be greater than 0")
	}
	if limit > 1000 {
		v.addError("limit", "must not exceed 1000")
	}
	if offset < 0 {
		v.addError("offset", "must be non-negative")
	}
}
