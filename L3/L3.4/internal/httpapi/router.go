package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter регистрирует все маршруты
func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Group(func(r chi.Router) {
		r.Get("/healthz", handler.Health)
		r.Get("/images", handler.ListImages)
		r.Post("/upload", handler.Upload)
		r.Get("/image/{id}", handler.GetImage)
		r.Delete("/image/{id}", handler.DeleteImage)
	})

	fileServer := http.FileServer(http.FS(assetsFS))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))
	r.Get("/", serveIndex)

	return r
}
