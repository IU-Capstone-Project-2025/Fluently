package routes

import (
	//	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterWordRoutes(r chi.Router) {
	r.Route("/words", func(r chi.Router) {
		// r.Get("/", handler.ListWords)
		// r.Get("/{id}", handler.GetWord)
		// r.Post("/", handler.CreateWord)
		// r.Put("/{id}", handler.UpdateWord)
		// r.Delete("/{id}", handler.DeleteWord)
		//
		// r.Get("/{id}/sentences", handler.ListSentences)
		// r.Post("/{id}/sentences", handler.CreateSentence)
	})
}
