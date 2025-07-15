package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
)

// HandleLessonCallback handles lesson-related callbacks
func (s *HandlerService) HandleLessonCallback(ctx context.Context, c tele.Context, userID int64, action string) error {
	switch action {
	case "start_word_set":
		return s.HandleStartWordSet(ctx, c, userID, fsm.StateShowingWordSet)
	case "continue":
		// Resume lesson from current state
		return s.HandleStartWordSet(ctx, c, userID, fsm.StateShowingWordSet)
	case "restart":
		// Clear progress and start new
		if err := s.stateManager.ClearLessonProgress(ctx, userID); err != nil {
			s.logger.Error("Failed to clear lesson progress", zap.Error(err))
		}
		return s.HandleNewLearningStart(ctx, c, userID, fsm.StateStart)
	case "ready_exercises":
		return s.HandleReadyForExercises(ctx, c, userID)
	case "start_exercises":
		return s.HandleStartExercises(ctx, c, userID)
	case "stats":
		return s.handleLessonStats(ctx, c, userID)
	case "final_stats":
		return s.handleFinalStats(ctx, c, userID)
	case "new":
		return s.HandleNewLearningStart(ctx, c, userID, fsm.StateStart)
	default:
		// Handle show_word callbacks with index
		if strings.HasPrefix(action, "show_word:") {
			indexStr := strings.TrimPrefix(action, "show_word:")
			wordIndex, err := strconv.Atoi(indexStr)
			if err != nil {
				return err
			}
			return s.HandleShowWord(ctx, c, userID, wordIndex)
		}
		return c.Send("❌ Неизвестная команда")
	}
}

// HandleExerciseCallback handles exercise-related callbacks
func (s *HandlerService) HandleExerciseCallback(ctx context.Context, c tele.Context, userID int64, action string) error {
	switch {
	case action == "next":
		return s.HandleExerciseNext(ctx, c, userID)
	case action == "skip":
		return s.HandleSkipExercise(ctx, c, userID)
	case action == "hint":
		return s.HandleExerciseHint(ctx, c, userID)
	case strings.HasPrefix(action, "pick_option:"):
		// Format: pick_option:index:option
		parts := strings.Split(action, ":")
		if len(parts) != 3 {
			return c.Send("❌ Неверный формат ответа")
		}

		optionIndex, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		selectedOption := parts[2]
		return s.HandlePickOptionAnswer(ctx, c, userID, optionIndex, selectedOption)
	case strings.HasPrefix(action, "translate_option:"):
		// Format: translate_option:index:option
		parts := strings.Split(action, ":")
		if len(parts) != 3 {
			return c.Send("❌ Неверный формат ответа")
		}

		optionIndex, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}

		selectedOption := parts[2]
		return s.HandlePickOptionAnswer(ctx, c, userID, optionIndex, selectedOption)
	default:
		return c.Send("❌ Неизвестная команда упражнения")
	}
}

// HandleAuthCallback handles authentication-related callbacks
func (s *HandlerService) HandleAuthCallback(ctx context.Context, c tele.Context, userID int64, action string) error {
	switch action {
	case "check_link":
		return s.handleCheckLinkStatus(ctx, c, userID)
	default:
		return c.Send("❌ Неизвестная команда авторизации")
	}
}

// HandleVoiceCallback handles voice-related callbacks
func (s *HandlerService) HandleVoiceCallback(ctx context.Context, c tele.Context, userID int64, action string) error {
	if strings.HasPrefix(action, "repeat:") {
		word := strings.TrimPrefix(action, "repeat:")
		return s.sendWordVoiceMessage(ctx, c, word)
	}

	return c.Send("❌ Неизвестная голосовая команда")
}

