package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
	"telegram-bot/internal/domain"
)

// SettingsMessageID stores the message ID for the settings message to update it
const SettingsMessageID = "settings_message_id"

// HandleSettingsCommand handles the /settings command
func (s *HandlerService) HandleSettingsCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Get current user progress for settings
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Set user state to settings
	if err := s.SetStateIfDifferent(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	// Send initial settings message
	return s.sendSettingsMessage(ctx, c, userID, userProgress, "")
}

// sendSettingsMessage sends or updates the settings message
func (s *HandlerService) sendSettingsMessage(ctx context.Context, c tele.Context, userID int64, userProgress *domain.UserProgress, statusMessage string) error {
	// Create settings message
	settingsText := "⚙️ *Настройки*\n\n" +
		fmt.Sprintf("🔤 Уровень CEFR: *%s*\n", formatCEFRLevel(userProgress.CEFRLevel)) +
		fmt.Sprintf("📚 Слов в день: *%d*\n", userProgress.WordsPerDay) +
		fmt.Sprintf("🔔 Уведомления: *%s*\n", formatNotificationTime(userProgress.NotificationTime))

	if statusMessage != "" {
		settingsText += "\n" + statusMessage
	}

	settingsText += "\n\nВыберите настройку для изменения:"

	// Create settings keyboard based on current state
	var keyboard *tele.ReplyMarkup

	// Check if we're in a specific settings sub-state
	currentState, err := s.stateManager.GetState(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get current state", zap.Error(err))
		currentState = fsm.StateSettings
	}

	switch currentState {
	case fsm.StateSettingsWordsPerDay, fsm.StateSettingsWordsPerDayInput:
		// Show word count options
		keyboard = &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{
					{Text: "5 слов", Data: "settings:words:5"},
					{Text: "10 слов", Data: "settings:words:10"},
					{Text: "15 слов", Data: "settings:words:15"},
				},
				{
					{Text: "20 слов", Data: "settings:words:20"},
					{Text: "25 слов", Data: "settings:words:25"},
					{Text: "30 слов", Data: "settings:words:30"},
				},
				{{Text: "Ввести вручную", Data: "settings:words:custom"}},
				{{Text: "Отмена", Data: "settings:back"}},
			},
		}
	case fsm.StateSettingsNotifications:
		// Show time options
		keyboard = &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{
					{Text: "08:00", Data: "settings:time:08:00"},
					{Text: "09:00", Data: "settings:time:09:00"},
					{Text: "10:00", Data: "settings:time:10:00"},
				},
				{
					{Text: "12:00", Data: "settings:time:12:00"},
					{Text: "15:00", Data: "settings:time:15:00"},
					{Text: "18:00", Data: "settings:time:18:00"},
				},
				{
					{Text: "20:00", Data: "settings:time:20:00"},
					{Text: "21:00", Data: "settings:time:21:00"},
					{Text: "22:00", Data: "settings:time:22:00"},
				},
				{{Text: "Ввести вручную", Data: "settings:time:custom"}},
				{{Text: "Отключить", Data: "settings:time:disabled"}},
				{{Text: "Отмена", Data: "settings:back"}},
			},
		}
	case fsm.StateSettingsCEFRLevel:
		// Show CEFR level options
		keyboard = &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{
					{Text: "A1 - Начинающий", Data: "settings:cefr:A1"},
					{Text: "A2 - Элементарный", Data: "settings:cefr:A2"},
				},
				{
					{Text: "B1 - Средний", Data: "settings:cefr:B1"},
					{Text: "B2 - Выше среднего", Data: "settings:cefr:B2"},
				},
				{
					{Text: "C1 - Продвинутый", Data: "settings:cefr:C1"},
					{Text: "C2 - В совершенстве", Data: "settings:cefr:C2"},
				},
				{{Text: "Пройти тест", Data: "settings:cefr:test"}},
				{{Text: "Ввести вручную", Data: "settings:cefr:custom"}},
				{{Text: "Отмена", Data: "settings:back"}},
			},
		}
	default:
		// Default settings keyboard
		keyboard = &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{{Text: "🔤 Уровень CEFR", Data: "settings:cefr_level"}},
				{{Text: "📚 Слов в день", Data: "settings:words_per_day"}},
				{{Text: "🔔 Уведомления", Data: "settings:notifications"}},
				{{Text: "Назад в главное меню", Data: "menu:main"}},
			},
		}
	}

	// Check if we have a stored message ID to edit
	if messageID, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataSettings); err == nil {
		if msgID, ok := messageID.(int); ok {
			// Try to edit existing message
			if _, err := c.Bot().Edit(&tele.Message{ID: msgID, Chat: c.Message().Chat}, settingsText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard); err == nil {
				return nil
			}
		}
	}

	// Send new message if editing failed or no stored message ID
	if err := c.Send(settingsText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard); err != nil {
		return err
	}

	// For now, we'll just send a new message each time since getting message ID is complex
	// In a production environment, you might want to implement a more sophisticated approach
	return nil
}

