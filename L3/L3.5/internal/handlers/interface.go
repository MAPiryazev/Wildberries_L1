package handlers

import "net/http"

// Handler описывает все HTTP-эндпоинты сервиса EventBooker.
type Handler interface {
	// POST /events — создать мероприятие
	CreateEvent(w http.ResponseWriter, r *http.Request)

	// GET /events — получить список всех мероприятий (+ свободные места)
	GetAllEvents(w http.ResponseWriter, r *http.Request)

	// GET /events/{id} — получить инфо о мероприятии
	GetEvent(w http.ResponseWriter, r *http.Request)

	// POST /events/{id}/book — забронировать место
	BookEvent(w http.ResponseWriter, r *http.Request)

	// POST /events/{id}/confirm — подтвердить (оплатить) бронь
	ConfirmBooking(w http.ResponseWriter, r *http.Request)

	// = WEB UI

	// GET /ui/events — список событий (страница)
	UIEventsList(w http.ResponseWriter, r *http.Request)

	// GET /ui/events/new — страница создания события
	UINewEventPage(w http.ResponseWriter, r *http.Request)

	// POST /ui/events/new — форма создания события
	UICreateEvent(w http.ResponseWriter, r *http.Request)

	// GET /ui/events/{id} — страница события
	UIEventPage(w http.ResponseWriter, r *http.Request)

	// POST /ui/events/{id}/book — форма бронирования
	UIBookEvent(w http.ResponseWriter, r *http.Request)

	// POST /ui/events/{id}/confirm — форма подтверждения
	UIConfirmBooking(w http.ResponseWriter, r *http.Request)
}
