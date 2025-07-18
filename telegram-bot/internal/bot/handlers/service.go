package handlers

import (
	"context"
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
		return nil, err
	}

	// Convert to domain model
	userProgress := &domain.UserProgress{
		UserID:           userID,
		CEFRLevel:        preferences.CEFRLevel,
		WordsPerDay:      preferences.WordsPerDay,
		NotificationTime: preferences.NotificationAt,
		LearnedWords:     0, // TODO: Get from backend stats
		CurrentStreak:    0, // TODO: Get from backend stats
		LongestStreak:    0, // TODO: Get from backend stats
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

	// Convert to backend format
	updateRequest := &api.UpdatePreferenceRequest{
		CEFRLevel:      progress.CEFRLevel,
		WordsPerDay:    progress.WordsPerDay,
		NotificationAt: progress.NotificationTime,
		Notifications:  progress.NotificationsEnabled(),
	}

	// Extract additional preferences
	if progress.Preferences != nil {
		if goal, ok := progress.Preferences["goal"].(string); ok {
			updateRequest.Goal = goal
		}
		if factEveryday, ok := progress.Preferences["fact_everyday"].(bool); ok {
			updateRequest.FactEveryday = factEveryday
		}
		if subscribed, ok := progress.Preferences["subscribed"].(bool); ok {
			updateRequest.Subscribed = subscribed
		}
		if avatarURL, ok := progress.Preferences["avatar_image_url"].(string); ok {
			updateRequest.AvatarImageURL = avatarURL
		}
	}

	// Update preferences in backend
	_, err = s.apiClient.UpdateUserPreferences(ctx, accessToken, updateRequest)
	if err != nil {
		s.logger.Error("Failed to update user preferences", zap.Int64("user_id", userID), zap.Error(err))
		return err
	}

	s.logger.Info("Successfully updated user progress", zap.Int64("user_id", userID))
	return nil
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
	case fsm.StateSettingsWordsPerDayInput:
		return s.HandleSettingsWordsPerDayInputMessage(ctx, c, userID, currentState)
	case fsm.StateSettingsTimeInput:
		return s.HandleSettingsTimeInputMessage(ctx, c, userID, currentState)
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
