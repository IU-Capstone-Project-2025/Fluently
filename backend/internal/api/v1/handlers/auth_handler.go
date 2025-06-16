package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

// Handlers combines multiple repositories needed for authentication
type Handlers struct {
	UserRepo     *postgres.UserRepository
	UserPrefRepo *postgres.PreferenceRepository
}

// GoogleAuthHandler godoc
// @Summary      Authenticates with Google
// @Description  Authenticates with Google using the authorization code flow
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        code query     string  true  "Authorization code"
// @Success      200  {object}  schemas.JwtResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      401  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /auth/google [get]
func (h *Handlers) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Getting idtoken and platform (for audience) from request
	var req struct {
		IDToken  string `json:"id_token"`
		Platform string `json:"platform"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	googleToken := req.IDToken
	platform := req.Platform

	cfg := config.GetConfig()
	var googleClientIDs = map[string]string{
		"ios":     fmt.Sprintf("%s.apps.googleusercontent.com", cfg.Google.IosClientID),
		"android": fmt.Sprintf("%s.apps.googleusercontent.com", cfg.Google.AndroidClientID),
		"web":     cfg.Google.WebClientID,
	}

	audience := googleClientIDs[platform]

	payload, err := idtoken.Validate(r.Context(), googleToken, audience)
	if err != nil {
		logger.Log.Error("Invalid token", zap.Error(err))
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Process claims
	claims := payload.Claims
	sub := claims["sub"].(string)
	email := claims["email"].(string)
	emailVerified := claims["email_verified"].(bool)
	name := claims["name"].(string)
	avatar := claims["picture"].(string)

	// Check if email is verified
	if !emailVerified {
		logger.Log.Error("Email not verified", zap.String("email", email))
		http.Error(w, "email not verified", http.StatusBadRequest)
		return
	}

	// Check if user exists
	user, err := h.UserRepo.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Firstly creating user preferences
			userID := uuid.New()
			prefID := uuid.New()
			userPreferences := models.Preference{
				ID:              prefID,
				Subscribed:      false,
				CEFRLevel:       0,
				FactEveryday:    false,
				Notifications:   true,
				NotificationsAt: nil,
				WordsPerDay:     10,
				Goal:            "Learn new words",
				AvatarImage:     []byte(avatar),
			}
			if err := h.UserPrefRepo.Create(r.Context(), &userPreferences); err != nil {
				logger.Log.Error("Failed to create user preferences", zap.Error(err))
				http.Error(w, "failed to create user preferences", http.StatusInternalServerError)
				return
			}

			logger.Log.Info("Creating new user with email: ", zap.String("email", email))
			newUser := &models.User{
				ID:           userID,
				Name:         name,
				Email:        email,
				PasswordHash: "",
				Provider:     "google",
				GoogleID:     sub,
				Role:         "user",
				IsActive:     true,
				PrefID:       &prefID,
				LastLoginAt:  time.Now(),
			}

			if err := h.UserRepo.Create(r.Context(), newUser); err != nil {
				logger.Log.Error("Failed to create user", zap.Error(err))
				http.Error(w, "failed to create user", http.StatusInternalServerError)
				return
			}
			user = newUser
		} else {
			logger.Log.Error("Failed to get user", zap.Error(err))
			http.Error(w, "failed to get user", http.StatusInternalServerError)
			return
		}
	} else {
		// Update last login time for existing user
		if err := h.UserRepo.UpdateLastLogin(r.Context(), user.ID); err != nil {
			logger.Log.Error("Failed to update last login time", zap.Error(err))
		}
	}

	// TODO: Generate and return JWT token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Authentication successful",
		"user_id": user.ID.String(),
	})
}
