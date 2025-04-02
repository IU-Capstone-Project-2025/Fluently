package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router) {
    r.Route("/users", func(r chi.Router) {
        r.Post("/", handler.CreateUser)
        r.Get("/{id}", handler.GetUser)
        r.Put("/{id}", handler.UpdateUser)
        r.Delete("/{id}", handler.DeleteUser)

        r.Get("/{id}/preferences", handler.GetUserPreferences)
        r.Put("/{id}/preferences", handler.UpdateUserPreferences)

        r.Get("/{id}/learned-words", handler.GetLearnedWords)
        r.Get("/{id}/learned-words/{word_id}", handler.GetLearnedWord)
        r.Post("/{id}/learned-words", handler.CreateLearnedWord)
        r.Put("/{id}/learned-words/{word_id}", handler.UpdateLearnedWord)
        r.Delete("/{id}/learned-words/{word_id}", handler.DeleteLearnedWord)
    })
}
