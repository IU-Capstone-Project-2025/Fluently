package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

// RegisterTelegramRoutes registers telegram routes
func RegisterTelegramRoutes(r chi.Router, h *handler.TelegramHandler) {
	// Public routes (no authentication required)
	r.Route("/telegram", func(r chi.Router) {
		r.Post("/create-link", h.CreateLinkToken)
		r.Post("/check-status", h.CheckLinkStatus)
	})

	// Protected routes
	r.Route("/api/v1/telegram", func(r chi.Router) {
		r.Post("/unlink", h.UnlinkTelegram)
	})

	// Routes for magic links
	r.Get("/link-google", h.LinkWithGoogle)
	r.Get("/link-google/callback", h.LinkGoogleCallback)
}
