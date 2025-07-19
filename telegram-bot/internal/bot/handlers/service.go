package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/config"
	"telegram-bot/internal/api"
	"telegram-bot/internal/bot/fsm"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/tasks"
	"telegram-bot/internal/utils"
)

// HandlerService provides common functionality for message handlers
type HandlerService struct {
	config       *config.Config
	redisClient  *redis.Client
	apiClient    *api.Client
	scheduler    *tasks.Scheduler
	bot          *tele.Bot
	stateManager *fsm.UserStateManager
	ttsService   *utils.TTSService
	logger       *zap.Logger
}

// NewHandlerService creates a new handler service
func NewHandlerService(
	cfg *config.Config,
	redisClient *redis.Client,
	apiClient *api.Client,
	scheduler *tasks.Scheduler,
	bot *tele.Bot,
	stateManager *fsm.UserStateManager,
	logger *zap.Logger,
) *HandlerService {
	// Initialize TTS service
	ttsService := utils.NewTTSService("./tmp/tts", logger)

	return &HandlerService{
		config:       cfg,
		redisClient:  redisClient,
		apiClient:    apiClient,
		scheduler:    scheduler,
		bot:          bot,
		stateManager: stateManager,
		ttsService:   ttsService,
		logger:       logger,
	}
}

// TransitionState is a convenience method for transitioning user state
func (s *HandlerService) TransitionState(ctx context.Context, userID int64, newState fsm.UserState) error {
	s.logger.With(zap.Int64("user_id", userID), zap.String("new_state", string(newState))).Debug("Transitioning state")

	return s.stateManager.SetState(ctx, userID, newState)
}

// SetStateIfDifferent sets the user state only if it's different from the current state
// This prevents "invalid state transition from X to X" errors
func (s *HandlerService) SetStateIfDifferent(ctx context.Context, userID int64, newState fsm.UserState) error {
	currentState, err := s.stateManager.GetState(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get current state", zap.Int64("user_id", userID), zap.Error(err))
		return err
	}

	// Only set state if it's different from current state
	if currentState != newState {
		s.logger.With(zap.Int64("user_id", userID), zap.String("from_state", string(currentState)), zap.String("to_state", string(newState))).Debug("Transitioning state")
		return s.stateManager.SetState(ctx, userID, newState)
	}

	s.logger.With(zap.Int64("user_id", userID), zap.String("state", string(newState))).Debug("State already set, skipping transition")
	return nil
}

