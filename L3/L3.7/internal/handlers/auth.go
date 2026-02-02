package handlers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/errors"
	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/service"
	jwt "github.com/golang-jwt/jwt/v5"
)

type loginRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type tokenResponse struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("failed to decode login request", "err", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Role == "" {
		h.log.Warn("missing required fields", "email", req.Email, "role", req.Role)
		respondError(w, http.StatusBadRequest, "email and role required")
		return
	}

	user, err := h.userService.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		h.log.Warn("user not found", "email", req.Email)
		respondError(w, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	if string(user.Role) != req.Role {
		h.log.Warn("role mismatch", "email", req.Email, "provided", req.Role, "stored", user.Role)
		respondError(w, http.StatusUnauthorized, "invalid role")
		return
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		h.log.Error("failed to sign token", "err", err)
		respondError(w, http.StatusInternalServerError, appErrors.ErrInternal)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenResponse{
		Token: tokenStr,
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	})
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := service.GetUserIDFromContext(ctx)
	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		h.log.Warn("user not found", "id", userID)
		respondError(w, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	resp := userResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role := service.GetRoleFromContext(ctx)
	if err := service.RequirePermission(role, "list_users"); err != nil {
		h.log.Warn("permission denied", "action", "list_users", "role", role)
		respondError(w, http.StatusForbidden, appErrors.ErrForbidden)
		return
	}

	users, err := h.userService.ListUsers(ctx, 100, 0)
	if err != nil {
		h.log.Error("failed to list users", "err", err)
		respondError(w, http.StatusInternalServerError, appErrors.ErrInternal)
		return
	}

	var resp []userResponse
	for _, u := range users {
		resp = append(resp, userResponse{
			ID:    u.ID.String(),
			Email: u.Email,
			Name:  u.Name,
			Role:  string(u.Role),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
