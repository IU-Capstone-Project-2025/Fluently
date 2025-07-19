package handlers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"

	"telegram-bot/internal/api"
	"telegram-bot/internal/bot/fsm"
)

// CEFR test questions organized by level groups
var cefrTestQuestions = map[int][]TestQuestion{
	1: { // A1 level
		{Word: "hello", Translation: "привет", Options: []string{"привет", "пока", "спасибо", "извините"}, Correct: 0},
		{Word: "cat", Translation: "кот", Options: []string{"собака", "кот", "птица", "рыба"}, Correct: 1},
	},
	2: { // A2 level
		{Word: "journey", Translation: "путешествие", Options: []string{"работа", "путешествие", "покупки", "учеба"}, Correct: 1},
		{Word: "weather", Translation: "погода", Options: []string{"время", "погода", "деньги", "здоровье"}, Correct: 1},
	},
	3: { // B1 level
		{Word: "accomplish", Translation: "выполнять", Options: []string{"начинать", "выполнять", "отменять", "планировать"}, Correct: 1},
		{Word: "advantage", Translation: "преимущество", Options: []string{"недостаток", "преимущество", "результат", "проблема"}, Correct: 1},
	},
	4: { // B2 level
		{Word: "substantial", Translation: "значительный", Options: []string{"незначительный", "значительный", "временный", "постоянный"}, Correct: 1},
		{Word: "elaborate", Translation: "детализировать", Options: []string{"упрощать", "детализировать", "скрывать", "копировать"}, Correct: 1},
	},
	5: { // C1-C2 level
		{Word: "inevitable", Translation: "неизбежный", Options: []string{"возможный", "неизбежный", "невероятный", "желательный"}, Correct: 1},
		{Word: "coherent", Translation: "связный", Options: []string{"разрозненный", "связный", "короткий", "длинный"}, Correct: 1},
	},
}

type TestQuestion struct {
	Word        string   `json:"word"`
	Translation string   `json:"translation"`
	Options     []string `json:"options"`
	Correct     int      `json:"correct"`
}

// HandleTestStartCallback handles the start of CEFR test
func (s *HandlerService) HandleTestStartCallback(ctx context.Context, c tele.Context, userID int64, currentState fsm.UserState) error {
	// Validate current state
	if currentState != fsm.StateVocabularyTest {
		s.logger.Warn("Invalid state for test start",
			zap.Int64("user_id", userID),
			zap.String("expected_state", string(fsm.StateVocabularyTest)),
			zap.String("actual_state", string(currentState)))
		return c.Send("Пожалуйста, начните с команды /test")
	}

	// Initialize test data
	testData := &fsm.CEFRTestData{
		CurrentGroup:   1,
		Answers:        make(map[string]string),
		CorrectAnswers: 0,
		StartTime:      time.Now(),
	}

	// Store test data
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataCEFRTest, testData); err != nil {
		s.logger.Error("Failed to store test data", zap.Error(err))
		return err
	}

	// Transition to first test group
	if err := s.stateManager.SetState(ctx, userID, fsm.StateTestGroup1); err != nil {
		s.logger.Error("Failed to set test group 1 state", zap.Error(err))
		return err
	}

	// Send first question
	return s.sendTestQuestion(ctx, c, userID, 1, 0)
}

