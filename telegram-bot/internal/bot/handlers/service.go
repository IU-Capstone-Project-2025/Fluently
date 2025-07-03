package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"fluently/telegram-bot/config"
	"fluently/telegram-bot/internal/api"
	"fluently/telegram-bot/internal/bot/fsm"
	"fluently/telegram-bot/internal/tasks"
	"fluently/telegram-bot/pkg/logger"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// HandlerService manages all bot handlers and state
type HandlerService struct {
	config      *config.Config
	redisClient *redis.Client
	apiClient   *api.Client
	scheduler   *tasks.Scheduler
	bot         *tele.Bot
}

// NewHandlerService creates a new handler service
func NewHandlerService(cfg *config.Config, redisClient *redis.Client, apiClient *api.Client, scheduler *tasks.Scheduler, bot *tele.Bot) *HandlerService {
	return &HandlerService{
		config:      cfg,
		redisClient: redisClient,
		apiClient:   apiClient,
		scheduler:   scheduler,
		bot:         bot,
	}
}

// ProcessMessage processes incoming messages based on user state
func (h *HandlerService) ProcessMessage(ctx context.Context, message *tele.Message) {
	if message.Sender == nil {
		logger.Log.Warn("Received message without sender")
		return
	}

	userID := message.Sender.ID
	logger.Log.Debug("Processing message",
		zap.Int64("user_id", userID),
		zap.String("text", message.Text))

	// Get or create user progress
	progress, err := h.getUserProgress(ctx, userID)
	if err != nil {
		logger.Log.Error("Failed to get user progress", zap.Error(err), zap.Int64("user_id", userID))
		h.sendErrorMessage(message.Chat, "Произошла ошибка. Попробуйте снова.")
		return
	}

	// Update last activity
	progress.LastActivity = time.Now()

	// Handle commands
	if message.Text != "" && strings.HasPrefix(message.Text, "/") {
		h.handleCommand(ctx, message, progress)
		return
	}

	// Handle message based on current state
	switch progress.State {
	case fsm.StateStart:
		h.handleStart(ctx, message, progress)
	case fsm.StateWelcome:
		h.handleWelcome(ctx, message, progress)
	case fsm.StateMethodExplanation:
		h.handleMethodExplanation(ctx, message, progress)
	case fsm.StateSpacedRepetition:
		h.handleSpacedRepetition(ctx, message, progress)
	case fsm.StateQuestionnaire:
		h.handleQuestionnaire(ctx, message, progress)
	case fsm.StateQuestionGoal:
		h.handleQuestionGoal(ctx, message, progress)
	case fsm.StateQuestionConfidence:
		h.handleQuestionConfidence(ctx, message, progress)
	case fsm.StateQuestionSerials:
		h.handleQuestionSerials(ctx, message, progress)
	case fsm.StateQuestionExperience:
		h.handleQuestionExperience(ctx, message, progress)
	case fsm.StateVocabularyTest:
		h.handleVocabularyTest(ctx, message, progress)
	case fsm.StateTestGroup1, fsm.StateTestGroup2, fsm.StateTestGroup3, fsm.StateTestGroup4, fsm.StateTestGroup5:
		h.handleVocabularyTestGroup(ctx, message, progress)
	case fsm.StateLevelDetermination:
		h.handleLevelDetermination(ctx, message, progress)
	case fsm.StatePlanCreation:
		h.handlePlanCreation(ctx, message, progress)
	case fsm.StateLessonStart:
		h.handleLessonStart(ctx, message, progress)
	case fsm.StateShowingFirstBlock:
		h.handleShowingFirstBlock(ctx, message, progress)
	case fsm.StateExerciseAfterBlock:
		h.handleExerciseAfterBlock(ctx, message, progress)
	case fsm.StateShowingIndividualWord:
		h.handleShowingIndividualWord(ctx, message, progress)
	case fsm.StateExerciseReview:
		h.handleExerciseReview(ctx, message, progress)
	case fsm.StateAudioDictation:
		h.handleAudioDictation(ctx, message, progress)
	case fsm.StateTranslationCheck:
		h.handleTranslationCheck(ctx, message, progress)
	case fsm.StateWaitingForAudio:
		h.handleWaitingForAudio(ctx, message, progress)
	case fsm.StateWaitingForTranslation:
		h.handleWaitingForTranslation(ctx, message, progress)
	case fsm.StateLessonComplete:
		h.handleLessonComplete(ctx, message, progress)
	case fsm.StateAccountLinking:
		h.handleAccountLinking(ctx, message, progress)
	case fsm.StateWaitingForLink:
		h.handleWaitingForLink(ctx, message, progress)
	case fsm.StateSettings:
		h.handleSettings(ctx, message, progress)
	default:
		logger.Log.Warn("Unhandled state", zap.String("state", string(progress.State)), zap.Int64("user_id", userID))
		h.sendMessage(message.Chat, "Извините, произошла ошибка. Введите /start для начала.")
	}

	// Save updated progress
	if err := h.saveUserProgress(ctx, progress); err != nil {
		logger.Log.Error("Failed to save user progress", zap.Error(err), zap.Int64("user_id", userID))
	}
}

