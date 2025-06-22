package utils

import (
	"time"

	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"

	"crypto/rand"
	"encoding/base64"

	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth

// InitJWTAuth initializes the JWT authenticator with the secret from config
func InitJWTAuth() {
	cfg := config.GetConfig()
	if cfg.Auth.JWTSecret == "" {
		panic("JWT_SECRET environment variable is required but not set")
	}
	TokenAuth = jwtauth.New("HS256", []byte(cfg.Auth.JWTSecret), nil)
}

// GenerateJWT creates a signed JWT string for the provided user using go-chi/jwtauth
func GenerateJWT(user *models.User) (string, error) {
	cfg := config.GetConfig()

	claims := map[string]interface{}{
		"sub":   user.ID.String(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(cfg.Auth.JWTExpiration).Unix(),
		"iat":   time.Now().Unix(),
	}

	_, tokenString, err := TokenAuth.Encode(claims)
	return tokenString, err
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
