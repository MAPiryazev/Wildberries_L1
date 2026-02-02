package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"shortener/internal/models"
	"shortener/internal/service"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()
	r.Post("/shorten", h.postShorten)
	r.Get("/s/{short}", h.getRedirect)
	r.Get("/analytics/{short}", h.getAnalytics)
	return r
}

type shortenRequest struct {
	OriginalURL string     `json:"original_url"`
	ClientID    *uuid.UUID `json:"client_id"`
	TTLSeconds  *int       `json:"ttl_seconds"`
}

type shortenResponse struct {
	ShortCode string `json:"short_code"`
}

func (h *Handler) postShorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.OriginalURL) == "" {
		http.Error(w, "original_url is required", http.StatusBadRequest)
		return
	}
	var expiresAt *time.Time
	if req.TTLSeconds != nil && *req.TTLSeconds > 0 {
		t := time.Now().Add(time.Duration(*req.TTLSeconds) * time.Second)
		expiresAt = &t
	}
	s, err := h.svc.CreateShort(r.Context(), req.OriginalURL, req.ClientID, expiresAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, shortenResponse{ShortCode: s.ShortCode}, http.StatusCreated)
}

func (h *Handler) getRedirect(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	s, err := h.svc.Resolve(r.Context(), short)
	if err != nil || s == nil {
		http.NotFound(w, r)
		return
	}
	// Аналитика клика (best-effort)
	if cid := extractClientID(r); cid != nil {
		_ = h.svc.RecordClick(r.Context(), models.ClickEvent{
			ShortCode: short,
			ClientID:  *cid,
			UserAgent: r.UserAgent(),
			IP:        clientIP(r),
			At:        time.Now(),
		})
	}
	http.Redirect(w, r, s.Original, http.StatusFound)
}

type analyticsResponse struct {
	Total int64             `json:"total"`
	Daily []models.AggPoint `json:"daily"`
	ByUA  []models.AggPoint `json:"by_user_agent"`
}

func (h *Handler) getAnalytics(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "short")
	total, err := h.svc.Count(r.Context(), short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	daily, err := h.svc.Daily(r.Context(), short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	byUA, err := h.svc.ByUserAgent(r.Context(), short)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, analyticsResponse{Total: total, Daily: daily, ByUA: byUA}, http.StatusOK)
}

// helpers

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func extractClientID(r *http.Request) *uuid.UUID {
	if cid := r.Header.Get("X-Client-Id"); cid != "" {
		if u, err := uuid.Parse(cid); err == nil {
			return &u
		}
	}
	if cid := r.URL.Query().Get("client_id"); cid != "" {
		if u, err := uuid.Parse(cid); err == nil {
			return &u
		}
	}
	return nil
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	host := r.RemoteAddr
	if idx := strings.LastIndex(host, ":"); idx > 0 {
		return host[:idx]
	}
	return host
}
