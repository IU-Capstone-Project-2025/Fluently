package handlers

import (
	"context"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
)

// HandleQuestionGoalMessage handles goal question messages
func (s *HandlerService) HandleQuestionGoalMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос о цели.")
}

// HandleQuestionConfidenceMessage handles confidence question messages
func (s *HandlerService) HandleQuestionConfidenceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос об уверенности.")
}

// HandleQuestionExperienceMessage handles experience question messages
func (s *HandlerService) HandleQuestionExperienceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос об опыте.")
}

// HandleQuestionWordsPerDayMessage handles words per day question messages
func (s *HandlerService) HandleQuestionWordsPerDayMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, выберите количество слов для изучения в день.")
}

// HandleQuestionNotificationsMessage handles notifications question messages
func (s *HandlerService) HandleQuestionNotificationsMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, выберите настройки уведомлений.")
}

// HandleQuestionNotificationTimeMessage handles notification time question messages
func (s *HandlerService) HandleQuestionNotificationTimeMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, выберите время для уведомлений.")
}

// HandleQuestionnaireStartCallback handles the questionnaire start callback
func (s *HandlerService) HandleQuestionnaireStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Validate current state
	if currentState != fsm.StateQuestionnaire {
		s.logger.Warn("Invalid state for questionnaire start",
			zap.Int64("user_id", userID),
			zap.String("expected_state", string(fsm.StateQuestionnaire)),
			zap.String("actual_state", string(currentState)))
		return c.Send("Пожалуйста, начните с команды /start")
	}

	// Transition to first question state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateQuestionGoal); err != nil {
		s.logger.Error("Failed to set question goal state", zap.Error(err))
		return err
	}

	// Send first question
	questionText := "🎯 *Первый вопрос*\n\n" +
		"Какая у тебя главная цель в изучении английского?"

	// Create answer options
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Работа/карьера", Data: "goal:work"}},
			{{Text: "Путешествия", Data: "goal:travel"}},
			{{Text: "Образование", Data: "goal:education"}},
			{{Text: "Общение", Data: "goal:communication"}},
		},
	}

	return c.Send(questionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleGoalCallback handles goal question callback
func (s *HandlerService) HandleGoalCallback(ctx context.Context, c tele.Context, userID int64, answer string) error {
	s.logger.Info("User answered goal question", zap.Int64("user_id", userID), zap.String("answer", answer))

	// Store goal answer
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataGoal, answer); err != nil {
		s.logger.Error("Failed to store goal answer", zap.Error(err))
	}

	// Transition to confidence question
	if err := s.stateManager.SetState(ctx, userID, fsm.StateQuestionConfidence); err != nil {
		s.logger.Error("Failed to set confidence state", zap.Error(err))
		return err
	}

	// Send confidence question
	questionText := "🤔 *Вопрос 2*\n\n" +
		"Как ты оцениваешь свой уровень английского?"

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Новичок", Data: "confidence:beginner"}},
			{{Text: "Базовый", Data: "confidence:elementary"}},
			{{Text: "Средний", Data: "confidence:intermediate"}},
			{{Text: "Продвинутый", Data: "confidence:advanced"}},
		},
	}

	return c.Send(questionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleConfidenceCallback handles confidence question callback
func (s *HandlerService) HandleConfidenceCallback(ctx context.Context, c tele.Context, userID int64, answer string) error {
	s.logger.Info("User answered confidence question", zap.Int64("user_id", userID), zap.String("answer", answer))

	// Store confidence level for later use
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataConfidence, answer); err != nil {
		s.logger.Error("Failed to store confidence level", zap.Error(err))
	}

	// Transition directly to experience question (skip serials)
	if err := s.stateManager.SetState(ctx, userID, fsm.StateQuestionExperience); err != nil {
		s.logger.Error("Failed to set experience state", zap.Error(err))
		return err
	}

	// Send experience question
	questionText := "🎓 *Последний вопрос*\n\n" +
		"Сколько времени ты изучаешь английский?"

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Только начинаю", Data: "experience:beginner"}},
			{{Text: "Меньше года", Data: "experience:less_year"}},
			{{Text: "1-3 года", Data: "experience:1_3_years"}},
			{{Text: "Больше 3 лет", Data: "experience:more_3_years"}},
		},
	}

	return c.Send(questionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleExperienceCallback handles experience question callback
func (s *HandlerService) HandleExperienceCallback(ctx context.Context, c tele.Context, userID int64, answer string) error {
	s.logger.Info("User answered experience question", zap.Int64("user_id", userID), zap.String("answer", answer))

	// Store experience answer
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataExperience, answer); err != nil {
		s.logger.Error("Failed to store experience answer", zap.Error(err))
	}

	// Transition to words per day question
	if err := s.stateManager.SetState(ctx, userID, fsm.StateQuestionWordsPerDay); err != nil {
		s.logger.Error("Failed to set words per day question state", zap.Error(err))
		return err
	}

	// Send words per day question
	questionText := "📚 *Сколько слов в день?*\n\n" +
		"Выберите количество новых слов, которые вы хотите изучать каждый день.\n\n" +
		"Рекомендуется начать с 10 слов для эффективного запоминания."

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "5 слов", Data: "words_per_day:5"}},
			{{Text: "10 слов (рекомендуется)", Data: "words_per_day:10"}},
			{{Text: "15 слов", Data: "words_per_day:15"}},
			{{Text: "20 слов", Data: "words_per_day:20"}},
		},
	}

	return c.Send(questionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleWordsPerDayCallback handles words per day question callback
func (s *HandlerService) HandleWordsPerDayCallback(ctx context.Context, c tele.Context, userID int64, answer string) error {
	s.logger.Info("User answered words per day question", zap.Int64("user_id", userID), zap.String("answer", answer))

	// Convert answer to integer
	var wordsPerDay int
	switch answer {
	case "5":
		wordsPerDay = 5
	case "10":
		wordsPerDay = 10
	case "15":
		wordsPerDay = 15
	case "20":
		wordsPerDay = 20
	default:
		wordsPerDay = 10 // Default fallback
	}

	// Store words per day answer
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataWordsPerDay, wordsPerDay); err != nil {
		s.logger.Error("Failed to store words per day answer", zap.Error(err))
	}

	// Transition to notifications question
	if err := s.stateManager.SetState(ctx, userID, fsm.StateQuestionNotifications); err != nil {
		s.logger.Error("Failed to set notifications question state", zap.Error(err))
		return err
	}

	// Send notifications question
	questionText := "🔔 *Уведомления*\n\n" +
		"Хотите получать ежедневные напоминания об изучении новых слов?\n\n" +
		"Это поможет сформировать полезную привычку!"

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "✅ Да, включить уведомления", Data: "notifications:enabled"}},
			{{Text: "❌ Нет, без уведомлений", Data: "notifications:disabled"}},
		},
	}

	return c.Send(questionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleNotificationsCallback handles notifications question callback
func (s *HandlerService) HandleNotificationsCallback(ctx context.Context, c tele.Context, userID int64, answer string) error {
	s.logger.Info("User answered notifications question", zap.Int64("user_id", userID), zap.String("answer", answer))

	// Convert answer to boolean
	notificationsEnabled := answer == "enabled"

	// Store notifications answer
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataNotifications, notificationsEnabled); err != nil {
		s.logger.Error("Failed to store notifications answer", zap.Error(err))
	}

	if notificationsEnabled {
		// If notifications are enabled, ask for time preference
		if err := s.stateManager.SetState(ctx, userID, fsm.StateQuestionNotificationTime); err != nil {
			s.logger.Error("Failed to set notification time question state", zap.Error(err))
			return err
		}

		// Send notification time question
		questionText := "⏰ *Время уведомлений*\n\n" +
			"В какое время дня вам удобно получать напоминания об изучении слов?"

		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{{Text: "🌅 Утром (9:00)", Data: "notification_time:09:00"}},
				{{Text: "🏢 Днем (14:00)", Data: "notification_time:14:00"}},
				{{Text: "🌆 Вечером (19:00)", Data: "notification_time:19:00"}},
				{{Text: "🌙 Поздно вечером (21:00)", Data: "notification_time:21:00"}},
			},
		}

		return c.Send(questionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
	} else {
		// If notifications are disabled, store default time and proceed to CEFR test
		if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataNotificationTime, "10:00"); err != nil {
			s.logger.Error("Failed to store default notification time", zap.Error(err))
		}

		return s.proceedToVocabularyTest(ctx, c, userID)
	}
}

