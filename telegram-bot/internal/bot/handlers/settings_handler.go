package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/api"
	"telegram-bot/internal/bot/fsm"
	"telegram-bot/internal/domain"
)

// SettingsMessageID stores the message ID for the settings message to update it
const SettingsMessageID = "settings_message_id"

// HandleSettingsCommand handles the /settings command
func (s *HandlerService) HandleSettingsCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Delete the previous message if it exists
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

	// Send initial settings message (will delete previous one if exists)
	return s.sendSettingsMessage(ctx, c, userID, userProgress, "")
}

// sendSettingsMessage sends or updates the settings message
func (s *HandlerService) sendSettingsMessage(ctx context.Context, c tele.Context, userID int64, userProgress *domain.UserProgress, statusMessage string) error {
	// Delete the previous message if it exists
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
	// Create settings message
	settingsText := "⚙️ *Настройки*\n\n" +
		fmt.Sprintf("🔤 Уровень CEFR: *%s*\n", formatCEFRLevel(userProgress.CEFRLevel)) +
		fmt.Sprintf("📚 Слов в день: *%d*\n", userProgress.WordsPerDay) +
		fmt.Sprintf("🔔 Уведомления: *%s*\n", formatNotificationTime(userProgress.NotificationTime))

	// Add goal display if available
	if userProgress.Preferences != nil {
		if goal, ok := userProgress.Preferences["goal"].(string); ok && goal != "" {
			settingsText += fmt.Sprintf("🎯 Цель изучения: *%s*\n", goal)
		}
	}

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
				{{Text: "Отмена", Data: "settings:back"}},
			},
		}
	case fsm.StateSettingsTopicSelection:
		// Show topic selection options
		return s.sendTopicSelectionMessage(ctx, c, userID)
	default:
		// Default settings keyboard
		keyboard = &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{{Text: "🔤 Уровень CEFR", Data: "settings:cefr_level"}},
				{{Text: "📚 Слов в день", Data: "settings:words_per_day"}},
				{{Text: "🔔 Уведомления", Data: "settings:notifications"}},
				{{Text: "🎯 Цель изучения", Data: "settings:goal_topic"}},
				{{Text: "Назад в главное меню", Data: "menu:back_to_main"}},
			},
		}
	}

	// Try to delete previous settings message if it exists
	if messageIDData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataSettingsMessageID); err == nil {
		// Handle JSON number unmarshaling (numbers come back as float64)
		var messageID int
		switch v := messageIDData.(type) {
		case int:
			messageID = v
		case float64:
			messageID = int(v)
		default:
			s.logger.Warn("Invalid message ID type for deletion", zap.Any("type", v))
			goto sendNewMessage
		}

		// Try to delete the previous message
		if deleteErr := c.Bot().Delete(&tele.Message{
			ID:   messageID,
			Chat: c.Message().Chat,
		}); deleteErr == nil {
			s.logger.Debug("Successfully deleted previous settings message",
				zap.Int64("user_id", userID),
				zap.Int("message_id", messageID))
		} else {
			// Only log as warning if it's not a "message not found" error
			if !strings.Contains(deleteErr.Error(), "message to delete not found") {
				s.logger.Warn("Failed to delete previous settings message",
					zap.Int64("user_id", userID),
					zap.Int("message_id", messageID),
					zap.Error(deleteErr))
			} else {
				s.logger.Debug("Previous settings message already deleted or not found",
					zap.Int64("user_id", userID),
					zap.Int("message_id", messageID))
			}
		}
	}
