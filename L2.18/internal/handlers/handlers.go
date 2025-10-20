package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"calendar/internal/models"
	"calendar/internal/service"
	"net/http"
)

// Handler описывает интерфейс всех HTTP-обработчиков для событий календаря.
type Handler interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
	EventsForDay(w http.ResponseWriter, r *http.Request)
	EventsForWeek(w http.ResponseWriter, r *http.Request)
	EventsForMonth(w http.ResponseWriter, r *http.Request)
}

// DefaultHandler реализует интерфейс Handler.
type DefaultHandler struct {
	svc service.Service
}

// NewDefaultHandler создаёт новый DefaultHandler с указанным сервисом.
func NewDefaultHandler(service service.Service) (*DefaultHandler, error) {
	if service == nil {
		return nil, fmt.Errorf("сервис, передаваемый в конструктор хендлера, не должен быть nil")
	}
	return &DefaultHandler{
		svc: service,
	}, nil
}

// CreateEvent обрабатывает POST /create_event и создаёт новое событие.
func (dh *DefaultHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		fmt.Println("фейлим1")
		return
	}

	err = validateEventInputHelper(&event)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		fmt.Println("фейлим2")
		return
	}

	created, err := dh.svc.CreateEvent(r.Context(), &event)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		fmt.Println("фейлим3")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"result": created})
}

// UpdateEvent обрабатывает POST /update_event и обновляет существующее событие.
func (dh *DefaultHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	err = validateEventInputHelper(&event)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	err = dh.svc.UpdateEvent(r.Context(), &event)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
}

// DeleteEvent обрабатывает POST /delete_event и удаляет событие по ID.
func (dh *DefaultHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	if err := dh.svc.DeleteEvent(r.Context(), id); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
}

// EventsForDay обрабатывает GET /events_for_day и возвращает события за конкретный день.
func (dh *DefaultHandler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	dateStr := r.URL.Query().Get("date")

	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, `{"error":"invalid date format"}`, http.StatusBadRequest)
		return
	}

	events, err := dh.svc.EventsForDay(r.Context(), userID, dateStr)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"result": events})
}

// EventsForWeek обрабатывает GET /events_for_week и возвращает события за неделю.
func (dh *DefaultHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	dateStr := r.URL.Query().Get("date")

	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, `{"error":"invalid date format"}`, http.StatusBadRequest)
		return
	}

	events, err := dh.svc.EventsForWeek(r.Context(), userID, dateStr)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"result": events})
}

// EventsForMonth обрабатывает GET /events_for_month и возвращает события за месяц.
func (dh *DefaultHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	dateStr := r.URL.Query().Get("date")

	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		http.Error(w, `{"error":"invalid date format"}`, http.StatusBadRequest)
		return
	}

	events, err := dh.svc.EventsForMonth(r.Context(), userID, dateStr)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"result": events})
}

// validateEventInputHelper проверяет базовую валидность полей события.
func validateEventInputHelper(event *models.Event) error {
	if event.UserID <= 0 {
		return fmt.Errorf("user_id must be positive")
	}
	if event.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if event.Date.IsZero() {
		return fmt.Errorf("date must be set")
	}
	return nil
}
