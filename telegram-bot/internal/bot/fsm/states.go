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

	// Stage 3: Learning Flow (New Implementation)
	// Main Learning States
	StateLessonStart      UserState = "lesson_start"
	StateLessonInProgress UserState = "lesson_in_progress"
	StateLessonComplete   UserState = "lesson_complete"

	// New Learning Flow: Show 3 words → Test → Repeat
	StateShowingWordSet     UserState = "showing_word_set"     // Showing current set of 3 words
	StateShowingWord1       UserState = "showing_word_1"       // Showing word 1 of 3
	StateShowingWord2       UserState = "showing_word_2"       // Showing word 2 of 3
	StateShowingWord3       UserState = "showing_word_3"       // Showing word 3 of 3
	StateReadyForExercises  UserState = "ready_for_exercises"  // Ready to start exercises for current set
	StateDoingExercises     UserState = "doing_exercises"      // Doing exercises for current set
	StateExerciseInProgress UserState = "exercise_in_progress" // Single exercise in progress
	StateSetComplete        UserState = "set_complete"         // Current set of 3 words complete
	StatePreparingNextSet   UserState = "preparing_next_set"   // Preparing next set of words

	// Exercise Types (new implementation)
	StatePickOptionSentence   UserState = "pick_option_sentence"   // Multiple choice with sentence template
	StateWriteWordTranslation UserState = "write_word_translation" // Write word from translation
	StateTranslateRuToEn      UserState = "translate_ru_to_en"     // Translate Russian to English
	StateWaitingForTextInput  UserState = "waiting_for_text_input" // Waiting for user text input

	// Legacy Learning Process (keeping for backward compatibility)
	StateShowingWords          UserState = "showing_words"
	StateShowingFirstBlock     UserState = "showing_first_block"
	StateShowingSecondBlock    UserState = "showing_second_block"
	StateShowingThirdBlock     UserState = "showing_third_block"
	StateExerciseAfterBlock    UserState = "exercise_after_block"
	StateShowingIndividualWord UserState = "showing_individual_word"
	StateExerciseReview        UserState = "exercise_review"

	// Legacy Exercise States
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
	{StateQuestionnaire, StateQuestionGoal}:            true,
	{StateQuestionGoal, StateQuestionConfidence}:       true,
	{StateQuestionConfidence, StateQuestionExperience}: true,
	{StateQuestionExperience, StateVocabularyTest}:     true,

	// Vocabulary test flow
	{StateVocabularyTest, StateTestGroup1}:         true,
	{StateTestGroup1, StateTestGroup2}:             true,
	{StateTestGroup2, StateTestGroup3}:             true,
	{StateTestGroup3, StateTestGroup4}:             true,
	{StateTestGroup4, StateTestGroup5}:             true,
	{StateTestGroup5, StateCEFRTestProcessing}:     true,
	{StateTestGroup5, StateCEFRTestResult}:         true, // Direct transition for completion
	{StateCEFRTestProcessing, StateCEFRTestResult}: true,
	{StateCEFRTestResult, StateLevelDetermination}: true,
	{StateLevelDetermination, StatePlanCreation}:   true,
	{StatePlanCreation, StateLessonStart}:          true,

	// Test skip flow
	{StateVocabularyTest, StateCEFRTestResult}: true, // Allow skipping test
	{StateTestGroup1, StateCEFRTestResult}:     true, // Allow skipping from any test stage
	{StateTestGroup2, StateCEFRTestResult}:     true,
	{StateTestGroup3, StateCEFRTestResult}:     true,
	{StateTestGroup4, StateCEFRTestResult}:     true,

	// Lesson flow - main states
	{StateLessonStart, StateLessonInProgress}:    true,
	{StateLessonInProgress, StateLessonComplete}: true,
	{StateLessonComplete, StateLessonStart}:      true,
	{StateStart, StateLessonStart}:               true, // Allow starting lesson from main menu
	{StateStart, StateLessonInProgress}:          true, // Allow continuing lesson from main menu

	// New Learning Flow: Show 3 words → Test → Repeat
	{StateLessonInProgress, StateShowingWordSet}:   true,
	{StateShowingWordSet, StateShowingWord1}:       true,
	{StateShowingWord1, StateShowingWord2}:         true,
	{StateShowingWord2, StateShowingWord3}:         true,
	{StateShowingWord3, StateReadyForExercises}:    true,
	{StateReadyForExercises, StateDoingExercises}:  true,
	{StateDoingExercises, StateExerciseInProgress}: true,
	{StateExerciseInProgress, StateDoingExercises}: true, // Next exercise
	{StateDoingExercises, StateSetComplete}:        true, // All exercises done
	{StateSetComplete, StatePreparingNextSet}:      true,
	{StatePreparingNextSet, StateShowingWordSet}:   true, // Next set of words
	{StateSetComplete, StateLessonComplete}:        true, // 10 words learned

	// Exercise flow for new implementation
	{StateExerciseInProgress, StatePickOptionSentence}:    true,
	{StateExerciseInProgress, StateWriteWordTranslation}:  true,
	{StateExerciseInProgress, StateTranslateRuToEn}:       true,
	{StateWriteWordTranslation, StateWaitingForTextInput}: true,
	{StateTranslateRuToEn, StateWaitingForTextInput}:      true,
	{StatePickOptionSentence, StateDoingExercises}:        true, // Back to exercise queue
	{StateWaitingForTextInput, StateDoingExercises}:       true, // Back to exercise queue

	// Legacy detailed learning flow (keeping for backward compatibility)
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

	// Additional authentication flows
	{StateWelcome, StateAccountLinking}:        true,
	{StateAccountLinking, StateStart}:          true,
	{StateAccountLinked, StateStart}:           true,
	{StateQuestionnaire, StateAccountLinking}:  true,
	{StateCEFRTestResult, StateAccountLinking}: true,
	{StateAccountLinking, StateQuestionnaire}:  true,

	// Welcome to lesson transitions (for authenticated users)
	{StateWelcome, StateLessonStart}:      true,
	{StateWelcome, StateLessonInProgress}: true,
	{StateWelcome, StateStart}:            true, // Allow transition from welcome to start
	{StateWelcome, StateQuestionnaire}:    true, // Allow transition from welcome to questionnaire for fast-track onboarding

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

	// Special case: allow transitions to lesson states from any state
	if to == StateLessonStart || to == StateLessonInProgress {
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
		StateLessonInProgress,
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
		StateLessonComplete,
		// New learning flow states
		StateShowingWordSet,
		StateShowingWord1,
		StateShowingWord2,
		StateShowingWord3,
		StateReadyForExercises,
		StateDoingExercises,
		StateExerciseInProgress,
		StateSetComplete,
		StatePreparingNextSet,
		StatePickOptionSentence,
		StateWriteWordTranslation,
		StateTranslateRuToEn,
		StateWaitingForTextInput,
		// Legacy states
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
	}

	return slices.Contains(lessonStates, state)
}

// IsExerciseState checks if the state is an exercise state
func IsExerciseState(state UserState) bool {
	exerciseStates := []UserState{
		// New exercise states
		StateDoingExercises,
		StateExerciseInProgress,
		StatePickOptionSentence,
		StateWriteWordTranslation,
		StateTranslateRuToEn,
		StateWaitingForTextInput,
		// Legacy exercise states
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