// HandleSettingsWordsPerDayCallback handles words per day settings callback
func (s *HandlerService) HandleSettingsWordsPerDayCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Set state to words per day selection
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettingsWordsPerDay); err != nil {
		s.logger.Error("Failed to set words per day state", zap.Error(err))
		return err
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	statusText := fmt.Sprintf("📚 *Слов в день*\n\nТекущее значение: *%d* слов\n\nВыберите новое количество или введите вручную:", userProgress.WordsPerDay)

	// Update the settings message
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusText)
}

// HandleSettingsNotificationsCallback handles notifications settings callback
func (s *HandlerService) HandleSettingsNotificationsCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Set state to notifications settings
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettingsNotifications); err != nil {
		s.logger.Error("Failed to set notifications state", zap.Error(err))
		return err
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	currentTime := formatNotificationTime(userProgress.NotificationTime)
	statusText := fmt.Sprintf("🔔 *Уведомления*\n\nТекущее время: *%s*\n\nВыберите время для ежедневных уведомлений:", currentTime)

	// Update the settings message
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusText)
}

// HandleSettingsCEFRLevelCallback handles CEFR level settings callback
func (s *HandlerService) HandleSettingsCEFRLevelCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Set state to CEFR level settings
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettingsCEFRLevel); err != nil {
		s.logger.Error("Failed to set CEFR level state", zap.Error(err))
		return err
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	currentLevel := formatCEFRLevel(userProgress.CEFRLevel)
	statusText := fmt.Sprintf("🔤 *Уровень CEFR*\n\nТекущий уровень: *%s*\n\nВыберите уровень или пройдите тест для определения:", currentLevel)

	// Update the settings message
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusText)
}

// HandleSettingsWordsPerDayInputMessage handles words per day input messages
func (s *HandlerService) HandleSettingsWordsPerDayInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	text := strings.TrimSpace(c.Text())

	// Parse the number
	wordsPerDay, err := strconv.Atoi(text)
	if err != nil || wordsPerDay < 1 || wordsPerDay > 100 {
		return c.Send("❌ Пожалуйста, введите число от 1 до 100.")
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Update words per day
	userProgress.WordsPerDay = wordsPerDay

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Return to settings with success message
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	statusMessage := fmt.Sprintf("✅ Количество слов в день изменено на *%d*", wordsPerDay)
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
}

// HandleSettingsTimeInputMessage handles time input messages
func (s *HandlerService) HandleSettingsTimeInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	text := strings.TrimSpace(c.Text())

	// Validate time format (HH:MM)
	if !isValidTimeFormat(text) {
		return c.Send("❌ Пожалуйста, введите время в формате ЧЧ:ММ (например, 09:30)")
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Update notification time
	userProgress.NotificationTime = text

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Return to settings with success message
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	statusMessage := fmt.Sprintf("✅ Время уведомлений изменено на *%s*", text)
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
}

// HandleSettingsCEFRLevelInputMessage handles CEFR level input messages
func (s *HandlerService) HandleSettingsCEFRLevelInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	text := strings.TrimSpace(c.Text())

	// Validate CEFR level
	validLevels := map[string]bool{"A1": true, "A2": true, "B1": true, "B2": true, "C1": true, "C2": true}
	if !validLevels[text] {
		return c.Send("❌ Неверный уровень CEFR. Используйте A1, A2, B1, B2, C1 или C2.")
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Update CEFR level
	userProgress.CEFRLevel = text

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Return to settings with success message
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	statusMessage := fmt.Sprintf("✅ Уровень CEFR изменен на *%s*", text)
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
}

// HandleSettingsMenuCallback handles settings menu callback
func (s *HandlerService) HandleSettingsMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// This is the same as the settings command
	return s.HandleSettingsCommand(ctx, c, userID, currentState)
}

// HandleSettingsBackCallback handles the back button in settings
func (s *HandlerService) HandleSettingsBackCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Return to main settings menu
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Clear any temporary data
	s.stateManager.ClearTempData(ctx, userID, SettingsMessageID)

	// Send clean settings message
	return s.sendSettingsMessage(ctx, c, userID, userProgress, "")
}

