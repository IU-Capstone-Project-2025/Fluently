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

// Word represents a vocabulary word with learning data
type Word struct {
	ID             int64     `json:"id"`
	Word           string    `json:"word"`
	Translation    string    `json:"translation"`
	Definition     string    `json:"definition"`
	Examples       []string  `json:"examples"`
	AudioURL       string    `json:"audio_url"`
	ImageURL       string    `json:"image_url"`
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

type LessonResponse struct {
	Lesson   LessonSchema   `json:"lesson"`
	Cards    []CardSchema   `json:"cards"`
	Progress ProgressSchema `json:"progress"`
	Sync     SyncSchema     `json:"sync"`
}
