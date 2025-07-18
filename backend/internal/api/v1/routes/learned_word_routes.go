package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterLearnedWordRoutes registers learned word routes
func RegisterLearnedWordRoutes(r chi.Router, h *handler.LearnedWordHandler) {
	r.Route("/users/{user_id}/learned-words", func(r chi.Router) {
		r.Get("/", h.GetLearnedWords)
		r.Get("/{word_id}", h.GetLearnedWord)
		r.Post("/", h.CreateLearnedWord)
		r.Put("/{word_id}", h.UpdateLearnedWord)
		r.Delete("/{word_id}", h.DeleteLearnedWord)
	})
}
