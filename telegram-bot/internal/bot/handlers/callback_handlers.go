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
	case "start":
		// Get current state for HandleLessonStartCallback
		currentState, err := s.stateManager.GetState(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to get current state", zap.Error(err))
			return err
		}
		return s.HandleLessonStartCallback(ctx, c, userID, currentState)
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
		// Handle already_know callbacks with index
		if strings.HasPrefix(action, "already_know:") {
			indexStr := strings.TrimPrefix(action, "already_know:")
			wordIndex, err := strconv.Atoi(indexStr)
			if err != nil {
				return err
			}
			return s.HandleWordAlreadyKnown(ctx, c, userID, wordIndex)
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
	case "existing_user":
		return s.HandleExistingUserAuth(ctx, c, userID)
	case "new_user":
		return s.HandleNewUserAuth(ctx, c, userID)
	case "register":
		return s.HandleRegisterAuth(ctx, c, userID)
	case "check_link":
		return s.handleCheckLinkStatus(ctx, c, userID)
	default:
		return c.Send("❌ Неизвестная команда авторизации")
	}
}

// HandleExistingUserAuth handles existing user authentication flow
func (s *HandlerService) HandleExistingUserAuth(ctx context.Context, c tele.Context, userID int64) error {
	// First check if user is already linked
	linkStatus, err := s.apiClient.CheckLinkStatus(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to check link status", zap.Error(err))
		return c.Send("❌ Произошла ошибка при проверке статуса связи. Попробуйте позже.")
	}

	if linkStatus.IsLinked {
		// User is already linked, proceed directly to authentication check
		s.logger.Info("User is already linked, proceeding to check authentication", zap.Int64("user_id", userID))
		return s.handleCheckLinkStatus(ctx, c, userID)
	}

	// Create link token for existing user who is not yet linked
	linkResponse, err := s.apiClient.CreateLinkToken(ctx, userID)
	if err != nil {
		// Handle the case where the account is already linked (409 error)
		if strings.Contains(err.Error(), "already linked") || strings.Contains(err.Error(), "409") {
			s.logger.Info("Account already linked, proceeding to check status", zap.Int64("user_id", userID))
			return s.handleCheckLinkStatus(ctx, c, userID)
		}
		s.logger.Error("Failed to create link token for existing user", zap.Error(err))
		return c.Send("❌ Произошла ошибка. Попробуйте позже.")
	}

	// Store linking data
	err = s.stateManager.StoreUserLinkingData(ctx, userID, linkResponse.Token, time.Hour)
	if err != nil {
		s.logger.Error("Failed to store linking data", zap.Error(err))
	}

	authText := fmt.Sprintf(
		"🔐 *Авторизация для существующего пользователя*\n\n"+
			"Для входа в ваш аккаунт необходимо пройти авторизацию через Google.\n\n"+
			"🔗 *Ссылка для авторизации:*\n[Нажмите здесь для входа](%s)\n\n"+
			"После авторизации вернитесь и нажмите \"Проверить связь\".",
		linkResponse.LinkURL,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🔄 Проверить связь", Data: "auth:check_link"},
				{Text: "❓ Помощь", Data: "help:auth"},
			},
		},
	}

	return c.Send(authText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleNewUserAuth handles new user authentication flow with onboarding
func (s *HandlerService) HandleNewUserAuth(ctx context.Context, c tele.Context, userID int64) error {
	// Check current state first
	currentState, err := s.stateManager.GetState(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get current state", zap.Error(err))
		return err
	}

	// Only set state to welcome if not already there
	if currentState != fsm.StateWelcome {
		if err := s.stateManager.SetState(ctx, userID, fsm.StateWelcome); err != nil {
			s.logger.Error("Failed to set welcome state", zap.Error(err))
			return err
		}
	}

	// Start onboarding process
	return s.HandleOnboardingStartCallback(ctx, c, userID, fsm.StateWelcome)
}

// HandleRegisterAuth handles user registration after CEFR test completion
func (s *HandlerService) HandleRegisterAuth(ctx context.Context, c tele.Context, userID int64) error {
	// First check if user is already linked
	linkStatus, err := s.apiClient.CheckLinkStatus(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to check link status", zap.Error(err))
		return c.Send("❌ Произошла ошибка при проверке статуса связи. Попробуйте позже.")
	}

	if linkStatus.IsLinked {
		// User is already linked, proceed directly to authentication check
		s.logger.Info("User is already linked during registration, proceeding to check authentication", zap.Int64("user_id", userID))
		return s.handleCheckLinkStatus(ctx, c, userID)
	}

	// Create link token for new user registration
	linkResponse, err := s.apiClient.CreateLinkToken(ctx, userID)
	if err != nil {
		// Handle the case where the account is already linked (409 error)
		if strings.Contains(err.Error(), "already linked") || strings.Contains(err.Error(), "409") {
			s.logger.Info("Account already linked during registration, proceeding to check status", zap.Int64("user_id", userID))
			return s.handleCheckLinkStatus(ctx, c, userID)
		}
		s.logger.Error("Failed to create link token for registration", zap.Error(err))
		return c.Send("❌ Произошла ошибка. Попробуйте позже.")
	}

	// Store linking data
	err = s.stateManager.StoreUserLinkingData(ctx, userID, linkResponse.Token, time.Hour)
	if err != nil {
		s.logger.Error("Failed to store linking data", zap.Error(err))
	}

	authText := fmt.Sprintf(
		"🔐 *Создание аккаунта*\n\n"+
			"Для сохранения прогресса создадим аккаунт через Google.\n\n"+
			"🎯 **Преимущества аккаунта:**\n"+
			"• Сохранение прогресса на всех устройствах\n"+
			"• Персональная статистика\n"+
			"• Синхронизация с веб-версией\n\n"+
			"🔗 *Ссылка для регистрации:*\n[Нажмите здесь для создания аккаунта](%s)\n\n"+
			"После регистрации вернитесь и нажмите \"Проверить связь\".",
		linkResponse.LinkURL,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🔄 Проверить связь", Data: "auth:check_link"},
				{Text: "⏭ Пропустить", Data: "lesson:start"},
			},
		},
	}

	return c.Send(authText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleStatsCallback handles statistics-related callbacks
func (s *HandlerService) HandleStatsCallback(ctx context.Context, c tele.Context, userID int64, action string) error {
	switch action {
	case "show":
		return s.HandleStatsCommand(ctx, c, userID, fsm.StateStart)
	case "overall":
		return s.HandleStatsCommand(ctx, c, userID, fsm.StateStart)
	default:
		return c.Send("❌ Неизвестная команда статистики")
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

	// Calculate statistics - exclude "already known" words from the count
	wellAnsweredWords := 0
	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 {
			// Check if this word was marked as "already known"
			if !wordProgress.AlreadyKnown {
				// This is a newly learned word, not "already known"
				wellAnsweredWords++
			}
		}
	}

	var accuracy float64
	newlyLearnedCount := progress.LearnedCount - progress.AlreadyKnownCount
	if newlyLearnedCount > 0 {
		accuracy = float64(wellAnsweredWords) / float64(newlyLearnedCount) * 100
	}

	statsText := fmt.Sprintf(
		"📊 *Статистика урока*\n\n"+
			"🎯 Прогресс: %d/%d слов\n"+
			"✅ Правильно: %d из %d\n"+
			"📈 Точность: %.1f%%\n"+
			"⏱ Время: %s\n",
		learnedCount,
		targetCount,
		wellAnsweredWords,
		targetCount, // Show correct answers vs target words (user_preferences.words_per_day)
		accuracy,
		duration,
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

	// Calculate statistics - exclude "already known" words from the count
	correctWords := 0
	newlyLearnedCorrectWords := 0

	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 {
			correctWords++
			// Check if this word was marked as "already known"
			if !wordProgress.AlreadyKnown {
				// This is a newly learned word, not "already known"
				newlyLearnedCorrectWords++
			}
		}
	}

	duration := s.formatDuration(time.Since(progress.StartTime))
	// Calculate accuracy based on newly learned words only
	newlyLearnedCount := progress.LearnedCount - progress.AlreadyKnownCount
	accuracy := float64(newlyLearnedCorrectWords) / float64(newlyLearnedCount) * 100

	// Show learned words with translations
	var learnedWordsText strings.Builder
	learnedWordsText.WriteString(fmt.Sprintf("📚 *За урок выучено %d слов:*\n\n", newlyLearnedCorrectWords))

	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 && !wordProgress.AlreadyKnown {
			learnedWordsText.WriteString(fmt.Sprintf("#%s - %s\n", wordProgress.Word, wordProgress.Translation))
		}
	}

	// Add information about retry words if any remain
	var retryInfo string
	if len(progress.RetryWords) > 0 {
		retryInfo = fmt.Sprintf("\n🔄 *Слова для повторения:* %d\n", len(progress.RetryWords))
	}

	finalStatsText := fmt.Sprintf(
		"🏆 *Финальная статистика*\n\n"+
			"✅ Слов выучено: %d\n"+
			"💡 Уже знал: %d слов\n"+
			"🎯 Правильно: %d из %d\n"+
			"📈 Точность: %.1f%%\n"+
			"⏱ Время урока: %s%s\n\n"+
			"%s",
		progress.LearnedCount,
		progress.AlreadyKnownCount,
		newlyLearnedCorrectWords,
		progress.LessonData.Lesson.WordsPerLesson,
		accuracy,
		duration,
		retryInfo,
		learnedWordsText.String(),
	)

	// Send progress to backend
	token, err := s.stateManager.GetJWTToken(ctx, userID)
	if err == nil {
		err = s.apiClient.SendLessonProgress(ctx, token, progress.WordsLearned, progress.BadlyAnsweredWords)
		if err != nil {
			s.logger.Error("Failed to send lesson progress to backend", zap.Error(err))
		}
	}

	// Clear lesson progress
	err = s.stateManager.ClearLessonProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to clear lesson progress", zap.Error(err))
	}

	// Reset state to start
	if err := s.stateManager.SetState(ctx, userID, fsm.StateStart); err != nil {
		s.logger.Error("Failed to reset state", zap.Error(err))
	}

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
	s.logger.Info("Checking link status", zap.Int64("user_id", userID))

	linkStatus, err := s.apiClient.CheckLinkStatus(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to check link status", zap.Int64("user_id", userID), zap.Error(err))
		return c.Send("❌ Не удалось проверить статус связи")
	}

	s.logger.Info("Link status result", zap.Int64("user_id", userID), zap.Bool("is_linked", linkStatus.IsLinked))

	if !linkStatus.IsLinked {
		return c.Send("🔗 Аккаунт еще не связан. Пожалуйста, завершите процесс авторизации по ссылке выше.")
	}

	// Account is linked, now we need to get JWT tokens
	s.logger.Info("Account is linked, attempting to get JWT tokens", zap.Int64("user_id", userID))

	// Get JWT tokens from the backend after successful linking
	jwtTokens, err := s.apiClient.GetJWTTokens(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get JWT tokens after linking", zap.Int64("user_id", userID), zap.Error(err))
		// For now, let's still proceed - this might be an API issue
		// In a real implementation, you'd handle this error appropriately
	} else {
		// Store JWT tokens
		s.logger.Info("Storing JWT tokens", zap.Int64("user_id", userID), zap.String("access_token_length", fmt.Sprintf("%d", len(jwtTokens.AccessToken))))

		// Store access token with 24 hour expiration (adjust as needed)
		err = s.stateManager.StoreJWTToken(ctx, userID, jwtTokens.AccessToken, 24*time.Hour)
		if err != nil {
			s.logger.Error("Failed to store JWT access token", zap.Int64("user_id", userID), zap.Error(err))
		} else {
			s.logger.Info("Successfully stored JWT access token", zap.Int64("user_id", userID))
		}

		// Also store using the new format if available
		if jwtTokens.RefreshToken != "" {
			err = s.stateManager.StoreJWTTokens(ctx, userID, jwtTokens.AccessToken, jwtTokens.RefreshToken, 24*time.Hour, 30*24*time.Hour)
			if err != nil {
				s.logger.Error("Failed to store JWT tokens", zap.Int64("user_id", userID), zap.Error(err))
			} else {
				s.logger.Info("Successfully stored JWT tokens", zap.Int64("user_id", userID))
			}
		}
	}

	// Clear linking data
	err = s.stateManager.ClearUserLinkingData(ctx, userID)
	if err != nil {
		s.logger.Warn("Failed to clear linking data", zap.Int64("user_id", userID), zap.Error(err))
	}

	// Check if user has completed onboarding
	isAuthenticated, hasCompletedOnboarding, err := s.GetUserAuthenticationStatus(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user authentication status", zap.Int64("user_id", userID), zap.Error(err))
		return err
	}

	s.logger.Info("Post-linking authentication status", zap.Int64("user_id", userID), zap.Bool("is_authenticated", isAuthenticated), zap.Bool("has_completed_onboarding", hasCompletedOnboarding))

	if isAuthenticated && hasCompletedOnboarding {
		// User is fully set up - show main menu
		successText := "✅ *Аккаунт успешно связан!*\n\n🎉 Добро пожаловать обратно! Теперь вы можете продолжить изучение."

		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{
					{Text: "🚀 Начать урок", Data: "lesson:start"},
					{Text: "🏠 Главное меню", Data: "menu:main"},
				},
			},
		}

		return c.Send(successText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
	} else {
		// User needs to complete onboarding
		successText := "✅ *Аккаунт успешно связан!*\n\n📋 Теперь давайте завершим настройку вашего профиля для эффективного обучения."

		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{
					{Text: "✅ Завершить настройку", Data: "questionnaire:start"},
					{Text: "🚀 Пропустить в урок", Data: "lesson:start"},
				},
			},
		}

		return c.Send(successText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
	}
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
