package handlers

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	// auth
	r.Get("/auth/me", h.GetCurrentUser)

	// users
	r.Get("/users", h.ListUsers)

	// items
	r.Get("/items", h.ListItems)
	r.Post("/items", h.CreateItem)
	r.Get("/items/{id}", h.GetItem)
	r.Put("/items/{id}", h.UpdateItem)
	r.Delete("/items/{id}", h.DeleteItem)
	r.Get("/items/{id}/history", h.GetItemHistory)
}
