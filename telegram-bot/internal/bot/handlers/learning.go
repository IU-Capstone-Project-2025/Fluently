package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
	"telegram-bot/internal/domain"
)

// HandleNewLearningStart initiates the new learning flow
func (s *HandlerService) HandleNewLearningStart(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Delete the previous message if it exists, but preserve lesson completion messages
	if c.Message() != nil {
		// Check if this is a lesson completion message (contains learned words list)
		messageText := c.Message().Text
		if messageText != "" && !strings.Contains(messageText, "🏆 *Урок завершен!*") {
			// Only delete if it's not a lesson completion message
			if err := c.Delete(); err != nil {
				// Only log as warning if it's not a "message not found" error
				if !strings.Contains(err.Error(), "message to delete not found") {
					s.logger.Warn("Failed to delete previous message", zap.Error(err))
				} else {
					s.logger.Debug("Previous message already deleted or not found", zap.Error(err))
				}
			}
		}
	}

	// Check if user is authenticated (has JWT token)
	token, err := s.stateManager.GetJWTToken(ctx, userID)
	if err != nil {
		return s.handleUnauthenticatedUser(ctx, c, userID)
	}

	// Check if user has an active lesson in progress
	hasActiveLesson, err := s.stateManager.HasActiveLessonProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to check active lesson", zap.Error(err))
		return err
	}

	if hasActiveLesson {
		// Resume existing lesson
		return s.resumeLesson(ctx, c, userID)
	}

	// Send thinking message and start typing indicator
	var lessonResponse *domain.LessonResponse
	err = s.withThinkingGifAndTyping(ctx, c, userID, "Генерирую персональный урок", func() error {
		// Generate new lesson from backend
		var generateErr error
		lessonResponse, generateErr = s.apiClient.GenerateLesson(ctx, token)
		return generateErr
	})

	if err != nil {
		s.logger.Error("Failed to generate lesson", zap.Error(err))

		// Check if this is a preferences-related error
		if strings.Contains(err.Error(), "failed to get preference") || strings.Contains(err.Error(), "preference not found") {
			s.logger.Warn("Lesson generation failed due to missing preferences, guiding user to setup", zap.Int64("user_id", userID))

			// Guide user to complete their profile setup
			message := "🔧 *Настройка профиля требуется*\n\n" +
				"Для создания персональных уроков необходимо завершить настройку профиля.\n\n" +
				"📝 Что нужно сделать:\n" +
				"• Определить ваш уровень английского\n" +
				"• Установить количество слов в день\n" +
				"• Настроить уведомления\n\n" +
				"Используйте команду /start для завершения настройки."

			return c.Send(message, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
		}

		// For other errors, show generic message
		return c.Send("❌ Не удалось получить урок. Попробуйте позже.")
	}

	// Initialize lesson progress
	progress := &domain.LessonProgress{
		LessonData:         lessonResponse,
		CurrentWordIndex:   0,
		CurrentPhase:       "showing_words",
		WordsInCurrentSet:  []domain.Card{},
		CurrentSetIndex:    0,
		ExerciseIndex:      0,
		WordsLearned:       []domain.WordProgress{},
		BadlyAnsweredWords: []domain.BadlyAnsweredWord{},
		StartTime:          time.Now(),
		LastActivity:       time.Now(),
		LearnedCount:       0,
		AlreadyKnownCount:  0,
	}

	// Store lesson progress
	err = s.stateManager.StoreLessonProgress(ctx, userID, progress)
	if err != nil {
		s.logger.Error("Failed to store lesson progress", zap.Error(err))
		return err
	}

	// Set state to lesson in progress
	if err := s.stateManager.SetState(ctx, userID, fsm.StateLessonInProgress); err != nil {
		s.logger.Error("Failed to set lesson state", zap.Error(err))
		return err
	}

	// Start the lesson with introduction
	return s.startNewLesson(ctx, c, userID, progress)
}

