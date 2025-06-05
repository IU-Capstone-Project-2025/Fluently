package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterLearnedWordRoutes(r chi.Router, h *handler.LearnedWordHandler) {
	r.Route("users/{user_id}/learned-words", func(r chi.Router) {
		r.Get("/", h.GetLearnedWord)                // ?user_id=
		r.Get("/{word_id}", h.GetLearnedWord)       // ?user_id=
		r.Post("/", h.CreateLearnedWord)            // в теле user_id
		r.Put("/{word_id}", h.UpdateLearnedWord)    // в теле user_id
		r.Delete("/{word_id}", h.DeleteLearnedWord) // ?user_id=
	})
}
