package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/models"
	tmplPkg "github.com/MAPiryazev/Wildberries_L1/tree/main/L3/L3.5/internal/templates"
	"github.com/go-chi/chi/v5"
)

// UIAdminEventsList — админская страница со списком событий и бронирований
func (h *DefaultHandler) UIAdminEventsList(w http.ResponseWriter, r *http.Request) {
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
		ID          int64
		Title       string
		StartTime   time.Time
		Capacity    int64
		FreePlaces  int64
		Bookings    []*models.Booking
		TotalBooked int64
	}

	var data []EventView
	for _, e := range events {
		free, err := h.svc.CountFreePlaces(r.Context(), e.ID)
		if err != nil {
			http.Error(w, "error counting free places: "+err.Error(), http.StatusInternalServerError)
			return
		}

		bookings, err := h.svc.GetBookingsByEventID(r.Context(), e.ID)
		if err != nil {
			http.Error(w, "error fetching bookings: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data = append(data, EventView{
			ID:          e.ID,
			Title:       e.Title,
			StartTime:   e.StartTime,
			Capacity:    e.Capacity,
			FreePlaces:  free,
			Bookings:    bookings,
			TotalBooked: int64(len(bookings)),
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplPkg.Templates.ExecuteTemplate(w, "admin_events.html", data); err != nil {
		http.Error(w, "error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UIAdminEventPage — админская страница одного события с деталями бронирований
func (h *DefaultHandler) UIAdminEventPage(w http.ResponseWriter, r *http.Request) {
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

	bookings, err := h.svc.GetBookingsByEventID(r.Context(), eventID)
	if err != nil {
		http.Error(w, "error fetching bookings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Event      *models.Event
		FreePlaces int64
		Bookings   []*models.Booking
	}{
		Event:      event,
		FreePlaces: freePlaces,
		Bookings:   bookings,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmplPkg.Templates.ExecuteTemplate(w, "admin_event.html", data); err != nil {
		http.Error(w, "error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UIAdminConfirmBooking — админ может подтвердить любую бронь
func (h *DefaultHandler) UIAdminConfirmBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form: "+err.Error(), http.StatusBadRequest)
		return
	}

	bookingIDStr := r.FormValue("booking_id")
	if bookingIDStr == "" {
		http.Error(w, "missing booking id", http.StatusBadRequest)
		return
	}

	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil || bookingID <= 0 {
		http.Error(w, "invalid booking id", http.StatusBadRequest)
		return
	}

	// Админ может подтверждать любую бронь
	if err := h.svc.UpdateBookingStatus(r.Context(), bookingID, "confirmed", true); err != nil {
		http.Error(w, "error confirming booking: "+err.Error(), http.StatusBadRequest)
		return
	}

	eventIDStr := r.FormValue("event_id")
	if eventIDStr != "" {
		http.Redirect(w, r, "/ui/admin/events/"+eventIDStr, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/ui/admin/events", http.StatusSeeOther)
	}
}
