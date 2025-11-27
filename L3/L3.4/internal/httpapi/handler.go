package httpapi

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.4/internal/service/image"
)

// Handler обрабатывает HTTP запросы
type Handler struct {
	service      image.Service
	maxListItems int
}

// HandlerConfig настройки HTTP обработчиков
type HandlerConfig struct {
	MaxListItems int
}

// NewHandler конструктор Handler
func NewHandler(service image.Service, cfg HandlerConfig) *Handler {
	maxItems := cfg.MaxListItems
	if maxItems <= 0 {
		maxItems = 50
	}
	return &Handler{
		service:      service,
		maxListItems: maxItems,
	}
}

// Upload обрабатывает POST /upload
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "failed to parse multipart form: %v", err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "file field is required: %v", err)
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(make([]byte, 512))
	}

	pt := models.ProcessingType(r.FormValue("processing_type"))

	info, err := h.service.Upload(r.Context(), file, fileHeader.Size, contentType, fileHeader.Filename, pt)
	if err != nil {
		respondError(w, http.StatusBadRequest, "upload failed: %v", err)
		return
	}

	respondJSON(w, http.StatusCreated, info)
}

// GetImage возвращает информацию либо редиректит на файл
func (h *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")
	if imageID == "" {
		respondError(w, http.StatusBadRequest, "image id is required")
		return
	}

	info, err := h.service.Get(r.Context(), imageID)
	if err != nil {
		respondError(w, http.StatusNotFound, "image not found: %v", err)
		return
	}

	variant := r.URL.Query().Get("variant")
	switch variant {
	case "":
		respondJSON(w, http.StatusOK, info)
	case "original":
		if info.OriginalURL == "" {
			respondError(w, http.StatusNotFound, "original image is not available yet")
			return
		}
		http.Redirect(w, r, info.OriginalURL, http.StatusTemporaryRedirect)
	case "processed":
		if info.ProcessedURL == "" {
			respondError(w, http.StatusNotFound, "processed image is not available")
			return
		}
		http.Redirect(w, r, info.ProcessedURL, http.StatusTemporaryRedirect)
	default:
		respondError(w, http.StatusBadRequest, "unknown variant: %s", variant)
	}
}

// DeleteImage удаляет изображение
func (h *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	imageID := chi.URLParam(r, "id")
	if imageID == "" {
		respondError(w, http.StatusBadRequest, "image id is required")
		return
	}

	if err := h.service.Delete(r.Context(), imageID); err != nil {
		respondError(w, http.StatusInternalServerError, "delete failed: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListImages возвращает список изображений
func (h *Handler) ListImages(w http.ResponseWriter, r *http.Request) {
	limit := h.maxListItems
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 && v <= 500 {
			limit = v
		}
	}

	items, err := h.service.List(r.Context(), limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "list images failed: %v", err)
		return
	}

	respondJSON(w, http.StatusOK, items)
}

// Health хелсчек
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