// GetUserProgress retrieves user progress from backend API
func (s *HandlerService) GetUserProgress(ctx context.Context, userID int64) (*domain.UserProgress, error) {
	// Check if user is authenticated
	if !s.stateManager.IsUserAuthenticated(ctx, userID) {
		// Return minimal user progress for unauthenticated users
		return &domain.UserProgress{
			UserID:           userID,
			CEFRLevel:        "",
			WordsPerDay:      10,
			NotificationTime: "10:00",
		}, nil
	}

	// Get access token
	accessToken, err := s.stateManager.GetValidAccessToken(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get access token", zap.Int64("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Get preferences from backend
	preferences, err := s.apiClient.GetUserPreferences(ctx, accessToken)
	if err != nil {
		s.logger.Error("Failed to get user preferences", zap.Int64("user_id", userID), zap.Error(err))

		// Check if this is a 404 error (preferences not found)
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "preference not found") {
			s.logger.Warn("User preferences not found, returning default preferences", zap.Int64("user_id", userID))

			// Return default preferences for authenticated users without saved preferences
			// This handles the case where a user is linked but preferences were never created
			return &domain.UserProgress{
				UserID:           userID,
				CEFRLevel:        "A1",    // Default level
				WordsPerDay:      10,      // Default words per day
				NotificationTime: "10:00", // Default notification time
				LearnedWords:     0,
				CurrentStreak:    0,
				LongestStreak:    0,
				LastActivity:     time.Now(),
				StartDate:        time.Now().Format("2006-01-02"),
				Preferences: map[string]interface{}{
					"goal":             "Learn new words",
					"fact_everyday":    false,
					"subscribed":       false,
					"avatar_image_url": "",
					"notifications":    false,
				},
			}, nil
		}

		// For other errors, still return the error
		return nil, err
	}

	// Convert to domain model
	userProgress := &domain.UserProgress{
		UserID:           userID,
		CEFRLevel:        preferences.CEFRLevel,
		WordsPerDay:      preferences.WordsPerDay,
		NotificationTime: "", // Will be set below after parsing
		LearnedWords:     0,  // TODO: Get from backend stats
		CurrentStreak:    0,  // TODO: Get from backend stats
		LongestStreak:    0,  // TODO: Get from backend stats
		LastActivity:     time.Now(),
		StartDate:        time.Now().Format("2006-01-02"),
		Preferences: map[string]interface{}{
			"goal":             preferences.Goal,
			"fact_everyday":    preferences.FactEveryday,
			"subscribed":       preferences.Subscribed,
			"avatar_image_url": preferences.AvatarImageURL,
			"notifications":    preferences.Notifications,
		},
	}

	// Parse notification time from ISO format to HH:MM format
	if preferences.NotificationAt != "" {
		s.logger.Debug("Parsing notification time from backend",
			zap.String("notification_at", preferences.NotificationAt))

		// Try to parse the ISO format time
		notificationTime, err := time.Parse(time.RFC3339, preferences.NotificationAt)
		if err != nil {
			// If RFC3339 fails, try other common formats
			notificationTime, err = time.Parse("2006-01-02T15:04:05Z", preferences.NotificationAt)
			if err != nil {
				// If that fails too, try parsing as HH:MM format (in case it's already in the right format)
				notificationTime, err = time.Parse("15:04", preferences.NotificationAt)
				if err != nil {
					s.logger.Warn("Failed to parse notification time from backend",
						zap.String("notification_at", preferences.NotificationAt),
						zap.Error(err))
					// Set default time if parsing fails
					userProgress.NotificationTime = "10:00"
				} else {
					// Already in HH:MM format
					userProgress.NotificationTime = preferences.NotificationAt
				}
			} else {
				// Successfully parsed ISO format, convert to HH:MM
				userProgress.NotificationTime = notificationTime.Format("15:04")
			}
		} else {
			// Successfully parsed RFC3339 format, convert to HH:MM
			userProgress.NotificationTime = notificationTime.Format("15:04")
		}

		s.logger.Debug("Successfully converted notification time",
			zap.String("from", preferences.NotificationAt),
			zap.String("to", userProgress.NotificationTime))
	} else {
		// No notification time set, use default
		userProgress.NotificationTime = "10:00"
		s.logger.Debug("No notification time from backend, using default",
			zap.String("default_time", userProgress.NotificationTime))
	}

	return userProgress, nil
}