// startNewLesson starts a new lesson with introduction
func (s *HandlerService) startNewLesson(ctx context.Context, c tele.Context, userID int64, progress *domain.LessonProgress) error {
	wordsPerLesson := progress.LessonData.Lesson.WordsPerLesson

	introText := fmt.Sprintf(
		"📚 *Персональный урок сгенерирован!*\n\n"+
			"🎯 Цель: выучить %d новых слов\n"+
			"📊 Уровень: %s\n\n"+
			"Готовы начать урок?",
		wordsPerLesson,
		progress.LessonData.Lesson.CEFRLevel,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🚀 Начать урок", Data: "lesson:start_word_set"},
				{Text: "📊 Статистика", Data: "lesson:stats"},
			},
		},
	}

	return c.Send(introText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// resumeLesson resumes an existing lesson
func (s *HandlerService) resumeLesson(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	if progress == nil {
		return fmt.Errorf("no lesson progress found")
	}

	learnedCount := progress.LearnedCount
	targetCount := progress.LessonData.Lesson.WordsPerLesson

	resumeText := fmt.Sprintf(
		"📖 *Продолжаем урок*\n\n"+
			"✅ Выучено слов: %d/%d\n"+
			"⏱ Время урока: %s\n\n"+
			"Продолжим с того места, где остановились?",
		learnedCount,
		targetCount,
		s.formatDuration(time.Since(progress.StartTime)),
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "▶️ Продолжить", Data: "lesson:continue"},
				{Text: "🔄 Начать заново", Data: "lesson:restart"},
			},
		},
	}

	return c.Send(resumeText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleStartWordSet starts showing a new set of 3 words
func (s *HandlerService) HandleStartWordSet(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Delete the previous message if it exists, but preserve lesson completion messages
	if c.Message() != nil {
		// Check if this is a lesson completion message (contains learned words list)
		messageText := c.Message().Text
		if messageText != "" && !strings.Contains(messageText, "🏆 *Урок завершен!*") {
			// Only delete if it's not a lesson completion message
			if err := c.Delete(); err != nil {
				// Only log as warning if it's not a "message not found" error
				if !strings.Contains(err.Error(), "message to delete not found") {
					s.logger.Warn("Failed to delete previous message", zap.Error(err))
				} else {
					s.logger.Debug("Previous message already deleted or not found", zap.Error(err))
				}
			}
		}
	}

	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	if progress == nil {
		return fmt.Errorf("no lesson progress found")
	}

	// Check if lesson is complete
	if progress.LearnedCount >= progress.LessonData.Lesson.WordsPerLesson {
		return s.completeLessonFlow(ctx, c, userID, progress)
	}

	// Prepare next set of 3 words
	nextSet, err := s.getNextWordSet(progress)
	if err != nil {
		return err
	}

	// Update progress with new word set
	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		p.WordsInCurrentSet = nextSet
		p.CurrentPhase = "showing_words"
		p.CurrentSetIndex++
		p.LastActivity = time.Now()
		return nil
	})
	if err != nil {
		return err
	}

	// Set state to showing word set
	if err := s.stateManager.SetState(ctx, userID, fsm.StateShowingWordSet); err != nil {
		return err
	}

	// Show set introduction
	setIntroText := fmt.Sprintf(
		"Сейчас изучим %d новых слов:\n",
		len(nextSet),
	)

	// Add each word to the text
	for i, word := range nextSet {
		setIntroText += fmt.Sprintf("%d️⃣ %s - %s\n", i+1, word.Word, word.Translation)
	}

	setIntroText += "\nНажмите \"Изучать\", чтобы посмотреть каждое слово подробно."

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "📖 Изучать слова", Data: "lesson:show_word:0"},
			},
		},
	}

	return c.Send(setIntroText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleShowWord shows a specific word with examples and voice
func (s *HandlerService) HandleShowWord(ctx context.Context, c tele.Context, userID int64, wordIndex int) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	if wordIndex >= len(progress.WordsInCurrentSet) {
		return fmt.Errorf("invalid word index")
	}

	word := progress.WordsInCurrentSet[wordIndex]

	// Set appropriate state
	var newState fsm.UserState
	switch wordIndex {
	case 0:
		newState = fsm.StateShowingWord1
	case 1:
		newState = fsm.StateShowingWord2
	case 2:
		newState = fsm.StateShowingWord3
	}

	if err := s.stateManager.SetState(ctx, userID, newState); err != nil {
		return err
	}

	// Format word information with bolded English word and quoted examples
	detailText := fmt.Sprintf(
		"Слово %d из %d\n"+
			"<b>%s</b> - %s\n"+
			"Examples:",
		wordIndex+1,
		len(progress.WordsInCurrentSet),
		word.Word,
		word.Translation,
	)

	// Add examples as quoted sentences using Telegram's quote format
	for i, sentence := range word.Sentences {
		if i >= 1 { // Limit to 1 examples
			break
		}
		detailText += fmt.Sprintf(
			"\n<blockquote>%s</blockquote>\n<blockquote>%s</blockquote>",
			sentence.Text,
			sentence.Translation,
		)
	}

	// Create navigation buttons
	var buttons [][]tele.InlineButton

	// Add "Already know" button
	buttons = append(buttons, []tele.InlineButton{
		{Text: "✅ Уже знаю", Data: fmt.Sprintf("lesson:already_know:%d", wordIndex)},
	})

	// Navigation buttons
	var navButtons []tele.InlineButton

	if wordIndex < len(progress.WordsInCurrentSet)-1 {
		navButtons = append(navButtons, tele.InlineButton{
			Text: "Следующее ▶️",
			Data: fmt.Sprintf("lesson:show_word:%d", wordIndex+1),
		})
	} else {
		// Last word - show "Ready for exercises" button
		navButtons = append(navButtons, tele.InlineButton{
			Text: "✅ К упражнениям",
			Data: "lesson:ready_exercises",
		})
	}

	if len(navButtons) > 0 {
		buttons = append(buttons, navButtons)
	}

	keyboard := &tele.ReplyMarkup{InlineKeyboard: buttons}

	// Generate and send voice message for the word pronunciation
	if err := s.sendWordVoiceMessage(ctx, c, word.Word); err != nil {
		s.logger.Warn("Failed to send voice message", zap.Error(err))
		// Continue with text-only version if voice fails
	}

	// Send the combined message with voice and text
	if err := c.Send(detailText, &tele.SendOptions{ParseMode: tele.ModeHTML}, keyboard); err != nil {
		return err
	}

	return nil
}