sendNewMessage:

	// Send new message
	msg, err := c.Bot().Send(c.Sender(), settingsText, &tele.SendOptions{
		ParseMode: tele.ModeMarkdown,
	}, keyboard)

	if err != nil {
		return err
	}

	// Store the new message ID for future deletion
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataSettingsMessageID, msg.ID); err != nil {
		s.logger.Warn("Failed to store settings message ID",
			zap.Int64("user_id", userID),
			zap.Error(err))
	} else {
		s.logger.Debug("Stored new settings message ID",
			zap.Int64("user_id", userID),
			zap.Int("message_id", msg.ID))
	}

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

	// Send thinking message and start typing indicator
	err = s.withThinkingGifAndTyping(ctx, c, userID, "Сохраняю настройки", func() error {
		// Save to backend
		return s.UpdateUserProgress(ctx, userID, userProgress)
	})

	if err != nil {
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

	// Parse time using flexible format parser
	parsedTime, err := ParseTimeFormat(text)
	if err != nil {
		s.logger.Warn("Invalid time format provided",
			zap.String("time_input", text),
			zap.Int64("user_id", userID),
			zap.Error(err))
		return c.Send("❌ Пожалуйста, введите время в формате ЧЧ:ММ (например, 09:30, 9 30, 9:30):")
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Update notification time with parsed format
	userProgress.NotificationTime = parsedTime

	// Send thinking message and start typing indicator
	err = s.withThinkingGifAndTyping(ctx, c, userID, "Сохраняю настройки", func() error {
		// Save to backend
		return s.UpdateUserProgress(ctx, userID, userProgress)
	})

	if err != nil {
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Return to settings with success message
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	statusMessage := fmt.Sprintf("✅ Время уведомлений изменено на *%s*", parsedTime)
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

	// Send thinking message and start typing indicator
	err = s.withThinkingGifAndTyping(ctx, c, userID, "Сохраняю настройки", func() error {
		// Save to backend
		return s.UpdateUserProgress(ctx, userID, userProgress)
	})

	if err != nil {
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

	// Send clean settings message (will delete previous one if exists)
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

	// Send thinking message
	thinkingMsg, err := s.sendThinkingMessage(ctx, c, userID, "Сохраняю настройки")
	if err != nil {
		s.logger.Error("Failed to send thinking message", zap.Error(err))
		// Continue without thinking message if it fails
	}

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		// Delete thinking message if it was sent
		if thinkingMsg != nil {
			if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
				s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
			}
		}
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Delete thinking message if it was sent
	if thinkingMsg != nil {
		if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
			s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
		}
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
	// Delete the previous message if it exists
	if c.Message() != nil {
		if err := c.Delete(); err != nil {
			// Log the error but don't fail the operation
			s.logger.Warn("Failed to delete previous message", zap.Error(err))
		}
	}
	s.logger.Debug("Processing settings time callback",
		zap.String("data", data),
		zap.Int64("user_id", userID))

	// Parse callback data: settings:time:value
	// For time values like "09:00", we need to handle the colon in time
	if !strings.HasPrefix(data, "settings:time:") {
		s.logger.Error("Invalid settings time callback format - missing prefix",
			zap.String("data", data))
		return c.Send("❌ Неверный формат данных.")
	}

	// Extract the time value after "settings:time:"
	value := strings.TrimPrefix(data, "settings:time:")

	s.logger.Debug("Extracted time value",
		zap.String("value", value),
		zap.Int64("user_id", userID))

	if value == "custom" {
		// Delete the previous message if it exists
		if c.Message() != nil {
			if err := c.Delete(); err != nil {
				// Log the error but don't fail the operation
				s.logger.Warn("Failed to delete previous message", zap.Error(err))
			}
		}
		// Set state to input mode (direct transition from settings)
		if err := s.stateManager.SetState(ctx, userID, fsm.StateSettingsTimeInput); err != nil {
			s.logger.Error("Failed to set time input state", zap.Error(err))
			return err
		}
		return c.Send("Когда поставить напоминалку?")
	}

	if value == "disabled" {
		// Delete the previous message if it exists
		if c.Message() != nil {
			if err := c.Delete(); err != nil {
				// Log the error but don't fail the operation
				s.logger.Warn("Failed to delete previous message", zap.Error(err))
			}
		}
		// Disable notifications
		userProgress, err := s.GetUserProgress(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to get user progress", zap.Error(err))
			return err
		}

		userProgress.NotificationTime = ""

		// Send thinking message
		thinkingMsg, err := s.sendThinkingMessage(ctx, c, userID, "Сохраняю настройки")
		if err != nil {
			s.logger.Error("Failed to send thinking message", zap.Error(err))
			// Continue without thinking message if it fails
		}

		// Save to backend
		if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
			// Delete thinking message if it was sent
			if thinkingMsg != nil {
				if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
					s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
				}
			}
			s.logger.Error("Failed to update user progress", zap.Error(err))
			return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
		}

		// Delete thinking message if it was sent
		if thinkingMsg != nil {
			if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
				s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
			}
		}

		// Return to settings with success message
		if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
			s.logger.Error("Failed to set settings state", zap.Error(err))
			return err
		}

		statusMessage := "✅ Уведомления отключены"
		return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
	}

	// Validate time format using flexible parser
	parsedTime, err := ParseTimeFormat(value)
	if err != nil {
		s.logger.Warn("Invalid time format provided in callback",
			zap.String("time_value", value),
			zap.Int64("user_id", userID),
			zap.Error(err))
		return c.Send("❌ Неверный формат времени.")
	}

	// Get current user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	// Update notification time with parsed format
	userProgress.NotificationTime = parsedTime

	// Send thinking message
	thinkingMsg, err := s.sendThinkingMessage(ctx, c, userID, "Сохраняю настройки")
	if err != nil {
		s.logger.Error("Failed to send thinking message", zap.Error(err))
		// Continue without thinking message if it fails
	}

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		// Delete thinking message if it was sent
		if thinkingMsg != nil {
			if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
				s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
			}
		}
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Delete thinking message if it was sent
	if thinkingMsg != nil {
		if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
			s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
		}
	}

	// Return to settings with success message
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	statusMessage := fmt.Sprintf("✅ Время уведомлений изменено на *%s*", parsedTime)
	return s.sendSettingsMessage(ctx, c, userID, userProgress, statusMessage)
}

