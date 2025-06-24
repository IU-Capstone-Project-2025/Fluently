package routes

import (
	handler "fluently/go-backend/internal/api/v1/handlers"

	"github.com/go-chi/chi/v5"
)

func RegisterTelegramRoutes(r chi.Router, h *handler.TelegramHandler) {
	// Публичные роуты (не требуют аутентификации)
	r.Route("/telegram", func(r chi.Router) {
		r.Post("/create-link", h.CreateLinkToken)
		r.Post("/check-status", h.CheckLinkStatus)
	})

	// Защищенные роуты
	r.Route("/api/v1/telegram", func(r chi.Router) {
		r.Post("/unlink", h.UnlinkTelegram)
	})

	// Роуты для магических ссылок (публичные)
	r.Get("/link-google", h.LinkWithGoogle)
	r.Get("/link-google/callback", h.LinkGoogleCallback)
}