// sendTestQuestion sends a test question for a specific group and question index
func (s *HandlerService) sendTestQuestion(ctx context.Context, c tele.Context, userID int64, group int, questionIndex int) error {
	questions := cefrTestQuestions[group]
	if questionIndex >= len(questions) {
		// Move to next group or finish test
		return s.handleGroupComplete(ctx, c, userID, group)
	}

	question := questions[questionIndex]

	// Create question text
	questionText := fmt.Sprintf(
		"📝 *Тест CEFR - Группа %d*\n\n"+
			"Вопрос %d из %d\n\n"+
			"Что означает слово: **%s**?",
		group,
		questionIndex+1,
		len(questions),
		question.Word,
	)

	// Create answer options
	var buttons [][]tele.InlineButton
	for i, option := range question.Options {
		buttonData := fmt.Sprintf("test_answer:%d:%d:%d", group, questionIndex, i)
		button := tele.InlineButton{
			Text: option,
			Data: buttonData,
		}
		buttons = append(buttons, []tele.InlineButton{button})
	}

	// Add "Don't know" button
	dontKnowButton := tele.InlineButton{
		Text: "🤷‍♂️ Не знаю",
		Data: fmt.Sprintf("test_dont_know:%d:%d", group, questionIndex),
	}
	buttons = append(buttons, []tele.InlineButton{dontKnowButton})

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: buttons,
	}

	return c.Send(questionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleTestAnswerCallback handles test answer callbacks
func (s *HandlerService) HandleTestAnswerCallback(ctx context.Context, c tele.Context, userID int64, group, questionIndex, answerIndex int) error {
	// Get test data
	testData, err := s.stateManager.GetCEFRTestData(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get test data", zap.Error(err))
		return err
	}

	// Check if answer is correct
	questions := cefrTestQuestions[group]
	if questionIndex < len(questions) {
		question := questions[questionIndex]
		answerKey := fmt.Sprintf("g%d_q%d", group, questionIndex)

		isCorrect := answerIndex == question.Correct
		if isCorrect {
			testData.CorrectAnswers++
			testData.Answers[answerKey] = "correct"
		} else {
			testData.Answers[answerKey] = "incorrect"
		}

		// Update test data
		if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataCEFRTest, testData); err != nil {
			s.logger.Error("Failed to update test data", zap.Error(err))
			return err
		}

		// Show answer feedback
		return s.showAnswerFeedback(ctx, c, userID, group, questionIndex, isCorrect, question)
	}

	// If question index is invalid, continue to next question
	return s.sendTestQuestion(ctx, c, userID, group, questionIndex+1)
}