// HandleSettingsCEFRCallback handles CEFR level selection callbacks
func (s *HandlerService) HandleSettingsCEFRCallback(ctx context.Context, c tele.Context, userID int64, data string) error {
	// Delete the previous message if it exists
	if c.Message() != nil {
		if err := c.Delete(); err != nil {
			// Log the error but don't fail the operation
			s.logger.Warn("Failed to delete previous message", zap.Error(err))
		}
	}

	parts := strings.Split(data, ":")
	if len(parts) != 3 {
		return c.Send("❌ Неверный формат данных.")
	}

	value := parts[2]

	if value == "test" {
		// Delete the previous message if it exists
		if c.Message() != nil {
			if err := c.Delete(); err != nil {
				// Log the error but don't fail the operation
				s.logger.Warn("Failed to delete previous message", zap.Error(err))
			}
		}
		// Start CEFR test
		if err := s.stateManager.SetState(ctx, userID, fsm.StateVocabularyTest); err != nil {
			s.logger.Error("Failed to set vocabulary test state", zap.Error(err))
			return err
		}
		return s.HandleTestStartCallback(ctx, c, userID, fsm.StateVocabularyTest)
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

	// Send thinking message
	thinkingMsg, err := s.sendThinkingMessage(ctx, c, userID, "Сохраняю настройки")
	if err != nil {
		s.logger.Error("Failed to send thinking message", zap.Error(err))
		// Continue without thinking message if it fails
	}

	// Save to backend
	if err := s.UpdateUserProgress(ctx, userID, userProgress); err != nil {
		// Delete thinking message if it was sent
		if thinkingMsg != nil {
			if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
				s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
			}
		}
		s.logger.Error("Failed to update user progress", zap.Error(err))
		return c.Send("❌ Не удалось сохранить настройки. Попробуйте позже.")
	}

	// Delete thinking message if it was sent
	if thinkingMsg != nil {
		if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
			s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
		}
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

// clearSettingsMessage deletes the settings message and clears the stored ID
func (s *HandlerService) clearSettingsMessage(ctx context.Context, c tele.Context, userID int64) {
	// Try to delete previous settings message if it exists
	if messageIDData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataSettingsMessageID); err == nil {
		// Handle JSON number unmarshaling (numbers come back as float64)
		var messageID int
		switch v := messageIDData.(type) {
		case int:
			messageID = v
		case float64:
			messageID = int(v)
		default:
			s.logger.Warn("Invalid message ID type for deletion", zap.Any("type", v))
			return
		}

		// Try to delete the previous message
		if err := c.Bot().Delete(&tele.Message{ID: messageID, Chat: c.Chat()}); err != nil {
			// Only log as warning if it's not a "message not found" error
			if !strings.Contains(err.Error(), "message to delete not found") {
				s.logger.Warn("Failed to delete previous settings message", zap.Error(err))
			} else {
				s.logger.Debug("Previous settings message already deleted or not found")
			}
		}
	}
}

