package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"

	"fluently/go-backend/internal/api/v1/handlers"
	"fluently/go-backend/internal/api/v1/routes"
	authMiddleware "fluently/go-backend/internal/middleware"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"
	"fluently/go-backend/pkg/logger"

	"go.uber.org/zap"
)

// FlexibleJWTVerifier is a custom JWT verifier that supports both "Bearer token" and just "token" formats
func flexibleJWTVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Log.Error("No Authorization header found")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
			return
		}

		// Extract token (support both "Bearer token" and just "token")
		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = authHeader
		}

		// Verify and parse the token
		token, err := utils.TokenAuth.Decode(tokenString)
		if err != nil {
			logger.Log.Error("JWT decode error", zap.Error(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
			return
		}

		if token == nil {
			logger.Log.Error("Invalid JWT token")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
			return
		}

		// Add token and claims to context
		ctx := jwtauth.NewContext(r.Context(), token, nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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

	// Initialize Telegram handler
	linkTokenRepo := postgres.NewLinkTokenRepository(db)
	telegramHandler := &handlers.TelegramHandler{
		UserRepo:      postgres.NewUserRepository(db),
		LinkTokenRepo: linkTokenRepo,
	}

	// Start cleanup task for expired tokens (every hour)
	utils.StartTokenCleanupTask(linkTokenRepo, time.Hour)

	// Public routes (NO AUTHENTICATION REQUIRED)
	routes.RegisterAuthRoutes(r, authHandlers)
	routes.RegisterTelegramRoutes(r, telegramHandler)

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Swagger documentation (public)
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Protected routes using flexible JWT authentication
	r.Route("/api/v1", func(r chi.Router) {
		// JWT authentication middleware (supports both "Bearer token" and "token" formats)
		r.Use(flexibleJWTVerifier)
		r.Use(authMiddleware.CustomAuthenticator)

		// Protected API routes
		routes.RegisterUserRoutes(r, &handlers.UserHandler{Repo: postgres.NewUserRepository(db)})
		routes.RegisterWordRoutes(r, &handlers.WordHandler{Repo: postgres.NewWordRepository(db)})
		routes.RegisterSentenceRoutes(r, &handlers.SentenceHandler{Repo: postgres.NewSentenceRepository(db)})
		routes.RegisterLearnedWordRoutes(r, &handlers.LearnedWordHandler{Repo: postgres.NewLearnedWordRepository(db)})
		routes.RegisterPreferencesRoutes(r, &handlers.PreferenceHandler{Repo: postgres.NewPreferenceRepository(db)})
		routes.RegisterPickOptionRoutes(r, &handlers.PickOptionHandler{Repo: postgres.NewPickOptionRepository(db)})
		routes.RegisterLessonRoutes(r, &handlers.LessonHandler{
			PreferenceRepo: postgres.NewPreferenceRepository(db),
			TopicRepo:      postgres.NewTopicRepository(db),
			SentenceRepo:   postgres.NewSentenceRepository(db),
			PickOptionRepo: postgres.NewPickOptionRepository(db),
			WordRepo:       postgres.NewWordRepository(db),
			Repo:           postgres.NewLessonRepository(db),
		})
	})
}
