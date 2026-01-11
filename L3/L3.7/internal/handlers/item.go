package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	appErrors "github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/errors"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/service"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := service.GetUserIDFromContext(ctx)
	role := service.GetRoleFromContext(ctx)

	if err := service.RequirePermission(role, "create_item"); err != nil {
		h.log.Warn("permission denied", "action", "create_item", "role", role)
		respondError(w, http.StatusForbidden, appErrors.ErrForbidden)
		return
	}

	var req models.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("failed to decode request", "err", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	v := validator.New()
	v.ValidateCreateItemRequest(&req)
	if !v.IsValid() {
		h.log.Warn("validation failed", "errors", v.ErrorMessage())
		respondError(w, http.StatusBadRequest, v.ErrorMessage())
		return
	}

	item, err := h.itemService.CreateItem(ctx, userID, &req)
	if err != nil {
		h.log.Error("failed to create item", "err", err)
		respondError(w, http.StatusInternalServerError, appErrors.ErrInternal)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role := service.GetRoleFromContext(ctx)

	if err := service.RequirePermission(role, "read_item"); err != nil {
		h.log.Warn("permission denied", "action", "read_item", "role", role)
		respondError(w, http.StatusForbidden, appErrors.ErrForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warn("invalid item id", "id", idStr)
		respondError(w, http.StatusBadRequest, appErrors.ErrInvalidItemID)
		return
	}

	item, err := h.itemService.GetItem(ctx, id)
	if err != nil {
		if errors.Is(err, errors.New(appErrors.ErrItemNotFound)) ||
			errors.Is(err, context.DeadlineExceeded) {
			h.log.Warn("item not found", "id", id)
			respondError(w, http.StatusNotFound, appErrors.ErrItemNotFound)
			return
		}
		h.log.Error("failed to get item", "err", err)
		respondError(w, http.StatusInternalServerError, appErrors.ErrInternal)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role := service.GetRoleFromContext(ctx)

	if err := service.RequirePermission(role, "read_item"); err != nil {
		h.log.Warn("permission denied", "action", "read_item", "role", role)
		respondError(w, http.StatusForbidden, appErrors.ErrForbidden)
		return
	}

	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	v := validator.New()
	v.ValidatePagination(limit, offset)
	if !v.IsValid() {
		h.log.Warn("invalid pagination", "errors", v.ErrorMessage())
		respondError(w, http.StatusBadRequest, v.ErrorMessage())
		return
	}

	items, err := h.itemService.ListItems(ctx, limit, offset)
	if err != nil {
		h.log.Error("failed to list items", "err", err)
		respondError(w, http.StatusInternalServerError, appErrors.ErrInternal)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := service.GetUserIDFromContext(ctx)
	role := service.GetRoleFromContext(ctx)

	if err := service.RequirePermission(role, "update_item"); err != nil {
		h.log.Warn("permission denied", "action", "update_item", "role", role)
		respondError(w, http.StatusForbidden, appErrors.ErrForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warn("invalid item id", "id", idStr)
		respondError(w, http.StatusBadRequest, appErrors.ErrInvalidItemID)
		return
	}

	var req models.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("failed to decode request", "err", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	v := validator.New()
	v.ValidateUpdateItemRequest(&req)
	if !v.IsValid() {
		h.log.Warn("validation failed", "errors", v.ErrorMessage())
		respondError(w, http.StatusBadRequest, v.ErrorMessage())
		return
	}

	item, err := h.itemService.UpdateItem(ctx, userID, id, &req)
	if err != nil {
		h.log.Error("failed to update item", "err", err)
		respondError(w, http.StatusInternalServerError, appErrors.ErrInternal)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := service.GetUserIDFromContext(ctx)
	role := service.GetRoleFromContext(ctx)

	if err := service.RequirePermission(role, "delete_item"); err != nil {
		h.log.Warn("permission denied", "action", "delete_item", "role", role)
		respondError(w, http.StatusForbidden, appErrors.ErrForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warn("invalid item id", "id", idStr)
		respondError(w, http.StatusBadRequest, appErrors.ErrInvalidItemID)
		return
	}

	if err := h.itemService.DeleteItem(ctx, userID, id); err != nil {
		h.log.Error("failed to delete item", "err", err)
		respondError(w, http.StatusInternalServerError, appErrors.ErrInternal)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetItemHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role := service.GetRoleFromContext(ctx)

	if err := service.RequirePermission(role, "view_history"); err != nil {
		h.log.Warn("permission denied", "action", "view_history", "role", role)
		respondError(w, http.StatusForbidden, appErrors.ErrForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warn("invalid item id", "id", idStr)
		respondError(w, http.StatusBadRequest, appErrors.ErrInvalidItemID)
		return
	}

	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	v := validator.New()
	v.ValidatePagination(limit, offset)
	if !v.IsValid() {
		h.log.Warn("invalid pagination", "errors", v.ErrorMessage())
		respondError(w, http.StatusBadRequest, v.ErrorMessage())
		return
	}

	history, err := h.itemService.GetItemHistory(ctx, id, limit, offset)
	if err != nil {
		h.log.Error("failed to get item history", "err", err)
		respondError(w, http.StatusInternalServerError, appErrors.ErrInternal)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(history)
}
