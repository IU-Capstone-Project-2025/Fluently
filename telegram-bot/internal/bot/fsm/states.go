// states.go
package fsm

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
	StateVocabularyTest     UserState = "vocabulary_test"
	StateTestGroup1         UserState = "test_group_1"
	StateTestGroup2         UserState = "test_group_2"
	StateTestGroup3         UserState = "test_group_3"
	StateTestGroup4         UserState = "test_group_4"
	StateTestGroup5         UserState = "test_group_5"
	StateLevelDetermination UserState = "level_determination"
	StatePlanCreation       UserState = "plan_creation"

	// Stage 3: Lesson Flow
	StateLessonStart           UserState = "lesson_start"
	StateShowingWords          UserState = "showing_words"
	StateShowingFirstBlock     UserState = "showing_first_block"
	StateExerciseAfterBlock    UserState = "exercise_after_block"
	StateShowingIndividualWord UserState = "showing_individual_word"
	StateExerciseReview        UserState = "exercise_review"

	// Exercise States
	StateAudioDictation        UserState = "audio_dictation"
	StateTranslationCheck      UserState = "translation_check"
	StateWaitingForAudio       UserState = "waiting_for_audio"
	StateWaitingForTranslation UserState = "waiting_for_translation"

	// Lesson Management
	StateLessonComplete   UserState = "lesson_complete"
	StateDailyProgress    UserState = "daily_progress"
	StateScheduleReminder UserState = "schedule_reminder"

	// Account Management
	StateAccountLinking UserState = "account_linking"
	StateWaitingForLink UserState = "waiting_for_link"
	StateAccountLinked  UserState = "account_linked"

	// Settings and Preferences
	StateSettings             UserState = "settings"
	StateChangeLevel          UserState = "change_level"
	StateChangeGoal           UserState = "change_goal"
	StateNotificationSettings UserState = "notification_settings"

	// Error and Recovery States
	StateError UserState = "error"
	StateRetry UserState = "retry"
	StateHelp  UserState = "help"
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
	{StateVocabularyTest, StateTestGroup1}:       true,
	{StateTestGroup1, StateTestGroup2}:           true,
	{StateTestGroup2, StateTestGroup3}:           true,
	{StateTestGroup3, StateTestGroup4}:           true,
	{StateTestGroup4, StateTestGroup5}:           true,
	{StateTestGroup5, StateLevelDetermination}:   true,
	{StateLevelDetermination, StatePlanCreation}: true,

	// Lesson flow
	{StatePlanCreation, StateLessonStart}:                 true,
	{StateLessonStart, StateShowingFirstBlock}:            true,
	{StateShowingFirstBlock, StateExerciseAfterBlock}:     true,
	{StateExerciseAfterBlock, StateShowingIndividualWord}: true,
	{StateShowingIndividualWord, StateExerciseReview}:     true,
	{StateExerciseReview, StateShowingIndividualWord}:     true,
	{StateExerciseReview, StateLessonComplete}:            true,
	{StateLessonComplete, StateLessonStart}:               true,

	// Exercise flows
	{StateExerciseAfterBlock, StateAudioDictation}:      true,
	{StateExerciseAfterBlock, StateTranslationCheck}:    true,
	{StateExerciseReview, StateAudioDictation}:          true,
	{StateExerciseReview, StateTranslationCheck}:        true,
	{StateAudioDictation, StateWaitingForAudio}:         true,
	{StateTranslationCheck, StateWaitingForTranslation}: true,
	{StateWaitingForAudio, StateExerciseReview}:         true,
	{StateWaitingForTranslation, StateExerciseReview}:   true,

	// Account management
	{StateStart, StateAccountLinking}:          true,
	{StateAccountLinking, StateWaitingForLink}: true,
	{StateWaitingForLink, StateAccountLinked}:  true,
	{StateAccountLinked, StateQuestionnaire}:   true,

	// Settings
	{StateLessonComplete, StateSettings}:       true,
	{StateSettings, StateChangeLevel}:          true,
	{StateSettings, StateChangeGoal}:           true,
	{StateSettings, StateNotificationSettings}: true,

	// Error handling - can transition from any state to error/help
	{StateError, StateRetry}: true,
	{StateHelp, StateStart}:  true,
}

// IsValidTransition checks if a state transition is valid
func IsValidTransition(from, to UserState) bool {
	return ValidTransitions[StateTransition{From: from, To: to}]
}

// GetInitialState returns the initial state for new users
func GetInitialState() UserState {
	return StateStart
}

// IsLessonState checks if the state is part of the lesson flow
func IsLessonState(state UserState) bool {
	lessonStates := []UserState{
		StateLessonStart,
		StateShowingWords,
		StateShowingFirstBlock,
		StateExerciseAfterBlock,
		StateShowingIndividualWord,
		StateExerciseReview,
		StateAudioDictation,
		StateTranslationCheck,
		StateWaitingForAudio,
		StateWaitingForTranslation,
		StateLessonComplete,
	}

	for _, ls := range lessonStates {
		if state == ls {
			return true
		}
	}
	return false
}

// IsExerciseState checks if the state is an exercise state
func IsExerciseState(state UserState) bool {
	exerciseStates := []UserState{
		StateAudioDictation,
		StateTranslationCheck,
		StateWaitingForAudio,
		StateWaitingForTranslation,
		StateExerciseReview,
	}

	for _, es := range exerciseStates {
		if state == es {
			return true
		}
	}
	return false
}
