package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterWordRoutes registers word routes
func RegisterWordRoutes(r chi.Router, h *handler.WordHandler) {
	r.Route("/words", func(r chi.Router) {
		r.Get("/", h.ListWords)
		r.Get("/{id}", h.GetWord)
		r.Post("/", h.CreateWord)
		r.Put("/{id}", h.UpdateWord)
		r.Delete("/{id}", h.DeleteWord)
	})
}
