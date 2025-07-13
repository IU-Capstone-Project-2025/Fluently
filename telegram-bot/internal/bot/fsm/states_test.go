package fsm

import (
	"testing"
)

func TestTransitionToLessonFromAnyState(t *testing.T) {
	// Test various source states transitioning to lesson states
	testCases := []struct {
		fromState UserState
		toState   UserState
		expected  bool
	}{
		// From any state to StateLessonStart should be allowed
		{StateStart, StateLessonStart, true},
		{StateWelcome, StateLessonStart, true},
		{StateSettings, StateLessonStart, true},
		{StateQuestionGoal, StateLessonStart, true},
		{StateVocabularyTest, StateLessonStart, true},
		{StateTestGroup1, StateLessonStart, true},
		{StateExerciseReview, StateLessonStart, true},

		// From any state to StateLessonInProgress should be allowed
		{StateStart, StateLessonInProgress, true},
		{StateWelcome, StateLessonInProgress, true},
		{StateSettings, StateLessonInProgress, true},
		{StateQuestionGoal, StateLessonInProgress, true},
		{StateVocabularyTest, StateLessonInProgress, true},
		{StateTestGroup1, StateLessonInProgress, true},
		{StateExerciseReview, StateLessonInProgress, true},

		// Test that other transitions still work as expected
		{StateStart, StateWelcome, true},
		{StateWelcome, StateMethodExplanation, true},

		// Test that invalid transitions are still blocked
		{StateWelcome, StateTestGroup1, false},
		{StateQuestionGoal, StateShowingWords, false},
	}

	for _, tc := range testCases {
		result := IsValidTransition(tc.fromState, tc.toState)
		if result != tc.expected {
			t.Errorf("IsValidTransition(%s, %s) = %v; expected %v", tc.fromState, tc.toState, result, tc.expected)
		}
	}
}

func TestErrorStatesStillWork(t *testing.T) {
	// Test that error state transitions still work from any state
	testStates := []UserState{
		StateStart,
		StateWelcome,
		StateQuestionGoal,
		StateLessonInProgress,
		StateShowingWords,
		StateExerciseReview,
	}

	for _, fromState := range testStates {
		// Test transitions to error states
		if !IsValidTransition(fromState, StateError) {
			t.Errorf("Expected transition from %s to StateError to be valid", fromState)
		}
		if !IsValidTransition(fromState, StateHelp) {
			t.Errorf("Expected transition from %s to StateHelp to be valid", fromState)
		}
		if !IsValidTransition(fromState, StateCancel) {
			t.Errorf("Expected transition from %s to StateCancel to be valid", fromState)
		}
	}
}

func TestMainStatesIncludeLessonStates(t *testing.T) {
	// Test that lesson states are considered main states
	if !isMainState(StateLessonStart) {
		t.Error("StateLessonStart should be considered a main state")
	}
	if !isMainState(StateLessonInProgress) {
		t.Error("StateLessonInProgress should be considered a main state")
	}

	// Test error recovery to lesson states
	if !IsValidTransition(StateErrorRecovery, StateLessonStart) {
		t.Error("Expected transition from StateErrorRecovery to StateLessonStart to be valid")
	}
	if !IsValidTransition(StateErrorRecovery, StateLessonInProgress) {
		t.Error("Expected transition from StateErrorRecovery to StateLessonInProgress to be valid")
	}
}
