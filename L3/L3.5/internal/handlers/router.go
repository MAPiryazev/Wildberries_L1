package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRouter настраивает все маршруты
func SetupRouter(h *DefaultHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes
	r.Route("/events", func(r chi.Router) {
		r.Post("/", h.CreateEvent)                // POST /events
		r.Get("/", h.GetAllEvents)                // GET /events
		r.Get("/{id}", h.GetEvent)                // GET /events/{id}
		r.Post("/{id}/book", h.BookEvent)         // POST /events/{id}/book
		r.Post("/{id}/confirm", h.ConfirmBooking) // POST /events/{id}/confirm
	})

	// UI routes - Admin
	r.Route("/ui/admin", func(r chi.Router) {
		r.Get("/events", h.UIAdminEventsList)                   // GET /ui/admin/events
		r.Get("/events/{id}", h.UIAdminEventPage)               // GET /ui/admin/events/{id}
		r.Post("/events/{id}/confirm", h.UIAdminConfirmBooking) // POST /ui/admin/events/{id}/confirm
	})

	// UI routes - User
	r.Route("/ui/user", func(r chi.Router) {
		r.Get("/events", h.UIUserEventsList)                   // GET /ui/user/events
		r.Get("/events/{id}", h.UIUserEventPage)               // GET /ui/user/events/{id}
		r.Post("/events/{id}/book", h.UIUserBookEvent)         // POST /ui/user/events/{id}/book
		r.Post("/events/{id}/confirm", h.UIUserConfirmBooking) // POST /ui/user/events/{id}/confirm
	})

	// UI routes - Legacy (для обратной совместимости)
	r.Route("/ui/events", func(r chi.Router) {
		r.Get("/", h.UIEventsList)                  // GET /ui/events
		r.Get("/new", h.UINewEventPage)             // GET /ui/events/new
		r.Post("/new", h.UICreateEvent)             // POST /ui/events/new
		r.Get("/{id}", h.UIEventPage)               // GET /ui/events/{id}
		r.Post("/{id}/book", h.UIBookEvent)         // POST /ui/events/{id}/book
		r.Post("/{id}/confirm", h.UIConfirmBooking) // POST /ui/events/{id}/confirm
	})

	// Root redirect to user events
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/user/events", http.StatusSeeOther)
	})

	return r
}
