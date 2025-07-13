package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"
	
	"github.com/go-chi/chi/v5"
)

func RegisterDayWordRoutes(r chi.Router, h *handler.DayWordHandler) {
	r.Get("/day-word", h.GetDayWord)
}