package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterPreferencesRoutes registers preferences routes
func RegisterPreferencesRoutes(r chi.Router, h *handler.PreferenceHandler) {
	r.Route("/preferences", func(r chi.Router) {
		r.Get("/", h.GetUserPreferences)                   // /preferences (gets user from context)
		r.Put("/", h.UpdateUserPreferences)                // /preferences (information from token)
		r.Delete("/", h.DeletePreference)                  // /preferences (gets user from context)
		r.Post("/user/{user_id}", h.CreateUserPreferences) // /preferences/user/{user_id}
	})
}
