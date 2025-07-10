package handlers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
)

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

// HandleLessonStartCallback handles lesson start callback
func (s *HandlerService) HandleLessonStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Начинаем урок...")
}

// HandleLessonLaterCallback handles lesson later callback
func (s *HandlerService) HandleLessonLaterCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Урок отложен.")
}

// HandleTestStartCallback handles test start callback
func (s *HandlerService) HandleTestStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Начинаем тест...")
}

// HandleWaitingForTranslationMessage handles translation waiting state
func (s *HandlerService) HandleWaitingForTranslationMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, предоставьте перевод.")
}

// HandleWaitingForAudioMessage handles audio waiting state
func (s *HandlerService) HandleWaitingForAudioMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, предоставьте аудио ответ.")
}

// HandleAudioExerciseResponse handles audio exercise responses
func (s *HandlerService) HandleAudioExerciseResponse(ctx context.Context, c tele.Context, userID int64, voice interface{}) error {
	return c.Send("Получен ответ на аудио упражнение.")
}

// HandleLearnMenuCallback handles learn menu callback
func (s *HandlerService) HandleLearnMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Меню обучения...")
}
