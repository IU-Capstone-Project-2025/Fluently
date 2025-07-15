package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterPreferencesRoutes(r chi.Router, h *handler.PreferenceHandler) {
	r.Route("/preferences/{user_id}", func(r chi.Router) {
		r.Post("/", h.CreateUserPreferences)
		r.Put("/", h.UpdateUserPreferences)
		r.Delete("/", h.DeletePreference)
	})

	r.Get("/", h.GetUserPreferences)
}
