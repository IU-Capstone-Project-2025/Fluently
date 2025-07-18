package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterDayWordRoutes registers day word routes
func RegisterDayWordRoutes(r chi.Router, h *handler.DayWordHandler) {
	r.Get("/day-word", h.GetDayWord) // get day word (using token from context)
}
