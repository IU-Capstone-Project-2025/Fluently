package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserProgress represents a user's learning progress and preferences
type UserProgress struct {
	UserID           int64                  `json:"user_id"`
	CEFRLevel        string                 `json:"cefr_level"`        // A1, A2, B1, B2, C1, C2
	WordsPerDay      int                    `json:"words_per_day"`     // Number of words to learn per day
	NotificationTime string                 `json:"notification_time"` // Time for daily notifications (HH:MM)
	LearnedWords     int                    `json:"learned_words"`     // Total learned words count
	CurrentStreak    int                    `json:"current_streak"`    // Current learning streak in days
	LongestStreak    int                    `json:"longest_streak"`    // Longest learning streak in days
	LastActivity     time.Time              `json:"last_activity"`     // Last activity timestamp
	StartDate        string                 `json:"start_date"`        // When user started learning (YYYY-MM-DD)
	Preferences      map[string]interface{} `json:"preferences"`       // Additional preferences
}

// New lesson structure matching the backend JSON format

// Lesson represents the lesson metadata
type Lesson struct {
	StartedAt      string `json:"started_at"`
	WordsPerLesson int    `json:"words_per_lesson"`
	TotalWords     int    `json:"total_words"`
	CEFRLevel      string `json:"cefr_level"`
}