// HandleNotificationTimeCallback handles notification time question callback
func (s *HandlerService) HandleNotificationTimeCallback(ctx context.Context, c tele.Context, userID int64, answer string) error {
	s.logger.Info("User answered notification time question", zap.Int64("user_id", userID), zap.String("answer", answer))

	// Store notification time answer
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataNotificationTime, answer); err != nil {
		s.logger.Error("Failed to store notification time answer", zap.Error(err))
	}

	return s.proceedToVocabularyTest(ctx, c, userID)
}

// proceedToVocabularyTest transitions to CEFR vocabulary test
func (s *HandlerService) proceedToVocabularyTest(ctx context.Context, c tele.Context, userID int64) error {
	// Transition to vocabulary test
	if err := s.stateManager.SetState(ctx, userID, fsm.StateVocabularyTest); err != nil {
		s.logger.Error("Failed to set vocabulary test state", zap.Error(err))
		return err
	}

	// Send completion message
	completionText := "🎉 *Отлично!*\n\n" +
		"Теперь давай определим твой точный уровень с помощью короткого теста.\n\n" +
		"Это поможет подобрать идеальные слова для изучения!"

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Начать тест", Data: "test:start"}},
			{{Text: "Пропустить тест", Data: "test:skip"}},
		},
	}

	return c.Send(completionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}
