package handlers

import (
	"context"
	"fmt"
	"strings"

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

	// Check user authentication status
	isAuthenticated, hasCompletedOnboarding, err := s.GetUserAuthenticationStatus(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user authentication status", zap.Error(err))
		return err
	}

	// Handle different user states
	if isAuthenticated && hasCompletedOnboarding {
		// User is fully set up - fast-track to main menu
		return s.showMainMenu(ctx, c, userID)
	} else if isAuthenticated && !hasCompletedOnboarding {
		// User is authenticated but hasn't completed onboarding - fast-track onboarding
		return s.showFastTrackOnboarding(ctx, c, userID)
	} else {
		// User is not authenticated - show initial welcome with auth options
		return s.showWelcomeWithAuthOptions(ctx, c, userID)
	}
}

// showWelcomeWithAuthOptions shows the welcome message with authentication options
func (s *HandlerService) showWelcomeWithAuthOptions(ctx context.Context, c tele.Context, userID int64) error {
	// Transition to welcome state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateWelcome); err != nil {
		s.logger.Error("Failed to set welcome state", zap.Error(err))
		return err
	}

	// Send welcome message
	welcomeText := fmt.Sprintf(
		"Привет, %s! 👋\n\n"+
			"Я помогу тебе выучить английский легко и весело!\n\n"+
			"Выбери, как продолжить:",
		c.Sender().FirstName,
	)

	// Create buttons for different flows
	existingUserBtn := &tele.InlineButton{
		Text: "У меня уже есть аккаунт",
		Data: "auth:existing_user",
	}
	newUserBtn := &tele.InlineButton{
		Text: "Начать",
		Data: "auth:new_user",
	}

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{*newUserBtn},
			{*existingUserBtn},
		},
	}

	return c.Send(welcomeText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// showMainMenu shows the main menu for authenticated users
func (s *HandlerService) showMainMenu(ctx context.Context, c tele.Context, userID int64) error {
	// Set state to start (only if different from current state)
	if err := s.SetStateIfDifferent(ctx, userID, fsm.StateStart); err != nil {
		s.logger.Error("Failed to set start state", zap.Error(err))
		return err
	}

	// Get user progress for personalized welcome
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	welcomeText := fmt.Sprintf(
		"🎉 *С возвращением!*\n\n"+
			"📊 Твой уровень: *%s*\n"+
			"📚 Слов в день: *%d*\n\n"+
			"Что будем изучать сегодня?",
		userProgress.CEFRLevel,
		userProgress.WordsPerDay,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🚀 Начать урок", Data: "lesson:start"},
				{Text: "📊 Статистика", Data: "stats:show"},
			},
			{
				{Text: "⚙️ Настройки", Data: "menu:settings"},
				{Text: "❓ Помощь", Data: "menu:help"},
			},
		},
	}

	return c.Send(welcomeText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// showFastTrackOnboarding shows fast-track onboarding for authenticated but incomplete users
func (s *HandlerService) showFastTrackOnboarding(ctx context.Context, c tele.Context, userID int64) error {
	// Set state to questionnaire
	if err := s.stateManager.SetState(ctx, userID, fsm.StateQuestionnaire); err != nil {
		s.logger.Error("Failed to set questionnaire state", zap.Error(err))
		return err
	}

	onboardingText := fmt.Sprintf(
		"✨ *Привет, %s!*\n\n"+
			"Осталось совсем немного - давай закончим настройку твоего профиля для эффективного обучения.\n\n"+
			"Это займет всего 2 минуты! 🕐",
		c.Sender().FirstName,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "✅ Завершить настройку", Data: "questionnaire:start"}},
			{{Text: "🚀 Пропустить в урок", Data: "lesson:start"}},
		},
	}

	return c.Send(onboardingText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleHelpCommand handles the /help command
func (s *HandlerService) HandleHelpCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	helpText := "🌟 *Справка по Fluently Bot* 🌟\n\n" +
		"Вот команды, которые вы можете использовать:\n\n" +
		"*/start* - Начать ваше путешествие в изучении языка\n" +
		"*/learn* - Начать сегодняшний урок\n" +
		"*/lesson* - Быстрый доступ к уроку\n" +
		"*/settings* - Настроить предпочтения обучения\n" +
		"*/test* - Пройти тест на определение уровня словарного запаса\n" +
		"*/stats* - Посмотреть статистику обучения\n" +
		"*/menu* - Вернуться в главное меню\n" +
		"*/help* - Показать это сообщение справки\n" +
		"*/cancel* - Отменить текущее действие\n\n" +
		"Нужна дополнительная помощь? Напишите свой вопрос, и я постараюсь помочь."

	return c.Send(helpText, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}

// HandleCancelCommand handles the /cancel command
func (s *HandlerService) HandleCancelCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Reset user to initial state
	if err := s.stateManager.ResetUserToInitial(ctx, userID); err != nil {
		s.logger.Error("Failed to reset user state", zap.Error(err))
		return err
	}

	// Send cancellation message
	cancelText := "❌ Действие отменено. Вы возвращены в главное меню.\n\n" +
		"Используйте /start чтобы начать заново или /help чтобы увидеть доступные команды."

	return c.Send(cancelText)
}

