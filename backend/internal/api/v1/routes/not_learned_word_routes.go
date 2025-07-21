package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterNotLearnedWordRoutes registers not learned word routes
func RegisterNotLearnedWordRoutes(r chi.Router, h *handler.NotLearnedWordHandler) {
	r.Route("/not-learned-words", func(r chi.Router) {
		r.Post("/", h.AddNotLearnedWord) // POST /api/v1/not-learned-words
		r.Get("/", h.GetNotLearnedWords) // GET /api/v1/not-learned-words
	})
}
