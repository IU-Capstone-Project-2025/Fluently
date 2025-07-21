package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/bot/fsm"
	"telegram-bot/internal/domain"
)

// replaceWordWithUnderscores replaces a word in a text with underscores
func replaceWordWithUnderscores(text, word string) string {
	lowerText := strings.ToLower(text)
	lowerWord := strings.ToLower(word)

	wordIndex := strings.Index(lowerText, lowerWord)
	if wordIndex == -1 {
		return text
	}

	// Replace the word with underscores, preserving original case
	originalWord := text[wordIndex : wordIndex+len(word)]
	underscores := strings.Repeat("_", len(originalWord))

	return text[:wordIndex] + underscores + text[wordIndex+len(word):]
}

// showPickOptionSentenceExercise displays a multiple choice exercise with sentence template
func (s *HandlerService) showPickOptionSentenceExercise(ctx context.Context, c tele.Context, userID int64, word domain.Card, exercise domain.Exercise) error {
	if err := s.stateManager.SetState(ctx, userID, fsm.StatePickOptionSentence); err != nil {
		return err
	}

	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	// Replace the word with underscores in the template
	processedTemplate := replaceWordWithUnderscores(exercise.Data.Template, word.Word)

	exerciseText := fmt.Sprintf(
		"Упражнение %d из %d\n\n*%s*\n\nВыберите правильный вариант для того чтобы вставить в предложение:\n\n",
		progress.ExerciseIndex+1,
		len(progress.WordsInCurrentSet),
		processedTemplate,
	)

	// Create option buttons
	var buttons [][]tele.InlineButton

	// Add option buttons
	for i, option := range exercise.Data.PickOptions {
		buttons = append(buttons, []tele.InlineButton{
			{
				Text: fmt.Sprintf("%s", option),
				Data: fmt.Sprintf("exercise:pick_option:%d:%s", i, option),
			},
		})
	}

	// Add hint button
	buttons = append(buttons, []tele.InlineButton{
		{Text: "💡 Подсказка", Data: "exercise:hint"},
	})

	keyboard := &tele.ReplyMarkup{InlineKeyboard: buttons}

	return c.Send(exerciseText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// showWriteWordTranslationExercise displays a text input exercise for writing word from translation
func (s *HandlerService) showWriteWordTranslationExercise(ctx context.Context, c tele.Context, userID int64, word domain.Card, exercise domain.Exercise) error {
	if err := s.stateManager.SetState(ctx, userID, fsm.StateWriteWordTranslation); err != nil {
		return err
	}

	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	exerciseText := fmt.Sprintf(
		"Упражнение %d из %d\n\n"+
			"Напишите английское слово по переводу:\n\n"+
			"Перевод: %s\n\n"+
			"Введите английское слово:",
		progress.ExerciseIndex+1,
		len(progress.WordsInCurrentSet),
		exercise.Data.Translation,
	)

	// Set state to waiting for text input
	if err := s.stateManager.SetState(ctx, userID, fsm.StateWaitingForTextInput); err != nil {
		return err
	}

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🔄 Пропустить", Data: "exercise:skip"},
				{Text: "💡 Подсказка", Data: "exercise:hint"},
			},
		},
	}

	return c.Send(exerciseText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// showTranslateRuToEnExercise displays a translation exercise from Russian to English
func (s *HandlerService) showTranslateRuToEnExercise(ctx context.Context, c tele.Context, userID int64, word domain.Card, exercise domain.Exercise) error {
	if err := s.stateManager.SetState(ctx, userID, fsm.StateTranslateRuToEn); err != nil {
		return err
	}

	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	// Check if it has pick options (multiple choice) or text input
	if len(exercise.Data.PickOptions) > 0 {
		return s.showTranslateRuToEnMultipleChoice(ctx, c, userID, word, exercise, progress)
	} else {
		return s.showTranslateRuToEnTextInput(ctx, c, userID, word, exercise, progress)
	}
}

// showTranslateRuToEnMultipleChoice displays Russian to English translation with multiple choice
func (s *HandlerService) showTranslateRuToEnMultipleChoice(ctx context.Context, c tele.Context, userID int64, word domain.Card, exercise domain.Exercise, progress *domain.LessonProgress) error {
	exerciseText := fmt.Sprintf(
		"Упражнение %d из %d\n\n"+
			"Переведите на английский:\n\n"+
			"%s\n\n"+
			"Выберите правильный перевод:",
		progress.ExerciseIndex+1,
		len(progress.WordsInCurrentSet),
		exercise.Data.Text,
	)

	// Create option buttons
	var buttons [][]tele.InlineButton
	for i, option := range exercise.Data.PickOptions {
		buttons = append(buttons, []tele.InlineButton{
			{
				Text: fmt.Sprintf("%s", option),
				Data: fmt.Sprintf("exercise:translate_option:%d:%s", i, option),
			},
		})
	}

	keyboard := &tele.ReplyMarkup{InlineKeyboard: buttons}

	return c.Send(exerciseText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// showTranslateRuToEnTextInput displays Russian to English translation with text input
func (s *HandlerService) showTranslateRuToEnTextInput(ctx context.Context, c tele.Context, userID int64, word domain.Card, exercise domain.Exercise, progress *domain.LessonProgress) error {
	exerciseText := fmt.Sprintf(
		"Упражнение %d из %d\n\n"+
			"Переведите на английский:\n\n"+
			"%s\n\n"+
			"Введите английский перевод:",
		progress.ExerciseIndex+1,
		len(progress.WordsInCurrentSet),
		exercise.Data.Text,
	)

	// Set state to waiting for text input
	if err := s.stateManager.SetState(ctx, userID, fsm.StateWaitingForTextInput); err != nil {
		return err
	}

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🔄 Пропустить", Data: "exercise:skip"},
				{Text: "💡 Подсказка", Data: "exercise:hint"},
			},
		},
	}

	return c.Send(exerciseText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandlePickOptionAnswer handles answer to pick option exercises
func (s *HandlerService) HandlePickOptionAnswer(ctx context.Context, c tele.Context, userID int64, optionIndex int, selectedOption string) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.Debug("Processing pick option answer",
		zap.Int64("user_id", userID),
		zap.String("current_phase", progress.CurrentPhase),
		zap.Int("exercise_index", progress.ExerciseIndex),
		zap.Int("retry_index", progress.RetryIndex),
		zap.Int("words_in_current_set", len(progress.WordsInCurrentSet)),
		zap.Int("retry_words_count", len(progress.RetryWords)))

	var currentWord domain.Card
	var exercise domain.Exercise

	// Check if we're in retry mode
	if progress.CurrentPhase == "retry" {
		// Validate retry index
		if progress.RetryIndex >= len(progress.RetryWords) {
			s.logger.Error("Retry index out of bounds",
				zap.Int64("user_id", userID),
				zap.Int("retry_index", progress.RetryIndex),
				zap.Int("retry_words_length", len(progress.RetryWords)))
			return fmt.Errorf("retry index out of bounds: %d >= %d", progress.RetryIndex, len(progress.RetryWords))
		}
		currentWord = progress.RetryWords[progress.RetryIndex]
	} else {
		// Validate exercise index
		if progress.ExerciseIndex >= len(progress.WordsInCurrentSet) {
			s.logger.Error("Exercise index out of bounds",
				zap.Int64("user_id", userID),
				zap.Int("exercise_index", progress.ExerciseIndex),
				zap.Int("words_in_current_set_length", len(progress.WordsInCurrentSet)))
			return fmt.Errorf("exercise index out of bounds: %d >= %d", progress.ExerciseIndex, len(progress.WordsInCurrentSet))
		}
		currentWord = progress.WordsInCurrentSet[progress.ExerciseIndex]
	}

	exercise = currentWord.Exercise
	isCorrect := selectedOption == exercise.Data.CorrectAnswer

	return s.processExerciseAnswer(ctx, c, userID, currentWord, isCorrect, selectedOption)
}

// HandleTextInputAnswer handles text input answers for exercises
func (s *HandlerService) HandleTextInputAnswer(ctx context.Context, c tele.Context, userID int64, userAnswer string) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.Debug("Processing text input answer",
		zap.Int64("user_id", userID),
		zap.String("current_phase", progress.CurrentPhase),
		zap.Int("exercise_index", progress.ExerciseIndex),
		zap.Int("retry_index", progress.RetryIndex),
		zap.Int("words_in_current_set", len(progress.WordsInCurrentSet)),
		zap.Int("retry_words_count", len(progress.RetryWords)))

	var currentWord domain.Card
	var exercise domain.Exercise

	// Check if we're in retry mode
	if progress.CurrentPhase == "retry" {
		// Validate retry index
		if progress.RetryIndex >= len(progress.RetryWords) {
			s.logger.Error("Retry index out of bounds",
				zap.Int64("user_id", userID),
				zap.Int("retry_index", progress.RetryIndex),
				zap.Int("retry_words_length", len(progress.RetryWords)))
			return fmt.Errorf("retry index out of bounds: %d >= %d", progress.RetryIndex, len(progress.RetryWords))
		}
		currentWord = progress.RetryWords[progress.RetryIndex]
	} else {
		// Validate exercise index
		if progress.ExerciseIndex >= len(progress.WordsInCurrentSet) {
			s.logger.Error("Exercise index out of bounds",
				zap.Int64("user_id", userID),
				zap.Int("exercise_index", progress.ExerciseIndex),
				zap.Int("words_in_current_set_length", len(progress.WordsInCurrentSet)))
			return fmt.Errorf("exercise index out of bounds: %d >= %d", progress.ExerciseIndex, len(progress.WordsInCurrentSet))
		}
		currentWord = progress.WordsInCurrentSet[progress.ExerciseIndex]
	}

	exercise = currentWord.Exercise

	// Clean and compare answers
	cleanUserAnswer := strings.ToLower(strings.TrimSpace(userAnswer))
	cleanCorrectAnswer := strings.ToLower(strings.TrimSpace(exercise.Data.CorrectAnswer))

	isCorrect := cleanUserAnswer == cleanCorrectAnswer

	return s.processExerciseAnswer(ctx, c, userID, currentWord, isCorrect, userAnswer)
}

