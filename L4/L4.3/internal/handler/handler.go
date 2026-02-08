package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"server-calendar/internal/service"
	"server-calendar/internal/storage/entity"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type EventHandler struct {
	svc       *service.CalendarService
	validator *validator.Validate
}

// NewEventHandler создаёт обработчик с бизнес-логикой
func NewEventHandler(svc *service.CalendarService) *EventHandler {
	return &EventHandler{
		svc:       svc,
		validator: validator.New(),
	}
}

// CreateEvent — создание нового события
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var e entity.Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		writeError(w, "Json invalid", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(&e)
	if err != nil {
		writeError(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	err = h.svc.CreateEvent(context.Background(), e)
	if err != nil {
		writeError(w, "Create error", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]string{"result": "created"})
}

// UpdateEvent — обновление события
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var e entity.Event

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		writeError(w, "Json invalid", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(&e)
	if err != nil {
		writeError(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	err = h.svc.UpdateEvent(context.Background(), e)
	if err != nil {
		writeError(w, "Update error", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]string{"result": "updated"})
}

// DeleteEvent — удаление события по id
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var e entity.Event

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		writeError(w, "Json invalid", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(&e)
	if err != nil {
		writeError(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	err = h.svc.DeleteEvent(context.Background(), e)
	if err != nil {
		writeError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, map[string]string{"result": "deleted"})
}

// EventsForDay — события за день
func (h *EventHandler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	events, err := h.svc.EventsForDay(context.Background(), entity.UserID(userID))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, events)
}

// EventsForWeek — события за неделю
func (h *EventHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	events, err := h.svc.EventsForWeek(context.Background(), entity.UserID(userID))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, events)
}

// EventsForMonth — события за месяц
func (h *EventHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	events, err := h.svc.EventsForMonth(context.Background(), entity.UserID(userID))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, events)
}

// writeJSON — утилита для возврата JSON-ответа
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