// HandleMenuCommand handles the /menu command
func (s *HandlerService) HandleMenuCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return s.HandleMainMenuCallback(ctx, c, userID, currentState)
}

// HandleWelcomeMessage handles welcome state messages
func (s *HandlerService) HandleWelcomeMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Добро пожаловать! Пожалуйста, используйте кнопки или команды для навигации.")
}

// HandleMethodExplanationMessage handles method explanation state messages
func (s *HandlerService) HandleMethodExplanationMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Объяснение методики. Пожалуйста, продолжите настройку.")
}

// HandleOnboardingStartCallback handles the onboarding start callback
func (s *HandlerService) HandleOnboardingStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Validate current state
	if currentState != fsm.StateWelcome {
		s.logger.Warn("Invalid state for onboarding start",
			zap.Int64("user_id", userID),
			zap.String("expected_state", string(fsm.StateWelcome)),
			zap.String("actual_state", string(currentState)))
		return c.Send("Пожалуйста, начните с команды /start")
	}

	// Transition to method explanation state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateMethodExplanation); err != nil {
		s.logger.Error("Failed to set method explanation state", zap.Error(err))
		return err
	}

	// Send method explanation message
	methodText := "🎯 *Как это работает?*\n\n" +
		"• 10 новых слов каждый день\n" +
		"• Только самые нужные слова\n" +
		"• Повторения в нужный момент\n\n" +
		"Просто и эффективно! 🚀"

	// Create continue button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Понятно!", Data: "onboarding:method"}},
		},
	}

	return c.Send(methodText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleOnboardingMethodCallback handles the transition from method explanation to spaced repetition
func (s *HandlerService) HandleOnboardingMethodCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Validate current state
	if currentState != fsm.StateMethodExplanation {
		s.logger.Warn("Invalid state for method callback",
			zap.Int64("user_id", userID),
			zap.String("expected_state", string(fsm.StateMethodExplanation)),
			zap.String("actual_state", string(currentState)))
		return c.Send("Пожалуйста, начните с команды /start")
	}

	// Transition to spaced repetition explanation state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSpacedRepetition); err != nil {
		s.logger.Error("Failed to set spaced repetition state", zap.Error(err))
		return err
	}

	// Send spaced repetition explanation message
	spacedRepetitionText := "🧠 *Секрет запоминания*\n\n" +
		"Показываю слово именно тогда, когда ты его почти забыл.\n\n" +
		"Так твой мозг запоминает навсегда! 💡"

	// Create continue button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Круто! Дальше", Data: "onboarding:questionnaire"}},
		},
	}

	return c.Send(spacedRepetitionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleOnboardingQuestionnaireCallback handles the transition to questionnaire
func (s *HandlerService) HandleOnboardingQuestionnaireCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Validate current state
	if currentState != fsm.StateSpacedRepetition {
		s.logger.Warn("Invalid state for questionnaire callback",
			zap.Int64("user_id", userID),
			zap.String("expected_state", string(fsm.StateSpacedRepetition)),
			zap.String("actual_state", string(currentState)))
		return c.Send("Пожалуйста, начните с команды /start")
	}

	// Transition to questionnaire state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateQuestionnaire); err != nil {
		s.logger.Error("Failed to set questionnaire state", zap.Error(err))
		return err
	}

	// Send questionnaire introduction message
	questionnaireText := "📋 *Расскажи о себе*\n\n" +
		"Пару быстрых вопросов, чтобы подобрать уроки именно для тебя.\n\n" +
		"Займет 1 минуту 🕐"

	// Create continue button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Поехали!", Data: "questionnaire:start"}},
		},
	}

	return c.Send(questionnaireText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleAccountLinkCallback handles account linking callback
func (s *HandlerService) HandleAccountLinkCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Получен callback аккаунта.")
}

// HandleMainMenuCallback handles main menu callback
func (s *HandlerService) HandleMainMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Clear settings message if we're coming from settings
	if fsm.IsSettingsState(currentState) {
		s.clearSettingsMessage(ctx, c, userID)
	}

	return s.showMainMenu(ctx, c, userID)
}

// HandleBackToMainMenuCallback handles "back to main menu" callback with message deletion
func (s *HandlerService) HandleBackToMainMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Delete the previous message when going back to main menu
	if c.Message() != nil {
		if err := c.Delete(); err != nil {
			// Only log as warning if it's not a "message not found" error
			if !strings.Contains(err.Error(), "message to delete not found") {
				s.logger.Warn("Failed to delete previous message", zap.Error(err))
			} else {
				s.logger.Debug("Previous message already deleted or not found", zap.Error(err))
			}
		}
	}

	// Clear settings message if we're coming from settings
	if fsm.IsSettingsState(currentState) {
		s.clearSettingsMessage(ctx, c, userID)
	}

	return s.showMainMenu(ctx, c, userID)
}

// HandleHelpMenuCallback handles help menu callback
func (s *HandlerService) HandleHelpMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return s.HandleHelpCommand(ctx, c, userID, currentState)
}

// HandleUnknownStateMessage handles unknown state messages
func (s *HandlerService) HandleUnknownStateMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Я не знаю, что делать в этом состоянии. Используйте /help для просмотра доступных команд.")
}

// HandleUnknownCallback handles unknown callbacks
func (s *HandlerService) HandleUnknownCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Получен неизвестный callback.")
}
