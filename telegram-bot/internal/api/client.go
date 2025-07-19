package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"telegram-bot/internal/domain"

	"telegram-bot/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Client represents the API client for backend communication
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AuthRequest represents authentication request
type AuthRequest struct {
	TelegramID int64 `json:"telegram_id"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	IsLinked  bool   `json:"is_linked"`
	LinkToken string `json:"link_token,omitempty"`
	LinkURL   string `json:"link_url,omitempty"`
}

// CreateLinkTokenRequest represents link token creation request
type CreateLinkTokenRequest struct {
	TelegramID int64 `json:"telegram_id"`
}

// CreateLinkTokenResponse represents link token creation response
type CreateLinkTokenResponse struct {
	Token     string `json:"token"`
	LinkURL   string `json:"link_url"`
	ExpiresAt string `json:"expires_at"`
}

// CheckLinkStatusRequest represents link status check request
type CheckLinkStatusRequest struct {
	TelegramID int64 `json:"telegram_id"`
}

// CheckLinkStatusResponse represents link status check response
type CheckLinkStatusResponse struct {
	IsLinked bool      `json:"is_linked"`
	User     *UserInfo `json:"user,omitempty"`
	Message  string    `json:"message"`
}

// UserInfo represents basic user information
type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// JWTResponse represents JWT token response from backend
type JWTResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// PreferenceResponse represents user preferences from backend
type PreferenceResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	CEFRLevel      string `json:"cefr_level"`
	WordsPerDay    int    `json:"words_per_day"`
	NotificationAt string `json:"notification_at"`
	Notifications  bool   `json:"notifications"`
	Goal           string `json:"goal"`
	FactEveryday   bool   `json:"fact_everyday"`
	Subscribed     bool   `json:"subscribed"`
	AvatarImageURL string `json:"avatar_image_url"`
}

// CreatePreferenceRequest represents user preferences creation request
type CreatePreferenceRequest struct {
	CEFRLevel      string     `json:"cefr_level"`
	FactEveryday   bool       `json:"fact_everyday"`
	Notifications  bool       `json:"notifications"`
	NotificationAt *time.Time `json:"notification_at,omitempty"`
	WordsPerDay    int        `json:"words_per_day"`
	Goal           string     `json:"goal"`
	Subscribed     bool       `json:"subscribed"`
	AvatarImageURL string     `json:"avatar_image_url"`
}

// UpdatePreferenceRequest represents user preferences update request
type UpdatePreferenceRequest struct {
	CEFRLevel      *string    `json:"cefr_level,omitempty"`
	FactEveryday   *bool      `json:"fact_everyday,omitempty"`
	Notifications  *bool      `json:"notifications,omitempty"`
	NotificationAt *time.Time `json:"notification_at,omitempty"`
	WordsPerDay    *int       `json:"words_per_day,omitempty"`
	Goal           *string    `json:"goal,omitempty"`
	Subscribed     *bool      `json:"subscribed,omitempty"`
	AvatarImageURL *string    `json:"avatar_image_url,omitempty"`
}

// WordProgressRequest represents word progress update request
type WordProgressRequest struct {
	UserID    string    `json:"user_id"`
	WordID    uuid.UUID `json:"word_id"`
	Correct   bool      `json:"correct"`
	TimeSpent int       `json:"time_spent"`
}

// ProgressUpdateRequest represents progress update request
type ProgressUpdateRequest struct {
	UserID             string                  `json:"user_id"`
	LessonID           uuid.UUID               `json:"lesson_id"`
	ExercisesCompleted int                     `json:"exercises_completed"`
	CorrectAnswers     int                     `json:"correct_answers"`
	TotalAttempts      int                     `json:"total_attempts"`
	TimeSpent          int                     `json:"time_spent"`
	WordsLearned       []uuid.UUID             `json:"words_learned"`
	ExerciseResults    []ExerciseResultRequest `json:"exercise_results"`
}

// ExerciseResultRequest represents exercise result
type ExerciseResultRequest struct {
	WordID       uuid.UUID `json:"word_id"`
	ExerciseType string    `json:"exercise_type"`
	Correct      bool      `json:"correct"`
	AttemptCount int       `json:"attempt_count"`
	TimeSpent    int       `json:"time_spent"`
}

// ProgressRequest represents progress request for backend API
type ProgressRequest struct {
	WordID          string    `json:"word_id"`
	LearnedAt       time.Time `json:"learned_at"`
	ConfidenceScore int       `json:"confidence_score"`
	CntReviewed     int       `json:"cnt_reviewed"`
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// doRequest performs HTTP request with proper headers
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "fluently-telegram-bot/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// doAuthenticatedRequest performs HTTP request with JWT authentication
func (c *Client) doAuthenticatedRequest(ctx context.Context, method, endpoint string, body interface{}, token string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "fluently-telegram-bot/1.0")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	logger.Log.Info("Sending request", zap.String("url", url), zap.String("method", method), zap.Any("body", body))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// parseResponse parses HTTP response
func (c *Client) parseResponse(resp *http.Response, dest interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Message)
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	if dest != nil {
		if err := json.Unmarshal(body, dest); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// CreateLinkToken creates a link token for Telegram account linking
func (c *Client) CreateLinkToken(ctx context.Context, telegramID int64) (*CreateLinkTokenResponse, error) {
	req := CreateLinkTokenRequest{TelegramID: telegramID}

	resp, err := c.doRequest(ctx, "POST", "/telegram/create-link", req)
	if err != nil {
		zap.L().With(zap.Int64("telegram_id", telegramID), zap.Error(err)).Error("Failed to create link token")
		return nil, err
	}

	var result CreateLinkTokenResponse
	if err := c.parseResponse(resp, &result); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse create link token response")
		return nil, err
	}

	zap.L().With(zap.Int64("telegram_id", telegramID)).Info("Successfully created link token")
	return &result, nil
}

// CheckLinkStatus checks if Telegram account is linked
func (c *Client) CheckLinkStatus(ctx context.Context, telegramID int64) (*CheckLinkStatusResponse, error) {
	req := CheckLinkStatusRequest{TelegramID: telegramID}

	resp, err := c.doRequest(ctx, "POST", "/telegram/check-status", req)
	if err != nil {
		zap.L().With(zap.Int64("telegram_id", telegramID), zap.Error(err)).Error("Failed to check link status")
		return nil, err
	}

	var result CheckLinkStatusResponse
	if err := c.parseResponse(resp, &result); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse check link status response")
		return nil, err
	}

	zap.L().With(zap.Int64("telegram_id", telegramID), zap.Bool("is_linked", result.IsLinked)).Debug("Checked link status")
	return &result, nil
}

// GetUserPreferences retrieves user preferences from backend
func (c *Client) GetUserPreferences(ctx context.Context, token string) (*PreferenceResponse, error) {
	resp, err := c.doAuthenticatedRequest(ctx, "GET", "/api/v1/preferences", nil, token)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to get user preferences")
		return nil, err
	}

	var result PreferenceResponse
	if err := c.parseResponse(resp, &result); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse user preferences response")
		return nil, err
	}

	zap.L().With(zap.String("user_id", result.UserID)).Debug("Successfully retrieved user preferences")
	return &result, nil
}

// CreateUserPreferences creates user preferences in backend
func (c *Client) CreateUserPreferences(ctx context.Context, token string, userID string, preferences *CreatePreferenceRequest) (*PreferenceResponse, error) {
	resp, err := c.doAuthenticatedRequest(ctx, "POST", "/api/v1/preferences/user/"+userID, preferences, token)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to create user preferences")
		return nil, err
	}

	var result PreferenceResponse
	if err := c.parseResponse(resp, &result); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse create preferences response")
		return nil, err
	}

	zap.L().With(zap.String("user_id", result.UserID)).Info("Successfully created user preferences")
	return &result, nil
}

// UpdateUserPreferences updates user preferences in backend
func (c *Client) UpdateUserPreferences(ctx context.Context, token string, preferences *UpdatePreferenceRequest) (*PreferenceResponse, error) {
	resp, err := c.doAuthenticatedRequest(ctx, "PUT", "/api/v1/preferences", preferences, token)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to update user preferences")
		return nil, err
	}

	var result PreferenceResponse
	if err := c.parseResponse(resp, &result); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse update preferences response")
		return nil, err
	}

	zap.L().With(zap.String("user_id", result.UserID)).Info("Successfully updated user preferences")
	return &result, nil
}

// GenerateLesson generates a new lesson for the user with JWT authentication
func (c *Client) GenerateLesson(ctx context.Context, token string) (*domain.LessonResponse, error) {
	resp, err := c.doAuthenticatedRequest(ctx, "GET", "/api/v1/lesson", nil, token)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to generate lesson")
		return nil, err
	}

	var result domain.LessonResponse
	if err := c.parseResponse(resp, &result); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse generate lesson response")
		return nil, err
	}

	zap.L().With(
		zap.String("cefr_level", result.Lesson.CEFRLevel),
		zap.Int("word_count", len(result.Cards)),
		zap.Int("words_per_lesson", result.Lesson.WordsPerLesson),
	).Info("Successfully generated lesson")
	return &result, nil
}

// UpdateWordProgress updates word learning progress
func (c *Client) UpdateWordProgress(ctx context.Context, userID string, wordID uuid.UUID, correct bool, timeSpent int) error {
	req := WordProgressRequest{
		UserID:    userID,
		WordID:    wordID,
		Correct:   correct,
		TimeSpent: timeSpent,
	}

	resp, err := c.doRequest(ctx, "POST", "/api/v1/words/progress", req)
	if err != nil {
		zap.L().With(zap.String("user_id", userID), zap.String("word_id", wordID.String()), zap.Error(err)).Error("Failed to update word progress")
		return err
	}

	if err := c.parseResponse(resp, nil); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse update word progress response")
		return err
	}

	zap.L().With(
		zap.String("user_id", userID),
		zap.String("word_id", wordID.String()),
		zap.Bool("correct", correct),
	).Debug("Successfully updated word progress")
	return nil
}

// UpdateProgress updates overall learning progress
func (c *Client) UpdateProgress(ctx context.Context, req ProgressUpdateRequest) error {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/progress/update", req)
	if err != nil {
		zap.L().With(zap.String("user_id", req.UserID), zap.Error(err)).Error("Failed to update progress")
		return err
	}

	if err := c.parseResponse(resp, nil); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse update progress response")
		return err
	}

	zap.L().With(
		zap.String("user_id", req.UserID),
		zap.String("lesson_id", req.LessonID.String()),
		zap.Int("exercises_completed", req.ExercisesCompleted),
	).Info("Successfully updated progress")
	return nil
}

// GetUserStats retrieves user learning statistics
func (c *Client) GetUserStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/users/"+userID+"/stats", nil)
	if err != nil {
		zap.L().With(zap.String("user_id", userID), zap.Error(err)).Error("Failed to get user stats")
		return nil, err
	}

	var result map[string]interface{}
	if err := c.parseResponse(resp, &result); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse get user stats response")
		return nil, err
	}

	zap.L().With(zap.String("user_id", userID)).Debug("Successfully retrieved user stats")
	return result, nil
}

// GetJWTTokens retrieves JWT tokens for an authenticated user
func (c *Client) GetJWTTokens(ctx context.Context, telegramID int64) (*JWTResponse, error) {
	req := AuthRequest{TelegramID: telegramID}

	resp, err := c.doRequest(ctx, "POST", "/telegram/get-tokens", req)
	if err != nil {
		zap.L().With(zap.Int64("telegram_id", telegramID), zap.Error(err)).Error("Failed to get JWT tokens")
		return nil, err
	}

	var result JWTResponse
	if err := c.parseResponse(resp, &result); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse JWT tokens response")
		return nil, err
	}

	zap.L().With(zap.Int64("telegram_id", telegramID)).Debug("Successfully retrieved JWT tokens")
	return &result, nil
}

// SendLessonProgress sends word progress data to backend after lesson completion
func (c *Client) SendLessonProgress(ctx context.Context, token string, progressData []domain.WordProgress) error {
	// Convert to backend format
	var progressRequests []ProgressRequest
	for _, progress := range progressData {
		progressRequests = append(progressRequests, ProgressRequest{
			WordID:          progress.Word, // Using word as ID for now
			LearnedAt:       progress.LearnedAt,
			ConfidenceScore: progress.ConfidenceScore,
			CntReviewed:     progress.CntReviewed,
		})
	}

	resp, err := c.doAuthenticatedRequest(ctx, "POST", "/api/v1/progress", progressRequests, token)
	if err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to send lesson progress")
		return err
	}

	if err := c.parseResponse(resp, nil); err != nil {
		zap.L().With(zap.Error(err)).Error("Failed to parse send lesson progress response")
		return err
	}

	zap.L().With(zap.Int("words_count", len(progressData))).Info("Successfully sent lesson progress")
	return nil
}
