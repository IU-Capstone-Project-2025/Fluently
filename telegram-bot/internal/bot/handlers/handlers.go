// handlers.go
package handlers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
)

// HandleStartCommand handles the /start command
func (s *HandlerService) HandleStartCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Reset user to initial state
	if err := s.stateManager.ResetUserToInitial(ctx, userID); err != nil {
		s.logger.Error("Failed to reset user state", zap.Error(err))
		return err
	}

	// Transition to welcome state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateWelcome); err != nil {
		s.logger.Error("Failed to set welcome state", zap.Error(err))
		return err
	}

	// Send welcome message
	welcomeText := fmt.Sprintf(
		"👋 Hello, %s!\n\n"+
			"Welcome to Fluently, your personal English language learning assistant.\n\n"+
			"I'll help you become fluent in English through personalized lessons, spaced repetition, and interactive exercises.",
		c.Sender().FirstName,
	)

	// Add "Get Started" button
	getStartedButton := &tele.InlineButton{
		Text: "Get Started",
		Data: "onboarding:start",
	}
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{*getStartedButton},
		},
	}

	return c.Send(welcomeText, keyboard)
}

// HandleHelpCommand handles the /help command
func (s *HandlerService) HandleHelpCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	helpText := "🌟 *Fluently Bot Help* 🌟\n\n" +
		"Here are the commands you can use:\n\n" +
		"*/start* - Begin your language learning journey\n" +
		"*/learn* - Start today's lesson\n" +
		"*/settings* - Configure your learning preferences\n" +
		"*/test* - Take a vocabulary test to determine your level\n" +
		"*/stats* - View your learning statistics\n" +
		"*/help* - Show this help message\n" +
		"*/cancel* - Cancel the current action\n\n" +
		"Need more help? Type your question and I'll try to assist you."

	return c.Send(helpText, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

// HandleSettingsCommand handles the /settings command
func (s *HandlerService) HandleSettingsCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Get current user progress for settings
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Set user state to settings
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	// Create settings message
	settingsText := "⚙️ *Settings*\n\n" +
		fmt.Sprintf("🔤 CEFR Level: *%s*\n", userProgress.CEFRLevel) +
		fmt.Sprintf("📚 Words per day: *%d*\n", userProgress.WordsPerDay) +
		fmt.Sprintf("🔔 Notifications: *%s*\n", formatNotificationTime(userProgress.NotificationTime)) +
		"\nSelect a setting to change:"

	// Create settings keyboard
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "🔤 CEFR Level", Data: "settings:cefr_level"}},
			{{Text: "📚 Words per day", Data: "settings:words_per_day"}},
			{{Text: "🔔 Notifications", Data: "settings:notifications"}},
			{{Text: "Back to Main Menu", Data: "menu:main"}},
		},
	}

	return c.Send(settingsText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleLearnCommand handles the /learn command
func (s *HandlerService) HandleLearnCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Check if the user has completed the onboarding process
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// If user hasn't completed onboarding, prompt them to start
	if userProgress.CEFRLevel == "" {
		startButton := &tele.InlineButton{
			Text: "Start Setup",
			Data: "onboarding:start",
		}
		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{*startButton},
			},
		}

		return c.Send("It looks like you haven't completed the initial setup yet. "+
			"Let's determine your English level first.", keyboard)
	}

	// Set user state to lesson start
	if err := s.stateManager.SetState(ctx, userID, fsm.StateLessonStart); err != nil {
		s.logger.Error("Failed to set lesson start state", zap.Error(err))
		return err
	}

	// Create new lesson data
	lessonData := &fsm.LessonData{
		StartTime: time.Now(),
	}

	// Store lesson data
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataLesson, lessonData); err != nil {
		s.logger.Error("Failed to store lesson data", zap.Error(err))
		return err
	}

	// Send lesson start message
	lessonText := fmt.Sprintf(
		"📚 *Today's Lesson*\n\n"+
			"Let's practice %d new words today.\n\n"+
			"Ready to start your daily lesson?",
		userProgress.WordsPerDay,
	)

	// Create lesson keyboard
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "Start Learning", Data: "lesson:start"},
				{Text: "Later", Data: "lesson:later"},
			},
		},
	}

	return c.Send(lessonText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleTestCommand handles the /test command
func (s *HandlerService) HandleTestCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Set user state to vocabulary test
	if err := s.stateManager.SetState(ctx, userID, fsm.StateVocabularyTest); err != nil {
		s.logger.Error("Failed to set vocabulary test state", zap.Error(err))
		return err
	}

	// Send test introduction message
	testText := "🧠 *CEFR Level Test*\n\n" +
		"This test will help determine your English proficiency level according to the CEFR scale.\n\n" +
		"You'll see a series of words. For each word, indicate if you know it well.\n\n" +
		"The test has 5 parts and takes about 5-10 minutes. Ready to begin?"

	// Create test keyboard
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "Start Test", Data: "test:start"},
				{Text: "Later", Data: "menu:main"},
			},
		},
	}

	return c.Send(testText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleStatsCommand handles the /stats command
func (s *HandlerService) HandleStatsCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Get user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Create stats message
	statsText := fmt.Sprintf(
		"📊 *Your Learning Statistics*\n\n"+
			"🔤 Current Level: *%s*\n"+
			"📚 Words per day: *%d*\n"+
			"🔥 Current streak: *%d days*\n"+
			"📖 Total words learned: *%d*\n"+
			"🎯 Lessons completed: *%d*\n"+
			"⏱ Total time spent: *%d minutes*\n",
		userProgress.CEFRLevel,
		userProgress.WordsPerDay,
		0, // streak days - placeholder
		0, // total words - placeholder
		0, // lessons completed - placeholder
		0, // total time - placeholder
	)

	// Create back button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Back to Main Menu", Data: "menu:main"}},
		},
	}

	return c.Send(statsText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleCancelCommand handles the /cancel command
func (s *HandlerService) HandleCancelCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Reset user to initial state
	if err := s.stateManager.ResetUserToInitial(ctx, userID); err != nil {
		s.logger.Error("Failed to reset user state", zap.Error(err))
		return err
	}

	// Send cancellation message
	cancelText := "❌ Action cancelled. You've been returned to the main menu.\n\n" +
		"Use /start to begin again or /help to see available commands."

	return c.Send(cancelText)
}