// UpdateUserProgress updates user progress through backend API
func (s *HandlerService) UpdateUserProgress(ctx context.Context, userID int64, progress *domain.UserProgress) error {
	// Check if user is authenticated
	if !s.stateManager.IsUserAuthenticated(ctx, userID) {
		return fmt.Errorf("user %d is not authenticated", userID)
	}

	// Get access token
	accessToken, err := s.stateManager.GetValidAccessToken(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get access token", zap.Int64("user_id", userID), zap.Error(err))
		return err
	}

	// Преобразуем в формат backend, только если значения существуют (не nil)
	updateRequest := &api.UpdatePreferenceRequest{}

	if progress.CEFRLevel != "" {
		updateRequest.CEFRLevel = &progress.CEFRLevel
	}
	if progress.WordsPerDay != 0 {
		updateRequest.WordsPerDay = &progress.WordsPerDay
	}
	if progress.NotificationTime != "" {
		s.logger.Debug("Parsing notification time for update",
			zap.String("notification_time", progress.NotificationTime))

		notificationTime, err := ParseTimeToTime(progress.NotificationTime)
		if err != nil {
			s.logger.Error("Failed to parse notification time", zap.String("notification_time", progress.NotificationTime), zap.Error(err))
			return err
		}
		updateRequest.NotificationAt = notificationTime
		s.logger.Debug("Successfully parsed notification time for update",
			zap.String("parsed_time", notificationTime.Format("15:04")))
	} else {
		// If notification time is empty, set notifications to false
		notifications := false
		updateRequest.Notifications = &notifications
		s.logger.Debug("Notification time is empty, setting notifications to false")
	}

	// Extract additional preferences
	if progress.Preferences != nil {
		if goal, ok := progress.Preferences["goal"].(string); ok {
			updateRequest.Goal = &goal
		}
		if factEveryday, ok := progress.Preferences["fact_everyday"].(bool); ok {
			updateRequest.FactEveryday = &factEveryday
		}
		if subscribed, ok := progress.Preferences["subscribed"].(bool); ok {
			updateRequest.Subscribed = &subscribed
		}
		if avatarURL, ok := progress.Preferences["avatar_image_url"].(string); ok {
			updateRequest.AvatarImageURL = &avatarURL
		}
	}

	// Update preferences in backend
	s.logger.Info("Updating user preferences", zap.Any("update_request", updateRequest))
	_, err = s.apiClient.UpdateUserPreferences(ctx, accessToken, updateRequest)
	if err != nil {
		s.logger.Error("Failed to update user preferences", zap.Int64("user_id", userID), zap.Error(err))

		// Check if this is a 404 error (preferences not found)
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "preference not found") {
			s.logger.Warn("User preferences not found, creating new preferences", zap.Int64("user_id", userID))

			// Get user ID from JWT token
			userIDFromToken, err := s.getUserIDFromToken(accessToken)
			if err != nil {
				s.logger.Error("Failed to get user ID from token", zap.Int64("user_id", userID), zap.Error(err))
				// Continue without backend sync
				return nil
			}

			// Try to create preferences instead
			createRequest := &api.CreatePreferenceRequest{
				CEFRLevel:      progress.CEFRLevel,
				WordsPerDay:    progress.WordsPerDay,
				Notifications:  progress.NotificationsEnabled(),
				Goal:           "Learn new words",
				FactEveryday:   false,
				Subscribed:     false,
				AvatarImageURL: "",
			}

			if progress.NotificationTime != "" {
				notificationTime, err := ParseTimeToTime(progress.NotificationTime)
				if err == nil {
					createRequest.NotificationAt = notificationTime
				} else {
					s.logger.Warn("Failed to parse notification time for create request",
						zap.String("notification_time", progress.NotificationTime),
						zap.Error(err))
					// Set notifications to false if time parsing fails
					createRequest.Notifications = false
				}
			} else {
				// If notification time is empty, set notifications to false
				createRequest.Notifications = false
			}

			// Extract additional preferences
			if progress.Preferences != nil {
				if goal, ok := progress.Preferences["goal"].(string); ok && goal != "" {
					createRequest.Goal = goal
				}
				if factEveryday, ok := progress.Preferences["fact_everyday"].(bool); ok {
					createRequest.FactEveryday = factEveryday
				}
				if subscribed, ok := progress.Preferences["subscribed"].(bool); ok {
					createRequest.Subscribed = subscribed
				}
				if avatarURL, ok := progress.Preferences["avatar_image_url"].(string); ok && avatarURL != "" {
					createRequest.AvatarImageURL = avatarURL
				}
			}

			// Create preferences
			if _, createErr := s.apiClient.CreateUserPreferences(ctx, accessToken, userIDFromToken, createRequest); createErr != nil {
				s.logger.Error("Failed to create user preferences", zap.Int64("user_id", userID), zap.Error(createErr))
				// Continue without backend sync
			} else {
				s.logger.Info("Successfully created user preferences", zap.Int64("user_id", userID))
			}
		} else {
			// For other errors, still return the error
			return err
		}
	}

	s.logger.Info("Successfully updated user progress", zap.Int64("user_id", userID))
	return nil
}

// getUserIDFromToken extracts user ID from JWT token
func (s *HandlerService) getUserIDFromToken(tokenString string) (string, error) {
	// Split the token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid JWT format")
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// Parse the JSON payload
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	// Extract user ID from 'sub' claim
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in JWT claims")
	}

	return userID, nil
}

