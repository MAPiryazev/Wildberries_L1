package handler

import (
	"net/http"

	"server-calendar/internal/service"
)

// NewRouter создаёт и настраивает HTTP-маршрутизатор.
func NewRouter(svc *service.CalendarService) *http.ServeMux {
	mux := http.NewServeMux()
	h := NewEventHandler(svc)

	mux.HandleFunc("/create_event", h.CreateEvent)
	mux.HandleFunc("/update_event", h.UpdateEvent)
	mux.HandleFunc("/delete_event", h.DeleteEvent)
	mux.HandleFunc("/events_for_day", h.EventsForDay)
	mux.HandleFunc("/events_for_week", h.EventsForWeek)
	mux.HandleFunc("/events_for_month", h.EventsForMonth)

	return mux
}