// showAnswerFeedback shows whether the answer was correct or incorrect
func (s *HandlerService) showAnswerFeedback(ctx context.Context, c tele.Context, userID int64, group, questionIndex int, isCorrect bool, question TestQuestion) error {
	var feedbackText string
	var emoji string

	if isCorrect {
		emoji = "✅"
		feedbackText = fmt.Sprintf(
			"%s *Правильно!*\n\n"+
				"**%s** = %s",
			emoji,
			question.Word,
			question.Translation,
		)
	} else {
		emoji = "❌"
		correctAnswer := question.Options[question.Correct]
		feedbackText = fmt.Sprintf(
			"%s *Неправильно*\n\n"+
				"**%s** = %s\n\n"+
				"Правильный ответ: **%s**",
			emoji,
			question.Word,
			question.Translation,
			correctAnswer,
		)
	}

	// Create continue button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Продолжить", Data: fmt.Sprintf("test_continue_next:%d:%d", group, questionIndex+1)}},
		},
	}

	return c.Send(feedbackText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleTestContinueNextCallback handles continuing to next question after feedback
func (s *HandlerService) HandleTestContinueNextCallback(ctx context.Context, c tele.Context, userID int64, group, nextQuestionIndex int) error {
	// Send next question
	return s.sendTestQuestion(ctx, c, userID, group, nextQuestionIndex)
}

// HandleTestDontKnowCallback handles "don't know" responses
func (s *HandlerService) HandleTestDontKnowCallback(ctx context.Context, c tele.Context, userID int64, group, questionIndex int) error {
	// Get test data
	testData, err := s.stateManager.GetCEFRTestData(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get test data", zap.Error(err))
		return err
	}

	// Mark answer as "don't know" (counts as incorrect)
	answerKey := fmt.Sprintf("g%d_q%d", group, questionIndex)
	testData.Answers[answerKey] = "dont_know"

	// Update test data
	if err := s.stateManager.StoreTempData(ctx, userID, fsm.TempDataCEFRTest, testData); err != nil {
		s.logger.Error("Failed to update test data", zap.Error(err))
		return err
	}

	// Show "don't know" feedback
	questions := cefrTestQuestions[group]
	if questionIndex < len(questions) {
		question := questions[questionIndex]
		return s.showDontKnowFeedback(ctx, c, userID, group, questionIndex, question)
	}

	// Count consecutive "don't know" answers in current group
	consecutiveDontKnow := 0

	for i := 0; i <= questionIndex; i++ {
		key := fmt.Sprintf("g%d_q%d", group, i)
		if answer, exists := testData.Answers[key]; exists && answer == "dont_know" {
			consecutiveDontKnow++
		} else {
			consecutiveDontKnow = 0 // Reset counter if there's a non-"don't know" answer
		}
	}

	// If user clicked "don't know" 2 times in a row in current group, offer to stop test
	if consecutiveDontKnow >= 2 {
		return s.offerToStopTest(ctx, c, userID, group)
	}

	// Otherwise, continue to next question
	return s.sendTestQuestion(ctx, c, userID, group, questionIndex+1)
}

// showDontKnowFeedback shows feedback for "don't know" answers
func (s *HandlerService) showDontKnowFeedback(ctx context.Context, c tele.Context, userID int64, group, questionIndex int, question TestQuestion) error {
	feedbackText := fmt.Sprintf(
		"🤷‍♂️ *Не знаешь? Ничего страшного!*\n\n"+
			"**%s** = %s\n\n"+
			"Теперь ты знаешь это слово! 📚",
		question.Word,
		question.Translation,
	)

	// Check if we need to offer stopping the test
	testData, err := s.stateManager.GetCEFRTestData(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get test data", zap.Error(err))
		return err
	}

	// Count consecutive "don't know" answers
	consecutiveDontKnow := 0
	for i := 0; i <= questionIndex; i++ {
		key := fmt.Sprintf("g%d_q%d", group, i)
		if answer, exists := testData.Answers[key]; exists && answer == "dont_know" {
			consecutiveDontKnow++
		} else {
			consecutiveDontKnow = 0
		}
	}

	if consecutiveDontKnow >= 2 {
		// Offer to stop test
		if _, err := c.Bot().Send(c.Sender(), feedbackText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}); err != nil {
			s.logger.Error("Failed to send feedback message", zap.Error(err))
		}
		return s.offerToStopTest(ctx, c, userID, group)
	}

	// Create continue button
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Продолжить", Data: fmt.Sprintf("test_continue_next:%d:%d", group, questionIndex+1)}},
		},
	}

	return c.Send(feedbackText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// offerToStopTest offers to stop the test and fix level at previous group
func (s *HandlerService) offerToStopTest(ctx context.Context, c tele.Context, userID int64, currentGroup int) error {
	// Determine suggested level based on previous performance
	suggestedLevel := "A1" // default
	if currentGroup > 1 {
		cefrLevels := []string{"A1", "A2", "B1", "B2", "C1"}
		suggestedLevel = cefrLevels[currentGroup-2] // previous group level
	}

	stopText := fmt.Sprintf(
		"🤔 *Похоже, вопросы стали сложными*\n\n"+
			"Ты ответил \"Не знаю\" на несколько вопросов подряд.\n\n"+
			"Могу предложить зафиксировать твой уровень как **%s** "+
			"или продолжить тест до конца.\n\n"+
			"Что выберешь?",
		suggestedLevel,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: fmt.Sprintf("Зафиксировать %s", suggestedLevel), Data: fmt.Sprintf("test_fix_level:%s", suggestedLevel)}},
			{{Text: "Продолжить тест", Data: fmt.Sprintf("test_continue:%d", currentGroup)}},
		},
	}

	return c.Send(stopText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleTestFixLevelCallback handles fixing user level without completing full test
func (s *HandlerService) HandleTestFixLevelCallback(ctx context.Context, c tele.Context, userID int64, level string) error {
	// Set user to start state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateStart); err != nil {
		s.logger.Error("Failed to set start state", zap.Error(err))
		return err
	}

	// Clear test data
	if err := s.stateManager.ClearTempData(ctx, userID, fsm.TempDataCEFRTest); err != nil {
		s.logger.Error("Failed to clear test data", zap.Error(err))
	}

	// Send completion message
	completionText := fmt.Sprintf(
		"✅ *Уровень зафиксирован!*\n\n"+
			"🎯 **Твой уровень: %s**\n\n"+
			"Отлично! Я буду подбирать слова, соответствующие твоему уровню.\n\n"+
			"Готов начать изучение?",
		level,
	)

	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{{Text: "Начать изучение", Data: "lesson:start"}},
			{{Text: "Настройки", Data: "menu:settings"}},
		},
	}

	return c.Send(completionText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
}

