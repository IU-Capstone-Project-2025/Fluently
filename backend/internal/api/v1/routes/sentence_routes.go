package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterSentenceRoutes registers sentence routes
func RegisterSentenceRoutes(r chi.Router, h *handler.SentenceHandler) {
	r.Get("/words/{word_id}/sentences", h.ListSentences)
	r.Post("/sentences", h.CreateSentence)
	r.Put("/sentences/{id}", h.UpdateSentence)
	r.Delete("/sentences/{id}", h.DeleteSentence)
}