// processExerciseAnswer processes the result of an exercise answer
func (s *HandlerService) processExerciseAnswer(ctx context.Context, c tele.Context, userID int64, word domain.Card, isCorrect bool, userAnswer string) error {
	var err error

	exercise := word.Exercise

	// Create feedback message
	var feedbackText string
	var emoji string

	if isCorrect {
		emoji = "✅"
		feedbackText = fmt.Sprintf(
			"%s *Правильно!*\n\n"+
				"🔤 %s - %s\n"+
				"✅ Ваш ответ: %s",
			emoji,
			word.Word,
			word.Translation,
			userAnswer,
		)

		// Check if this is a retry exercise
		progress, err := s.stateManager.GetLessonProgress(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to get lesson progress", zap.Error(err))
		} else if progress.CurrentPhase == "retry" {
			// This is a retry exercise - update the existing word progress and remove from retry queue
			err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
				// Find and update the existing word progress
				for i, wp := range p.WordsLearned {
					if wp.WordID == word.WordID {
						p.WordsLearned[i].ConfidenceScore = 100
						p.WordsLearned[i].CntReviewed++
						break
					}
				}

				// Remove word from retry queue
				for i, retryWord := range p.RetryWords {
					if retryWord.WordID == word.WordID {
						p.RetryWords = append(p.RetryWords[:i], p.RetryWords[i+1:]...)
						break
					}
				}

				p.LastActivity = time.Now()
				return nil
			})
			if err != nil {
				s.logger.Error("Failed to update retry word progress", zap.Error(err))
			}
		} else {
			// Regular exercise - add new word progress
			wordProgress := domain.WordProgress{
				Word:            word.Word,
				Translation:     word.Translation,
				WordID:          word.WordID,
				LearnedAt:       time.Now(),
				ConfidenceScore: 100,
				CntReviewed:     1,
			}

			err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
				p.WordsLearned = append(p.WordsLearned, wordProgress)
				p.LearnedCount++
				p.LastActivity = time.Now()
				return nil
			})
			if err != nil {
				s.logger.Error("Failed to add word progress", zap.Error(err))
			}
		}
	} else {
		emoji = "❌"
		feedbackText = fmt.Sprintf(
			"%s *Неправильно*\n\n"+
				"🔤 %s - %s\n"+
				"❌ Ваш ответ: %s\n"+
				"✅ Правильный ответ: %s",
			emoji,
			word.Word,
			word.Translation,
			userAnswer,
			exercise.Data.CorrectAnswer,
		)

		// Check if this is a retry exercise
		progress, err := s.stateManager.GetLessonProgress(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to get lesson progress", zap.Error(err))
		} else if progress.CurrentPhase == "retry" {
			// This is a retry exercise - update existing word progress and keep in retry queue
			err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
				// Find and update the existing word progress
				for i, wp := range p.WordsLearned {
					if wp.WordID == word.WordID {
						p.WordsLearned[i].ConfidenceScore = 0
						p.WordsLearned[i].CntReviewed++
						break
					}
				}

				// Move word to end of retry queue for another attempt
				for i, retryWord := range p.RetryWords {
					if retryWord.WordID == word.WordID {
						// Remove from current position
						p.RetryWords = append(p.RetryWords[:i], p.RetryWords[i+1:]...)
						// Add to end of queue
						p.RetryWords = append(p.RetryWords, word)
						break
					}
				}

				p.LastActivity = time.Now()
				return nil
			})
			if err != nil {
				s.logger.Error("Failed to update retry word progress", zap.Error(err))
			}
		} else {
			// Regular exercise - add new word progress
			wordProgress := domain.WordProgress{
				Word:            word.Word,
				Translation:     word.Translation,
				WordID:          word.WordID,
				LearnedAt:       time.Now(),
				ConfidenceScore: 0,
				CntReviewed:     0,
			}

			// Add to badly answered words list
			badlyAnsweredWord := domain.BadlyAnsweredWord{
				WordID: word.WordID,
			}

			err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
				p.WordsLearned = append(p.WordsLearned, wordProgress)
				p.BadlyAnsweredWords = append(p.BadlyAnsweredWords, badlyAnsweredWord)
				// Add word to retry queue for later practice
				p.RetryWords = append(p.RetryWords, word)
				p.LearnedCount++
				p.LastActivity = time.Now()
				return nil
			})
			if err != nil {
				s.logger.Error("Failed to add word progress", zap.Error(err))
			}
		}
	}

	// Update exercise index based on current phase
	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		if p.CurrentPhase == "retry" {
			p.RetryIndex++
		} else {
			p.ExerciseIndex++
		}
		p.LastActivity = time.Now()
		return nil
	})
	if err != nil {
		return err
	}

	// Send feedback
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "➡️ Продолжить", Data: "exercise:next"},
			},
		},
	}

	return c.Send(feedbackText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleExerciseNext handles moving to the next exercise
func (s *HandlerService) HandleExerciseNext(ctx context.Context, c tele.Context, userID int64) error {
	return s.showNextExercise(ctx, c, userID)
}

// completeCurrentSet handles completion of the current set of words
func (s *HandlerService) completeCurrentSet(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	// Check if lesson is complete
	if progress.LearnedCount >= progress.LessonData.Lesson.WordsPerLesson {
		// Lesson is complete - go directly to final statistics
		return s.completeLessonFlow(ctx, c, userID, progress)
	}

	// Set state to set complete
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSetComplete); err != nil {
		return err
	}

	// Continue with next set
	// Calculate how many words are left
	wordsLearnedInLesson := progress.LearnedCount - progress.AlreadyKnownCount
	wordsLeft := progress.LessonData.Lesson.WordsPerLesson - wordsLearnedInLesson
	nextSetSize := 3
	if wordsLeft < 3 {
		nextSetSize = wordsLeft
	}

	completionText := fmt.Sprintf(
		"✅ *Набор слов завершен!*\n\n"+
			"🎯 Вы успешно изучили %d слов в этом наборе.\n\n"+
			"Продолжайте изучение или посмотрите статистику:",
		len(progress.WordsInCurrentSet),
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: fmt.Sprintf("➡️ Следующие %d слова", nextSetSize), Data: "lesson:start_word_set"},
				{Text: "📊 Статистика", Data: "lesson:stats"},
			},
		},
	}

	return c.Send(completionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// completeLessonFlow handles completion of the entire lesson
func (s *HandlerService) completeLessonFlow(ctx context.Context, c tele.Context, userID int64, progress *domain.LessonProgress) error {
	// Check if there are words to retry
	if len(progress.RetryWords) > 0 {
		return s.startRetryPhase(ctx, c, userID)
	}

	// Set state to lesson complete
	if err := s.stateManager.SetState(ctx, userID, fsm.StateLessonComplete); err != nil {
		return err
	}

	// Calculate final statistics - exclude "already known" words from the count
	wellAnsweredWords := 0
	alreadyKnownCorrectWords := 0

	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 {
			// Check if this word was marked as "already known"
			if wordProgress.AlreadyKnown {
				// This is an "already known" word
				alreadyKnownCorrectWords++
			} else {
				// This is a newly learned word
				wellAnsweredWords++
			}
		}
	}

	duration := time.Since(progress.StartTime)
	// Calculate accuracy based on newly learned words only
	accuracy := float64(wellAnsweredWords) / float64(progress.LearnedCount-progress.AlreadyKnownCount) * 100

	// Build list of learned words
	var learnedWordsList strings.Builder
	learnedWordsList.WriteString(fmt.Sprintf("📚 *За урок выучено %d слов:*\n\n", wellAnsweredWords))

	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 && !wordProgress.AlreadyKnown {
			learnedWordsList.WriteString(fmt.Sprintf("#%s - %s\n", wordProgress.Word, wordProgress.Translation))
		}
	}

	finalText := fmt.Sprintf(
		"🏆 *Урок завершен!*\n\n"+
			"📊 *Финальная статистика:*\n"+
			"✅ Слов выучено: %d\n"+
			"💡 Уже знал: %d слов\n"+
			"🎯 Правильно: %d из %d\n"+
			"📈 Точность: %.1f%%\n"+
			"⏱️ Время урока: %s\n\n"+
			"%s\n"+
			"🎉 Отличная работа! Продолжайте изучение!",
		progress.LearnedCount,
		progress.AlreadyKnownCount,
		wellAnsweredWords,                         // Only newly learned words
		progress.LessonData.Lesson.WordsPerLesson, // Show correct answers vs target words
		accuracy,
		s.formatDuration(duration),
		learnedWordsList.String(),
	)

	// Send progress to backend
	token, err := s.stateManager.GetJWTToken(ctx, userID)
	if err == nil {
		err = s.apiClient.SendLessonProgress(ctx, token, progress.WordsLearned, progress.BadlyAnsweredWords)
		if err != nil {
			s.logger.Error("Failed to send lesson progress to backend", zap.Error(err))
		}
	}

	// Clear lesson progress
	err = s.stateManager.ClearLessonProgress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to clear lesson progress", zap.Error(err))
	}

	// Reset state to start
	if err := s.stateManager.SetState(ctx, userID, fsm.StateStart); err != nil {
		s.logger.Error("Failed to reset state", zap.Error(err))
	}

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🚀 Новый урок", Data: "lesson:new"},
				{Text: "📊 Общая статистика", Data: "stats:overall"},
			},
		},
	}

	return c.Send(finalText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleSkipExercise handles skipping an exercise
func (s *HandlerService) HandleSkipExercise(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	var currentWord domain.Card

	// Check if we're in retry mode
	if progress.CurrentPhase == "retry" {
		// Validate retry index
		if progress.RetryIndex >= len(progress.RetryWords) {
			return fmt.Errorf("retry index out of bounds: %d >= %d", progress.RetryIndex, len(progress.RetryWords))
		}
		currentWord = progress.RetryWords[progress.RetryIndex]
	} else {
		// Validate exercise index
		if progress.ExerciseIndex >= len(progress.WordsInCurrentSet) {
			return fmt.Errorf("exercise index out of bounds: %d >= %d", progress.ExerciseIndex, len(progress.WordsInCurrentSet))
		}
		currentWord = progress.WordsInCurrentSet[progress.ExerciseIndex]
	}

	// Mark word as skipped (low confidence)
	wordProgress := domain.WordProgress{
		Word:            currentWord.Word,
		Translation:     currentWord.Translation,
		WordID:          currentWord.WordID,
		LearnedAt:       time.Now(),
		ConfidenceScore: 25, // Low but not zero for skipped
		CntReviewed:     0,
	}

	// Add to badly answered words list
	badlyAnsweredWord := domain.BadlyAnsweredWord{
		WordID: currentWord.WordID,
	}

	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		p.WordsLearned = append(p.WordsLearned, wordProgress)
		p.BadlyAnsweredWords = append(p.BadlyAnsweredWords, badlyAnsweredWord)
		p.LearnedCount++
		p.LastActivity = time.Now()
		return nil
	})
	if err != nil {
		s.logger.Error("Failed to add skipped word progress", zap.Error(err))
	}

	// Update exercise index based on current phase
	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		if p.CurrentPhase == "retry" {
			p.RetryIndex++
		} else {
			p.ExerciseIndex++
		}
		p.LastActivity = time.Now()
		return nil
	})
	if err != nil {
		return err
	}

	skipText := fmt.Sprintf(
		"⏭ *Упражнение пропущено*\n\n"+
			"🔤 %s - %s\n\n"+
			"Рекомендуем повторить это слово позже.",
		currentWord.Word,
		currentWord.Translation,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "➡️ Продолжить", Data: "exercise:next"},
			},
		},
	}

	return c.Send(skipText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleExerciseHint provides a hint for the current exercise
func (s *HandlerService) HandleExerciseHint(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	var currentWord domain.Card
	var exercise domain.Exercise

	// Check if we're in retry mode
	if progress.CurrentPhase == "retry" {
		// Validate retry index
		if progress.RetryIndex >= len(progress.RetryWords) {
			return fmt.Errorf("retry index out of bounds: %d >= %d", progress.RetryIndex, len(progress.RetryWords))
		}
		currentWord = progress.RetryWords[progress.RetryIndex]
	} else {
		// Validate exercise index
		if progress.ExerciseIndex >= len(progress.WordsInCurrentSet) {
			return fmt.Errorf("exercise index out of bounds: %d >= %d", progress.ExerciseIndex, len(progress.WordsInCurrentSet))
		}
		currentWord = progress.WordsInCurrentSet[progress.ExerciseIndex]
	}

	exercise = currentWord.Exercise

	var hintText string

	// Provide different hints based on exercise type
	switch exercise.Type {
	case "pick_option_sentence":
		// For pick option sentence, show the sentence translation and word meaning
		// Try to find the sentence translation from the word's sentences
		sentenceTranslation := "Перевод недоступен"
		for _, sentence := range currentWord.Sentences {
			if sentence.Text == exercise.Data.Template {
				sentenceTranslation = sentence.Translation
				break
			}
		}

		// If no exact match found, use the first available sentence translation
		if sentenceTranslation == "Перевод недоступен" && len(currentWord.Sentences) > 0 {
			sentenceTranslation = currentWord.Sentences[0].Translation
		}

		hintText = fmt.Sprintf("💡 *Подсказка:*\n\n"+
			"*Предложение:* %s\n"+
			"*Перевод предложения:* %s\n\n"+
			"*Слово:* %s - %s\n"+
			"*Правильный ответ:* %s",
			exercise.Data.Template,
			sentenceTranslation,
			currentWord.Word,
			currentWord.Translation,
			exercise.Data.CorrectAnswer)
	case "write_word_from_translation":
		word := exercise.Data.CorrectAnswer
		if len(word) > 3 {
			hintText = fmt.Sprintf("💡 *Подсказка:*\n\nСлово начинается на \"%s\" и содержит %d букв",
				strings.ToUpper(string(word[0])), len(word))
		} else {
			hintText = fmt.Sprintf("💡 *Подсказка:*\n\nСлово содержит %d букв", len(word))
		}
	case "translate_ru_to_en":
		word := exercise.Data.CorrectAnswer
		if len(word) > 3 {
			hintText = fmt.Sprintf("💡 *Подсказка:*\n\nПеревод начинается на \"%s\" и содержит %d букв",
				strings.ToUpper(string(word[0])), len(word))
		} else {
			hintText = fmt.Sprintf("💡 *Подсказка:*\n\nПеревод содержит %d букв", len(word))
		}
	default:
		hintText = "💡 *Подсказка:*\n\nВнимательно прочитайте предложение и подумайте о контексте."
	}

	return c.Send(hintText, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
}
