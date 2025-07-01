package handlers

import (
	"time"

	"fluently/telegram-bot/internal/bot"

	tele "gopkg.in/telebot.v3"
)

// UserData represents user data stored in Redis
type UserData struct {
	Name     string    `json:"name"`
	Language string    `json:"language"`
	LastSeen time.Time `json:"last_seen"`
}

// ExampleRedisHandler demonstrates how to use Redis in handlers
func ExampleRedisHandler(c tele.Context) error {
	userID := c.Sender().ID
	redisService := bot.NewRedisService()

	// Example: Store user data
	userData := UserData{
		Name:     c.Sender().FirstName,
		Language: c.Sender().LanguageCode,
		LastSeen: time.Now(),
	}

	err := redisService.SetUserData(userID, userData, 24*time.Hour)
	if err != nil {
		return c.Send("Error storing user data")
	}

	// Example: Retrieve user data
	var retrievedData UserData
	err = redisService.GetUserData(userID, &retrievedData)
	if err != nil {
		return c.Send("Error retrieving user data")
	}

	// Example: Set user state
	err = redisService.SetUserState(userID, "waiting_for_input", time.Hour)
	if err != nil {
		return c.Send("Error setting user state")
	}

	// Example: Get user state
	state, err := redisService.GetUserState(userID)
	if err != nil {
		return c.Send("Error getting user state")
	}

	return c.Send("Redis operations completed successfully! State: " + state)
}
