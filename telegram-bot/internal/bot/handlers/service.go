package handlers

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/config"
	"telegram-bot/internal/api"
	"telegram-bot/internal/bot/fsm"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/tasks"
)

// HandlerService provides common functionality for message handlers
type HandlerService struct {
	config       *config.Config
	redisClient  *redis.Client
	apiClient    *api.Client
	scheduler    *tasks.Scheduler
	bot          *tele.Bot
	stateManager *fsm.UserStateManager
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
	return &HandlerService{
		config:       cfg,
		redisClient:  redisClient,
		apiClient:    apiClient,
		scheduler:    scheduler,
		bot:          bot,
		stateManager: stateManager,
		logger:       logger,
	}
}

// TransitionState is a convenience method for transitioning user state
func (s *HandlerService) TransitionState(ctx context.Context, userID int64, newState fsm.UserState) error {
	s.logger.With(zap.Int64("user_id", userID), zap.String("new_state", string(newState))).Debug("Transitioning state")

	return s.stateManager.SetState(ctx, userID, newState)
}

// GetUserProgress (from API) for now - mock
func (s *HandlerService) GetUserProgress(ctx context.Context, userID int64) (*domain.UserProgress, error) {
	// For now, return a default user progress
	// This should be implemented to call the actual API
	return &domain.UserProgress{
		UserID:           userID,
		CEFRLevel:        "A1",
		WordsPerDay:      10,
		NotificationTime: "10:00",
	}, nil
}

// UpdateUserProgress update user progress through API
func (s *HandlerService) UpdateUserProgress(ctx context.Context, userID int64, progress *domain.UserProgress) error {
	// For now, this is a placeholder
	// This should be implemented to call the actual API
	return nil
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
	case fsm.StateQuestionSerials:
		return s.HandleQuestionSerialsMessage(ctx, c, userID, currentState)
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

	// Route based on callback data prefix and current state
	switch data {
	case "onboarding:start":
		return s.HandleOnboardingStartCallback(ctx, c, userID, currentState)
	case "onboarding:method":
		return s.HandleOnboardingMethodCallback(ctx, c, userID, currentState)
	case "onboarding:questionnaire":
		return s.HandleOnboardingQuestionnaireCallback(ctx, c, userID, currentState)
	case "test:start":
		return s.HandleTestStartCallback(ctx, c, userID, currentState)
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