// HandleTestContinueCallback handles continuing test after "don't know" warning
func (s *HandlerService) HandleTestContinueCallback(ctx context.Context, c tele.Context, userID int64, group int) error {
	// Get test data to find current question index
	testData, err := s.stateManager.GetCEFRTestData(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get test data", zap.Error(err))
		return err
	}

	// Find the next question index in current group
	questions := cefrTestQuestions[group]
	nextQuestionIndex := 0

	// Count how many questions in current group have been answered
	for i := 0; i < len(questions); i++ {
		answerKey := fmt.Sprintf("g%d_q%d", group, i)
		if _, exists := testData.Answers[answerKey]; exists {
			nextQuestionIndex = i + 1
		} else {
			break
		}
	}

	// Continue with next question
	return s.sendTestQuestion(ctx, c, userID, group, nextQuestionIndex)
}

// handleGroupComplete handles completion of a test group
func (s *HandlerService) handleGroupComplete(ctx context.Context, c tele.Context, userID int64, completedGroup int) error {
	if completedGroup >= 5 {
		// Test complete
		return s.completeTest(ctx, c, userID)
	}

	nextGroup := completedGroup + 1
	nextState := getTestGroupState(nextGroup)

	// Transition to next group
	if err := s.stateManager.SetState(ctx, userID, nextState); err != nil {
		s.logger.Error("Failed to set next test group state", zap.Error(err))
		return err
	}

	// Send progress message
	progressText := fmt.Sprintf(
		"✅ *Группа %d завершена!*\n\n"+
			"Переходим к группе %d из 5...",
		completedGroup,
		nextGroup,
	)

	if _, err := c.Bot().Send(c.Sender(), progressText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}); err != nil {
		s.logger.Error("Failed to send progress message", zap.Error(err))
	}

	// Send first question of next group
	return s.sendTestQuestion(ctx, c, userID, nextGroup, 0)
}