// Placeholder handler functions - these need to be implemented
func (s *HandlerService) HandleWelcomeMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Welcome! Please use the buttons or commands to navigate.")
}

func (s *HandlerService) HandleMethodExplanationMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Method explanation. Please continue with the setup.")
}

func (s *HandlerService) HandleQuestionGoalMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Please answer the goal question.")
}

func (s *HandlerService) HandleQuestionConfidenceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Please answer the confidence question.")
}

func (s *HandlerService) HandleQuestionSerialsMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Please answer the serials question.")
}

func (s *HandlerService) HandleQuestionExperienceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Please answer the experience question.")
}

func (s *HandlerService) HandleSettingsWordsPerDayInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Please enter the number of words per day.")
}

func (s *HandlerService) HandleSettingsTimeInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Please enter the notification time.")
}

func (s *HandlerService) HandleWaitingForTranslationMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Please provide the translation.")
}

func (s *HandlerService) HandleWaitingForAudioMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Please provide the audio response.")
}

func (s *HandlerService) HandleUnknownStateMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("I'm not sure what to do in this state. Use /help for available commands.")
}

// Callback handlers - placeholder implementations
func (s *HandlerService) HandleOnboardingStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Starting onboarding...")
}

func (s *HandlerService) HandleOnboardingMethodCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Method callback received.")
}

func (s *HandlerService) HandleTestStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Starting test...")
}

func (s *HandlerService) HandleLessonStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Starting lesson...")
}

func (s *HandlerService) HandleLessonLaterCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Lesson postponed.")
}

func (s *HandlerService) HandleSettingsWordsPerDayCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Words per day setting...")
}

func (s *HandlerService) HandleSettingsNotificationsCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Notifications setting...")
}

func (s *HandlerService) HandleSettingsCEFRLevelCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("CEFR level setting...")
}

func (s *HandlerService) HandleMainMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Main menu...")
}

func (s *HandlerService) HandleSettingsMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Settings menu...")
}

func (s *HandlerService) HandleLearnMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Learn menu...")
}

func (s *HandlerService) HandleHelpMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Help menu...")
}

func (s *HandlerService) HandleUnknownCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Unknown callback received.")
}

func (s *HandlerService) HandleAudioExerciseResponse(ctx context.Context, c tele.Context, userID int64, voice interface{}) error {
	return c.Send("Audio exercise response received.")
}

func formatNotificationTime(timeStr string) string {
	if timeStr == "" {
		return "Disabled"
	}
	return timeStr
}
