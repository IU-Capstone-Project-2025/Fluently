// router.go
package router

import (
	"fluently/telegram-bot/internal/bot/handlers"

	tele "gopkg.in/telebot.v3"
)

func Register(b *tele.Bot) {
	b.Handle("/start", handlers.StartHandler)
}
