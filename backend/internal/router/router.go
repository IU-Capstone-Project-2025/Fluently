package router

import (
	"fmt"
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

		// Trim any whitespace that might cause issues
		tokenString = strings.TrimSpace(tokenString)

		// WORKAROUND: Handle malformed client tokens like "ServerToken(accessToken=eyJ...)"
		if strings.HasPrefix(tokenString, "ServerToken(accessToken=") {
			logger.Log.Warn("Client sending malformed token format - applying workaround",
				zap.String("malformed_prefix", tokenString[:min(30, len(tokenString))]),
			)

			// Extract the actual JWT token from ServerToken(accessToken=TOKEN)
			start := strings.Index(tokenString, "accessToken=")
			if start != -1 {
				start += len("accessToken=")
				end := strings.LastIndex(tokenString, ")")
				if end != -1 && end > start {
					extractedToken := tokenString[start:end]
					logger.Log.Info("Extracted JWT token from malformed format",
						zap.String("extracted_prefix", extractedToken[:min(20, len(extractedToken))]),
						zap.Int("extracted_length", len(extractedToken)),
					)
					tokenString = extractedToken
				} else {
					logger.Log.Error("Failed to extract token from malformed format - missing closing parenthesis")
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error": "malformed token format"}`))
					return
				}
			} else {
				logger.Log.Error("Failed to extract token from malformed format - missing accessToken=")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "malformed token format"}`))
				return
			}
		}

		// Log token details for debugging (first 20 chars only for security)
		tokenPrefix := tokenString
		if len(tokenString) > 20 {
			tokenPrefix = tokenString[:20] + "..."
		}
		logger.Log.Debug("JWT token received",
			zap.String("token_prefix", tokenPrefix),
			zap.Int("token_length", len(tokenString)),
			zap.String("authorization_header", authHeader[:min(50, len(authHeader))]),
		)

		// Check if token looks like a valid JWT (should have 3 parts separated by dots)
		parts := strings.Split(tokenString, ".")
		if len(parts) != 3 {
			logger.Log.Error("Invalid JWT format - should have 3 parts separated by dots",
				zap.Int("parts_count", len(parts)),
				zap.String("token_prefix", tokenPrefix),
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid token format"}`))
			return
		}

		// Verify each part contains valid base64
		for i, part := range parts {
			if len(part) == 0 {
				logger.Log.Error("JWT part is empty",
					zap.Int("part_index", i),
					zap.String("token_prefix", tokenPrefix),
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error": "invalid token format"}`))
				return
			}

			// Check for invalid characters in base64
			for j, char := range part {
				if !isValidBase64Char(char) {
					logger.Log.Error("Invalid character found in JWT part",
						zap.Int("part_index", i),
						zap.Int("char_position", j),
						zap.String("invalid_char", string(char)),
						zap.Int("char_code", int(char)),
						zap.String("token_prefix", tokenPrefix),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error": "invalid token encoding"}`))
					return
				}
			}
		}

		// Verify and parse the token
		token, err := utils.TokenAuth.Decode(tokenString)
		if err != nil {
			logger.Log.Error("JWT decode error",
				zap.Error(err),
				zap.String("token_prefix", tokenPrefix),
				zap.Int("token_length", len(tokenString)),
				zap.Strings("token_parts_lengths", getPartsLengths(parts)),
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
			return
		}

		// Check if token is nil
		if token == nil {
			logger.Log.Error("Invalid JWT token - token is nil",
				zap.String("token_prefix", tokenPrefix),
			)
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

// Helper function to check if a character is valid in base64
func isValidBase64Char(char rune) bool {
	return (char >= 'A' && char <= 'Z') ||
		(char >= 'a' && char <= 'z') ||
		(char >= '0' && char <= '9') ||
		char == '+' || char == '/' || char == '-' || char == '_' || char == '='
}

// Helper function to get lengths of JWT parts
func getPartsLengths(parts []string) []string {
	lengths := make([]string, len(parts))
	for i, part := range parts {
		lengths[i] = fmt.Sprintf("%d", len(part))
	}
	return lengths
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// InitRoutes initializes routes
func InitRoutes(db *gorm.DB, r *chi.Mux) {
	// Initialize JWT auth
	utils.InitJWTAuth()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
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

	// Initialize link token repository for cleanup task
	linkTokenRepo := postgres.NewLinkTokenRepository(db)

	// Start cleanup task for expired tokens (every hour)
	utils.StartTokenCleanupTask(linkTokenRepo, time.Hour)

	// Public routes (NO AUTHENTICATION REQUIRED)
	routes.RegisterAuthRoutes(r, authHandlers)

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

	userRepo := postgres.NewUserRepository(db)
	wordRepo := postgres.NewWordRepository(db)
	sentenceRepo := postgres.NewSentenceRepository(db)
	learnedWordRepo := postgres.NewLearnedWordRepository(db)
	preferenceRepo := postgres.NewPreferenceRepository(db)
	pickOptionRepo := postgres.NewPickOptionRepository(db)
	topicRepo := postgres.NewTopicRepository(db)
	lessonRepo := postgres.NewLessonRepository(db)
	chatHistoryRepo := postgres.NewChatHistoryRepository(db)
	notLearnedWordRepo := postgres.NewNotLearnedWordRepository(db)

	thesaurusClient := utils.NewThesaurusClient(utils.ThesaurusClientConfig{})
	llmClient := utils.NewLLMClient(utils.LLMClientConfig{})
	distractorClient := utils.NewDistractorClient(utils.DistractorClientConfig{})

	chatHistoryHandler := &handlers.ChatHistoryHandler{Repo: chatHistoryRepo}

	// Initialize Telegram handler with all required repositories
	telegramHandler := &handlers.TelegramHandler{
		UserRepo:         userRepo,
		LinkTokenRepo:    linkTokenRepo,
		RefreshTokenRepo: postgres.NewRefreshTokenRepository(db),
	}

	// Register telegram routes now that handler is initialized
	routes.RegisterTelegramRoutes(r, telegramHandler)

	// Protected routes using flexible JWT authentication
	r.Route("/api/v1", func(r chi.Router) {
		// JWT authentication middleware (supports both "Bearer token" and "token" formats)
		r.Use(flexibleJWTVerifier)
		r.Use(authMiddleware.CustomAuthenticator)

		// Protected API routes
		routes.RegisterUserRoutes(r, &handlers.UserHandler{Repo: userRepo})
		routes.RegisterWordRoutes(r, &handlers.WordHandler{Repo: wordRepo})
		routes.RegisterSentenceRoutes(r, &handlers.SentenceHandler{Repo: sentenceRepo})
		routes.RegisterLearnedWordRoutes(r, &handlers.LearnedWordHandler{Repo: learnedWordRepo})
		routes.RegisterPreferencesRoutes(r, &handlers.PreferenceHandler{Repo: preferenceRepo})
		routes.RegisterPickOptionRoutes(r, &handlers.PickOptionHandler{Repo: pickOptionRepo})
		routes.RegisterTopicRoutes(r, &handlers.TopicHandler{Repo: topicRepo})
		routes.RegisterProgressRoutes(r, &handlers.ProgressHandler{
			WordRepo:           wordRepo,
			LearnedWordRepo:    learnedWordRepo,
			NotLearnedWordRepo: notLearnedWordRepo,
		})
		routes.RegisterDayWordRoutes(r, &handlers.DayWordHandler{
			WordRepo:        wordRepo,
			PreferenceRepo:  preferenceRepo,
			TopicRepo:       topicRepo,
			SentenceRepo:    sentenceRepo,
			PickOptionRepo:  pickOptionRepo,
			LearnedWordRepo: learnedWordRepo,
		})
		routes.RegisterLessonRoutes(r, &handlers.LessonHandler{
			PreferenceRepo:  preferenceRepo,
			TopicRepo:       topicRepo,
			SentenceRepo:    sentenceRepo,
			PickOptionRepo:  pickOptionRepo,
			WordRepo:        wordRepo,
			Repo:            lessonRepo,
			LearnedWordRepo: learnedWordRepo,
			ThesaurusClient: thesaurusClient,
		})

		// --- new AI-related routes ---
		chatHandler := &handlers.ChatHandler{
			Redis:              utils.Redis(),
			HistoryRepo:        chatHistoryRepo,
			LLMClient:          llmClient,
			LearnedWordRepo:    learnedWordRepo,
			NotLearnedWordRepo: notLearnedWordRepo,
			WordRepo:           wordRepo,
			TopicRepo:          topicRepo,
		}
		distractorHandler := &handlers.DistractorHandler{
			Client: distractorClient,
		}
		thesaurusHandler := &handlers.ThesaurusHandler{
			Client: thesaurusClient,
		}

		routes.RegisterChatRoutes(r, chatHandler, chatHistoryHandler)
		routes.RegisterDistractorRoutes(r, distractorHandler)
		routes.RegisterThesaurusRoutes(r, thesaurusHandler)
	})
}
