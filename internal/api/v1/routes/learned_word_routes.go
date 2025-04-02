package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterLearnedWordRoutes(r chi.Router) {
    r.Route("/learned-words", func(r chi.Router) {
        r.Get("/", handler.GetLearnedWords)                       // ?user_id=
        r.Get("/{word_id}", handler.GetLearnedWord)              // ?user_id=
        r.Post("/", handler.CreateLearnedWord)                   // в теле user_id
        r.Put("/{word_id}", handler.UpdateLearnedWord)           // в теле user_id
        r.Delete("/{word_id}", handler.DeleteLearnedWord)        // ?user_id=
    })
}