// sendTopicSelectionMessage sends the topic selection interface
func (s *HandlerService) sendTopicSelectionMessage(ctx context.Context, c tele.Context, userID int64) error {
	// Check if user is authenticated
	if !s.stateManager.IsUserAuthenticated(ctx, userID) {
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Необходимо связать аккаунт для изменения цели")
	}

	// Get access token
	accessToken, err := s.stateManager.GetValidAccessToken(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get access token for topic selection", zap.Error(err))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка авторизации")
	}

	// Send thinking message and start typing indicator
	var topics []api.TopicResponse
	err = s.withThinkingGifAndTyping(ctx, c, userID, "Загружаю доступные темы", func() error {
		// Get topics from API
		var topicsErr error
		topics, topicsErr = s.apiClient.GetTopics(ctx, accessToken)
		return topicsErr
	})

	if err != nil {
		s.logger.Error("Failed to get topics", zap.Error(err))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка загрузки тем")
	}

	// Remove duplicates and create unique topics list
	uniqueTopics := make([]string, 0)
	seen := make(map[string]bool)

	for _, topic := range topics {
		if !seen[topic.Title] {
			seen[topic.Title] = true
			uniqueTopics = append(uniqueTopics, topic.Title)
		}
	}

	// Create topic selection data
	topicData := &fsm.TopicSelectionData{
		Topics:        uniqueTopics,
		CurrentPage:   0,
		TopicsPerPage: 5,
		TotalPages:    (len(uniqueTopics) + 4) / 5, // Ceiling division
	}

	// Store topic selection data
	if err := s.stateManager.StoreTopicSelectionData(ctx, userID, topicData); err != nil {
		s.logger.Error("Failed to store topic selection data", zap.Error(err))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка сохранения данных")
	}

	// Send topic selection message
	return s.sendTopicSelectionPage(ctx, c, userID, topicData)
}

// sendTopicSelectionPage sends a specific page of topics
func (s *HandlerService) sendTopicSelectionPage(ctx context.Context, c tele.Context, userID int64, topicData *fsm.TopicSelectionData) error {
	// Delete previous settings message if it exists
	if messageIDData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataSettingsMessageID); err == nil {
		// Handle JSON number unmarshaling (numbers come back as float64)
		var messageID int
		switch v := messageIDData.(type) {
		case int:
			messageID = v
		case float64:
			messageID = int(v)
		default:
			s.logger.Warn("Invalid message ID type for deletion", zap.Any("type", v))
			goto sendNewMessage
		}

		// Try to delete the previous message
		if err := c.Bot().Delete(&tele.Message{ID: messageID, Chat: c.Chat()}); err != nil {
			// Only log as warning if it's not a "message not found" error
			if !strings.Contains(err.Error(), "message to delete not found") {
				s.logger.Warn("Failed to delete previous settings message", zap.Error(err))
			} else {
				s.logger.Debug("Previous settings message already deleted or not found",
					zap.Int64("user_id", userID),
					zap.Int("message_id", messageID))
			}
		} else {
			s.logger.Debug("Successfully deleted previous settings message",
				zap.Int64("user_id", userID),
				zap.Int("message_id", messageID))
		}
	}
sendNewMessage:

	// Calculate start and end indices for current page
	start := topicData.CurrentPage * topicData.TopicsPerPage
	end := start + topicData.TopicsPerPage
	if end > len(topicData.Topics) {
		end = len(topicData.Topics)
	}

	// Get topics for current page
	pageTopics := topicData.Topics[start:end]

	// Create message text
	messageText := "🎯 *Выберите цель изучения*\n\n"
	for i, topic := range pageTopics {
		messageText += fmt.Sprintf("%d. %s\n", i+1, topic)
	}
	messageText += fmt.Sprintf("\nСтраница %d из %d", topicData.CurrentPage+1, topicData.TotalPages)

	// Create keyboard
	var keyboard [][]tele.InlineButton

	// Add topic buttons (5 per row)
	for i, topic := range pageTopics {
		row := i / 5

		// Ensure we have enough rows
		for len(keyboard) <= row {
			keyboard = append(keyboard, []tele.InlineButton{})
		}

		keyboard[row] = append(keyboard[row], tele.InlineButton{
			Text: fmt.Sprintf("%d", i+1),
			Data: fmt.Sprintf("settings:topic:select:%s", topic),
		})
	}

	// Add navigation buttons
	var navRow []tele.InlineButton

	if topicData.CurrentPage > 0 {
		navRow = append(navRow, tele.InlineButton{
			Text: "⬅️ Назад",
			Data: "settings:topic:prev",
		})
	}

	if topicData.CurrentPage < topicData.TotalPages-1 {
		navRow = append(navRow, tele.InlineButton{
			Text: "Вперед ➡️",
			Data: "settings:topic:next",
		})
	}

	if len(navRow) > 0 {
		keyboard = append(keyboard, navRow)
	}

	// Add cancel button
	keyboard = append(keyboard, []tele.InlineButton{
		{Text: "Отмена", Data: "settings:back"},
	})

	// Send message
	msg, err := c.Bot().Send(c.Chat(), messageText, &tele.ReplyMarkup{
		InlineKeyboard: keyboard,
	}, tele.ModeMarkdown)
	if err != nil {
		return err
	}

	// Store message ID for later deletion
	return s.stateManager.StoreTempData(ctx, userID, fsm.TempDataSettingsMessageID, msg.ID)
}