// Sentence represents a sentence with its translation
type Sentence struct {
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

// ExerciseData represents different types of exercise data
type ExerciseData struct {
	Template      string   `json:"template,omitempty"` // For pick_option_sentence
	CorrectAnswer string   `json:"correct_answer"`
	PickOptions   []string `json:"pick_options,omitempty"` // For multiple choice exercises
	Translation   string   `json:"translation,omitempty"`  // For write_word_from_translation
	Text          string   `json:"text,omitempty"`         // For translate_ru_to_en
}

// Exercise represents an exercise for a word
type Exercise struct {
	Type string       `json:"type"`
	Data ExerciseData `json:"data"`
}

// Card represents a word card with all its learning data
type Card struct {
	WordID      string     `json:"word_id"`
	Word        string     `json:"word"`
	Translation string     `json:"translation"`
	Topic       string     `json:"topic"`
	Subtopic    string     `json:"subtopic"`
	Sentences   []Sentence `json:"sentences"`
	Exercise    Exercise   `json:"exercise"`
}

// LessonResponse represents the complete lesson response from backend
type LessonResponse struct {
	Lesson Lesson `json:"lesson"`
	Cards  []Card `json:"cards"`
}

// Learning progress tracking models

// WordProgress represents progress for a single word
type WordProgress struct {
	Word            string    `json:"word"`
	Translation     string    `json:"translation"` // Add translation for display
	WordID          string    `json:"word_id"`     // Add WordID for backend API
	LearnedAt       time.Time `json:"learned_at"`
	ConfidenceScore int       `json:"confidence_score"`
	CntReviewed     int       `json:"cnt_reviewed"`
	AlreadyKnown    bool      `json:"already_known"` // Flag to mark words as "already known"
}

// BadlyAnsweredWord represents a word that was answered incorrectly
type BadlyAnsweredWord struct {
	WordID string `json:"word_id"`
}

// LessonProgress represents overall lesson progress stored in Redis
type LessonProgress struct {
	LessonData         *LessonResponse     `json:"lesson_data"`
	CurrentWordIndex   int                 `json:"current_word_index"`
	CurrentPhase       string              `json:"current_phase"` // "showing_words", "exercises", "completed", "retry"
	WordsInCurrentSet  []Card              `json:"words_in_current_set"`
	CurrentSetIndex    int                 `json:"current_set_index"`
	ExerciseIndex      int                 `json:"exercise_index"`
	WordsLearned       []WordProgress      `json:"words_learned"`
	BadlyAnsweredWords []BadlyAnsweredWord `json:"badly_answered_words"` // Words answered incorrectly
	RetryWords         []Card              `json:"retry_words"`          // Words that need to be retried
	RetryIndex         int                 `json:"retry_index"`          // Current index in retry queue
	StartTime          time.Time           `json:"start_time"`
	LastActivity       time.Time           `json:"last_activity"`
	LearnedCount       int                 `json:"learned_count"`       // Count of words actually learned (goal: 10)
	AlreadyKnownCount  int                 `json:"already_known_count"` // Count of words marked as already known
}

// Legacy models - keeping for backward compatibility
type Word struct {
	ID             int64     `json:"id"`
	Word           string    `json:"word"`
	Translation    string    `json:"translation"`
	Examples       []string  `json:"examples"`
	AudioURL       string    `json:"audio_url"`
	CEFRLevel      string    `json:"cefr_level"`      // A1, A2, B1, B2, C1, C2
	NextReview     time.Time `json:"next_review"`     // When to review next
	ReviewCount    int       `json:"review_count"`    // How many times reviewed
	CorrectCount   int       `json:"correct_count"`   // How many times answered correctly
	IncorrectCount int       `json:"incorrect_count"` // How many times answered incorrectly
	Learned        bool      `json:"learned"`         // Whether the word is considered learned
	LearningScore  float64   `json:"learning_score"`  // Score based on spaced repetition algorithm
}

// NewUserProgress creates a new UserProgress with default values
func NewUserProgress(userID int64) *UserProgress {
	return &UserProgress{
		UserID:           userID,
		CEFRLevel:        "",
		WordsPerDay:      5,       // Default to 5 words per day
		NotificationTime: "10:00", // Default to 10:00 AM
		LearnedWords:     0,
		CurrentStreak:    0,
		LongestStreak:    0,
		LastActivity:     time.Now(),
		StartDate:        time.Now().Format("2006-01-02"),
		Preferences:      map[string]interface{}{},
	}
}

// NotificationsEnabled checks if the user has enabled notifications
func (up *UserProgress) NotificationsEnabled() bool {
	return up.NotificationTime != ""
}

// GetPreferenceString returns a string preference value
func (up *UserProgress) GetPreferenceString(key string, defaultValue string) string {
	if val, ok := up.Preferences[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

// GetPreferenceBool returns a boolean preference value
func (up *UserProgress) GetPreferenceBool(key string, defaultValue bool) bool {
	if val, ok := up.Preferences[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}

// UpdateStreak updates the streak based on the last activity
func (up *UserProgress) UpdateStreak(now time.Time) {
	yesterday := now.AddDate(0, 0, -1)
	lastActivityDate := up.LastActivity.Format("2006-01-02")
	yesterdayDate := yesterday.Format("2006-01-02")
	todayDate := now.Format("2006-01-02")

	if lastActivityDate == todayDate {
		// Already updated today, no change needed
		return
	} else if lastActivityDate == yesterdayDate {
		// Activity yesterday, increment streak
		up.CurrentStreak++
	} else {
		// Activity gap, reset streak
		up.CurrentStreak = 1
	}

	// Update longest streak if current one is longer
	if up.CurrentStreak > up.LongestStreak {
		up.LongestStreak = up.CurrentStreak
	}

	up.LastActivity = now
}

// GetLearningSchedule generates a learning schedule for the next few days
func (up *UserProgress) GetLearningSchedule(days int) []int {
	schedule := make([]int, days)
	wordsPerDay := up.WordsPerDay

	for i := 0; i < days; i++ {
		schedule[i] = wordsPerDay
	}

	return schedule
}

// GetEffectiveLevel gets the effective CEFR level, returning a default if not set
func (up *UserProgress) GetEffectiveLevel() string {
	if up.CEFRLevel == "" {
		return "A1" // Default to beginner level
	}
	return up.CEFRLevel
}

// CreateReviewWord creates a review word from a word
func CreateReviewWord(word *Word) *Word {
	review := *word
	review.NextReview = calculateNextReview(word)
	return &review
}

// Helper function to calculate when to review a word next
func calculateNextReview(word *Word) time.Time {
	// Simple spaced repetition algorithm
	// 1 day + (reviewCount * factor)
	reviewCount := word.ReviewCount
	if reviewCount <= 0 {
		reviewCount = 1
	}

	// Calculate interval based on success rate
	var factor float64 = 1.0
	if word.CorrectCount+word.IncorrectCount > 0 {
		successRate := float64(word.CorrectCount) / float64(word.CorrectCount+word.IncorrectCount)
		factor = 0.5 + successRate*2.0 // Range from 0.5 to 2.5
	}

	// Calculate days to add (exponential spacing)
	daysToAdd := int(float64(reviewCount) * factor)
	if daysToAdd < 1 {
		daysToAdd = 1
	} else if daysToAdd > 30 {
		daysToAdd = 30 // Cap at 30 days
	}

	return time.Now().AddDate(0, 0, daysToAdd)
}

// Legacy schemas - keeping for backward compatibility but not used in new learning logic
type SentenceSchema struct {
	SentenceID  uuid.UUID `json:"sentence_id"`
	Text        string    `json:"text"`
	Translation string    `json:"translation"`
}

type ExerciseDataPickOption struct {
	WordID        uuid.UUID `json:"word_id"`
	Template      string    `json:"template"`
	Options       []string  `json:"options"`
	CorrectAnswer string    `json:"correct_answer"`
}

type ExerciseSchema struct {
	ExerciseID uuid.UUID              `json:"exercise_id"`
	Type       string                 `json:"type"`
	Data       ExerciseDataPickOption `json:"data"`
}

type CardSchema struct {
	WordID        uuid.UUID        `json:"word_id"`
	Word          string           `json:"word"`
	Translation   string           `json:"translation"`
	Transcription string           `json:"transcription"`
	CEFRLevel     string           `json:"cefr_level"`
	IsNew         bool             `json:"is_new"`
	Topic         string           `json:"topic"`
	Subtopic      string           `json:"subtopic"`
	Sentences     []SentenceSchema `json:"sentences"`
	Exercise      ExerciseSchema   `json:"exercise"`
}

type LessonSchema struct {
	LessonID       uuid.UUID `json:"lesson_id"`
	UserID         uuid.UUID `json:"user_id"`
	StartedAt      string    `json:"started_at"`
	WordsPerLesson int       `json:"words_per_lesson"`
	TotalWords     int       `json:"total_words"`
}

type ProgressSchema struct {
	CurrentCardIndex int           `json:"current_card_index"`
	WordsShown       []uuid.UUID   `json:"words_shown"`
	ExerciseShown    []uuid.UUID   `json:"exercise_shown"`
	Exercises        []interface{} `json:"exercises"`
}

type SyncSchema struct {
	Dirty        bool   `json:"dirty"`
	LastSyncedAt string `json:"last_synced_at"`
}
