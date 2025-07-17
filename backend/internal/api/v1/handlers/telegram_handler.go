package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

// TelegramHandler handles the telegram endpoint
type TelegramHandler struct {
	UserRepo      *postgres.UserRepository
	LinkTokenRepo *postgres.LinkTokenRepository
}

// generateLinkToken generates a random link token
func generateLinkToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateLinkToken godoc
// @Summary      Create link token for Telegram account linking
// @Description  Create a magic link for linking Telegram ID with Google account
// @Tags         telegram
// @Accept       json
// @Produce      json
// @Param        request  body      schemas.TelegramLinkRequest  true  "Telegram ID"
// @Success      200  {object}  schemas.TelegramLinkResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      409  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /telegram/create-link [post]
func (h *TelegramHandler) CreateLinkToken(w http.ResponseWriter, r *http.Request) {
	var req schemas.TelegramLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("Invalid request body", zap.Error(err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Check if Telegram ID is already linked
	if _, err := h.UserRepo.GetByTelegramID(r.Context(), req.TelegramID); err == nil {
		logger.Log.Warn("Telegram ID already linked", zap.Int64("telegram_id", req.TelegramID))
		http.Error(w, "telegram account already linked", http.StatusConflict)
		return
	}

	// Delete old tokens for this Telegram ID
	if existingTokens, err := h.LinkTokenRepo.GetActiveTelegramTokens(r.Context(), req.TelegramID); err == nil {
		for _, token := range existingTokens {
			h.LinkTokenRepo.MarkAsUsed(r.Context(), token.ID)
		}
	}

	// Generate link token
	token, err := generateLinkToken()
	if err != nil {
		logger.Log.Error("Failed to generate link token", zap.Error(err))
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// Create link token
	linkToken := &models.LinkToken{
		Token:      token,
		TelegramID: req.TelegramID,
		Used:       false,
		ExpiresAt:  time.Now().Add(15 * time.Minute), // 15 минут на связывание
	}

	if err := h.LinkTokenRepo.Create(r.Context(), linkToken); err != nil {
		logger.Log.Error("Failed to create link token", zap.Error(err))
		http.Error(w, "failed to create link token", http.StatusInternalServerError)
		return
	}

	// Format link URL
	ExternalHostName := "fluently-app.ru"
	linkURL := fmt.Sprintf("https://%s/link-google?token=%s", ExternalHostName, token)

	resp := schemas.TelegramLinkResponse{
		Token:     token,
		LinkURL:   linkURL,
		ExpiresAt: linkToken.ExpiresAt.Format(time.RFC3339),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CheckLinkStatus godoc
// @Summary      Check status of Telegram account linking
// @Description  Check if Telegram ID is linked to any account
// @Tags         telegram
// @Accept       json
// @Produce      json
// @Param        request  body      schemas.TelegramLinkStatusRequest  true  "Telegram ID"
// @Success      200  {object}  schemas.TelegramLinkStatusResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /telegram/check-status [post]
func (h *TelegramHandler) CheckLinkStatus(w http.ResponseWriter, r *http.Request) {
	var req schemas.TelegramLinkStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetByTelegramID(r.Context(), req.TelegramID)
	var resp schemas.TelegramLinkStatusResponse

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp = schemas.TelegramLinkStatusResponse{
				IsLinked: false,
				Message:  "Telegram account not linked",
			}
		} else {
			logger.Log.Error("Failed to check telegram status", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		resp = schemas.TelegramLinkStatusResponse{
			IsLinked: true,
			User: &schemas.UserBasic{
				ID:    user.ID.String(),
				Name:  user.Name,
				Email: user.Email,
			},
			Message: "Telegram account successfully linked",
		}
	}

	logger.Log.Info("CheckLinkStatus: ", zap.Any("resp", resp))

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// LinkWithGoogle godoc
// @Summary      Link Telegram account with Google account
// @Description  Process magic link and link Telegram ID with Google account
// @Tags         telegram
// @Accept       json
// @Produce      json
// @Param        token   query      string  true  "Link token"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      401  {object}  schemas.ErrorResponse
// @Failure      410  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /link-google [get]
func (h *TelegramHandler) LinkWithGoogle(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token parameter required", http.StatusBadRequest)
		return
	}

	// Fund link token
	linkToken, err := h.LinkTokenRepo.GetByToken(r.Context(), token)
	if err != nil {
		logger.Log.Error("Link token not found", zap.Error(err))
		http.Error(w, "invalid or expired token", http.StatusGone)
		return
	}

	// Check expiration and usage
	if linkToken.Used || linkToken.ExpiresAt.Before(time.Now()) {
		logger.Log.Warn("Link token expired or used",
			zap.Bool("used", linkToken.Used),
			zap.Time("expires_at", linkToken.ExpiresAt))
		http.Error(w, "token expired or already used", http.StatusGone)
		return
	}

	/*
	* If user is not authenticated, redirect to Google OAuth
	* then connect his account with Telegram ID
	*
	* Here should be check if user is authenticated
	* If user is not authenticated, redirect to Google OAuth
	 */

	// HTML example for redirect
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Связывание аккаунта</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .container { max-width: 400px; margin: 0 auto; }
        .btn { background: #4285f4; color: white; padding: 12px 24px; 
               text-decoration: none; border-radius: 4px; display: inline-block; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔗 Связывание аккаунтов</h1>
        <p>Для связывания вашего Telegram аккаунта с Google, войдите в систему:</p>
        <a href="/auth/google?redirect_uri=%s&state=%s" class="btn">
            Войти через Google
        </a>
    </div>
</body>
</html>`

	// Form redirect_uri and state
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	redirectURI := fmt.Sprintf("%s://%s/link-google/callback", scheme, r.Host)
	state := token // use token as state

	// Send response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, html, redirectURI, state)
}

// LinkGoogleCallback godoc
// @Summary      Callback for linking with Google
// @Description  Process callback from Google OAuth and finish linking
// @Tags         telegram
// @Accept       json
// @Produce      json
// @Param        code    query      string  true  "OAuth code"
// @Param        state   query      string  true  "OAuth state (contains link token)"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /link-google/callback [get]
func (h *TelegramHandler) LinkGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Check parameters - use state as token
	if code == "" || state == "" {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}

	token := state // token is passed as state parameter

	// Check link token
	linkToken, err := h.LinkTokenRepo.GetByToken(r.Context(), token)
	if err != nil || linkToken.Used || linkToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "invalid or expired token", http.StatusGone)
		return
	}

	// Here should be OAuth code exchange and get user info
	// Use existing logic from auth_handler.go
	cfg := config.GetConfig()

	// Configure OAuth config for callback
	oauthCfg := config.GoogleOAuthConfig()
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	oauthCfg.RedirectURL = fmt.Sprintf("%s://%s/link-google/callback", scheme, r.Host)

	// Exchange code for token
	oauthToken, err := oauthCfg.Exchange(r.Context(), code)
	if err != nil {
		logger.Log.Error("OAuth code exchange failed", zap.Error(err))
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	// Get ID token
	rawIDToken, ok := oauthToken.Extra("id_token").(string)
	if !ok {
		logger.Log.Error("ID token missing in OAuth response")
		http.Error(w, "authorization failed", http.StatusInternalServerError)
		return
	}

	// Validate ID token
	audiences := []string{
		cfg.Google.WebClientID,
		fmt.Sprintf("%s.apps.googleusercontent.com", cfg.Google.IosClientID),
		fmt.Sprintf("%s.apps.googleusercontent.com", cfg.Google.AndroidClientID),
	}

	var payload *idtoken.Payload
	for _, aud := range audiences {
		payload, err = idtoken.Validate(r.Context(), rawIDToken, aud)
		if err == nil {
			break
		}
	}
	if err != nil {
		logger.Log.Error("Invalid ID token", zap.Error(err))
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	// Get user info
	claims := payload.Claims
	userEmail := claims["email"].(string)
	emailVerified := claims["email_verified"].(bool)

	if !emailVerified {
		logger.Log.Error("Email not verified", zap.String("email", userEmail))
		http.Error(w, "email not verified", http.StatusBadRequest)
		return
	}

	// Find user by email
	user, err := h.UserRepo.GetByEmail(r.Context(), userEmail)
	if err != nil {
		logger.Log.Error("User not found", zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Link Telegram ID with user
	if err := h.UserRepo.LinkTelegramID(r.Context(), user.ID, linkToken.TelegramID); err != nil {
		logger.Log.Error("Failed to link telegram ID", zap.Error(err))
		http.Error(w, "failed to link account", http.StatusInternalServerError)
		return
	}

	// Mark token as used
	if err := h.LinkTokenRepo.MarkAsUsed(r.Context(), linkToken.ID); err != nil {
		logger.Log.Error("Failed to mark token as used", zap.Error(err))
	}

	// Send success HTML response
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Успешно!</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .container { max-width: 400px; margin: 0 auto; }
        .success { color: #28a745; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="success">✅ Готово!</h1>
        <p>Ваш Telegram аккаунт успешно связан с Google аккаунтом.</p>
        <p>Теперь можете вернуться в Telegram бот.</p>
    </div>
</body>
</html>`

	// Send HTML response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

// UnlinkTelegram godoc
// @Summary      Unlink Telegram account
// @Description  Delete the link between Telegram ID and user account
// @Tags         telegram
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      schemas.TelegramUnlinkRequest  true  "Telegram ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/telegram/unlink [post]
func (h *TelegramHandler) UnlinkTelegram(w http.ResponseWriter, r *http.Request) {
	var req schemas.TelegramUnlinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by Telegram ID
	user, err := h.UserRepo.GetByTelegramID(r.Context(), req.TelegramID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "telegram account not linked", http.StatusNotFound)
			return
		}
		logger.Log.Error("Failed to find user", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Unlink Telegram ID
	if err := h.UserRepo.UnlinkTelegramID(r.Context(), user.ID); err != nil {
		logger.Log.Error("Failed to unlink telegram", zap.Error(err))
		http.Error(w, "failed to unlink account", http.StatusInternalServerError)
		return
	}

	// Send success JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Telegram account successfully unlinked",
	})
}
