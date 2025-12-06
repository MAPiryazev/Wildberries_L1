package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/models"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/service"
)

type DefaultHandler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *DefaultHandler {
	return &DefaultHandler{svc: svc}
}

func (h *DefaultHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string `json:"title"`
		StartTime string `json:"start_time"`
		Capacity  int64  `json:"capacity"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "wrong method selected for POST operation", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	event := &models.Event{
		Title:    input.Title,
		Capacity: input.Capacity,
	}

	// Парсим время
	t, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		http.Error(w, "invalid start_time format, must be RFC3339", http.StatusBadRequest)
		return
	}
	event.StartTime = t

	// Проверяем админский флаг из query-параметра
	isAdmin := r.URL.Query().Get("admin") == "1"

	created, err := h.svc.CreateEvent(r.Context(), event, isAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *DefaultHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "wrong method selected for GET operation", http.StatusMethodNotAllowed)
		return
	}

	events, err := h.svc.GetAllEvents(r.Context())
	if err != nil {
		http.Error(w, "error fetching events: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем свободные места для каждого события
	type EventWithFree struct {
		*models.Event
		FreePlaces int64 `json:"free_places"`
	}

	var out []EventWithFree
	for _, e := range events {
		free, err := h.svc.CountFreePlaces(r.Context(), e.ID)
		if err != nil {
			http.Error(w, "error counting free places: "+err.Error(), http.StatusInternalServerError)
			return
		}
		out = append(out, EventWithFree{
			Event:      e,
			FreePlaces: free,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *DefaultHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "wrong method selected for GET operation", http.StatusMethodNotAllowed)
		return
	}

	// Читаем ID события из query-параметра (для простоты)
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "missing event id", http.StatusBadRequest)
		return
	}

	eventID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || eventID <= 0 {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	event, err := h.svc.GetEventByID(r.Context(), eventID)
	if err != nil {
		http.Error(w, "error fetching event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	free, err := h.svc.CountFreePlaces(r.Context(), event.ID)
	if err != nil {
		http.Error(w, "error counting free places: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type EventWithFree struct {
		*models.Event
		FreePlaces int64 `json:"free_places"`
	}

	out := EventWithFree{
		Event:      event,
		FreePlaces: free,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *DefaultHandler) BookEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method selected for POST operation", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID события из query-параметра
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "missing event id", http.StatusBadRequest)
		return
	}
	eventID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || eventID <= 0 {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из query-параметра
	userParam := r.URL.Query().Get("user_id")
	if userParam == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userParam, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// Для упрощения считаем, что бронь истекает через 1 час
	expiresAt := time.Now().Add(time.Hour)

	booking := &models.Booking{
		EventID:   eventID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	created, err := h.svc.CreateBooking(r.Context(), booking)
	if err != nil {
		http.Error(w, "error creating booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *DefaultHandler) ConfirmBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method selected for POST operation", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID события из query-параметра
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "missing event id", http.StatusBadRequest)
		return
	}
	eventID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || eventID <= 0 {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя
	userParam := r.URL.Query().Get("user_id")
	if userParam == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userParam, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	// Проверяем админский флаг (для Confirm можно ставить через admin=true, чтобы админ мог подтверждать чужие брони)
	isAdmin := r.URL.Query().Get("admin") == "1"

	booking, err := h.svc.GetBooking(r.Context(), eventID, userID)
	if err != nil {
		http.Error(w, "error fetching booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	if booking.Status == "confirmed" {
		http.Error(w, "booking already confirmed", http.StatusBadRequest)
		return
	}

	err = h.svc.UpdateBookingStatus(r.Context(), booking.ID, "confirmed", isAdmin)
	if err != nil {
		http.Error(w, "error confirming booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "confirmed",
	})
}
