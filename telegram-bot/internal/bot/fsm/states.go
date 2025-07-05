package fsm

import "slices"

// UserState represents the current state of the user in the learning flow
type UserState string

const (
	// Stage 1: Onboarding
	StateStart             UserState = "start"
	StateWelcome           UserState = "welcome"
	StateMethodExplanation UserState = "method_explanation"
	StateSpacedRepetition  UserState = "spaced_repetition"

	// Stage 2: Personalization
	StateQuestionnaire      UserState = "questionnaire"
	StateQuestionGoal       UserState = "question_goal"
	StateQuestionConfidence UserState = "question_confidence"
	StateQuestionSerials    UserState = "question_serials"
	StateQuestionExperience UserState = "question_experience"

	// CEFR Test Flow
	StateVocabularyTest     UserState = "vocabulary_test"
	StateTestGroup1         UserState = "test_group_1"
	StateTestGroup2         UserState = "test_group_2"
	StateTestGroup3         UserState = "test_group_3"
	StateTestGroup4         UserState = "test_group_4"
	StateTestGroup5         UserState = "test_group_5"
	StateCEFRTestProcessing UserState = "cefr_test_processing"
	StateCEFRTestResult     UserState = "cefr_test_result"
	StateLevelDetermination UserState = "level_determination"
	StatePlanCreation       UserState = "plan_creation"

	// Stage 3: Learning Flow
	// Main Learning States
	StateLessonStart      UserState = "lesson_start"
	StateLessonInProgress UserState = "lesson_in_progress"
	StateLessonComplete   UserState = "lesson_complete"

	// Detailed Learning Process
	StateShowingWords          UserState = "showing_words"
	StateShowingFirstBlock     UserState = "showing_first_block"
	StateShowingSecondBlock    UserState = "showing_second_block"
	StateShowingThirdBlock     UserState = "showing_third_block"
	StateExerciseAfterBlock    UserState = "exercise_after_block"
	StateShowingIndividualWord UserState = "showing_individual_word"
	StateExerciseReview        UserState = "exercise_review"

	// Exercise States
	StateAudioDictation        UserState = "audio_dictation"
	StateTranslationCheck      UserState = "translation_check"
	StateMultipleChoiceCheck   UserState = "multiple_choice_check"
	StateWaitingForAudio       UserState = "waiting_for_audio"
	StateWaitingForTranslation UserState = "waiting_for_translation"
	StateWaitingForChoice      UserState = "waiting_for_choice"

	// Lesson Management
	StateDailyProgress    UserState = "daily_progress"
	StateScheduleReminder UserState = "schedule_reminder"

	// Settings Flow
	StateSettings                 UserState = "settings"
	StateSettingsWordsPerDay      UserState = "settings_words_per_day"
	StateSettingsWordsPerDayInput UserState = "settings_words_per_day_input"
	StateSettingsNotifications    UserState = "settings_notifications"
	StateSettingsTimeSelection    UserState = "settings_time_selection"
	StateSettingsTimeInput        UserState = "settings_time_input"
	StateSettingsCEFRLevel        UserState = "settings_cefr_level"
	StateSettingsLanguage         UserState = "settings_language"

	// Account Management
	StateAccountLinking UserState = "account_linking"
	StateWaitingForLink UserState = "waiting_for_link"
	StateAccountLinked  UserState = "account_linked"

	// Error and Recovery States
	StateError         UserState = "error"
	StateErrorRecovery UserState = "error_recovery"
	StateRetry         UserState = "retry"
	StateHelp          UserState = "help"
	StateCancel        UserState = "cancel"
)

// StateTransition represents a valid state transition
type StateTransition struct {
	From UserState
	To   UserState
}

