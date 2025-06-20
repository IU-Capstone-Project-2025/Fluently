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

	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

// Handlers combines multiple repositories needed for authentication
type Handlers struct {
	UserRepo         *postgres.UserRepository
	UserPrefRepo     *postgres.PreferenceRepository
	RefreshTokenRepo *postgres.RefreshTokenRepository
}

func generateRandomState() (string, error) {
	return utils.GenerateRefreshToken() // reuse secure random generator
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
// @Router       /auth/google [post]
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

	resp, err := h.generateTokens(user, w, r)
	if err != nil {
		logger.Log.Error("Failed to generate tokens", zap.Error(err))
		http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handlers) generateTokens(user *models.User, w http.ResponseWriter, r *http.Request) (schemas.JwtResponse, error) {
	// Generate JWT token
	tokenString, err := utils.GenerateJWT(user)
	if err != nil {
		logger.Log.Error("Failed to generate JWT", zap.Error(err))
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return schemas.JwtResponse{}, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		logger.Log.Error("Failed to generate refresh token", zap.Error(err))
		http.Error(w, "failed to generate refresh token", http.StatusInternalServerError)
		return schemas.JwtResponse{}, err
	}

	refreshTokenModel := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	if err := h.RefreshTokenRepo.Create(r.Context(), refreshTokenModel); err != nil {
		logger.Log.Error("Failed to store refresh token", zap.Error(err))
		http.Error(w, "failed to store refresh token", http.StatusInternalServerError)
		return schemas.JwtResponse{}, err
	}

	resp := schemas.JwtResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(config.GetConfig().Auth.JWTExpiration.Seconds()),
	}

	return resp, nil
}

// LoginHandler godoc
// @Summary      Login with email & password
// @Description  Authenticates user and returns JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      schemas.LoginRequest  true  "Email & Password"
// @Success      200  {object}  schemas.JwtResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      401  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /auth/login [post]
func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req schemas.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		logger.Log.Error("User not found", zap.Error(err))
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	resp, err := h.generateTokens(user, w, r)
	if err != nil {
		logger.Log.Error("Failed to generate tokens", zap.Error(err))
		http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RegisterHandler godoc
// @Summary      Register with email & password
// @Description  Creates a user, hashes password, returns JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      schemas.RegisterRequest  true  "Registration data"
// @Success      201  {object}  schemas.JwtResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      409  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /auth/register [post]
func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req schemas.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	if _, err := h.UserRepo.GetByEmail(r.Context(), req.Email); err == nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error("Failed to hash password", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Provider:     "password",
		Role:         "user",
		IsActive:     true,
		LastLoginAt:  time.Now(),
	}

	if err := h.UserRepo.Create(r.Context(), user); err != nil {
		logger.Log.Error("Failed to create user", zap.Error(err))
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	resp, err := h.generateTokens(user, w, r)
	if err != nil {
		logger.Log.Error("Failed to generate tokens", zap.Error(err))
		http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// RefreshTokenHandler godoc
// @Summary      Refresh access token
// @Description  Rotates refresh token and returns new access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refreshToken body object true "Refresh Token"
// @Success      200  {object}  schemas.JwtResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      401  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /auth/refresh [post]
func (h *Handlers) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if r.Body == nil {
		logger.Log.Error("Empty request")
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	rt, err := h.RefreshTokenRepo.GetByToken(r.Context(), req.RefreshToken)
	if err != nil {
		if rt.Revoked {
			logger.Log.Error("Invalid refresh token", zap.Error(err))
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		} else if rt.ExpiresAt.Before(time.Now()) {
			logger.Log.Error("Refresh token expired", zap.Error(err))
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		logger.Log.Error("Invalid refresh token", zap.Error(err))
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	user, err := h.UserRepo.GetByID(r.Context(), rt.UserID)
	if err != nil {
		logger.Log.Error("User not found", zap.Error(err))
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// revoke old refresh token
	if err := h.RefreshTokenRepo.Revoke(r.Context(), rt.ID); err != nil {
		logger.Log.Error("Could not revoke token", zap.Error(err))
		http.Error(w, "could not revoke token", http.StatusInternalServerError)
		return
	}

	resp, err := h.generateTokens(user, w, r)
	if err != nil {
		logger.Log.Error("Failed to generate tokens", zap.Error(err))
		http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GoogleAuthRedirectHandler godoc
// @Summary      Redirects to Google OAuth consent screen
// @Description  Initiates Google OAuth 2.0 authorization code flow
// @Tags         auth
// @Produce      json
// @Success      302  {string}  string "Redirect"
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /auth/google [get]
func (h *Handlers) GoogleAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		logger.Log.Error("Failed to generate state", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Store state in a secure cookie for later verification
	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	oauthCfg := config.GoogleOAuthConfig()
	url := oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallbackHandler godoc
// @Summary      Handles Google OAuth callback
// @Description  Exchanges authorization code for tokens and signs user in
// @Tags         auth
// @Produce      json
// @Success      200  {object}  schemas.JwtResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      401  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /auth/google/callback [get]
func (h *Handlers) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Verify state parameter
	state := r.URL.Query().Get("state")
	cookie, err := r.Cookie("oauthstate")
	if err != nil || cookie.Value != state {
		logger.Log.Error("Invalid OAuth state", zap.Error(err))
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code not found", http.StatusBadRequest)
		return
	}

	oauthCfg := config.GoogleOAuthConfig()
	token, err := oauthCfg.Exchange(r.Context(), code)
	if err != nil {
		logger.Log.Error("Code exchange failed", zap.Error(err))
		http.Error(w, "code exchange failed", http.StatusUnauthorized)
		return
	}

	// Extract id_token from token response
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		logger.Log.Error("ID token missing in OAuth response")
		http.Error(w, "id_token missing", http.StatusInternalServerError)
		return
	}

	// Reuse existing logic to process id_token
	processGoogleIDToken(h, w, r, rawIDToken)
}

func processGoogleIDToken(h *Handlers, w http.ResponseWriter, r *http.Request, googleToken string) {
	cfg := config.GetConfig()
	// For validation we accept any of the configured client IDs as audience
	audiences := []string{
		cfg.Google.WebClientID,
		fmt.Sprintf("%s.apps.googleusercontent.com", cfg.Google.IosClientID),
		fmt.Sprintf("%s.apps.googleusercontent.com", cfg.Google.AndroidClientID),
	}

	var payload *idtoken.Payload
	var err error
	for _, aud := range audiences {
		payload, err = idtoken.Validate(r.Context(), googleToken, aud)
		if err == nil {
			break
		}
	}
	if err != nil {
		logger.Log.Error("Invalid id_token", zap.Error(err))
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	claims := payload.Claims
	sub := claims["sub"].(string)
	email := claims["email"].(string)
	emailVerified := claims["email_verified"].(bool)
	name := claims["name"].(string)
	avatar := claims["picture"].(string)

	if !emailVerified {
		logger.Log.Error("Email not verified", zap.String("email", email))
		http.Error(w, "email not verified", http.StatusBadRequest)
		return
	}

	// Check if user exists
	user, err := h.UserRepo.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new user and preferences
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
		// Update last login time
		if err := h.UserRepo.UpdateLastLogin(r.Context(), user.ID); err != nil {
			logger.Log.Error("Failed to update last login time", zap.Error(err))
		}
	}

	resp, err := h.generateTokens(user, w, r)
	if err != nil {
		logger.Log.Error("Failed to generate tokens", zap.Error(err))
		http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
