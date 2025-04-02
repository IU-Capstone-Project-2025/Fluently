package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterPreferencesRoutes(r chi.Router) {
    r.Route("/users/{id}/preferences", func(r chi.Router) {
        r.Get("/", handler.GetUserPreferences)
        r.Put("/", handler.UpdateUserPreferences)
    })
}
