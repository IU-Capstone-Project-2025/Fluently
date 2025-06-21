package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
	"fluently/go-backend/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

// SwaggerOAuthCallbackHandler godoc
// @Summary      Handles Swagger OAuth callback
// @Description  Exchanges authorization code for tokens and redirects to Swagger with token
// @Tags         auth
// @Produce      html
// @Success      302  {string}  string "Redirect to Swagger UI"
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /swagger/oauth2-redirect.html [get]
func (h *Handlers) SwaggerOAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code not found", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")

	oauthCfg := config.GoogleOAuthConfig()
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	oauthCfg.RedirectURL = fmt.Sprintf("%s://%s/swagger/oauth2-redirect.html", scheme, r.Host)

	token, err := oauthCfg.Exchange(r.Context(), code)
	if err != nil {
		logger.Log.Error("Code exchange failed (Swagger)", zap.Error(err))
		http.Error(w, "code exchange failed", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "id_token missing", http.StatusInternalServerError)
		return
	}

	resp, err := h.processGoogleIDTokenForSwagger(r, rawIDToken)
	if err != nil {
		logger.Log.Error("Token processing failed (Swagger)", zap.Error(err))
		http.Error(w, "token processing failed", http.StatusInternalServerError)
		return
	}

	// Redirect to Swagger's static oauth2-redirect.html with token in hash fragment
	redirectTarget := fmt.Sprintf("/swagger/oauth2-redirect.html#access_token=%s&token_type=Bearer&expires_in=%d&state=%s", resp.AccessToken, resp.ExpiresIn, state)
	http.Redirect(w, r, redirectTarget, http.StatusTemporaryRedirect)
}

// processGoogleIDTokenForSwagger processes the Google id_token and returns JwtResponse
func (h *Handlers) processGoogleIDTokenForSwagger(r *http.Request, googleToken string) (*schemas.JwtResponse, error) {
	cfg := config.GetConfig()
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
		return nil, err
	}

	claims := payload.Claims
	sub := claims["sub"].(string)
	email := claims["email"].(string)
	emailVerified := claims["email_verified"].(bool)
	name := claims["name"].(string)
	avatar := claims["picture"].(string)

	if !emailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	user, err := h.UserRepo.GetByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create new user
			uid := uuid.New()
			prefID := uuid.New()
			pref := models.Preference{
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
			if err := h.UserPrefRepo.Create(r.Context(), &pref); err != nil {
				return nil, err
			}
			user = &models.User{
				ID:          uid,
				Name:        name,
				Email:       email,
				Provider:    "google",
				GoogleID:    sub,
				Role:        "user",
				IsActive:    true,
				PrefID:      &prefID,
				LastLoginAt: time.Now(),
			}
			if err := h.UserRepo.Create(r.Context(), user); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		_ = h.UserRepo.UpdateLastLogin(r.Context(), user.ID)
	}

	access, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, err
	}
	refresh, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	rtModel := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refresh,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	if err := h.RefreshTokenRepo.Create(r.Context(), rtModel); err != nil {
		return nil, err
	}

	return &schemas.JwtResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    int64(cfg.Auth.JWTExpiration.Seconds()),
	}, nil
}