// IsUserAuthenticated checks if user has valid authentication
func (s *HandlerService) IsUserAuthenticated(ctx context.Context, userID int64) bool {
	return s.stateManager.IsUserAuthenticated(ctx, userID)
}

// HasUserCompletedOnboarding checks if user has completed the onboarding process
func (s *HandlerService) HasUserCompletedOnboarding(ctx context.Context, userID int64) (bool, error) {
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		return false, err
	}

	// User has completed onboarding if they have:
	// 1. CEFR level set
	// 2. Words per day preference set (non-zero)
	// 3. Notification preference set (either enabled or disabled)
	return userProgress.CEFRLevel != "" &&
		userProgress.WordsPerDay > 0 &&
		userProgress.NotificationTime != "", nil
}

// GetUserAuthenticationStatus returns detailed authentication status
func (s *HandlerService) GetUserAuthenticationStatus(ctx context.Context, userID int64) (bool, bool, error) {
	isAuthenticated := s.IsUserAuthenticated(ctx, userID)
	s.logger.Debug("User authentication check", zap.Int64("user_id", userID), zap.Bool("is_authenticated", isAuthenticated))

	hasCompletedOnboarding := false

	if isAuthenticated {
		completed, err := s.HasUserCompletedOnboarding(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to check onboarding completion", zap.Int64("user_id", userID), zap.Error(err))
			return isAuthenticated, false, err
		}
		hasCompletedOnboarding = completed
		s.logger.Debug("User onboarding check", zap.Int64("user_id", userID), zap.Bool("has_completed_onboarding", hasCompletedOnboarding))
	}

	s.logger.Info("User authentication status", zap.Int64("user_id", userID), zap.Bool("is_authenticated", isAuthenticated), zap.Bool("has_completed_onboarding", hasCompletedOnboarding))
	return isAuthenticated, hasCompletedOnboarding, nil
}

// HandleTextMessage handles text messages based on state
func (s *HandlerService) HandleTextMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	text := c.Text()

	s.logger.With(zap.Int64("user_id", userID), zap.String("state", string(currentState)), zap.String("text", text)).Debug("Processing text message")

	// Route based on current state
	switch currentState {
	case fsm.StateWelcome:
		return s.HandleWelcomeMessage(ctx, c, userID, currentState)
	case fsm.StateMethodExplanation:
		return s.HandleMethodExplanationMessage(ctx, c, userID, currentState)
	case fsm.StateQuestionGoal:
		return s.HandleQuestionGoalMessage(ctx, c, userID, currentState)
	case fsm.StateQuestionConfidence:
		return s.HandleQuestionConfidenceMessage(ctx, c, userID, currentState)
	case fsm.StateQuestionExperience:
		return s.HandleQuestionExperienceMessage(ctx, c, userID, currentState)
	case fsm.StateQuestionWordsPerDay:
		return s.HandleQuestionWordsPerDayMessage(ctx, c, userID, currentState)
	case fsm.StateQuestionNotifications:
		return s.HandleQuestionNotificationsMessage(ctx, c, userID, currentState)
	case fsm.StateQuestionNotificationTime:
		return s.HandleQuestionNotificationTimeMessage(ctx, c, userID, currentState)
	case fsm.StateSettingsWordsPerDayInput:
		return s.HandleSettingsWordsPerDayInputMessage(ctx, c, userID, currentState)
	case fsm.StateSettingsTimeInput:
		return s.HandleSettingsTimeInputMessage(ctx, c, userID, currentState)
	case fsm.StateSettingsCEFRLevel:
		return s.HandleSettingsCEFRLevelInputMessage(ctx, c, userID, currentState)
	case fsm.StateWaitingForTranslation:
		return s.HandleWaitingForTranslationMessage(ctx, c, userID, currentState)
	case fsm.StateWaitingForAudio:
		return s.HandleWaitingForAudioMessage(ctx, c, userID, currentState)
	case fsm.StateWaitingForTextInput:
		// New learning flow: handle exercise text input
		return s.HandleTextInputAnswer(ctx, c, userID, text)
	default:
		return s.HandleUnknownStateMessage(ctx, c, userID, currentState)
	}
}

