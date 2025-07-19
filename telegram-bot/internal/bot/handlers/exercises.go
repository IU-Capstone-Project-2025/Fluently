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

	return text[:wordIndex] + strings.Repeat("_", len(word)) + text[wordIndex+len(word):]
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
		"Упражнение %d из %d\n\n"+
			"Выберите правильный вариант:\n\n"+
			"%s\n\n"+
			"Выберите правильный ответ:",
		progress.ExerciseIndex+1,
		len(progress.WordsInCurrentSet),
		processedTemplate,
	)

	// Create option buttons
	var buttons [][]tele.InlineButton
	for i, option := range exercise.Data.PickOptions {
		buttons = append(buttons, []tele.InlineButton{
			{
				Text: fmt.Sprintf("%s", option),
				Data: fmt.Sprintf("exercise:pick_option:%d:%s", i, option),
			},
		})
	}

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

	currentWord := progress.WordsInCurrentSet[progress.ExerciseIndex]
	exercise := currentWord.Exercise
	isCorrect := selectedOption == exercise.Data.CorrectAnswer

	return s.processExerciseAnswer(ctx, c, userID, currentWord, isCorrect, selectedOption)
}

// HandleTextInputAnswer handles text input answers for exercises
func (s *HandlerService) HandleTextInputAnswer(ctx context.Context, c tele.Context, userID int64, userAnswer string) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	currentWord := progress.WordsInCurrentSet[progress.ExerciseIndex]
	exercise := currentWord.Exercise

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

		// Add word to learned words if correct
		wordProgress := domain.WordProgress{
			Word:            word.Word,
			LearnedAt:       time.Now(),
			ConfidenceScore: 100,
			CntReviewed:     1,
		}

		err = s.stateManager.AddWordProgress(ctx, userID, wordProgress)
		if err != nil {
			s.logger.Error("Failed to add word progress", zap.Error(err))
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

		// Add word with low confidence if incorrect
		wordProgress := domain.WordProgress{
			Word:            word.Word,
			LearnedAt:       time.Now(),
			ConfidenceScore: 0,
			CntReviewed:     0,
		}

		err = s.stateManager.AddWordProgress(ctx, userID, wordProgress)
		if err != nil {
			s.logger.Error("Failed to add word progress", zap.Error(err))
		}
	}

	// Update exercise index
	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		p.ExerciseIndex++
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

// completeCurrentSet handles completion of the current set of 3 words
func (s *HandlerService) completeCurrentSet(ctx context.Context, c tele.Context, userID int64) error {
	progress, err := s.stateManager.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	// Set state to set complete
	if err := s.stateManager.SetState(ctx, userID, fsm.StateSetComplete); err != nil {
		return err
	}

	// Calculate set statistics
	correctCount := 0
	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 {
			correctCount++
		}
	}

	setCompleteText := fmt.Sprintf(
		"🎉 *Набор завершен!*\n\n"+
			"📊 Результаты:\n"+
			"✅ Правильно: %d из %d\n"+
			"📈 Точность: %.1f%%\n\n"+
			"Изученные слова:\n"+
			"• %s\n"+
			"• %s\n"+
			"• %s",
		correctCount,
		len(progress.WordsInCurrentSet),
		float64(correctCount)/float64(len(progress.WordsInCurrentSet))*100,
		progress.WordsInCurrentSet[0].Word,
		progress.WordsInCurrentSet[1].Word,
		progress.WordsInCurrentSet[2].Word,
	)

	// Check if lesson is complete
	if progress.LearnedCount >= progress.LessonData.Lesson.WordsPerLesson {
		setCompleteText += "\n\n🏆 *Поздравляем! Урок завершен!*"

		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{
					{Text: "📊 Финальная статистика", Data: "lesson:final_stats"},
				},
			},
		}

		err = c.Send(setCompleteText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
		if err != nil {
			return err
		}

		return s.completeLessonFlow(ctx, c, userID, progress)
	}

	// Continue with next set
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "➡️ Следующие 3 слова", Data: "lesson:start_word_set"},
				{Text: "📊 Статистика", Data: "lesson:stats"},
			},
		},
	}

	return c.Send(setCompleteText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// completeLessonFlow handles completion of the entire lesson
func (s *HandlerService) completeLessonFlow(ctx context.Context, c tele.Context, userID int64, progress *domain.LessonProgress) error {
	// Set state to lesson complete
	if err := s.stateManager.SetState(ctx, userID, fsm.StateLessonComplete); err != nil {
		return err
	}

	// Calculate final statistics
	totalWords := len(progress.WordsLearned)
	correctWords := 0
	for _, wordProgress := range progress.WordsLearned {
		if wordProgress.ConfidenceScore > 0 {
			correctWords++
		}
	}

	duration := time.Since(progress.StartTime)
	accuracy := float64(correctWords) / float64(totalWords) * 100

	finalText := fmt.Sprintf(
		"🏆 *Урок завершен!*\n\n"+
			"📊 *Финальная статистика:*\n"+
			"✅ Слов выучено: %d\n"+
			"🎯 Правильных ответов: %d из %d\n"+
			"📈 Точность: %.1f%%\n"+
			"⏱ Время урока: %s\n\n"+
			"🎉 Отличная работа! Продолжайте изучение!",
		progress.LearnedCount,
		correctWords,
		totalWords,
		accuracy,
		s.formatDuration(duration),
	)

	// Send progress to backend
	token, err := s.stateManager.GetJWTToken(ctx, userID)
	if err == nil {
		err = s.apiClient.SendLessonProgress(ctx, token, progress.WordsLearned)
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

	currentWord := progress.WordsInCurrentSet[progress.ExerciseIndex]

	// Mark word as skipped (low confidence)
	wordProgress := domain.WordProgress{
		Word:            currentWord.Word,
		LearnedAt:       time.Now(),
		ConfidenceScore: 25, // Low but not zero for skipped
		CntReviewed:     0,
	}

	err = s.stateManager.AddWordProgress(ctx, userID, wordProgress)
	if err != nil {
		s.logger.Error("Failed to add skipped word progress", zap.Error(err))
	}

	// Update exercise index
	err = s.stateManager.UpdateLessonProgress(ctx, userID, func(p *domain.LessonProgress) error {
		p.ExerciseIndex++
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

	currentWord := progress.WordsInCurrentSet[progress.ExerciseIndex]
	exercise := currentWord.Exercise

	var hintText string

	// Provide different hints based on exercise type
	switch exercise.Type {
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