// HandleSettingsWordsCallback handles words per day selection callbacks
func (s *HandlerService) HandleSettingsWordsCallback(ctx context.Context, c tele.Context, userID int64, data string) error {
	parts := strings.Split(data, ":")
	if len(parts) != 3 {
		return c.Send("❌ Неверный формат данных.")
	}

	value := parts[2]

	if value == "custom" {
		// Set state to input mode (direct transition from settings)
		if err := s.stateManager.SetState(ctx, userID, fsm.StateSettingsWordsPerDayInput); err != nil {
			s.logger.Error("Failed to set words per day input state", zap.Error(err))
			return err
		}
		return c.Send("📝 Введите количество слов в день (от 1 до 100):")
	}

	// Parse the number
	wordsPerDay, err := strconv.Atoi(value)
	if err != nil || wordsPerDay < 1 || wordsPerDay > 100 {
		return c.Send("❌ Неверное количество слов.")
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Update words per day
	userProgress.WordsPerDay = wordsPerDay

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Return to settings with success message
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	statusMessage := fmt.Sprintf("✅ Количество слов в день изменено на *%d*", wordsPerDay)
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
}

// HandleSettingsTimeCallback handles notification time selection callbacks
func (s *HandlerService) HandleSettingsTimeCallback(ctx context.Context, c tele.Context, userID int64, data string) error {
	parts := strings.Split(data, ":")
	if len(parts) != 3 {
		return c.Send("❌ Неверный формат данных.")
	}

	value := parts[2]

	if value == "custom" {
		// Set state to input mode (direct transition from settings)
		if err := s.stateManager.SetState(ctx, userID, fsm.StateSettingsTimeInput); err != nil {
			s.logger.Error("Failed to set time input state", zap.Error(err))
			return err
		}
		return c.Send("📝 Введите время уведомлений в формате ЧЧ:ММ (например, 09:30):")
	}

	if value == "disabled" {
		// Disable notifications
		userProgress, err := s.GetUserProgress(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to get user progress", zap.Error(err))
			return err
		}

		userProgress.NotificationTime = ""

		// Save to backend
		if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
			s.logger.Error("Failed to update user progress", zap.Error(err))
			return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
		}

		// Return to settings with success message
		if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
			s.logger.Error("Failed to set settings state", zap.Error(err))
			return err
		}

		statusMessage := "✅ Уведомления отключены"
		return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
	}

	// Validate time format
	if !isValidTimeFormat(value) {
		return c.Send("❌ Неверный формат времени.")
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Update notification time
	userProgress.NotificationTime = value

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Return to settings with success message
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	statusMessage := fmt.Sprintf("✅ Время уведомлений изменено на *%s*", value)
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
}

// HandleSettingsCEFRCallback handles CEFR level selection callbacks
func (s *HandlerService) HandleSettingsCEFRCallback(ctx context.Context, c tele.Context, userID int64, data string) error {
	parts := strings.Split(data, ":")
	if len(parts) != 3 {
		return c.Send("❌ Неверный формат данных.")
	}

	value := parts[2]

	if value == "test" {
		// Start CEFR test
		if err := s.stateManager.SetState(ctx, userID, fsm.StateVocabularyTest); err != nil {
			s.logger.Error("Failed to set vocabulary test state", zap.Error(err))
			return err
		}
		return s.HandleTestStartCallback(ctx, c, userID, fsm.StateVocabularyTest)
	}

	if value == "custom" {
		// Set state to input mode (direct transition from settings)
		if err := s.stateManager.SetState(ctx, userID, fsm.StateSettingsCEFRLevel); err != nil {
			s.logger.Error("Failed to set CEFR level input state", zap.Error(err))
			return err
		}
		return c.Send("📝 Введите уровень CEFR (A1, A2, B1, B2, C1, C2):")
	}

	// Validate CEFR level
	validLevels := map[string]bool{"A1": true, "A2": true, "B1": true, "B2": true, "C1": true, "C2": true}
	if !validLevels[value] {
		return c.Send("❌ Неверный уровень CEFR. Используйте A1, A2, B1, B2, C1 или C2.")
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Update CEFR level
	userProgress.CEFRLevel = value

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Return to settings with success message
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	statusMessage := fmt.Sprintf("✅ Уровень CEFR изменен на *%s*", value)
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
}

// Helper functions

// formatNotificationTime formats notification time string
func formatNotificationTime(timeStr string) string {
	if timeStr == "" {
		return "Отключены"
	}
	return timeStr
}

// formatCEFRLevel formats CEFR level string
func formatCEFRLevel(level string) string {
	if level == "" {
		return "Не установлен"
	}
	return level
}

// isValidTimeFormat checks if the time string is in HH:MM format
func isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 5 || timeStr[2] != ':' {
		return false
	}

	hour, err1 := strconv.Atoi(timeStr[:2])
	minute, err2 := strconv.Atoi(timeStr[3:])

	if err1 != nil || err2 != nil {
		return false
	}

	return hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59
}
