package routes

import (
	"fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterThesaurusRoutes(r chi.Router, h *handlers.ThesaurusHandler) {
	r.Post("/thesaurus/recommend", h.Recommend) // POST /api/v1/thesaurus/recommend
}
