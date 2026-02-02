package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/models"
	tmplPkg "github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/templates"
	"github.com/go-chi/chi/v5"
)

// DefaultHandler уже должен быть объявлен в другом файле:
// type DefaultHandler struct { svc service.Service }
// func NewHandler(svc service.Service) *DefaultHandler { return &DefaultHandler{svc: svc} }

// helper: tries to obtain event id from query/form/path
func getEventIDFromRequest(r *http.Request) (int64, error) {
	// 1) query param ?id=
	if idq := r.URL.Query().Get("id"); idq != "" {
		return strconv.ParseInt(idq, 10, 64)
	}
	// 2) form value "id"
	if err := r.ParseForm(); err == nil {
		if idf := r.FormValue("id"); idf != "" {
			return strconv.ParseInt(idf, 10, 64)
		}
	}

	// 3) try to parse from path like /ui/events/{id}/book or /ui/events/{id}/confirm or /ui/events/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// expected: ["ui","events","{id}","book"]
	if len(parts) >= 3 {
		// find segment that looks like a number among the parts
		for _, p := range parts {
			if n, err := strconv.ParseInt(p, 10, 64); err == nil && n > 0 {
				return n, nil
			}
		}
	}

	return 0, http.ErrNoLocation
}

// UIEventsList отображает список всех событий с количеством свободных мест
func (h *DefaultHandler) UIEventsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "wrong method selected for GET operation", http.StatusMethodNotAllowed)
		return
	}

	if tmplPkg.Templates == nil {
		http.Error(w, "templates not loaded", http.StatusInternalServerError)
		return
	}

	events, err := h.svc.GetAllEvents(r.Context())
	if err != nil {
		http.Error(w, "error fetching events: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Для каждого события добавляем свободные места
	type EventView struct {
		ID         int64
		Title      string
		StartTime  time.Time
		Capacity   int64
		FreePlaces int64
	}

	var data []EventView
	for _, e := range events {
		free, err := h.svc.CountFreePlaces(r.Context(), e.ID)
		if err != nil {
			http.Error(w, "error counting free places: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, EventView{
			ID:         e.ID,
			Title:      e.Title,
			StartTime:  e.StartTime,
			Capacity:   e.Capacity,
			FreePlaces: free,
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplPkg.Templates.ExecuteTemplate(w, "events.html", data); err != nil {
		http.Error(w, "error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UINewEventPage отображает форму создания нового события
func (h *DefaultHandler) UINewEventPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "wrong method selected for GET operation", http.StatusMethodNotAllowed)
		return
	}

	if tmplPkg.Templates == nil {
		http.Error(w, "templates not loaded", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplPkg.Templates.ExecuteTemplate(w, "new_event.html", nil); err != nil {
		http.Error(w, "error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UICreateEvent обрабатывает форму создания события
func (h *DefaultHandler) UICreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong method selected for POST operation", http.StatusMethodNotAllowed)
		return
	}

	// parse form values
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	startTimeStr := r.FormValue("start_time")
	capacityStr := r.FormValue("capacity")
	adminFlag := r.FormValue("admin") == "1"

	if title == "" || startTimeStr == "" || capacityStr == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	capacity, err := strconv.ParseInt(capacityStr, 10, 64)
	if err != nil || capacity <= 0 {
		http.Error(w, "invalid capacity", http.StatusBadRequest)
		return
	}

	t, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		http.Error(w, "invalid start_time format, must be RFC3339", http.StatusBadRequest)
		return
	}

	event := &models.Event{
		Title:     title,
		StartTime: t,
		Capacity:  capacity,
	}

	_, err = h.svc.CreateEvent(r.Context(), event, adminFlag)
	if err != nil {
		// CreateEvent возвращает ошибку, если не админ или данные неверны
		http.Error(w, "error creating event: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Перенаправляем обратно на список
	http.Redirect(w, r, "/ui/events", http.StatusSeeOther)
}

// UIEventPage — показывает страницу одного события (рендерит форму бронирования/подтверждения)
func (h *DefaultHandler) UIEventPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	if tmplPkg.Templates == nil {
		http.Error(w, "templates not loaded", http.StatusInternalServerError)
		return
	}

	// Пробуем получить из path параметра, если нет - из query
	eventIDStr := chi.URLParam(r, "id")
	var eventID int64
	var err error
	if eventIDStr != "" {
		eventID, err = strconv.ParseInt(eventIDStr, 10, 64)
	} else {
		eventID, err = getEventIDFromRequest(r)
	}
	if err != nil || eventID <= 0 {
		http.Error(w, "missing or invalid event id", http.StatusBadRequest)
		return
	}

	event, err := h.svc.GetEventByID(r.Context(), eventID)
	if err != nil {
		http.Error(w, "event not found: "+err.Error(), http.StatusNotFound)
		return
	}

	freePlaces, err := h.svc.CountFreePlaces(r.Context(), eventID)
	if err != nil {
		http.Error(w, "error counting free places: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Event      *models.Event
		FreePlaces int64
	}{
		Event:      event,
		FreePlaces: freePlaces,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Для отображения одного события используем шаблон booking.html (в нём есть форма для user_id)
	if err := tmplPkg.Templates.ExecuteTemplate(w, "booking.html", data); err != nil {
		http.Error(w, "error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UIBookEvent — обрабатывает submission формы бронирования (форма отправляет user_id)
func (h *DefaultHandler) UIBookEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// Пробуем получить из path параметра, если нет - из query/form
	eventIDStr := chi.URLParam(r, "id")
	var eventID int64
	var err error
	if eventIDStr != "" {
		eventID, err = strconv.ParseInt(eventIDStr, 10, 64)
	} else {
		eventID, _ = getEventIDFromRequest(r)
	}

	// гарантируем чтение формы
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
		return
	}
	// если не получили id ранее, попробуем из формы
	if eventID == 0 {
		if idf := r.FormValue("id"); idf != "" {
			eventID, _ = strconv.ParseInt(idf, 10, 64)
		}
	}

	if eventID <= 0 {
		http.Error(w, "missing event id", http.StatusBadRequest)
		return
	}

	// user id из формы (booking.html использует name="user_id")
	userIDStr := r.FormValue("user_id")
	if userIDStr == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	booking := &models.Booking{
		EventID:   eventID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(15 * time.Minute), // срок подтверждения — 15 минут
	}

	_, err = h.svc.CreateBooking(r.Context(), booking)
	if err != nil {
		http.Error(w, "error creating booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	// после успешного бронирования — редиректим на страницу события (чтобы можно было подтвердить)
	http.Redirect(w, r, "/ui/events/"+strconv.FormatInt(eventID, 10), http.StatusSeeOther)
}

// UIConfirmBooking — форма подтверждения брони; принимает event id + user_id (или booking id через booking param)
func (h *DefaultHandler) UIConfirmBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse form or query
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Prefer booking id if provided
	bookingIDStr := r.FormValue("booking")
	if bookingIDStr == "" {
		bookingIDStr = r.URL.Query().Get("booking")
	}

	if bookingIDStr != "" {
		bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
		if err != nil || bookingID <= 0 {
			http.Error(w, "invalid booking id", http.StatusBadRequest)
			return
		}
		// пользователь подтверждает свою бронь
		if err := h.svc.UpdateBookingStatus(r.Context(), bookingID, "confirmed", false); err != nil {
			http.Error(w, "error confirming booking: "+err.Error(), http.StatusBadRequest)
			return
		}
		// redirect back to events list
		http.Redirect(w, r, "/ui/events", http.StatusSeeOther)
		return
	}

	// If no booking id — require event id + user_id
	eventID, _ := getEventIDFromRequest(r)
	if eventID == 0 {
		if idf := r.FormValue("id"); idf != "" {
			eventID, _ = strconv.ParseInt(idf, 10, 64)
		}
	}
	if eventID <= 0 {
		http.Error(w, "missing event id", http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("user_id")
	if userIDStr == "" {
		userIDStr = r.URL.Query().Get("user_id")
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "missing or invalid user id", http.StatusBadRequest)
		return
	}

	booking, err := h.svc.GetBooking(r.Context(), eventID, userID)
	if err != nil {
		http.Error(w, "error fetching booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	if booking == nil {
		http.Error(w, "booking not found", http.StatusBadRequest)
		return
	}

	if booking.Status == "confirmed" {
		http.Error(w, "booking already confirmed", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateBookingStatus(r.Context(), booking.ID, "confirmed", false); err != nil {
		http.Error(w, "error confirming booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/ui/events/"+strconv.FormatInt(eventID, 10), http.StatusSeeOther)
}