// ValidTransitions defines the allowed state transitions
var ValidTransitions = map[StateTransition]bool{
	// Onboarding flow
	{StateStart, StateWelcome}:                      true,
	{StateWelcome, StateMethodExplanation}:          true,
	{StateMethodExplanation, StateSpacedRepetition}: true,
	{StateSpacedRepetition, StateQuestionnaire}:     true,

	// Questionnaire flow
	{StateQuestionnaire, StateQuestionGoal}:         true,
	{StateQuestionGoal, StateQuestionConfidence}:    true,
	{StateQuestionConfidence, StateQuestionSerials}: true,
	{StateQuestionSerials, StateQuestionExperience}: true,
	{StateQuestionExperience, StateVocabularyTest}:  true,

	// Vocabulary test flow
	{StateVocabularyTest, StateTestGroup1}:         true,
	{StateTestGroup1, StateTestGroup2}:             true,
	{StateTestGroup2, StateTestGroup3}:             true,
	{StateTestGroup3, StateTestGroup4}:             true,
	{StateTestGroup4, StateTestGroup5}:             true,
	{StateTestGroup5, StateCEFRTestProcessing}:     true,
	{StateCEFRTestProcessing, StateCEFRTestResult}: true,
	{StateCEFRTestResult, StateLevelDetermination}: true,
	{StateLevelDetermination, StatePlanCreation}:   true,
	{StatePlanCreation, StateLessonStart}:          true,

	// Lesson flow - main states
	{StateLessonStart, StateLessonInProgress}:    true,
	{StateLessonInProgress, StateLessonComplete}: true,
	{StateLessonComplete, StateLessonStart}:      true,

	// Detailed learning flow
	{StateLessonInProgress, StateShowingFirstBlock}:       true,
	{StateShowingFirstBlock, StateExerciseAfterBlock}:     true,
	{StateExerciseAfterBlock, StateShowingSecondBlock}:    true,
	{StateShowingSecondBlock, StateExerciseAfterBlock}:    true,
	{StateExerciseAfterBlock, StateShowingThirdBlock}:     true,
	{StateShowingThirdBlock, StateExerciseAfterBlock}:     true,
	{StateExerciseAfterBlock, StateShowingIndividualWord}: true,
	{StateShowingIndividualWord, StateExerciseReview}:     true,
	{StateExerciseReview, StateShowingIndividualWord}:     true,
	{StateExerciseReview, StateLessonComplete}:            true,

	// Exercise flows
	{StateExerciseAfterBlock, StateAudioDictation}:      true,
	{StateExerciseAfterBlock, StateTranslationCheck}:    true,
	{StateExerciseAfterBlock, StateMultipleChoiceCheck}: true,
	{StateExerciseReview, StateAudioDictation}:          true,
	{StateExerciseReview, StateTranslationCheck}:        true,
	{StateExerciseReview, StateMultipleChoiceCheck}:     true,
	{StateAudioDictation, StateWaitingForAudio}:         true,
	{StateTranslationCheck, StateWaitingForTranslation}: true,
	{StateMultipleChoiceCheck, StateWaitingForChoice}:   true,
	{StateWaitingForAudio, StateExerciseReview}:         true,
	{StateWaitingForTranslation, StateExerciseReview}:   true,
	{StateWaitingForChoice, StateExerciseReview}:        true,

	// Account management
	{StateStart, StateAccountLinking}:          true,
	{StateAccountLinking, StateWaitingForLink}: true,
	{StateWaitingForLink, StateAccountLinked}:  true,
	{StateAccountLinked, StateQuestionnaire}:   true,

	// Settings flow - main menu transitions
	{StateStart, StateSettings}:                 true,
	{StateLessonComplete, StateSettings}:        true,
	{StateSettings, StateSettingsWordsPerDay}:   true,
	{StateSettings, StateSettingsNotifications}: true,
	{StateSettings, StateSettingsCEFRLevel}:     true,
	{StateSettings, StateSettingsLanguage}:      true,

	// Settings - Words Per Day flow
	{StateSettingsWordsPerDay, StateSettingsWordsPerDayInput}: true,
	{StateSettingsWordsPerDayInput, StateSettings}:            true,
	{StateSettingsWordsPerDay, StateSettings}:                 true,

	// Settings - Notifications flow
	{StateSettingsNotifications, StateSettingsTimeSelection}: true,
	{StateSettingsTimeSelection, StateSettingsTimeInput}:     true,
	{StateSettingsTimeInput, StateSettingsNotifications}:     true,
	{StateSettingsTimeSelection, StateSettingsNotifications}: true,
	{StateSettingsNotifications, StateSettings}:              true,

	// Settings - CEFR Level flow
	{StateSettingsCEFRLevel, StateVocabularyTest}: true,
	{StateCEFRTestResult, StateSettingsCEFRLevel}: true,
	{StateSettingsCEFRLevel, StateSettings}:       true,

	// Settings - Language flow
	{StateSettingsLanguage, StateSettings}: true,

	// Common transitions to/from lesson flow
	{StateSettings, StateLessonStart}: true,
	{StateLessonComplete, StateStart}: true,

	// Error handling - transitions from any state to error/help and recovery
	{StateError, StateErrorRecovery}: true,
	{StateErrorRecovery, StateRetry}: true,
	{StateRetry, StateStart}:         true,
	{StateHelp, StateStart}:          true,
	{StateCancel, StateStart}:        true,
}

