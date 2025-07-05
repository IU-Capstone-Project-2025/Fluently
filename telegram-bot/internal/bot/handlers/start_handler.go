package handlers

import (
	"context"
	"fmt"

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
		"Привет, %s!\n\n"+
			"Добро пожаловать в Fluently — платформа для изучения *разговорного* английского языка.\n\n"+
			"Я помогу тебе научиться говорить на английском *свободно*",
		c.Sender().FirstName,
	)

	// Add "Get Started" button
	startBtn := &tele.InlineButton{
		Text: "Начать",
		Data: "onboarding:start",
	}
	alreadyHaveAccount := &tele.InlineButton{
		Text: "У меня уже есть аккаунт",
		Data: "account:link",
	}
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{*startBtn, *alreadyHaveAccount},
		},
	}

	// Send the message
	if _, err := s.bot.Send(c.Sender(), welcomeText, &tele.SendOptions{ParseMode: tele.ModeMarkdown, ReplyMarkup: keyboard}); err != nil {
		s.logger.Error("Failed to send welcome message", zap.Error(err))
		return err
	}

	// User should now be in StateWelcome, waiting for them to click "Начать"
	return nil
}

// HandleHelpCommand handles the /help command
func (s *HandlerService) HandleHelpCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	helpText := "🌟 *Справка по Fluently Bot* 🌟\n\n" +
		"Вот команды, которые вы можете использовать:\n\n" +
		"*/start* - Начать ваше путешествие в изучении языка\n" +
		"*/learn* - Начать сегодняшний урок\n" +
		"*/settings* - Настроить предпочтения обучения\n" +
		"*/test* - Пройти тест на определение уровня словарного запаса\n" +
		"*/stats* - Посмотреть статистику обучения\n" +
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
	methodText := "🎯 *Методика изучения Fluently*\n\n" +
		"Наша методика основана на принципах интервальных повторений и фокусе на самых важных словах.\n\n" +
		"🔑 **Ключевые принципы:**\n" +
		"• Изучение 10 слов в день\n" +
		"• Фокус на 80-90% самых используемых слов\n" +
		"• Интервальные повторения для закрепления\n" +
		"• Персонализированный подход\n\n" +
		"Готовы продолжить настройку?"

	// Create continue button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Продолжить", Data: "onboarding:method"}},
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
	spacedRepetitionText := "🧠 *Интервальные повторения*\n\n" +
		"Это научно обоснованная методика, которая помогает перенести информацию из кратковременной памяти в долговременную.\n\n" +
		"🔄 **Как это работает:**\n" +
		"• Повторение через увеличивающиеся интервалы\n" +
		"• Адаптация к вашему темпу обучения\n" +
		"• Оптимизация времени для каждого слова\n" +
		"• Максимальная эффективность запоминания\n\n" +
		"Теперь давайте узнаем о ваших целях и предпочтениях!"

	// Create continue button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Начать опрос", Data: "onboarding:questionnaire"}},
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
	questionnaireText := "📋 *Анкета для персонализации*\n\n" +
		"Для создания идеального плана обучения нам нужно узнать несколько вещей о вас.\n\n" +
		"Это займет всего 2-3 минуты, но поможет сделать ваше обучение максимально эффективным.\n\n" +
		"Готовы начать?"

	// Create continue button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Да, начинаем!", Data: "questionnaire:start"}},
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
	return c.Send("Главное меню...")
}

// HandleHelpMenuCallback handles help menu callback
func (s *HandlerService) HandleHelpMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Меню помощи...")
}

// HandleUnknownStateMessage handles unknown state messages
func (s *HandlerService) HandleUnknownStateMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Я не знаю, что делать в этом состоянии. Используйте /help для просмотра доступных команд.")
}

// HandleUnknownCallback handles unknown callbacks
func (s *HandlerService) HandleUnknownCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Получен неизвестный callback.")
}
