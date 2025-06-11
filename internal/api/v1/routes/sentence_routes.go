package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterSentenceRoutes(r chi.Router, h *handler.SentenceHandler) {
	r.Put("/sentences/{id}", h.UpdateSentence)
	r.Delete("/sentences/{id}", h.DeleteSentence)
}