// HandleWordAlreadyKnown handles when user marks a word as already known
func (s *HandlerService) HandleWordAlreadyKnown(ctx context.Context, c tele.Context, userID int64, wordIndex int) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	if wordIndex >= len(progress.WordsInCurrentSet) {
		return fmt.Errorf("invalid word index")
	}

	word := progress.WordsInCurrentSet[wordIndex]

	// Set state to word already known
	if err := s.stateManager.SetState(ctx, userID, fsm.StateWordAlreadyKnown); err != nil {
		return err
	}

	// Add word to learned words with high confidence (since user already knows it)
	// But don't increment LearnedCount since it doesn't count toward daily limit
	wordProgress := domain.WordProgress{
		Word:            word.Word,
		Translation:     word.Translation,
		WordID:          word.WordID,
		LearnedAt:       time.Now(),
		ConfidenceScore: 100, // High confidence since user already knows it
		CntReviewed:     1,
		AlreadyKnown:    true, // Mark as already known
	}

	// Add to WordsLearned but don't increment LearnedCount
	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		p.WordsLearned = append(p.WordsLearned, wordProgress)
		// Note: LearnedCount is NOT incremented for "already known" words
		p.AlreadyKnownCount++ // Track words marked as already known
		p.LastActivity = time.Now()
		return nil
	})
	if err != nil {
		s.logger.Error("Failed to add word progress", zap.Error(err))
	}

	// Try to replace the word with a new one from the pool
	replacementErr := s.replaceWordInCurrentSet(ctx, userID, wordIndex)

	var confirmText string
	var buttons [][]tele.InlineButton

	if replacementErr != nil {
		// No more replacement words available
		confirmText = fmt.Sprintf(
			"✅ *Отлично!*\n\n"+
				"Слово **%s** отмечено как уже известное.\n"+
				"Больше слов для замены нет.\n"+
				"Оно не будет засчитано в дневной лимит слов.",
			word.Word,
		)

		// Continue to next word or exercises
		if wordIndex < 2 {
			// Not the last word in set
			buttons = append(buttons, []tele.InlineButton{
				{Text: "▶️ Следующее слово", Data: fmt.Sprintf("lesson:show_word:%d", wordIndex+1)},
			})
		} else {
			// Last word in set - go to exercises
			buttons = append(buttons, []tele.InlineButton{
				{Text: "✅ К упражнениям", Data: "lesson:ready_exercises"},
			})
		}
	} else {
		// Word was successfully replaced
		// Get the updated progress to show the new word
		updatedProgress, err := s.stateManager.GetLessonProgress(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to get updated progress", zap.Error(err))
			return err
		}

		replacementWord := updatedProgress.WordsInCurrentSet[wordIndex]

		confirmText = fmt.Sprintf(
			"✅ *Отлично!*\n\n"+
				"Слово **%s** отмечено как уже известное.\n"+
				"Заменено на новое слово: **%s**\n"+
				"Оно не будет засчитано в дневной лимит слов.",
			word.Word,
			replacementWord.Word,
		)

		// Show the replacement word
		buttons = append(buttons, []tele.InlineButton{
			{Text: "👀 Посмотреть новое слово", Data: fmt.Sprintf("lesson:show_word:%d", wordIndex)},
		})
	}

	keyboard := &tele.ReplyMarkup{InlineKeyboard: buttons}

	return c.Send(confirmText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleReadyForExercises transitions to exercise phase
func (s *HandlerService) HandleReadyForExercises(ctx context.Context, c tele.Context, userID int64) error {
	// Set state to ready for exercises
	if err := s.stateManager.SetState(ctx, userID, fsm.StateReadyForExercises); err != nil {
		return err
	}

	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	readyText := fmt.Sprintf(
		"Время упражнений!\n\n"+
			"Проверим, как вы усвоили эти %d слова:\n",
		len(progress.WordsInCurrentSet),
	)

	// Add each word to the text
	for _, word := range progress.WordsInCurrentSet {
		readyText += fmt.Sprintf("• %s\n", word.Word)
	}

	readyText += fmt.Sprintf("\nБудет %d упражнения. Готовы?", len(progress.WordsInCurrentSet))

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🚀 Начать упражнения", Data: "lesson:start_exercises"},
			},
		},
	}

	return c.Send(readyText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleStartExercises begins the exercise phase
func (s *HandlerService) HandleStartExercises(ctx context.Context, c tele.Context, userID int64) error {
	// Set state to doing exercises
	if err := s.stateManager.SetState(ctx, userID, fsm.StateDoingExercises); err != nil {
		return err
	}

	// Update progress
	err := s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		p.CurrentPhase = "exercises"
		p.ExerciseIndex = 0
		p.LastActivity = time.Now()
		return nil
	})
	if err != nil {
		return err
	}

	// Start first exercise
	return s.showNextExercise(ctx, c, userID)
}