// HandleCallback handles callback queries based on state
func (s *HandlerService) HandleCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	callback := c.Callback()
	if callback == nil {
		return fmt.Errorf("no callback data")
	}

	data := callback.Data

	s.logger.With(zap.Int64("user_id", userID), zap.String("state", string(currentState)), zap.String("data", data)).Debug("Processing callback")

	// Always respond to callback to remove loading state
	defer func() {
		if err := c.Respond(); err != nil {
			s.logger.Error("Failed to respond to callback", zap.Error(err))
		}
	}()

	// Check for new learning flow callbacks
	if strings.HasPrefix(data, "lesson:") {
		action := strings.TrimPrefix(data, "lesson:")
		return s.HandleLessonCallback(ctx, c, userID, action)
	}

	if strings.HasPrefix(data, "exercise:") {
		action := strings.TrimPrefix(data, "exercise:")
		return s.HandleExerciseCallback(ctx, c, userID, action)
	}

	if strings.HasPrefix(data, "auth:") {
		action := strings.TrimPrefix(data, "auth:")
		return s.HandleAuthCallback(ctx, c, userID, action)
	}

	if strings.HasPrefix(data, "voice:") {
		action := strings.TrimPrefix(data, "voice:")
		return s.HandleVoiceCallback(ctx, c, userID, action)
	}

	if strings.HasPrefix(data, "stats:") {
		action := strings.TrimPrefix(data, "stats:")
		return s.HandleStatsCallback(ctx, c, userID, action)
	}

	if strings.HasPrefix(data, "help:") {
		action := strings.TrimPrefix(data, "help:")
		if action == "auth" {
			return s.handleHelpAuth(ctx, c, userID)
		}
	}

	// Check for test answer callbacks
	if strings.HasPrefix(data, "test_answer:") {
		return s.handleTestAnswerCallback(ctx, c, userID, data)
	}

	// Check for test "don't know" callbacks
	if strings.HasPrefix(data, "test_dont_know:") {
		return s.handleTestDontKnowCallback(ctx, c, userID, data)
	}

	// Check for test fix level callbacks
	if strings.HasPrefix(data, "test_fix_level:") {
		level := strings.TrimPrefix(data, "test_fix_level:")
		return s.HandleTestFixLevelCallback(ctx, c, userID, level)
	}

	// Check for test continue callbacks
	if strings.HasPrefix(data, "test_continue:") {
		return s.handleTestContinueCallback(ctx, c, userID, data)
	}

	// Check for test continue next callbacks (after answer feedback)
	if strings.HasPrefix(data, "test_continue_next:") {
		return s.handleTestContinueNextCallback(ctx, c, userID, data)
	}

	// Check for questionnaire answer callbacks first
	if strings.HasPrefix(data, "goal:") {
		answer := strings.TrimPrefix(data, "goal:")
		return s.HandleGoalCallback(ctx, c, userID, answer)
	}
	if strings.HasPrefix(data, "confidence:") {
		answer := strings.TrimPrefix(data, "confidence:")
		return s.HandleConfidenceCallback(ctx, c, userID, answer)
	}
	if strings.HasPrefix(data, "experience:") {
		answer := strings.TrimPrefix(data, "experience:")
		return s.HandleExperienceCallback(ctx, c, userID, answer)
	}
	if strings.HasPrefix(data, "words_per_day:") {
		answer := strings.TrimPrefix(data, "words_per_day:")
		return s.HandleWordsPerDayCallback(ctx, c, userID, answer)
	}
	if strings.HasPrefix(data, "notifications:") {
		answer := strings.TrimPrefix(data, "notifications:")
		return s.HandleNotificationsCallback(ctx, c, userID, answer)
	}
	if strings.HasPrefix(data, "notification_time:") {
		answer := strings.TrimPrefix(data, "notification_time:")
		return s.HandleNotificationTimeCallback(ctx, c, userID, answer)
	}

	// Handle settings callbacks
	if strings.HasPrefix(data, "settings:words:") {
		return s.HandleSettingsWordsCallback(ctx, c, userID, data)
	}
	if strings.HasPrefix(data, "settings:time:") {
		return s.HandleSettingsTimeCallback(ctx, c, userID, data)
	}
	if strings.HasPrefix(data, "settings:cefr:") {
		return s.HandleSettingsCEFRCallback(ctx, c, userID, data)
	}
	if data == "settings:back" {
		return s.HandleSettingsBackCallback(ctx, c, userID, currentState)
	}

	// Route based on callback data prefix and current state
	switch data {
	case "onboarding:start":
		return s.HandleOnboardingStartCallback(ctx, c, userID, currentState)
	case "onboarding:method":
		return s.HandleOnboardingMethodCallback(ctx, c, userID, currentState)
	case "onboarding:questionnaire":
		return s.HandleOnboardingQuestionnaireCallback(ctx, c, userID, currentState)
	case "questionnaire:start":
		return s.HandleQuestionnaireStartCallback(ctx, c, userID, currentState)
	case "test:start":
		return s.HandleTestStartCallback(ctx, c, userID, currentState)
	case "test:skip":
		return s.HandleTestSkipCallback(ctx, c, userID, currentState)
	case "lesson:start":
		return s.HandleLessonStartCallback(ctx, c, userID, currentState)
	case "lesson:later":
		return s.HandleLessonLaterCallback(ctx, c, userID, currentState)
	case "settings:words_per_day":
		return s.HandleSettingsWordsPerDayCallback(ctx, c, userID, currentState)
	case "settings:notifications":
		return s.HandleSettingsNotificationsCallback(ctx, c, userID, currentState)
	case "settings:cefr_level":
		return s.HandleSettingsCEFRLevelCallback(ctx, c, userID, currentState)
	case "menu:main":
		return s.HandleMainMenuCallback(ctx, c, userID, currentState)
	case "menu:settings":
		return s.HandleSettingsMenuCallback(ctx, c, userID, currentState)
	case "menu:learn":
		return s.HandleLearnMenuCallback(ctx, c, userID, currentState)
	case "menu:help":
		return s.HandleHelpMenuCallback(ctx, c, userID, currentState)
	case "account:link":
		return s.HandleAccountLinkCallback(ctx, c, userID, currentState)
	default:
		return s.HandleUnknownCallback(ctx, c, userID, currentState)
	}
}

