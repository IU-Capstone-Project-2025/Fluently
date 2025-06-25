// start.go
package handlers

import (
	"log"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// Стартовая команда /start
func StartHandler(c tele.Context) error {
	// Отправляем приветствие
	intro := strings.Join([]string{
		"Добро пожаловать! Это приложение Fluently!",
	}, "\n")

	if err := c.Send(intro, &tele.SendOptions{
		ParseMode: tele.ModeHTML,
	}); err != nil {
		log.Printf("Start: Failed to send greeting message: %v", err)
		return err
	}

	return nil
}