// showNextExercise displays the next exercise for the current word set
func (s *HandlerService) showNextExercise(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	// Check if we're in retry mode
	if progress.CurrentPhase == "retry" {
		return s.showNextRetryExercise(ctx, c, userID)
	}

	// Check if all exercises for current set are complete
	if progress.ExerciseIndex >= len(progress.WordsInCurrentSet) {
		return s.completeCurrentSet(ctx, c, userID)
	}

	// Get current word and its exercise
	currentWord := progress.WordsInCurrentSet[progress.ExerciseIndex]
	exercise := currentWord.Exercise

	// Set state to exercise in progress
	if err := s.stateManager.SetState(ctx, userID, fsm.StateExerciseInProgress); err != nil {
		return err
	}

	// Handle different exercise types
	switch exercise.Type {
	case "pick_option_sentence":
		return s.showPickOptionSentenceExercise(ctx, c, userID, currentWord, exercise)
	case "write_word_from_translation":
		return s.showWriteWordTranslationExercise(ctx, c, userID, currentWord, exercise)
	case "translate_ru_to_en":
		return s.showTranslateRuToEnExercise(ctx, c, userID, currentWord, exercise)
	default:
		return fmt.Errorf("unknown exercise type: %s", exercise.Type)
	}
}

// Helper function to get next set of words (up to 3, but respecting daily limit)
func (s *HandlerService) getNextWordSet(progress *domain.LessonProgress) ([]domain.Card, error) {
	// Calculate how many words we've already learned in this lesson
	wordsLearnedInLesson := progress.LearnedCount - progress.AlreadyKnownCount

	// Calculate how many words are left to reach the daily goal
	wordsLeft := progress.LessonData.Lesson.WordsPerLesson - wordsLearnedInLesson

	if wordsLeft <= 0 {
		return nil, fmt.Errorf("daily word limit reached")
	}

	// Determine how many words to show in this set (max 3, but not more than words left)
	wordsInThisSet := 3
	if wordsLeft < 3 {
		wordsInThisSet = wordsLeft
	}

	// Calculate start index based on words already shown in this lesson
	// We need to track how many words we've already shown, not just learned
	wordsShownInLesson := progress.CurrentSetIndex * 3 // Each set has 3 words
	cards := progress.LessonData.Cards

	// Start from the next available word after what we've already shown
	startIndex := wordsShownInLesson

	if startIndex >= len(cards) {
		return nil, fmt.Errorf("no more words available")
	}

	// Find the next available words that haven't been used yet
	var selectedWords []domain.Card
	wordsFound := 0

	for i := startIndex; i < len(cards) && wordsFound < wordsInThisSet; i++ {
		wordID := cards[i].WordID

		// Check if this word is already in the current set
		isInCurrentSet := false
		for _, currentWord := range progress.WordsInCurrentSet {
			if currentWord.WordID == wordID {
				isInCurrentSet = true
				break
			}
		}

		// Check if this word was already learned or shown
		isAlreadyUsed := false
		for _, learnedWord := range progress.WordsLearned {
			if learnedWord.WordID == wordID {
				isAlreadyUsed = true
				break
			}
		}

		// If word is not in current set and not already used, add it
		if !isInCurrentSet && !isAlreadyUsed {
			selectedWords = append(selectedWords, cards[i])
			wordsFound++
		}
	}

	if len(selectedWords) == 0 {
		return nil, fmt.Errorf("no more words available")
	}

	return selectedWords, nil
}