// HandleSettingsGoalTopicCallback handles the goal topic selection callback
func (s *HandlerService) HandleSettingsGoalTopicCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Set state to topic selection
	if err := s.SetStateIfDifferent(ctx, userID, fsm.StateSettingsTopicSelection); err != nil {
		s.logger.Error("Failed to set topic selection state", zap.Error(err))
		return err
	}

	// Send topic selection message
	return s.sendTopicSelectionMessage(ctx, c, userID)
}

// HandleSettingsTopicCallback handles topic selection callbacks
func (s *HandlerService) HandleSettingsTopicCallback(ctx context.Context, c tele.Context, userID int64, data string) error {
	// Parse the callback data
	parts := strings.Split(data, ":")
	if len(parts) < 3 {
		s.logger.Error("Invalid topic callback format", zap.String("data", data))
		return fmt.Errorf("invalid callback format")
	}

	action := parts[2]

	switch action {
	case "select":
		if len(parts) < 4 {
			s.logger.Error("Invalid topic select callback format", zap.String("data", data))
			return fmt.Errorf("invalid callback format")
		}
		// Extract the selected topic (parts[3] and beyond, joined with ":")
		selectedTopic := strings.Join(parts[3:], ":")
		return s.handleTopicSelection(ctx, c, userID, selectedTopic)
	case "prev":
		return s.handleTopicNavigation(ctx, c, userID, -1)
	case "next":
		return s.handleTopicNavigation(ctx, c, userID, 1)
	default:
		s.logger.Error("Unknown topic action", zap.String("action", action))
		return fmt.Errorf("unknown topic action")
	}
}

// handleTopicSelection handles when a user selects a topic
func (s *HandlerService) handleTopicSelection(ctx context.Context, c tele.Context, userID int64, selectedTopic string) error {
	// Check if user is authenticated
	if !s.stateManager.IsUserAuthenticated(ctx, userID) {
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Необходимо связать аккаунт для изменения цели")
	}

	// Get access token
	accessToken, err := s.stateManager.GetValidAccessToken(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get access token for topic selection", zap.Error(err))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка авторизации")
	}

	// Send thinking message
	thinkingMsg, err := s.sendThinkingMessage(ctx, c, userID, "Обновляю цель изучения")
	if err != nil {
		s.logger.Error("Failed to send thinking message", zap.Error(err))
		// Continue without thinking message if it fails
	}

	// Update user preferences with the selected topic
	updateRequest := &api.UpdatePreferenceRequest{
		Goal: &selectedTopic,
	}

	_, err = s.apiClient.UpdateUserPreferences(ctx, accessToken, updateRequest)

	// Delete thinking message if it was sent
	if thinkingMsg != nil {
		if deleteErr := s.deleteMessage(ctx, c, thinkingMsg.ID); deleteErr != nil {
			s.logger.Warn("Failed to delete thinking message", zap.Error(deleteErr))
		}
	}

	if err != nil {
		s.logger.Error("Failed to update user preferences with topic", zap.Error(err))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка обновления цели")
	}

	// Delete the topic selection message
	if messageIDData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataSettingsMessageID); err == nil {
		// Handle JSON number unmarshaling (numbers come back as float64)
		var messageID int
		switch v := messageIDData.(type) {
		case int:
			messageID = v
		case float64:
			messageID = int(v)
		default:
			s.logger.Warn("Invalid message ID type for deletion", zap.Any("type", v))
			goto continueAfterDeletion
		}

		if err := c.Bot().Delete(&tele.Message{ID: messageID, Chat: c.Chat()}); err != nil {
			s.logger.Warn("Failed to delete topic selection message", zap.Error(err))
		}
	}
continueAfterDeletion:

	// Clear topic selection data
	if err := s.stateManager.ClearTempData(ctx, userID, fsm.TempDataTopicSelection); err != nil {
		s.logger.Warn("Failed to clear topic selection data", zap.Error(err))
	}

	// Set state back to settings
	if err := s.SetStateIfDifferent(ctx, userID, fsm.StateSettings); err != nil {
		s.logger.Error("Failed to set settings state", zap.Error(err))
		return err
	}

	// Get updated user progress
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress after topic selection", zap.Error(err))
		return err
	}

	// Send updated settings message
	return s.sendSettingsMessage(ctx, c, userID, userProgress, fmt.Sprintf("✅ Цель обновлена: *%s*", selectedTopic))
}

