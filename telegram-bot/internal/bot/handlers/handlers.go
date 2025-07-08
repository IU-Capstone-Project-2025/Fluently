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
		"👋 Привет, %s!\n\n"+
			"Добро пожаловать в Fluently — ваш персональный помощник для изучения английского языка.\n\n"+
			"Я помогу вам стать свободно говорящим на английском через персонализированные уроки, интервальные повторения и интерактивные упражнения.",
		c.Sender().FirstName,
	)

	// Add "Get Started" button
	getStartedButton := &tele.InlineButton{
		Text: "Начать",
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
	settingsText := "⚙️ *Настройки*\n\n" +
		fmt.Sprintf("🔤 Уровень CEFR: *%s*\n", userProgress.CEFRLevel) +
		fmt.Sprintf("📚 Слов в день: *%d*\n", userProgress.WordsPerDay) +
		fmt.Sprintf("🔔 Уведомления: *%s*\n", formatNotificationTime(userProgress.NotificationTime)) +
		"\nВыберите настройку для изменения:"

	// Create settings keyboard
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "🔤 Уровень CEFR", Data: "settings:cefr_level"}},
			{{Text: "📚 Слов в день", Data: "settings:words_per_day"}},
			{{Text: "🔔 Уведомления", Data: "settings:notifications"}},
			{{Text: "Назад в главное меню", Data: "menu:main"}},
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
			Text: "Начать настройку",
			Data: "onboarding:start",
		}
		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{*startButton},
			},
		}

		return c.Send("Похоже, вы еще не завершили первоначальную настройку. "+
			"Давайте сначала определим ваш уровень английского.", keyboard)
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
		"📚 *Сегодняшний урок*\n\n"+
			"Давайте изучим %d новых слов сегодня.\n\n"+
			"Готовы начать ваш ежедневный урок?",
		userProgress.WordsPerDay,
	)

	// Create lesson keyboard
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "Начать изучение", Data: "lesson:start"},
				{Text: "Позже", Data: "lesson:later"},
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
	testText := "🧠 *Тест уровня CEFR*\n\n" +
		"Этот тест поможет определить ваш уровень владения английским языком согласно шкале CEFR.\n\n" +
		"Вы увидите серию слов. Для каждого слова укажите, хорошо ли вы его знаете.\n\n" +
		"Тест состоит из 5 частей и займет около 5-10 минут. Готовы начать?"

	// Create test keyboard
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "Начать тест", Data: "test:start"},
				{Text: "Позже", Data: "menu:main"},
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
		"📊 *Ваша статистика обучения*\n\n"+
			"🔤 Текущий уровень: *%s*\n"+
			"📚 Слов в день: *%d*\n"+
			"🔥 Текущая серия: *%d дней*\n"+
			"📖 Всего слов изучено: *%d*\n"+
			"🎯 Уроков завершено: *%d*\n"+
			"⏱ Общее время обучения: *%d минут*\n",
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
			{{Text: "Назад в главное меню", Data: "menu:main"}},
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
	cancelText := "❌ Действие отменено. Вы возвращены в главное меню.\n\n" +
		"Используйте /start чтобы начать заново или /help чтобы увидеть доступные команды."

	return c.Send(cancelText)
}

// Placeholder handler functions - these need to be implemented
func (s *HandlerService) HandleWelcomeMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Добро пожаловать! Пожалуйста, используйте кнопки или команды для навигации.")
}

func (s *HandlerService) HandleMethodExplanationMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Объяснение методики. Пожалуйста, продолжите настройку.")
}

func (s *HandlerService) HandleQuestionGoalMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос о цели.")
}

func (s *HandlerService) HandleQuestionConfidenceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос об уверенности.")
}

func (s *HandlerService) HandleQuestionSerialsMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос о сериалах.")
}

func (s *HandlerService) HandleQuestionExperienceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос об опыте.")
}

func (s *HandlerService) HandleSettingsWordsPerDayInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, введите количество слов в день.")
}

func (s *HandlerService) HandleSettingsTimeInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, введите время уведомлений.")
}

func (s *HandlerService) HandleWaitingForTranslationMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, предоставьте перевод.")
}

func (s *HandlerService) HandleWaitingForAudioMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, предоставьте аудио ответ.")
}

func (s *HandlerService) HandleUnknownStateMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Я не знаю, что делать в этом состоянии. Используйте /help для просмотра доступных команд.")
}

// Callback handlers - placeholder implementations
func (s *HandlerService) HandleOnboardingStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Начинаем знакомство...")
}

func (s *HandlerService) HandleOnboardingMethodCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Получен callback методики.")
}

func (s *HandlerService) HandleTestStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Начинаем тест...")
}

func (s *HandlerService) HandleLessonStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Начинаем урок...")
}

func (s *HandlerService) HandleLessonLaterCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Урок отложен.")
}

func (s *HandlerService) HandleSettingsWordsPerDayCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Настройка слов в день...")
}

func (s *HandlerService) HandleSettingsNotificationsCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Настройка уведомлений...")
}

func (s *HandlerService) HandleSettingsCEFRLevelCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Настройка уровня CEFR...")
}

func (s *HandlerService) HandleMainMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Главное меню...")
}

func (s *HandlerService) HandleSettingsMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Меню настроек...")
}

func (s *HandlerService) HandleLearnMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Меню обучения...")
}

func (s *HandlerService) HandleHelpMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Меню помощи...")
}

func (s *HandlerService) HandleUnknownCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Получен неизвестный callback.")
}

func (s *HandlerService) HandleAudioExerciseResponse(ctx context.Context, c tele.Context, userID int64, voice interface{}) error {
	return c.Send("Получен ответ на аудио упражнение.")
}

func formatNotificationTime(timeStr string) string {
	if timeStr == "" {
		return "Отключены"
	}
	return timeStr
}
