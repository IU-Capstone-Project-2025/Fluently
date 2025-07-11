package handlers

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
)

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

// HandleSettingsWordsPerDayCallback handles words per day settings callback
func (s *HandlerService) HandleSettingsWordsPerDayCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Настройка слов в день...")
}

// HandleSettingsNotificationsCallback handles notifications settings callback
func (s *HandlerService) HandleSettingsNotificationsCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Настройка уведомлений...")
}

// HandleSettingsCEFRLevelCallback handles CEFR level settings callback
func (s *HandlerService) HandleSettingsCEFRLevelCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Настройка уровня CEFR...")
}

// HandleSettingsWordsPerDayInputMessage handles words per day input messages
func (s *HandlerService) HandleSettingsWordsPerDayInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, введите количество слов в день.")
}

// HandleSettingsTimeInputMessage handles time input messages
func (s *HandlerService) HandleSettingsTimeInputMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, введите время уведомлений.")
}

// HandleSettingsMenuCallback handles settings menu callback
func (s *HandlerService) HandleSettingsMenuCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Меню настроек...")
}

// formatNotificationTime formats notification time string
func formatNotificationTime(timeStr string) string {
	if timeStr == "" {
		return "Отключены"
	}
	return timeStr
}