// startRetryPhase begins the retry phase for words that were answered incorrectly
func (s *HandlerService) startRetryPhase(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.Info("Starting retry phase",
		zap.Int64("user_id", userID),
		zap.Int("retry_words_count", len(progress.RetryWords)))

	// Set phase to retry and reset retry index
	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		p.CurrentPhase = "retry"
		p.RetryIndex = 0
		p.LastActivity = time.Now()
		return nil
	})
	if err != nil {
		return err
	}

	// Send retry phase message
	retryText := fmt.Sprintf(
		"🔄 *Повторение слов*\n\n"+
			"📝 У вас есть %d слов, которые нужно повторить.\n\n"+
			"Давайте закрепим знания!",
		len(progress.RetryWords),
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🚀 Начать повторение", Data: "exercise:next"},
			},
		},
	}

	return c.Send(retryText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// showNextRetryExercise displays the next retry exercise
func (s *HandlerService) showNextRetryExercise(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	// Check if all retry exercises are complete
	if progress.RetryIndex >= len(progress.RetryWords) {
		return s.completeRetryPhase(ctx, c, userID)
	}

	// Get current retry word and its exercise
	currentWord := progress.RetryWords[progress.RetryIndex]
	exercise := currentWord.Exercise

	// Set state to exercise in progress
	if err := s.stateManager.SetState(ctx, userID, fsm.StateExerciseInProgress); err != nil {
		return err
	}

	// Handle different exercise types
	switch exercise.Type {
	case "pick_option_sentence":
		return s.showPickOptionSentenceExercise(ctx, c, userID, currentWord, exercise)
	case "write_word_from_translation":
		return s.showWriteWordTranslationExercise(ctx, c, userID, currentWord, exercise)
	case "translate_ru_to_en":
		return s.showTranslateRuToEnExercise(ctx, c, userID, currentWord, exercise)
	default:
		return fmt.Errorf("unknown exercise type: %s", exercise.Type)
	}
}

// completeRetryPhase handles completion of the retry phase
func (s *HandlerService) completeRetryPhase(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.Info("Completing retry phase",
		zap.Int64("user_id", userID),
		zap.Int("remaining_retry_words", len(progress.RetryWords)))

	// Set phase back to completed
	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		p.CurrentPhase = "completed"
		p.LastActivity = time.Now()
		return nil
	})
	if err != nil {
		return err
	}

	// Send completion message
	completionText := fmt.Sprintf(
		"✅ *Повторение завершено!*\n\n"+
			"🎯 Вы повторили %d слов.\n\n"+
			"Теперь давайте завершим урок!",
		len(progress.RetryWords),
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🏆 Завершить урок", Data: "lesson:final_stats"},
			},
		},
	}

	return c.Send(completionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// Helper function to get replacement word for "already known" words
func (s *HandlerService) getReplacementWord(progress *domain.LessonProgress, currentWordIndex int) (*domain.Card, error) {
	// Start looking for replacement words after all words that will be shown in regular sets
	// Calculate how many words will be shown in total
	totalWordsToShow := progress.LessonData.Lesson.WordsPerLesson

	// Start looking from the end of all regular words to avoid conflicts
	startIndex := totalWordsToShow

	// Check if we have more words available
	if startIndex >= len(progress.LessonData.Cards) {
		return nil, fmt.Errorf("no more replacement words available")
	}

	// Find the next available word that hasn't been learned yet
	for i := startIndex; i < len(progress.LessonData.Cards); i++ {
		wordID := progress.LessonData.Cards[i].WordID

		// Check if this word is already in the current set
		isInCurrentSet := false
		for _, currentWord := range progress.WordsInCurrentSet {
			if currentWord.WordID == wordID {
				isInCurrentSet = true
				break
			}
		}

		// Check if this word was already learned or shown
		isAlreadyUsed := false
		for _, learnedWord := range progress.WordsLearned {
			if learnedWord.WordID == wordID {
				isAlreadyUsed = true
				break
			}
		}

		// If word is not in current set and not already used, use it
		if !isInCurrentSet && !isAlreadyUsed {
			replacementCard := progress.LessonData.Cards[i]
			return &replacementCard, nil
		}
	}

	return nil, fmt.Errorf("no more replacement words available")
}

// Helper function to replace a word in the current set
func (s *HandlerService) replaceWordInCurrentSet(ctx context.Context, userID int64, wordIndex int) error {
	return s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		// Get replacement word
		replacementWord, err := s.getReplacementWord(p, wordIndex)
		if err != nil {
			return err
		}

		// Replace the word at the specified index
		if wordIndex < len(p.WordsInCurrentSet) {
			p.WordsInCurrentSet[wordIndex] = *replacementWord
		}

		p.LastActivity = time.Now()
		return nil
	})
}