// handleTopicNavigation handles topic page navigation
func (s *HandlerService) handleTopicNavigation(ctx context.Context, c tele.Context, userID int64, direction int) error {
	// Get current topic selection data
	topicData, err := s.stateManager.GetTopicSelectionData(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get topic selection data", zap.Error(err))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка загрузки данных")
	}

	if topicData == nil {
		s.logger.Error("No topic selection data found")
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Данные не найдены")
	}

	// Calculate new page
	newPage := topicData.CurrentPage + direction
	if newPage < 0 || newPage >= topicData.TotalPages {
		s.logger.Error("Invalid page navigation", zap.Int("current_page", topicData.CurrentPage), zap.Int("direction", direction))
		return fmt.Errorf("invalid page navigation")
	}

	// Update current page
	topicData.CurrentPage = newPage

	// Store updated data
	if err := s.stateManager.StoreTopicSelectionData(ctx, userID, topicData); err != nil {
		s.logger.Error("Failed to store updated topic selection data", zap.Error(err))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка сохранения данных")
	}

	// Get current message ID
	messageIDData, err := s.stateManager.GetTempData(ctx, userID, fsm.TempDataSettingsMessageID)
	if err != nil {
		s.logger.Error("Failed to get message ID", zap.Error(err))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка получения ID сообщения")
	}

	// Handle JSON number unmarshaling (numbers come back as float64)
	var messageID int
	switch v := messageIDData.(type) {
	case int:
		messageID = v
	case float64:
		messageID = int(v)
	default:
		s.logger.Error("Invalid message ID type", zap.Any("type", v))
		return s.sendSettingsMessage(ctx, c, userID, &domain.UserProgress{}, "❌ Ошибка типа ID сообщения")
	}

	// Update the existing message instead of deleting and recreating
	return s.updateTopicSelectionMessage(ctx, c, userID, topicData, messageID)
}

// updateTopicSelectionMessage updates an existing topic selection message
func (s *HandlerService) updateTopicSelectionMessage(ctx context.Context, c tele.Context, userID int64, topicData *fsm.TopicSelectionData, messageID int) error {
	// Calculate start and end indices for current page
	start := topicData.CurrentPage * topicData.TopicsPerPage
	end := start + topicData.TopicsPerPage
	if end > len(topicData.Topics) {
		end = len(topicData.Topics)
	}

	// Get topics for current page
	pageTopics := topicData.Topics[start:end]

	// Create message text
	messageText := "🎯 *Выберите цель изучения*\n\n"
	for i, topic := range pageTopics {
		messageText += fmt.Sprintf("%d. %s\n", i+1, topic)
	}
	messageText += fmt.Sprintf("\nСтраница %d из %d", topicData.CurrentPage+1, topicData.TotalPages)

	// Create keyboard
	var keyboard [][]tele.InlineButton

	// Add topic buttons (5 per row)
	for i, topic := range pageTopics {
		row := i / 5

		// Ensure we have enough rows
		for len(keyboard) <= row {
			keyboard = append(keyboard, []tele.InlineButton{})
		}

		keyboard[row] = append(keyboard[row], tele.InlineButton{
			Text: fmt.Sprintf("%d", i+1),
			Data: fmt.Sprintf("settings:topic:select:%s", topic),
		})
	}

	// Add navigation buttons
	var navRow []tele.InlineButton

	if topicData.CurrentPage > 0 {
		navRow = append(navRow, tele.InlineButton{
			Text: "⬅️ Назад",
			Data: "settings:topic:prev",
		})
	}

	if topicData.CurrentPage < topicData.TotalPages-1 {
		navRow = append(navRow, tele.InlineButton{
			Text: "Вперед ➡️",
			Data: "settings:topic:next",
		})
	}

	if len(navRow) > 0 {
		keyboard = append(keyboard, navRow)
	}

	// Add cancel button
	keyboard = append(keyboard, []tele.InlineButton{
		{Text: "Отмена", Data: "settings:back"},
	})

	// Update the existing message
	_, err := c.Bot().Edit(
		&tele.Message{ID: messageID, Chat: c.Chat()},
		messageText,
		&tele.ReplyMarkup{InlineKeyboard: keyboard},
		tele.ModeMarkdown,
	)

	return err
}
