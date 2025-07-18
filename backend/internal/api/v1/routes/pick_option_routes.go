package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterPickOptionRoutes registers pick option routes
func RegisterPickOptionRoutes(r chi.Router, h *handler.PickOptionHandler) {
	r.Route("/pick-options", func(r chi.Router) {
		r.Post("/", h.CreatePickOption)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetPickOption)
			r.Put("/", h.UpdatePickOption)
			r.Delete("/", h.DeletePickOption)
		})
	})

	r.Get("/words/{word_id}/pick-options", h.ListPickOptions)
}
