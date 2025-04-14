package routes

import (
<<<<<<< HEAD
	//	handler "fluently/go-backend/internal/api/v1/handlers"
=======
	handler "fluently/go-backend/internal/api/v1/handlers"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/service"
>>>>>>> 7de7f04 (Add all word logic)

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

<<<<<<< HEAD
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
=======
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
>>>>>>> 7de7f04 (Add all word logic)
}
