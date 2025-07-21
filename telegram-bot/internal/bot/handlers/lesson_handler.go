package handlers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/api"
	"telegram-bot/internal/bot/fsm"
)

// HandleLearnCommand handles the /learn command with new learning flow
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

	// Start new learning flow
	return s.HandleNewLearningStart(ctx, c, userID, currentState)
}

// HandleLessonCommand handles the /lesson command (same as /learn for quick testing)
func (s *HandlerService) HandleLessonCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// For quick testing, use the same logic as /learn
	return s.HandleLearnCommand(ctx, c, userID, currentState)
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
	s.logger.Info("HandleLessonStartCallback called", zap.Int64("user_id", userID), zap.String("current_state", string(currentState)))

	// For users in welcome state, first transition to start state
	if currentState == fsm.StateWelcome {
		if err := s.stateManager.SetState(ctx, userID, fsm.StateStart); err != nil {
			s.logger.Error("Failed to set start state from welcome", zap.Error(err))
			return err
		}
		currentState = fsm.StateStart
	}

	// Check if user is authenticated and has completed onboarding
	isAuthenticated, hasCompletedOnboarding, err := s.GetUserAuthenticationStatus(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user authentication status", zap.Error(err))
		return err
	}

	if !isAuthenticated {
		// User is not authenticated, redirect to authentication
		return c.Send("🔐 Для начала уроков необходимо войти в аккаунт. Используйте /start для регистрации или входа.")
	}

	if !hasCompletedOnboarding {
		// User is authenticated but hasn't completed onboarding
		userProgress, err := s.GetUserProgress(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to get user progress", zap.Error(err))
			return err
		}

		if userProgress.CEFRLevel == "" {
			// User hasn't set CEFR level, redirect to onboarding
			return c.Send("📚 Сначала необходимо завершить настройку профиля. Используйте /start для продолжения.")
		}
	}

	// User is ready for lessons, start the lesson flow
	return s.HandleLearnCommand(ctx, c, userID, currentState)
}

// HandleLessonLaterCallback handles lesson later callback
func (s *HandlerService) HandleLessonLaterCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Урок отложен.")
}

// HandleTestSkipCallback handles test skip callback
func (s *HandlerService) HandleTestSkipCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Get the confidence level from questionnaire to determine CEFR level
	confidenceLevel, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataConfidence)
	if err != nil {
		s.logger.Error("Failed to get confidence level", zap.Error(err))
		// Default to beginner if we can't get confidence level
		confidenceLevel = "beginner"
	}

	// Map confidence level to CEFR level
	cefrLevel := s.mapConfidenceToCEFR(confidenceLevel.(string))

	// Set the CEFR level based on user's self-assessment
	if err := s.stateManager.SetState(ctx, userID, fsm.StateCEFRTestResult); err != nil {
		s.logger.Error("Failed to set CEFR test result state", zap.Error(err))
		return err
	}

	// Store the determined CEFR level
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataCEFRTest, cefrLevel); err != nil {
		s.logger.Error("Failed to store CEFR level", zap.Error(err))
	}

	// Save all preferences to the backend (if user is authenticated)
	token, err := s.stateManager.GetJWTToken(ctx, userID)
	if err == nil {
		// User is authenticated, build complete preferences from questionnaire answers
		preferences, err := s.buildCompletePreferencesFromQuestionnaire(ctx, userID, cefrLevel)
		if err != nil {
			s.logger.Error("Failed to build complete preferences", zap.Error(err))
			// Fallback to just CEFR level
			preferences = &api.UpdatePreferenceRequest{
				CEFRLevel: &cefrLevel,
			}
		}

		// Send thinking message
		thinkingMsg, err := s.sendThinkingMessage(ctx, c, userID, "Сохраняю настройки")
		if err != nil {
			s.logger.Error("Failed to send thinking message", zap.Error(err))
			// Continue without thinking message if it fails
		}

		if _, err := s.apiClient.UpdateUserPreferences(ctx, token, preferences); err != nil {
			// Delete thinking message if it was sent
			if thinkingMsg != nil {
				if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
					s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
				}
			}
			s.logger.Error("Failed to update user preferences", zap.Error(err))
		} else {
			// Delete thinking message if it was sent
			if thinkingMsg != nil {
				if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
					s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
				}
			}
			s.logger.Info("Successfully updated complete user preferences",
				zap.Int64("user_id", userID),
				zap.String("cefr_level", cefrLevel),
				zap.Int("words_per_day", *preferences.WordsPerDay),
				zap.Bool("notifications", *preferences.Notifications),
				zap.String("goal", *preferences.Goal))
		}
	}

	// Send completion message with assigned level and preferences summary
	wordsPerDay := 10 // Default
	if wordsPerDayData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataWordsPerDay); err == nil {
		// Handle JSON number unmarshaling (numbers come back as float64)
		switch v := wordsPerDayData.(type) {
		case int:
			wordsPerDay = v
		case float64:
			wordsPerDay = int(v)
		}
	}

	notifications := false // Default
	if notificationsData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataNotifications); err == nil {
		if notificationsValue, ok := notificationsData.(bool); ok {
			notifications = notificationsValue
		}
	}

	notificationStatus := "отключены"
	if notifications {
		notificationStatus = "включены"
	}

	completionText := fmt.Sprintf(
		"🎉 *Добро пожаловать в Fluently!*\n\n"+
			"📊 **Ваши настройки:**\n"+
			"• Уровень: *%s*\n"+
			"• Слов в день: *%d*\n"+
			"• Уведомления: *%s*\n\n"+
			"Настройка завершена! Теперь ты можешь начать изучение.\n\n"+
			"Используй /learn чтобы начать свой первый урок!",
		cefrLevel,
		wordsPerDay,
		notificationStatus,
	)

	// Create main menu keyboard
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Начать изучение", Data: "lesson:start"}},
			{{Text: "Настройки", Data: "menu:settings"}},
		},
	}

	// Send the completion message and transition to start state
	if err := c.Send(completionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard); err != nil {
		return err
	}

	// Clear all questionnaire temp data now that onboarding is complete
	tempDataTypes := []fsm.TempDataType{
		fsm.TempDataGoal,
		fsm.TempDataExperience,
		fsm.TempDataConfidence,
		fsm.TempDataWordsPerDay,
		fsm.TempDataNotifications,
		fsm.TempDataNotificationTime,
	}

	for _, dataType := range tempDataTypes {
		if err := s.stateManager.ClearTempData(ctx, userID, dataType); err != nil {
			s.logger.Error("Failed to clear questionnaire temp data",
				zap.String("data_type", string(dataType)), zap.Error(err))
		}
	}

	// Set final state to start (onboarding complete)
	return s.stateManager.SetState(ctx, userID, fsm.StateStart)
}