// ProcessCallback processes callback queries
func (h *HandlerService) ProcessCallback(ctx context.Context, callback *tele.Callback) {
	if callback.Sender == nil {
		logger.Log.Warn("Received callback without sender")
		return
	}

	userID := callback.Sender.ID
	logger.Log.Debug("Processing callback",
		zap.Int64("user_id", userID),
		zap.String("data", callback.Data))

	// Get user progress
	progress, err := h.getUserProgress(ctx, userID)
	if err != nil {
		logger.Log.Error("Failed to get user progress", zap.Error(err), zap.Int64("user_id", userID))
		h.bot.Respond(callback, &tele.CallbackResponse{Text: "Произошла ошибка"})
		return
	}

	// Handle callback based on data
	switch {
	case strings.HasPrefix(callback.Data, "start_lesson"):
		h.handleStartLessonCallback(ctx, callback, progress)
	case strings.HasPrefix(callback.Data, "vocab_test_"):
		h.handleVocabTestCallback(ctx, callback, progress)
	case strings.HasPrefix(callback.Data, "exercise_"):
		h.handleExerciseCallback(ctx, callback, progress)
	case strings.HasPrefix(callback.Data, "link_account"):
		h.handleLinkAccountCallback(ctx, callback, progress)
	case strings.HasPrefix(callback.Data, "settings_"):
		h.handleSettingsCallback(ctx, callback, progress)
	default:
		logger.Log.Warn("Unhandled callback", zap.String("data", callback.Data))
		h.bot.Respond(callback, &tele.CallbackResponse{Text: "Неизвестная команда"})
	}

	// Save updated progress
	if err := h.saveUserProgress(ctx, progress); err != nil {
		logger.Log.Error("Failed to save user progress", zap.Error(err), zap.Int64("user_id", userID))
	}
}

// ProcessInlineQuery processes inline queries
func (h *HandlerService) ProcessInlineQuery(ctx context.Context, query *tele.Query) {
	// Handle inline queries if needed
	logger.Log.Debug("Processing inline query", zap.String("query", query.Text))
}

// handleCommand handles bot commands
func (h *HandlerService) handleCommand(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	command := strings.ToLower(strings.TrimPrefix(message.Text, "/"))

	switch command {
	case "start":
		h.handleStart(ctx, message, progress)
	case "help":
		h.handleHelp(ctx, message, progress)
	case "stats":
		h.handleStats(ctx, message, progress)
	case "settings":
		h.handleSettingsCommand(ctx, message, progress)
	case "lesson":
		h.handleLessonCommand(ctx, message, progress)
	case "cancel":
		h.handleCancel(ctx, message, progress)
	default:
		h.sendMessage(message.Chat, "Неизвестная команда. Введите /help для получения справки.")
	}
}

// State handlers
func (h *HandlerService) handleStart(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	if progress.State != fsm.StateStart {
		progress.UpdateState(fsm.StateStart)
	}

	// Check if user is already linked
	linkStatus, err := h.apiClient.CheckLinkStatus(ctx, progress.TelegramID)
	if err == nil && linkStatus.IsLinked {
		progress.GoogleLinked = true
		h.sendWelcomeMessage(ctx, message, progress)
		return
	}

	// Send welcome message with link option
	text := `🙂 Привет! Это Fluently!

Я помогу тебе улучшить словарный запас английского языка всего за 10 минут в день.

🎯 Что делает этот бот особенным:
• Фокус на самых важных словах (80-90% речи)
• Особая техника запоминания
• Всего 10 слов в день
• Научно обоснованный метод

Для полного доступа к функциям рекомендую связать аккаунт с Google.`

	keyboard := &tele.ReplyMarkup{}
	linkBtn := keyboard.Data("🔗 Связать с Google", "link_account")
	continueBtn := keyboard.Data("▶️ Продолжить без связывания", "continue_without_link")
	keyboard.Inline(
		keyboard.Row(linkBtn),
		keyboard.Row(continueBtn),
	)

	h.sendMessageWithKeyboard(message.Chat, text, keyboard)
	progress.UpdateState(fsm.StateAccountLinking)
}

