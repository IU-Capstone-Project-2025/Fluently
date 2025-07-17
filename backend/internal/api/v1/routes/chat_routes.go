package routes

import (
	"fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterChatRoutes(r chi.Router, h *handlers.ChatHandler) {
	r.Route("/chat", func(r chi.Router) {
		r.Post("/", h.Chat)             // POST /api/v1/chat
		r.Post("/finish", h.FinishChat) // POST /api/v1/chat/finish
	})
}
