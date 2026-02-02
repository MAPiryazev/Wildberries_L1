package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/models"
	tmplPkg "github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/templates"
	"github.com/go-chi/chi/v5"
)

// UIUserEventsList — пользовательская страница со списком событий
func (h *DefaultHandler) UIUserEventsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method allowed", http.StatusMethodNotAllowed)
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
	if err := tmplPkg.Templates.ExecuteTemplate(w, "user_events.html", data); err != nil {
		http.Error(w, "error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UIUserEventPage — пользовательская страница одного события
func (h *DefaultHandler) UIUserEventPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	if tmplPkg.Templates == nil {
		http.Error(w, "templates not loaded", http.StatusInternalServerError)
		return
	}

	eventIDStr := chi.URLParam(r, "id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
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
	if err := tmplPkg.Templates.ExecuteTemplate(w, "user_event.html", data); err != nil {
		http.Error(w, "error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UIUserBookEvent — обработка бронирования для пользователя
func (h *DefaultHandler) UIUserBookEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	eventIDStr := chi.URLParam(r, "id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil || eventID <= 0 {
		http.Error(w, "missing or invalid event id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
		return
	}

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

	http.Redirect(w, r, "/ui/user/events/"+strconv.FormatInt(eventID, 10), http.StatusSeeOther)
}

// UIUserConfirmBooking — пользователь подтверждает свою бронь
func (h *DefaultHandler) UIUserConfirmBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	eventIDStr := chi.URLParam(r, "id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil || eventID <= 0 {
		http.Error(w, "missing or invalid event id", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("user_id")
	if userIDStr == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
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

	// Пользователь может подтверждать только свою бронь (isAdmin = false)
	if err := h.svc.UpdateBookingStatus(r.Context(), booking.ID, "confirmed", false); err != nil {
		http.Error(w, "error confirming booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/ui/user/events/"+strconv.FormatInt(eventID, 10), http.StatusSeeOther)
}