// handleTestAnswerCallback parses and handles test answer callbacks
func (s *HandlerService) handleTestAnswerCallback(ctx context.Context, c tele.Context, userID int64, callbackData string) error {
	// Parse callback data: test_answer:group:question:answer
	parts := strings.Split(callbackData, ":")
	if len(parts) != 4 {
		s.logger.Error("Invalid test answer callback format", zap.String("data", callbackData))
		return fmt.Errorf("invalid callback format")
	}

	group, err := strconv.Atoi(parts[1])
	if err != nil {
		s.logger.Error("Invalid group number", zap.String("group", parts[1]), zap.Error(err))
		return fmt.Errorf("invalid group number")
	}

	questionIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		s.logger.Error("Invalid question index", zap.String("question", parts[2]), zap.Error(err))
		return fmt.Errorf("invalid question index")
	}

	answerIndex, err := strconv.Atoi(parts[3])
	if err != nil {
		s.logger.Error("Invalid answer index", zap.String("answer", parts[3]), zap.Error(err))
		return fmt.Errorf("invalid answer index")
	}

	return s.HandleTestAnswerCallback(ctx, c, userID, group, questionIndex, answerIndex)
}

// handleTestDontKnowCallback parses and handles test "don't know" callbacks
func (s *HandlerService) handleTestDontKnowCallback(ctx context.Context, c tele.Context, userID int64, callbackData string) error {
	// Parse callback data: test_dont_know:group:question
	parts := strings.Split(callbackData, ":")
	if len(parts) != 3 {
		s.logger.Error("Invalid test dont know callback format", zap.String("data", callbackData))
		return fmt.Errorf("invalid callback format")
	}

	group, err := strconv.Atoi(parts[1])
	if err != nil {
		s.logger.Error("Invalid group number", zap.String("group", parts[1]), zap.Error(err))
		return fmt.Errorf("invalid group number")
	}

	questionIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		s.logger.Error("Invalid question index", zap.String("question", parts[2]), zap.Error(err))
		return fmt.Errorf("invalid question index")
	}

	return s.HandleTestDontKnowCallback(ctx, c, userID, group, questionIndex)
}