func (h *HandlerService) handleWelcome(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	h.sendWelcomeMessage(ctx, message, progress)
}

func (h *HandlerService) sendWelcomeMessage(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	text := `✨ Добро пожаловать в Fluently!

🎯 Наша цель: увеличить твой словарный запас на 1000+ слов за 3 месяца.

Как это работает:
• Не нужно учить весь словарь
• Фокус на самых частых словах
• Особая техника запоминания
• 10 минут в день = результат

Слова взяты из авторитетных словарей Oxford, Macmillan, Longman с профессиональным переводом.`

	keyboard := &tele.ReplyMarkup{}
	nextBtn := keyboard.Data("Далее ➡️", "method_explanation")
	keyboard.Inline(keyboard.Row(nextBtn))

	h.sendMessageWithKeyboard(message.Chat, text, keyboard)
	progress.UpdateState(fsm.StateMethodExplanation)
}

// Helper methods
func (h *HandlerService) getUserProgress(ctx context.Context, userID int64) (*fsm.UserProgress, error) {
	key := fmt.Sprintf("user_progress:%d", userID)

	data, err := h.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		// Create new user progress
		progress := fsm.CreateNewUserProgress(userID)
		return progress, nil
	} else if err != nil {
		return nil, err
	}

	// Parse existing progress
	progress, err := fsm.FromJSON([]byte(data))
	if err != nil {
		// If parsing fails, create new progress
		logger.Log.Warn("Failed to parse user progress, creating new", zap.Error(err), zap.Int64("user_id", userID))
		progress = fsm.CreateNewUserProgress(userID)
	}

	return progress, nil
}

func (h *HandlerService) saveUserProgress(ctx context.Context, progress *fsm.UserProgress) error {
	key := fmt.Sprintf("user_progress:%d", progress.TelegramID)

	data, err := progress.ToJSON()
	if err != nil {
		return err
	}

	return h.redisClient.Set(ctx, key, data, 24*time.Hour).Err()
}

func (h *HandlerService) sendMessage(chat *tele.Chat, text string) {
	if _, err := h.bot.Send(chat, text, &tele.SendOptions{ParseMode: tele.ModeHTML}); err != nil {
		logger.Log.Error("Failed to send message", zap.Error(err))
	}
}

func (h *HandlerService) sendMessageWithKeyboard(chat *tele.Chat, text string, keyboard *tele.ReplyMarkup) {
	if _, err := h.bot.Send(chat, text, &tele.SendOptions{
		ParseMode:   tele.ModeHTML,
		ReplyMarkup: keyboard,
	}); err != nil {
		logger.Log.Error("Failed to send message with keyboard", zap.Error(err))
	}
}

func (h *HandlerService) sendErrorMessage(chat *tele.Chat, text string) {
	h.sendMessage(chat, "❌ "+text)
}

// Health check methods
func (h *HandlerService) CheckRedisHealth() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.redisClient.Ping(ctx).Err() == nil
}

func (h *HandlerService) CheckAPIHealth() bool {
	// Implement API health check
	return true
}

// Placeholder handlers for other states (to be implemented)
func (h *HandlerService) handleMethodExplanation(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for method explanation
}

func (h *HandlerService) handleSpacedRepetition(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for spaced repetition explanation
}

func (h *HandlerService) handleQuestionnaire(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for questionnaire
}

func (h *HandlerService) handleQuestionGoal(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for goal question
}

func (h *HandlerService) handleQuestionConfidence(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for confidence question
}

func (h *HandlerService) handleQuestionSerials(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for serials question
}

func (h *HandlerService) handleQuestionExperience(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for experience question
}

func (h *HandlerService) handleVocabularyTest(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for vocabulary test
}

func (h *HandlerService) handleVocabularyTestGroup(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for vocabulary test groups
}

func (h *HandlerService) handleLevelDetermination(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for level determination
}

func (h *HandlerService) handlePlanCreation(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for plan creation
}

func (h *HandlerService) handleLessonStart(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for lesson start
}

func (h *HandlerService) handleShowingFirstBlock(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for showing first block
}

func (h *HandlerService) handleExerciseAfterBlock(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for exercise after block
}

func (h *HandlerService) handleShowingIndividualWord(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for showing individual word
}

