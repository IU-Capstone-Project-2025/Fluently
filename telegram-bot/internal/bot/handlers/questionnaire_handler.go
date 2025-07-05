package handlers

import (
	"context"

	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
)

// HandleQuestionGoalMessage handles goal question messages
func (s *HandlerService) HandleQuestionGoalMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос о цели.")
}

// HandleQuestionConfidenceMessage handles confidence question messages
func (s *HandlerService) HandleQuestionConfidenceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос об уверенности.")
}

// HandleQuestionSerialsMessage handles serials question messages
func (s *HandlerService) HandleQuestionSerialsMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос о сериалах.")
}

// HandleQuestionExperienceMessage handles experience question messages
func (s *HandlerService) HandleQuestionExperienceMessage(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	return c.Send("Пожалуйста, ответьте на вопрос об опыте.")
}