// mapConfidenceToCEFR maps user confidence level to CEFR level
func (s *HandlerService) mapConfidenceToCEFR(confidenceLevel string) string {
	switch confidenceLevel {
	case "beginner":
		return "A1"
	case "elementary":
		return "A2"
	case "intermediate":
		return "B1"
	case "advanced":
		return "C1"
	default:
		return "A1" // Default to beginner
	}
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

// buildCompletePreferencesFromQuestionnaire collects all questionnaire answers and builds a complete preferences request
func (s *HandlerService) buildCompletePreferencesFromQuestionnaire(ctx context.Context, userID int64, cefrLevel string) (*api.UpdatePreferenceRequest, error) {
	preferences := &api.UpdatePreferenceRequest{
		CEFRLevel: &cefrLevel,
	}

	goal := "Изучение английского" // Default goal
	factEveryday := false
	subscribed := false
	avatarImageURL := ""
	wordsPerDay := 10 // Default words per day
	notifications := false
	var notificationTime *time.Time

	// Get goal from temp data
	if goalData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataGoal); err == nil {
		if goalValue, ok := goalData.(string); ok {
			// Map goal values to more descriptive text
			switch goalValue {
			case "work":
				goal = "Работа и карьера"
			case "travel":
				goal = "Путешествия"
			case "education":
				goal = "Образование"
			case "communication":
				goal = "Общение"
			default:
				goal = "Изучение английского"
			}
		}
	}

	// Get words per day from temp data
	if wordsPerDayData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataWordsPerDay); err == nil {
		// Handle JSON number unmarshaling (numbers come back as float64)
		switch v := wordsPerDayData.(type) {
		case int:
			wordsPerDay = v
		case float64:
			wordsPerDay = int(v)
		}
	}

	// Get notifications setting from temp data
	if notificationsData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataNotifications); err == nil {
		if notificationsValue, ok := notificationsData.(bool); ok {
			notifications = notificationsValue
		}
	}

	// Get notification time from temp data
	if notificationTimeData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataNotificationTime); err == nil {
		if notificationTimeStr, ok := notificationTimeData.(string); ok && notificationTimeStr != "" {
			// Parse the time string to time.Time using flexible parser
			if parsedTime, err := ParseTimeToTime(notificationTimeStr); err == nil {
				notificationTime = parsedTime
			} else {
				s.logger.Warn("Failed to parse notification time from questionnaire",
					zap.String("notification_time", notificationTimeStr),
					zap.Error(err))
			}
		}
	}

	// Assign all values to preferences
	preferences.Goal = &goal
	preferences.FactEveryday = &factEveryday
	preferences.Subscribed = &subscribed
	preferences.AvatarImageURL = &avatarImageURL
	preferences.WordsPerDay = &wordsPerDay
	preferences.Notifications = &notifications
	if notificationTime != nil {
		preferences.NotificationAt = notificationTime
	}

	return preferences, nil
}
