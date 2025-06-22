package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/pkg/logger"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

// CustomAuthenticator is a custom version of jwtauth.Authenticator that returns JSON errors
// instead of plain text and adds user information to context
func CustomAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err != nil {
			logger.Log.Error("JWT authentication error", zap.Error(err))
			writeJSONError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if token == nil {
			logger.Log.Error("Invalid JWT token")
			writeJSONError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Check whitelist if configured
		whitelist := config.GetConfig().Swagger.AllowedEmails
		if len(whitelist) > 0 {
			email, _ := claims["email"].(string)
			if !whitelist[email] {
				logger.Log.Error("Email not whitelisted", zap.String("email", email))
				writeJSONError(w, "forbidden", http.StatusForbidden)
				return
			}
		}

		// Extract user ID from claims
		userIDStr, ok := claims["sub"].(string)
		if !ok {
			logger.Log.Error("Invalid user ID in token")
			writeJSONError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.Log.Error("Invalid UUID in token", zap.Error(err))
			writeJSONError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Create user context with additional info from claims
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)

		user := &models.User{
			ID:    userID,
			Email: email,
			Role:  role,
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the user from the context
func GetUserFromContext(ctx context.Context) *models.User {
	if user, ok := ctx.Value(UserContextKey).(*models.User); ok {
		return user
	}
	return nil
}

// writeJSONError writes a JSON error response
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
