package routes

import (
	"fluently/go-backend/internal/api/v1/handlers"
	"fluently/go-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, h *handlers.Handlers) {
	r.Route("/auth", func(r chi.Router) {
		// Add debug middleware to auth routes
		r.Use(middleware.DebugAuthRoutes)

		r.Post("/login", h.LoginHandler)
		r.Post("/register", h.RegisterHandler)
		r.Post("/google", h.GoogleAuthHandler)
		r.Get("/google", h.GoogleAuthRedirectHandler)
		r.Get("/google/callback", h.GoogleCallbackHandler)
		r.Post("/refresh", h.RefreshTokenHandler)
		// r.Post("/logout", h.LogoutHandler)
		// r.Post("/forgot-password", h.ForgotPasswordHandler)
		// r.Post("/reset-password", h.ResetPasswordHandler)
	})
}
