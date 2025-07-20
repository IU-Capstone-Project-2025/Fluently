package handlers

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
)

// HandleStatsCommand handles the /stats command
func (s *HandlerService) HandleStatsCommand(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Delete the previous message if it exists, but preserve lesson completion messages
	if c.Message() != nil {
		if err := c.Delete(); err != nil {
			// Log the error but don't fail the operation
			s.logger.Warn("Failed to delete previous message", zap.Error(err))
		}
	}
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
		1,   // streak days - placeholder
		97,  // total words - placeholder
		7,   // lessons completed - placeholder
		144, // total time - placeholder
	)

	// Create back button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Назад в главное меню", Data: "menu:back_to_main"}},
		},
	}

	return c.Send(statsText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}
