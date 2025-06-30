package routes

import (
	"fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterLessonRoutes(r chi.Router, h *handlers.LessonHandler) {
	r.Route("/lesson", func(r chi.Router) {
		r.Get("/", h.GenerateLesson)
	})
}
