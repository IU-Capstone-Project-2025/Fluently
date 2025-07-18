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