func (h *HandlerService) handleExerciseReview(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for exercise review
}

func (h *HandlerService) handleAudioDictation(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for audio dictation
}

func (h *HandlerService) handleTranslationCheck(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for translation check
}

func (h *HandlerService) handleWaitingForAudio(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for waiting for audio
}

func (h *HandlerService) handleWaitingForTranslation(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for waiting for translation
}

func (h *HandlerService) handleLessonComplete(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for lesson complete
}

func (h *HandlerService) handleAccountLinking(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for account linking
}

func (h *HandlerService) handleWaitingForLink(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for waiting for link
}

func (h *HandlerService) handleSettings(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for settings
}

func (h *HandlerService) handleHelp(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for help
}

func (h *HandlerService) handleStats(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for stats
}

func (h *HandlerService) handleSettingsCommand(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for settings command
}

func (h *HandlerService) handleLessonCommand(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for lesson command
}

func (h *HandlerService) handleCancel(ctx context.Context, message *tele.Message, progress *fsm.UserProgress) {
	// Implementation for cancel
}

// Callback handlers
func (h *HandlerService) handleStartLessonCallback(ctx context.Context, callback *tele.Callback, progress *fsm.UserProgress) {
	// Implementation for start lesson callback
}

func (h *HandlerService) handleVocabTestCallback(ctx context.Context, callback *tele.Callback, progress *fsm.UserProgress) {
	// Implementation for vocab test callback
}

func (h *HandlerService) handleExerciseCallback(ctx context.Context, callback *tele.Callback, progress *fsm.UserProgress) {
	// Implementation for exercise callback
}

func (h *HandlerService) handleLinkAccountCallback(ctx context.Context, callback *tele.Callback, progress *fsm.UserProgress) {
	// Create link token
	linkTokenResp, err := h.apiClient.CreateLinkToken(ctx, progress.TelegramID)
	if err != nil {
		logger.Log.Error("Failed to create link token", zap.Error(err))
		h.bot.Respond(callback, &tele.CallbackResponse{Text: "Ошибка создания ссылки"})
		return
	}

	text := fmt.Sprintf(`🔗 Для связывания аккаунтов перейдите по ссылке:

%s

Ссылка действительна 15 минут.`, linkTokenResp.LinkURL)

	h.bot.Edit(callback.Message, text)
	h.bot.Respond(callback, &tele.CallbackResponse{Text: "Ссылка создана"})
}

func (h *HandlerService) handleSettingsCallback(ctx context.Context, callback *tele.Callback, progress *fsm.UserProgress) {
	// TODO: Implement settings callback handling
	logger.Log.Info("Settings callback", zap.String("data", callback.Data))
	h.bot.Respond(callback, &tele.CallbackResponse{Text: "Настройки будут доступны в следующих версиях"})
}

// Task Handlers for Asynq

// HandleLessonReminderTask handles lesson reminder tasks
func (h *HandlerService) HandleLessonReminderTask(ctx context.Context, task *asynq.Task) error {
	payload := task.Payload()
	logger.Log.Info("Handling lesson reminder task", zap.ByteString("payload", payload))
	// TODO: Implement lesson reminder logic
	return nil
}

// HandleDailyNotificationTask handles daily notification tasks
func (h *HandlerService) HandleDailyNotificationTask(ctx context.Context, task *asynq.Task) error {
	payload := task.Payload()
	logger.Log.Info("Handling daily notification task", zap.ByteString("payload", payload))
	// TODO: Implement daily notification logic
	return nil
}

// HandleGenerateLessonTask handles lesson generation tasks
func (h *HandlerService) HandleGenerateLessonTask(ctx context.Context, task *asynq.Task) error {
	payload := task.Payload()
	logger.Log.Info("Handling generate lesson task", zap.ByteString("payload", payload))
	// TODO: Implement lesson generation logic
	return nil
}

// HandleSyncProgressTask handles progress synchronization tasks
func (h *HandlerService) HandleSyncProgressTask(ctx context.Context, task *asynq.Task) error {
	payload := task.Payload()
	logger.Log.Info("Handling sync progress task", zap.ByteString("payload", payload))
	// TODO: Implement progress sync logic
	return nil
}

// HandleCleanupSessionsTask handles session cleanup tasks
func (h *HandlerService) HandleCleanupSessionsTask(ctx context.Context, task *asynq.Task) error {
	payload := task.Payload()
	logger.Log.Info("Handling cleanup sessions task", zap.ByteString("payload", payload))
	// TODO: Implement session cleanup logic
	return nil
}
