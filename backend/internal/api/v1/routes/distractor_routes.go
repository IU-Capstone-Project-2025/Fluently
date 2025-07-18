package routes

import (
	"fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterDistractorRoutes(r chi.Router, h *handlers.DistractorHandler) {
	r.Post("/distractors", h.Generate) // POST /api/v1/distractors
}
