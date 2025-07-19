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

	// Generate new lesson from backend
	lessonResponse, err := s.apiClient.GenerateLesson(ctx, token)
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
		LessonData:        lessonResponse,
		CurrentWordIndex:  0,
		CurrentPhase:      "showing_words",
		WordsInCurrentSet: []domain.Card{},
		CurrentSetIndex:   0,
		ExerciseIndex:     0,
		WordsLearned:      []domain.WordProgress{},
		StartTime:         time.Now(),
		LastActivity:      time.Now(),
		LearnedCount:      0,
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
		"📚 *Набор слов #%d*\n\n"+
			"Сейчас изучим 3 новых слова:\n"+
			"1️⃣ %s - %s\n"+
			"2️⃣ %s - %s\n"+
			"3️⃣ %s - %s\n\n"+
			"Нажмите \"Изучать\", чтобы посмотреть каждое слово подробно.",
		progress.CurrentSetIndex,
		nextSet[0].Word, nextSet[0].Translation,
		nextSet[1].Word, nextSet[1].Translation,
		nextSet[2].Word, nextSet[2].Translation,
	)

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
		if i >= 2 { // Limit to 2 examples
			break
		}
		detailText += fmt.Sprintf(
			"\n<blockquote>%s</blockquote>\n%s",
			sentence.Text,
			sentence.Translation,
		)
	}

	// Create navigation buttons
	var buttons [][]tele.InlineButton

	// Navigation buttons
	var navButtons []tele.InlineButton

	if wordIndex < 2 {
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

	// Add repeat voice button
	buttons = append(buttons, []tele.InlineButton{
		{Text: "🔊 Повторить произношение", Data: fmt.Sprintf("voice:repeat:%s", word.Word)},
	})

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
			"Проверим, как вы усвоили эти 3 слова:\n"+
			"• %s\n"+
			"• %s\n"+
			"• %s\n\n"+
			"Будет %d упражнения. Готовы?",
		progress.WordsInCurrentSet[0].Word,
		progress.WordsInCurrentSet[1].Word,
		progress.WordsInCurrentSet[2].Word,
		len(progress.WordsInCurrentSet), // 3 exercises
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				{Text: "🚀 Начать упражнения", Data: "lesson:start_exercises"},
				{Text: "📖 Повторить слова", Data: "lesson:show_word:0"},
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

// Helper function to get next set of 3 words
func (s *HandlerService) getNextWordSet(progress *domain.LessonProgress) ([]domain.Card, error) {
	startIndex := progress.CurrentSetIndex * 3
	cards := progress.LessonData.Cards

	if startIndex >= len(cards) {
		return nil, fmt.Errorf("no more words available")
	}

	endIndex := startIndex + 3
	if endIndex > len(cards) {
		endIndex = len(cards)
	}

	return cards[startIndex:endIndex], nil
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