// handleTestContinueCallback parses and handles test continue callbacks
func (s *HandlerService) handleTestContinueCallback(ctx context.Context, c tele.Context, userID int64, callbackData string) error {
	// Parse callback data: test_continue:group
	parts := strings.Split(callbackData, ":")
	if len(parts) != 2 {
		s.logger.Error("Invalid test continue callback format", zap.String("data", callbackData))
		return fmt.Errorf("invalid callback format")
	}

	group, err := strconv.Atoi(parts[1])
	if err != nil {
		s.logger.Error("Invalid group number", zap.String("group", parts[1]), zap.Error(err))
		return fmt.Errorf("invalid group number")
	}

	return s.HandleTestContinueCallback(ctx, c, userID, group)
}

// handleTestContinueNextCallback parses and handles test continue next callbacks
func (s *HandlerService) handleTestContinueNextCallback(ctx context.Context, c tele.Context, userID int64, callbackData string) error {
	// Parse callback data: test_continue_next:group:nextQuestionIndex
	parts := strings.Split(callbackData, ":")
	if len(parts) != 3 {
		s.logger.Error("Invalid test continue next callback format", zap.String("data", callbackData))
		return fmt.Errorf("invalid callback format")
	}

	group, err := strconv.Atoi(parts[1])
	if err != nil {
		s.logger.Error("Invalid group number", zap.String("group", parts[1]), zap.Error(err))
		return fmt.Errorf("invalid group number")
	}

	nextQuestionIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		s.logger.Error("Invalid next question index", zap.String("nextQuestion", parts[2]), zap.Error(err))
		return fmt.Errorf("invalid next question index")
	}

	return s.HandleTestContinueNextCallback(ctx, c, userID, group, nextQuestionIndex)
}

// HandleVoiceMessage handles voice messages
func (s *HandlerService) HandleVoiceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	voice := c.Message().Voice
	if voice == nil {
		return fmt.Errorf("no voice message")
	}

	s.logger.With(zap.Int64("user_id", userID), zap.String("state", string(currentState)), zap.Int("duration", voice.Duration)).Debug("Processing voice message")

	// Handle voice based on state
	if currentState == fsm.StateWaitingForAudio {
		return s.HandleAudioExerciseResponse(ctx, c, userID, voice)
	}

	// For other states, provide guidance
	return c.Send("Голосовые сообщения поддерживаются во время аудио упражнений. Используйте /learn чтобы начать урок.")
}

// HandleAudioMessage handles audio messages
func (s *HandlerService) HandleAudioMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	audio := c.Message().Audio
	if audio == nil {
		return fmt.Errorf("no audio message")
	}

	s.logger.With(zap.Int64("user_id", userID), zap.String("state", string(currentState)), zap.Int("duration", audio.Duration)).Debug("Processing audio message")

	// Similar to voice handling
	if currentState == fsm.StateWaitingForAudio {
		return s.HandleAudioExerciseResponse(ctx, c, userID, audio)
	}

	return c.Send("Аудио сообщения поддерживаются во время аудио упражнений. Используйте /learn чтобы начать урок.")
}

// HandlePhotoMessage handles photo messages
func (s *HandlerService) HandlePhotoMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	photo := c.Message().Photo
	if photo == nil {
		return fmt.Errorf("no photo message")
	}

	s.logger.With(zap.Int64("user_id", userID), zap.String("state", string(currentState))).Debug("Processing photo message")

	// For now, photos aren't part of the learning flow
	return c.Send("Используйте /help чтобы увидеть доступные команды.")
}
