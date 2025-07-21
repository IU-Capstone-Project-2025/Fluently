package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/config"
	"telegram-bot/internal/api"
	"telegram-bot/internal/bot/fsm"
	"telegram-bot/internal/bot/handlers"
	"telegram-bot/internal/tasks"
)

// TelegramBot wraps the bot API and provides additional functionality
type TelegramBot struct {
	bot            *tele.Bot
	stateManager   *fsm.UserStateManager
	apiClient      *api.Client
	logger         *zap.Logger
	redisClient    *redis.Client
	handlerService *handlers.HandlerService
	config         *config.Config
}

// NewTelegramBot creates a new Telegram bot instance
func NewTelegramBot(cfg *config.Config, redisClient *redis.Client, apiClient *api.Client, scheduler *tasks.Scheduler, logger *zap.Logger) (*TelegramBot, error) {
	// Create bot settings for long polling
	settings := tele.Settings{
		Token:  cfg.Bot.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	logger.Info("Using long polling")

	// Create bot instance
	bot, err := tele.NewBot(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	// Create state manager
	stateManager := fsm.NewUserStateManager(redisClient)

	// Create handler service
	handlerService := handlers.NewHandlerService(cfg, redisClient, apiClient, scheduler, bot, stateManager, logger)

	telegramBot := &TelegramBot{
		bot:            bot,
		stateManager:   stateManager,
		apiClient:      apiClient,
		logger:         logger,
		redisClient:    redisClient,
		handlerService: handlerService,
		config:         cfg,
	}

	// Setup handlers
	telegramBot.setupHandlers()

	return telegramBot, nil
}

// setupHandlers configures all bot handlers
func (tb *TelegramBot) setupHandlers() {
	// Command handlers
	tb.bot.Handle("/start", tb.handleStart)
	tb.bot.Handle("/help", tb.handleHelp)
	tb.bot.Handle("/settings", tb.handleSettings)
	tb.bot.Handle("/learn", tb.handleLearn)
	tb.bot.Handle("/lesson", tb.handleLesson)
	tb.bot.Handle("/test", tb.handleTest)
	tb.bot.Handle("/stats", tb.handleStats)
	tb.bot.Handle("/cancel", tb.handleCancel)
	tb.bot.Handle("/menu", tb.handleMenu)

	// Message handler for all text messages
	tb.bot.Handle(tele.OnText, tb.handleMessage)

	// Callback handler for all callback queries
	tb.bot.Handle(tele.OnCallback, tb.handleCallback)

	// Set error handler
	tb.bot.Use(tb.errorMiddleware)
}

// errorMiddleware handles errors in middleware
func (tb *TelegramBot) errorMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		err := next(c)
		if err != nil {
			tb.logger.Error("Handler error occurred", zap.Error(err))

			// Try to send error message to user
			if sendErr := c.Send("⚠️ Что-то пошло не так. Пожалуйста, попробуйте еще раз или используйте /cancel для сброса."); sendErr != nil {
				tb.logger.Error("Failed to send error message to user", zap.Error(sendErr))
			}
		}
		return err
	}
}

// handleStart handles the /start command
func (tb *TelegramBot) handleStart(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleStartCommand(ctx, c, userID, currentState)
}

// handleHelp handles the /help command
func (tb *TelegramBot) handleHelp(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleHelpCommand(ctx, c, userID, currentState)
}

// handleSettings handles the /settings command
func (tb *TelegramBot) handleSettings(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleSettingsCommand(ctx, c, userID, currentState)
}

// handleLearn handles the /learn command
func (tb *TelegramBot) handleLearn(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleLearnCommand(ctx, c, userID, currentState)
}

// handleLesson handles the /lesson command
func (tb *TelegramBot) handleLesson(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleLessonCommand(ctx, c, userID, currentState)
}

// handleTest handles the /test command
func (tb *TelegramBot) handleTest(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleTestCommand(ctx, c, userID, currentState)
}

// handleStats handles the /stats command
func (tb *TelegramBot) handleStats(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleStatsCommand(ctx, c, userID, currentState)
}

// handleCancel handles the /cancel command
func (tb *TelegramBot) handleCancel(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleCancelCommand(ctx, c, userID, currentState)
}

// handleMenu handles the /menu command
func (tb *TelegramBot) handleMenu(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleMenuCommand(ctx, c, userID, currentState)
}

// handleMessage handles all text messages
func (tb *TelegramBot) handleMessage(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleTextMessage(ctx, c, userID, currentState)
}

// handleCallback handles all callback queries
func (tb *TelegramBot) handleCallback(c tele.Context) error {
	ctx := context.Background()
	userID := c.Sender().ID
	currentState, err := tb.stateManager.GetState(ctx, userID)
	if err != nil {
		tb.logger.Error("Failed to get user state")
		currentState = fsm.StateStart
	}

	return tb.handlerService.HandleCallback(ctx, c, userID, currentState)
}

// Start starts the bot
func (tb *TelegramBot) Start() error {
	tb.logger.Info("Starting Telegram bot...")

	// Test Redis connection
	ctx := context.Background()
	if err := tb.redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	tb.logger.Info("Redis connection established")

	// Start the bot
	tb.bot.Start()
	return nil
}

// Stop stops the bot gracefully
func (tb *TelegramBot) Stop() {
	tb.logger.Info("Stopping Telegram bot...")
	tb.bot.Stop()
}