// Helper function to format duration
func (s *HandlerService) formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	if minutes < 1 {
		return "меньше минуты"
	} else if minutes == 1 {
		return "1 минута"
	} else if minutes < 60 {
		return fmt.Sprintf("%d минут", minutes)
	} else {
		hours := minutes / 60
		mins := minutes % 60
		if mins == 0 {
			return fmt.Sprintf("%d час(ов)", hours)
		}
		return fmt.Sprintf("%d час(ов) %d минут", hours, mins)
	}
}

// handleUnauthenticatedUser handles users without JWT tokens
func (s *HandlerService) handleUnauthenticatedUser(ctx context.Context, c tele.Context, userID int64) error {
	// Check if user has completed onboarding (questionnaire + CEFR test)
	userProgress, err := s.GetUserProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user progress", zap.Error(err))
		return err
	}

	if userProgress.CEFRLevel == "" {
		// User hasn't completed onboarding - redirect to onboarding
		return s.redirectToOnboarding(ctx, c, userID)
	}

	// User has completed onboarding but isn't authenticated - offer authentication
	linkResponse, err := s.apiClient.CreateLinkToken(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to create link token", zap.Error(err))
		return c.Send("❌ Произошла ошибка. Попробуйте позже.")
	}

	// Store linking data
	err = s.stateManager.StoreUserLinkingData(ctx, userID, linkResponse.Token, time.Hour)
	if err != nil {
		s.logger.Error("Failed to store linking data", zap.Error(err))
	}

	authText := fmt.Sprintf(
		"🔐 *Требуется авторизация*\n\n"+
			"Для доступа к персональным урокам необходимо связать ваш аккаунт Telegram с аккаунтом Google.\n\n"+
			"🎯 **Это позволит:**\n"+
			"• Сохранить ваш прогресс (уровень %s)\n"+
			"• Получить персональные уроки\n"+
			"• Синхронизировать данные между устройствами\n\n"+
			"🔗 *Ссылка для авторизации:*\n[Нажмите здесь для авторизации](%s)\n\n"+
			"После авторизации вернитесь и нажмите \"Проверить связь\".",
		userProgress.CEFRLevel,
		linkResponse.LinkURL,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🔄 Проверить связь", Data: "auth:check_link"},
				{Text: "❓ Помощь", Data: "help:auth"},
			},
		},
	}

	return c.Send(authText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// redirectToOnboarding redirects user to complete onboarding first
func (s *HandlerService) redirectToOnboarding(ctx context.Context, c tele.Context, userID int64) error {
	onboardingText := fmt.Sprintf(
		"👋 *Привет, %s!*\n\n"+
			"Перед началом изучения давайте сначала настроим твой профиль.\n\n"+
			"📋 **Что нужно сделать:**\n"+
			"• Ответить на пару вопросов\n"+
			"• Пройти тест уровня CEFR\n"+
			"• Создать аккаунт для сохранения прогресса\n\n"+
			"Займет всего 3-5 минут! 🕐",
		c.Sender().FirstName,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🚀 Начать настройку", Data: "auth:new_user"},
				{Text: "🔗 У меня есть аккаунт", Data: "auth:existing_user"},
			},
		},
	}

	return c.Send(onboardingText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// sendWordVoiceMessage generates and sends a voice message for a word
func (s *HandlerService) sendWordVoiceMessage(ctx context.Context, c tele.Context, word string) error {
	// Generate voice message
	audioData, err := s.ttsService.GenerateWordVoiceMessage(word)
	if err != nil {
		return fmt.Errorf("failed to generate voice message: %w", err)
	}

	// Validate audio data
	if err := s.ttsService.ValidateAudioData(audioData); err != nil {
		return fmt.Errorf("invalid audio data: %w", err)
	}

	// Create temporary file for the voice message
	tempFile, err := s.ttsService.CreateVoiceMessageFromBytes(audioData, word)
	if err != nil {
		return fmt.Errorf("failed to create voice file: %w", err)
	}

	// Clean up temporary file after sending
	defer func() {
		if err := os.Remove(tempFile); err != nil {
			s.logger.Warn("Failed to clean up temp voice file", zap.Error(err))
		}
	}()

	// Send voice message
	voice := &tele.Voice{File: tele.FromDisk(tempFile), Caption: word}
	return c.Send(voice)
}