// IsValidTransition checks if a state transition is valid
func IsValidTransition(from, to UserState) bool {
	// Special case: allow transitions to error state from any state
	if to == StateError || to == StateHelp || to == StateCancel {
		return true
	}

	// Special case: allow transitions from error recovery to most main states
	if from == StateErrorRecovery && isMainState(to) {
		return true
	}

	return ValidTransitions[StateTransition{From: from, To: to}]
}

// isMainState checks if a state is one of the main entry points
func isMainState(state UserState) bool {
	mainStates := []UserState{
		StateStart,
		StateSettings,
		StateLessonStart,
		StateVocabularyTest,
	}

	for _, ms := range mainStates {
		if state == ms {
			return true
		}
	}
	return false
}

// GetInitialState returns the initial state for new users
func GetInitialState() UserState {
	return StateStart
}

// IsLessonState checks if the state is part of the lesson flow
func IsLessonState(state UserState) bool {
	lessonStates := []UserState{
		StateLessonStart,
		StateLessonInProgress,
		StateShowingWords,
		StateShowingFirstBlock,
		StateShowingSecondBlock,
		StateShowingThirdBlock,
		StateExerciseAfterBlock,
		StateShowingIndividualWord,
		StateExerciseReview,
		StateAudioDictation,
		StateTranslationCheck,
		StateMultipleChoiceCheck,
		StateWaitingForAudio,
		StateWaitingForTranslation,
		StateWaitingForChoice,
		StateLessonComplete,
	}

	return slices.Contains(lessonStates, state)
}

// IsExerciseState checks if the state is an exercise state
func IsExerciseState(state UserState) bool {
	exerciseStates := []UserState{
		StateAudioDictation,
		StateTranslationCheck,
		StateMultipleChoiceCheck,
		StateWaitingForAudio,
		StateWaitingForTranslation,
		StateWaitingForChoice,
		StateExerciseReview,
	}

	return slices.Contains(exerciseStates, state)
}

// IsSettingsState checks if the state is part of the settings flow
func IsSettingsState(state UserState) bool {
	settingsStates := []UserState{
		StateSettings,
		StateSettingsWordsPerDay,
		StateSettingsWordsPerDayInput,
		StateSettingsNotifications,
		StateSettingsTimeSelection,
		StateSettingsTimeInput,
		StateSettingsCEFRLevel,
		StateSettingsLanguage,
	}

	return slices.Contains(settingsStates, state)
}

// IsCEFRTestState checks if the state is part of the CEFR test flow
func IsCEFRTestState(state UserState) bool {
	testStates := []UserState{
		StateVocabularyTest,
		StateTestGroup1,
		StateTestGroup2,
		StateTestGroup3,
		StateTestGroup4,
		StateTestGroup5,
		StateCEFRTestProcessing,
		StateCEFRTestResult,
		StateLevelDetermination,
	}

	return slices.Contains(testStates, state)
}
