package middleware

import (
	"net/http"

	"fluently/go-backend/pkg/logger"

	"go.uber.org/zap"
)

// DebugAuthRoutes logs all requests to auth routes for debugging
func DebugAuthRoutes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info("Auth route accessed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Any("headers", r.Header),
		)
		next.ServeHTTP(w, r)
	})
}