// completeTest handles test completion and determines CEFR level
func (s *HandlerService) completeTest(ctx context.Context, c tele.Context, userID int64) error {
	// Get test data
	testData, err := s.stateManager.GetCEFRTestData(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get test data", zap.Error(err))
		return err
	}

	testData.EndTime = time.Now()

	// Calculate results for each group
	groupScores := make(map[int]int)
	for group := 1; group <= 5; group++ {
		score := 0
		questions := cefrTestQuestions[group]
		for questionIndex := 0; questionIndex < len(questions); questionIndex++ {
			answerKey := fmt.Sprintf("g%d_q%d", group, questionIndex)
			if answer, exists := testData.Answers[answerKey]; exists && answer == "correct" {
				score++
			}
		}
		groupScores[group] = score
	}

	// Determine CEFR level based on scores
	cefrLevel := determineCEFRLevel(groupScores)

	// Transition to test result state
	if err := s.stateManager.SetState(ctx, userID, fsm.StateCEFRTestResult); err != nil {
		s.logger.Error("Failed to set test result state", zap.Error(err))
		return err
	}

	// Create result message
	resultText := fmt.Sprintf(
		"🎉 *Тест завершен!*\n\n"+
			"📊 **Твои результаты:**\n"+
			"🥉 A1 уровень: %d/2\n"+
			"🥈 A2 уровень: %d/2\n"+
			"🥇 B1 уровень: %d/2\n"+
			"🏆 B2 уровень: %d/2\n"+
			"👑 C1-C2 уровень: %d/2\n\n"+
			"🎯 **Твой уровень: %s**\n\n"+
			"Отлично! Теперь я знаю, какие слова тебе подходят.",
		groupScores[1], groupScores[2], groupScores[3], groupScores[4], groupScores[5],
		cefrLevel,
	)

	// Clear test data and questionnaire temp data
	if err := s.stateManager.ClearTempData(ctx, userID, fsm.TempDataCEFRTest); err != nil {
		s.logger.Error("Failed to clear test data", zap.Error(err))
	}

	// Clear all questionnaire temp data now that onboarding is complete
	tempDataTypes := []fsm.TempDataType{
		fsm.TempDataGoal,
		fsm.TempDataExperience,
		fsm.TempDataConfidence,
		fsm.TempDataWordsPerDay,
		fsm.TempDataNotifications,
		fsm.TempDataNotificationTime,
	}

	for _, dataType := range tempDataTypes {
		if err := s.stateManager.ClearTempData(ctx, userID, dataType); err != nil {
			s.logger.Error("Failed to clear questionnaire temp data",
				zap.String("data_type", string(dataType)), zap.Error(err))
		}
	}

	// Check if user is authenticated
	isAuthenticated := s.IsUserAuthenticated(ctx, userID)

	if !isAuthenticated {
		// New user - need to authenticate to save progress
		resultText += "\n\n🔐 **Для сохранения прогресса нужно создать аккаунт**\n\n" +
			"Регистрация через Google займет всего 30 секунд и позволит сохранить твой прогресс на всех устройствах."

		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{{Text: "🔐 Создать аккаунт", Data: "auth:register"}},
				{{Text: "🔑 У меня уже есть аккаунт", Data: "auth:existing_user"}},
				{{Text: "⏭ Попробовать без аккаунта", Data: "test:skip"}},
			},
		}

		return c.Send(resultText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
	} else {
		// Authenticated user - update complete preferences with test results
		token, err := s.stateManager.GetJWTToken(ctx, userID)
		if err == nil {
			// Build complete preferences from questionnaire answers
			preferences, err := s.buildCompletePreferencesFromQuestionnaire(ctx, userID, cefrLevel)
			if err != nil {
				s.logger.Error("Failed to build complete preferences", zap.Error(err))
				// Fallback to just CEFR level
				preferences = &api.UpdatePreferenceRequest{
					CEFRLevel: cefrLevel,
				}
			}

			if _, err := s.apiClient.UpdateUserPreferences(ctx, token, preferences); err != nil {
				s.logger.Error("Failed to update user preferences", zap.Error(err))
			} else {
				s.logger.Info("Successfully updated complete user preferences from test",
					zap.Int64("user_id", userID),
					zap.String("cefr_level", cefrLevel),
					zap.Int("words_per_day", preferences.WordsPerDay),
					zap.Bool("notifications", preferences.Notifications),
					zap.String("goal", preferences.Goal))
			}
		}

		// Create completion keyboard
		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{{Text: "Начать изучение", Data: "lesson:start"}},
				{{Text: "Настройки", Data: "menu:settings"}},
			},
		}

		// Add preferences summary to result text for authenticated users
		wordsPerDayData, _ := s.stateManager.GetTempData(ctx, userID, fsm.TempDataWordsPerDay)
		wordsPerDay, _ := wordsPerDayData.(int)
		if wordsPerDay == 0 {
			wordsPerDay = 10 // Default
		}

		notificationsData, _ := s.stateManager.GetTempData(ctx, userID, fsm.TempDataNotifications)
		notifications, _ := notificationsData.(bool)

		notificationStatus := "отключены"
		if notifications {
			notificationStatus = "включены"
		}

		resultText += fmt.Sprintf(
			"\n\n📊 **Ваши настройки:**\n"+
				"• Слов в день: *%d*\n"+
				"• Уведомления: *%s*\n\n"+
				"Настройка завершена! Теперь можете начать изучение.",
			wordsPerDay,
			notificationStatus,
		)

		// Set user back to start state
		if err := s.stateManager.SetState(ctx, userID, fsm.StateStart); err != nil {
			s.logger.Error("Failed to set start state", zap.Error(err))
		}

		return c.Send(resultText, &tele.SendOptions{ParseMode: tele.ModeMarkdown}, keyboard)
	}
}

// determineCEFRLevel determines CEFR level based on group scores
func determineCEFRLevel(groupScores map[int]int) string {
	// Simple algorithm: find the highest level where user got at least 1/2 correct
	cefrLevels := []string{"A1", "A2", "B1", "B2", "C1"}

	maxLevel := "A1" // default minimum level

	for group := 1; group <= 5; group++ {
		if score := groupScores[group]; score >= 1 {
			maxLevel = cefrLevels[group-1]
		}
	}

	return maxLevel
}

// getTestGroupState returns the FSM state for a test group
func getTestGroupState(group int) fsm.UserState {
	switch group {
	case 1:
		return fsm.StateTestGroup1
	case 2:
		return fsm.StateTestGroup2
	case 3:
		return fsm.StateTestGroup3
	case 4:
		return fsm.StateTestGroup4
	case 5:
		return fsm.StateTestGroup5
	default:
		return fsm.StateVocabularyTest
	}
}