// handleLessonStats shows current lesson statistics
func (s *HandlerService) handleLessonStats(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil || progress == nil {
		return c.Send("❌ Нет активного урока")
	}

	learnedCount := progress.LearnedCount
	targetCount := progress.LessonData.Lesson.WordsPerLesson
	duration := s.formatDuration(time.Since(progress.StartTime))

	correctWords := 0
	totalWords := len(progress.WordsLearned)
	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 {
			correctWords++
		}
	}

	var accuracy float64
	if totalWords > 0 {
		accuracy = float64(correctWords) / float64(totalWords) * 100
	}

	statsText := fmt.Sprintf(
		"📊 *Статистика урока*\n\n"+
			"🎯 Прогресс: %d/%d слов\n"+
			"✅ Правильных ответов: %d из %d\n"+
			"📈 Точность: %.1f%%\n"+
			"⏱ Время: %s\n"+
			"📚 Текущий набор: #%d",
		learnedCount,
		targetCount,
		correctWords,
		totalWords,
		accuracy,
		duration,
		progress.CurrentSetIndex,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "▶️ Продолжить урок", Data: "lesson:continue"},
			},
		},
	}

	return c.Send(statsText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// handleFinalStats shows final lesson statistics
func (s *HandlerService) handleFinalStats(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil || progress == nil {
		return c.Send("❌ Нет данных о завершенном уроке")
	}

	totalWords := len(progress.WordsLearned)
	correctWords := 0
	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 {
			correctWords++
		}
	}

	duration := s.formatDuration(time.Since(progress.StartTime))
	accuracy := float64(correctWords) / float64(totalWords) * 100

	// Show learned words
	var learnedWordsText strings.Builder
	learnedWordsText.WriteString("📚 *Изученные слова:*\n")

	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 {
			learnedWordsText.WriteString(fmt.Sprintf("✅ %s\n", wordProgress.Word))
		} else {
			learnedWordsText.WriteString(fmt.Sprintf("❌ %s\n", wordProgress.Word))
		}
	}

	finalStatsText := fmt.Sprintf(
		"🏆 *Финальная статистика*\n\n"+
			"✅ Слов выучено: %d\n"+
			"🎯 Правильных ответов: %d из %d\n"+
			"📈 Точность: %.1f%%\n"+
			"⏱ Время урока: %s\n\n"+
			"%s",
		progress.LearnedCount,
		correctWords,
		totalWords,
		accuracy,
		duration,
		learnedWordsText.String(),
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🚀 Новый урок", Data: "lesson:new"},
				{Text: "🏠 Главное меню", Data: "menu:main"},
			},
		},
	}

	return c.Send(finalStatsText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// handleCheckLinkStatus checks if user's Google account is linked
func (s *HandlerService) handleCheckLinkStatus(ctx context.Context, c tele.Context, userID int64) error {
	linkStatus, err := s.apiClient.CheckLinkStatus(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to check link status", zap.Error(err))
		return c.Send("❌ Не удалось проверить статус связи")
	}

	if !linkStatus.IsLinked {
		return c.Send("🔗 Аккаунт еще не связан. Пожалуйста, завершите процесс авторизации по ссылке выше.")
	}

	// Store JWT token if available
	// Note: You would typically get the JWT token from the link status response
	// For now, we'll assume the token is available from a separate authentication flow

	return c.Send("✅ Аккаунт успешно связан! Теперь вы можете начать изучение.\n\nИспользуйте /learn для начала урока.")
}

// handleHelpAuth provides help information about authentication
func (s *HandlerService) handleHelpAuth(ctx context.Context, c tele.Context, userID int64) error {
	helpText := "❓ *Помощь по авторизации*\n\n" +
		"Для использования бота необходимо связать ваш аккаунт Telegram с аккаунтом Google.\n\n" +
		"*Шаги:*\n" +
		"1. Нажмите на ссылку авторизации\n" +
		"2. Войдите в свой аккаунт Google\n" +
		"3. Разрешите доступ приложению\n" +
		"4. Вернитесь в бота и нажмите \"Проверить связь\"\n\n" +
		"*Безопасность:*\n" +
		"• Мы не храним ваши пароли\n" +
		"• Доступ используется только для обучения\n" +
		"• Вы можете отозвать доступ в любое время"

	return c.Send(helpText, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}
