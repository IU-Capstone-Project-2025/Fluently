package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
)

// HandleLearningStart handles the start of a learning session
func (s *HandlerService) HandleLearningStart(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Create sample words for the lesson as map[string]interface{}
	words := []map[string]interface{}{
		{
			"id":          int64(1),
			"word":        "hello",
			"translation": "привет",
			"definition":  "used as a greeting",
			"examples":    []string{"Hello, how are you?", "She said hello to everyone."},
			"cefr_level":  "A1",
		},
		{
			"id":          int64(2),
			"word":        "world",
			"translation": "мир",
			"definition":  "the earth and all the people, places, and things on it",
			"examples":    []string{"The world is beautiful.", "People from around the world came to visit."},
			"cefr_level":  "A1",
		},
		{
			"id":          int64(3),
			"word":        "beautiful",
			"translation": "красивый",
			"definition":  "having beauty; pleasing to the senses",
			"examples":    []string{"She is beautiful.", "What a beautiful day!"},
			"cefr_level":  "A2",
		},
	}

	// Create lesson data
	lessonData := &fsm.LessonData{
		Words:            words,
		CurrentWordIndex: 0,
		StartTime:        time.Now(),
	}

	// Store lesson data
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataLesson, lessonData); err != nil {
		log.Printf("Failed to store lesson data: %v", err)
		return err
	}

	// Set state to showing words
	if err := s.stateManager.SetState(ctx, userID, fsm.StateShowingWords); err != nil {
		log.Printf("Failed to set learning state: %v", err)
		return err
	}

	// Send first word
	return s.sendCurrentWord(ctx, c, userID, lessonData)
}

// sendCurrentWord sends the current word to the user
func (s *HandlerService) sendCurrentWord(ctx context.Context, c tele.Context, userID int64, lessonData *fsm.LessonData) error {
	if lessonData.CurrentWordIndex >= len(lessonData.Words) {
		return s.completeLearningSession(ctx, c, userID, lessonData)
	}

	word := lessonData.Words[lessonData.CurrentWordIndex]

	// Extract word data
	wordText := word["word"].(string)
	definition := word["definition"].(string)
	examples := word["examples"].([]string)

	// Create word presentation
	wordMessage := fmt.Sprintf(
		"📚 *Word %d of %d*\n\n"+
			"🔤 **%s**\n"+
			"🔊 /%s/\n"+
			"📝 %s\n\n"+
			"💭 Examples:\n%s",
		lessonData.CurrentWordIndex+1,
		len(lessonData.Words),
		wordText,
		wordText, // placeholder for pronunciation
		definition,
		strings.Join(examples, "\n"),
	)

	// Create interaction buttons
	wordID := word["id"].(int64)
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "Show Translation", Data: fmt.Sprintf("word:translation:%d", wordID)},
				{Text: "I Know This", Data: fmt.Sprintf("word:know:%d", wordID)},
			},
			{
				{Text: "Next Word", Data: fmt.Sprintf("word:next:%d", wordID)},
			},
		},
	}

	return c.Send(wordMessage, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// completeLearningSession completes the learning session
func (s *HandlerService) completeLearningSession(ctx context.Context, c tele.Context, userID int64, lessonData *fsm.LessonData) error {
	// Calculate session statistics
	completedWords := len(lessonData.Words)

	// Create completion message
	completionText := fmt.Sprintf(
		"🎉 *Lesson Complete!*\n\n"+
			"📊 **Session Summary:**\n"+
			"✅ Words reviewed: %d\n"+
			"🎯 Progress: %d/%d\n\n"+
			"Great job! Keep up the good work!",
		completedWords,
		completedWords,
		len(lessonData.Words),
	)

	// Create post-lesson buttons
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "Review Again", Data: "lesson:review"},
				{Text: "New Lesson", Data: "lesson:new"},
			},
			{
				{Text: "Main Menu", Data: "menu:main"},
			},
		},
	}

	// Reset state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateStart); err != nil {
		log.Printf("Failed to reset state: %v", err)
	}

	// Clear lesson data
	if err := s.stateManager.ClearTempData(ctx, userID, fsm.TempDataLesson); err != nil {
		log.Printf("Failed to clear lesson data: %v", err)
	}

	return c.Send(completionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleWordCallback handles word-related callbacks
func (s *HandlerService) HandleWordCallback(ctx context.Context, c tele.Context, userID int64, action string, wordID int64) error {
	// Get current lesson data
	lessonData, err := s.stateManager.GetLessonData(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get lesson data: %w", err)
	}

	// Find the current word
	if lessonData.CurrentWordIndex >= len(lessonData.Words) {
		return fmt.Errorf("invalid word index")
	}

	word := lessonData.Words[lessonData.CurrentWordIndex]

	switch action {
	case "translation":
		return s.showWordTranslation(ctx, c, userID, word)
	case "know":
		return s.markWordAsKnown(ctx, c, userID, word, lessonData)
	case "next":
		return s.moveToNextWord(ctx, c, userID, lessonData)
	default:
		return fmt.Errorf("unknown word action: %s", action)
	}
}

// showWordTranslation shows the translation of the current word
func (s *HandlerService) showWordTranslation(ctx context.Context, c tele.Context, userID int64, word map[string]interface{}) error {
	wordText := word["word"].(string)
	translation := word["translation"].(string)

	translationText := fmt.Sprintf(
		"🔤 **%s**\n"+
			"🌐 Translation: **%s**\n\n"+
			"Did this help you understand the word?",
		wordText,
		translation,
	)

	wordID := word["id"].(int64)
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "Yes, I understand", Data: fmt.Sprintf("word:understand:%d", wordID)},
				{Text: "Still confused", Data: fmt.Sprintf("word:confused:%d", wordID)},
			},
		},
	}

	return c.Send(translationText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// markWordAsKnown marks a word as known by the user
func (s *HandlerService) markWordAsKnown(ctx context.Context, c tele.Context, userID int64, word map[string]interface{}, lessonData *fsm.LessonData) error {
	wordText := word["word"].(string)
	// Update word learning score
	// This would normally update the database
	log.Printf("User %d marked word '%s' as known", userID, wordText)

	// Move to next word
	return s.moveToNextWord(ctx, c, userID, lessonData)
}

// moveToNextWord moves to the next word in the lesson
func (s *HandlerService) moveToNextWord(ctx context.Context, c tele.Context, userID int64, lessonData *fsm.LessonData) error {
	lessonData.CurrentWordIndex++

	// Update lesson data
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataLesson, lessonData); err != nil {
		return fmt.Errorf("failed to update lesson data: %w", err)
	}

	// Send next word or complete session
	return s.sendCurrentWord(ctx, c, userID, lessonData)
}
