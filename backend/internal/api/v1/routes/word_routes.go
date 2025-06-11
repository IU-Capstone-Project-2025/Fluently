package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/service"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func RegisterWordRoutes(r chi.Router, db *gorm.DB) {
	wordRepo := postgres.NewWordPostgres(db)
	WordService := service.NewWordService(wordRepo)
	wordHandler := handler.NewWordHandler(WordService)

	r.Route("/words", func(r chi.Router) {
		r.Get("/", wordHandler.ListWords)
		r.Get("/{id}", wordHandler.GetWord)
		r.Post("/", wordHandler.CreateWord)
		r.Put("/{id}", wordHandler.UpdateWord)
		r.Delete("/{id}", wordHandler.DeleteWord)

		r.Get("/{id}/sentences", handler.ListSentences)
		r.Post("/{id}/sentences", handler.CreateSentence)
	})
}
