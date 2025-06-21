package middleware

import (
	"net/http"
	"strings"

	"fluently/go-backend/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

// SwaggerAuthMiddleware checks JWT token and ensures the email is allowed to use Swagger.
// Allowed emails are taken from the environment variable SWAGGER_ALLOWED_EMAILS (comma-separated list).
// If the variable is empty, Swagger will be accessible without additional restrictions.
func SwaggerAuthMiddleware(next http.Handler) http.Handler {
	// Retrieve whitelist from application configuration (set in config.Init).
	whitelist := config.GetConfig().Swagger.AllowedEmails

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip whitelist check if variable not set or empty
		if len(whitelist) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().Auth.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		emailRaw, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !whitelist[emailRaw] {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
