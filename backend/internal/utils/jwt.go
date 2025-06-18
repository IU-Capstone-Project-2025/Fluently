package utils

import (
	"time"

	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"

	"github.com/golang-jwt/jwt/v4"
)

// JwtCustomClaims represents the custom claims used in the access token.
// It embeds jwt.RegisteredClaims to get exp, iat, etc. encoded following RFC 7519.
type JwtCustomClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT string for the provided user using the
// HS256 signing method and the secret configured in the application config.
func GenerateJWT(user *models.User) (string, error) {
	cfg := config.GetConfig()

	expiresAt := time.Now().Add(cfg.Auth.JWTExpiration)

	claims := JwtCustomClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.JWTSecret))
}
