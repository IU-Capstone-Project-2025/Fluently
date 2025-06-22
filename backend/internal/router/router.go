package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"

	"fluently/go-backend/internal/api/v1/handlers"
	"fluently/go-backend/internal/api/v1/routes"
	authMiddleware "fluently/go-backend/internal/middleware"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"
)

func InitRoutes(db *gorm.DB, r *chi.Mux) {
	// Initialize JWT auth
	utils.InitJWTAuth()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // или конкретные
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authHandlers := &handlers.Handlers{
		UserRepo:         postgres.NewUserRepository(db),
		UserPrefRepo:     postgres.NewPreferenceRepository(db),
		RefreshTokenRepo: postgres.NewRefreshTokenRepository(db),
	}

	// Public routes (NO AUTHENTICATION REQUIRED)
	routes.RegisterAuthRoutes(r, authHandlers)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Custom authenticated Swagger UI (recommended)
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-auth", http.StatusMovedPermanently)
	})

	// Serve custom Swagger UI with authentication
	r.Get("/swagger-auth", func(w http.ResponseWriter, r *http.Request) {
		// Check if file exists
		filePath := filepath.Join("static", "swagger-auth.html")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Fallback to regular swagger if custom file doesn't exist
			http.Redirect(w, r, "/swagger/index.html", http.StatusTemporaryRedirect)
			return
		}
		http.ServeFile(w, r, filePath)
	})

	// Original Swagger documentation (fallback)
	r.Get("/swagger/index.html", httpSwagger.WrapHandler)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Protected routes using go-chi/jwtauth
	r.Route("/api/v1", func(r chi.Router) {
		// JWT authentication middleware
		r.Use(jwtauth.Verifier(utils.TokenAuth))
		r.Use(authMiddleware.CustomAuthenticator)

		// Protected API routes
		routes.RegisterUserRoutes(r, &handlers.UserHandler{Repo: postgres.NewUserRepository(db)})
		routes.RegisterWordRoutes(r, &handlers.WordHandler{Repo: postgres.NewWordRepository(db)})
		routes.RegisterSentenceRoutes(r, &handlers.SentenceHandler{Repo: postgres.NewSentenceRepository(db)})
		routes.RegisterLearnedWordRoutes(r, &handlers.LearnedWordHandler{Repo: postgres.NewLearnedWordRepository(db)})
		routes.RegisterPreferencesRoutes(r, &handlers.PreferenceHandler{Repo: postgres.NewPreferenceRepository(db)})
	})
}
