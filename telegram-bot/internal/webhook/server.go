package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fluently/telegram-bot/config"
	"fluently/telegram-bot/internal/bot/handlers"
	"fluently/telegram-bot/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	tele "gopkg.in/telebot.v3"
)

// Server represents the webhook server
type Server struct {
	router      *chi.Mux
	config      *config.Config
	rateLimiter *rate.Limiter
	handlers    *handlers.HandlerService
}

// TelegramUpdate represents incoming Telegram update
type TelegramUpdate struct {
	UpdateID int                   `json:"update_id"`
	Message  *tele.Message         `json:"message,omitempty"`
	Callback *tele.Callback        `json:"callback_query,omitempty"`
	Query    *tele.Query           `json:"inline_query,omitempty"`
	Chat     *tele.ChatJoinRequest `json:"chat_join_request,omitempty"`
}

// NewServer creates a new webhook server
func NewServer(cfg *config.Config, handlerService *handlers.HandlerService) *Server {
	r := chi.NewRouter()

	// Rate limiter: 100 requests per second with burst of 200
	limiter := rate.NewLimiter(100, 200)

	server := &Server{
		router:      r,
		config:      cfg,
		rateLimiter: limiter,
		handlers:    handlerService,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

// setupMiddleware configures middleware
func (s *Server) setupMiddleware() {
	// Basic middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(30 * time.Second))

	// CORS middleware
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Rate limiting middleware
	s.router.Use(s.rateLimitMiddleware)

	// Webhook secret validation middleware
	s.router.Use(s.webhookSecretMiddleware)

	// Request size limiting
	s.router.Use(s.limitRequestSizeMiddleware)
}

// setupRoutes configures routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.Get("/health", s.healthHandler)
	s.router.Get("/ready", s.readinessHandler)

	// Webhook endpoint
	s.router.Post(s.config.Webhook.Path, s.webhookHandler)

	// Metrics endpoint (optional)
	s.router.Get("/metrics", s.metricsHandler)
}

// rateLimitMiddleware implements rate limiting
func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.rateLimiter.Allow() {
			logger.Log.Warn("Rate limit exceeded", zap.String("remote_addr", r.RemoteAddr))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// webhookSecretMiddleware validates webhook secret
func (s *Server) webhookSecretMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip secret validation for health checks
		if r.URL.Path == "/health" || r.URL.Path == "/ready" || r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		if s.config.Webhook.Secret != "" {
			secretToken := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if secretToken != s.config.Webhook.Secret {
				logger.Log.Warn("Invalid webhook secret", zap.String("remote_addr", r.RemoteAddr))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// limitRequestSizeMiddleware limits request body size
func (s *Server) limitRequestSizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, s.config.Webhook.MaxBodySize)
		next.ServeHTTP(w, r)
	})
}

// webhookHandler handles incoming Telegram updates
func (s *Server) webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Log incoming request for debugging
	logger.Log.Debug("Received webhook request",
		zap.String("method", r.Method),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("content_type", r.Header.Get("Content-Type")))

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("Failed to read request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Log body size for debugging
	logger.Log.Debug("Request body received", zap.Int("size", len(body)))

	// Parse Telegram update
	var update TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		logger.Log.Error("Failed to parse update",
			zap.Error(err),
			zap.String("body", string(body)),
			zap.Int("body_length", len(body)))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Log successful parsing
	logger.Log.Debug("Successfully parsed update", zap.Int("update_id", update.UpdateID))

	// Immediately return 200 OK to Telegram
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))

	// Process update asynchronously
	go s.processUpdate(r.Context(), &update)
}

// processUpdate processes Telegram update asynchronously
func (s *Server) processUpdate(ctx context.Context, update *TelegramUpdate) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("Panic in update processing", zap.Any("panic", r))
		}
	}()

	logger.Log.Debug("Processing update", zap.Int("update_id", update.UpdateID))

	// Process different types of updates
	if update.Message != nil {
		s.handlers.ProcessMessage(ctx, update.Message)
	} else if update.Callback != nil {
		s.handlers.ProcessCallback(ctx, update.Callback)
	} else if update.Query != nil {
		s.handlers.ProcessInlineQuery(ctx, update.Query)
	} else {
		logger.Log.Debug("Unhandled update type", zap.Int("update_id", update.UpdateID))
	}
}

// healthHandler handles health check requests
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("Health check requested", zap.String("remote_addr", r.RemoteAddr))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
		"service":   "telegram-bot-webhook",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode health response", zap.Error(err))
	}
}

// readinessHandler handles readiness check requests
func (s *Server) readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Check if services are ready (Redis, API, etc.)
	ready := s.checkServicesReady()

	w.Header().Set("Content-Type", "application/json")

	if ready {
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":    "ready",
			"timestamp": time.Now().Unix(),
		}
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		response := map[string]interface{}{
			"status":    "not ready",
			"timestamp": time.Now().Unix(),
		}
		json.NewEncoder(w).Encode(response)
	}
}

// metricsHandler handles metrics requests
func (s *Server) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In production, you would implement proper metrics collection
	response := map[string]interface{}{
		"requests_total": 0,
		"errors_total":   0,
		"uptime":         time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// checkServicesReady checks if all required services are ready
func (s *Server) checkServicesReady() bool {
	// Check Redis connection
	if !s.handlers.CheckRedisHealth() {
		return false
	}

	// Check API connection
	if !s.handlers.CheckAPIHealth() {
		return false
	}

	return true
}

// Start starts the webhook server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Webhook.Host, s.config.Webhook.Port)

	logger.Log.Info("Starting webhook server",
		zap.String("addr", addr),
		zap.String("webhook_path", s.config.Webhook.Path))

	if s.config.Webhook.CertFile != "" && s.config.Webhook.KeyFile != "" {
		// HTTPS server
		return http.ListenAndServeTLS(addr, s.config.Webhook.CertFile, s.config.Webhook.KeyFile, s.router)
	} else {
		// HTTP server
		return http.ListenAndServe(addr, s.router)
	}
}

// GetRouter returns the router for testing
func (s *Server) GetRouter() *chi.Mux {
	return s.router
}

// validateWebhookSignature validates webhook signature (for enhanced security)
func (s *Server) validateWebhookSignature(body []byte, signature string) bool {
	if s.config.Webhook.Secret == "" {
		return true // No secret configured
	}

	mac := hmac.New(sha256.New, []byte(s.config.Webhook.Secret))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
