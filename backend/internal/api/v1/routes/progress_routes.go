package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterProgressRoutes registers progress routes
func RegisterProgressRoutes(r chi.Router, h *handler.ProgressHandler) {
	r.Post("/progress", h.UpdateUserProgress) // update user progress (using token from context)
}
