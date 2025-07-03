// memory.go
package fsm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UserProgress represents the user's learning progress
type UserProgress struct {
	UserID       int64     `json:"user_id"`
	TelegramID   int64     `json:"telegram_id"`
	GoogleLinked bool      `json:"google_linked"`
	State        UserState `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`

	// Onboarding data
	OnboardingComplete bool `json:"onboarding_complete"`

	// Questionnaire data
	Goal         string `json:"goal"`
	Confidence   int    `json:"confidence"`
	SerialHabits string `json:"serial_habits"`
	Experience   string `json:"experience"`

	// Vocabulary test
	VocabularyTest  VocabTestData `json:"vocabulary_test"`
	DeterminedLevel string        `json:"determined_level"`

	// Lesson progress
	CurrentLesson LessonProgress `json:"current_lesson"`
	DailyStats    DailyStats     `json:"daily_stats"`
	TotalStats    TotalStats     `json:"total_stats"`
}

// VocabTestData stores vocabulary test results
type VocabTestData struct {
	Group1Score   int  `json:"group1_score"`
	Group2Score   int  `json:"group2_score"`
	Group3Score   int  `json:"group3_score"`
	Group4Score   int  `json:"group4_score"`
	Group5Score   int  `json:"group5_score"`
	TestCompleted bool `json:"test_completed"`
	CurrentGroup  int  `json:"current_group"`
}

// LessonProgress tracks current lesson state
type LessonProgress struct {
	LessonID           uuid.UUID   `json:"lesson_id"`
	WordsShown         []WordShown `json:"words_shown"`
	CurrentWordIndex   int         `json:"current_word_index"`
	ExercisesCompleted int         `json:"exercises_completed"`
	LessonStarted      time.Time   `json:"lesson_started"`
	LessonComplete     bool        `json:"lesson_complete"`

	// Exercise state
	CurrentExercise Exercise   `json:"current_exercise"`
	ExerciseQueue   []Exercise `json:"exercise_queue"`

	// Words currently being learned in this lesson
	CurrentWords    []LessonWord `json:"current_words"`
	FirstBlockShown bool         `json:"first_block_shown"`
}

// WordShown represents a word that has been shown to the user
type WordShown struct {
	WordID        uuid.UUID `json:"word_id"`
	Word          string    `json:"word"`
	Translation   string    `json:"translation"`
	ShownAt       time.Time `json:"shown_at"`
	TimesRepeated int       `json:"times_repeated"`
	LastCorrect   bool      `json:"last_correct"`
}

// LessonWord represents a word in the current lesson
type LessonWord struct {
	WordID       uuid.UUID  `json:"word_id"`
	Word         string     `json:"word"`
	Translation  string     `json:"translation"`
	CEFRLevel    string     `json:"cefr_level"`
	Sentences    []Sentence `json:"sentences"`
	AudioURL     string     `json:"audio_url"`
	PartOfSpeech string     `json:"part_of_speech"`
	IsNew        bool       `json:"is_new"`
}

// Sentence represents a sentence with translation
type Sentence struct {
	SentenceID  uuid.UUID `json:"sentence_id"`
	Text        string    `json:"text"`
	Translation string    `json:"translation"`
	AudioURL    string    `json:"audio_url"`
}

// Exercise represents an exercise
type Exercise struct {
	ExerciseID    uuid.UUID `json:"exercise_id"`
	Type          string    `json:"type"` // "audio_dictation", "translation"
	WordID        uuid.UUID `json:"word_id"`
	Word          string    `json:"word"`
	Translation   string    `json:"translation"`
	AudioURL      string    `json:"audio_url"`
	CorrectAnswer string    `json:"correct_answer"`
	AttemptCount  int       `json:"attempt_count"`
	MaxAttempts   int       `json:"max_attempts"`
	CreatedAt     time.Time `json:"created_at"`
}

// DailyStats tracks daily learning statistics
type DailyStats struct {
	Date               string `json:"date"`
	NewWordsLearned    int    `json:"new_words_learned"`
	ExercisesCompleted int    `json:"exercises_completed"`
	CorrectAnswers     int    `json:"correct_answers"`
	TotalAttempts      int    `json:"total_attempts"`
	TimeSpent          int    `json:"time_spent"` // in minutes
	LessonsCompleted   int    `json:"lessons_completed"`
}

// TotalStats tracks overall learning statistics
type TotalStats struct {
	TotalWordsLearned int       `json:"total_words_learned"`
	TotalLessons      int       `json:"total_lessons"`
	TotalExercises    int       `json:"total_exercises"`
	TotalTimeSpent    int       `json:"total_time_spent"` // in minutes
	StreakDays        int       `json:"streak_days"`
	LastLessonDate    time.Time `json:"last_lesson_date"`
	LevelAchieved     string    `json:"level_achieved"`
}

// SessionData stores temporary session information
type SessionData struct {
	UserID        int64                  `json:"user_id"`
	TempData      map[string]interface{} `json:"temp_data"`
	LastMessageID int                    `json:"last_message_id"`
	CallbackData  string                 `json:"callback_data"`
	ExpiresAt     time.Time              `json:"expires_at"`
}

// CreateNewUserProgress creates a new user progress instance
func CreateNewUserProgress(telegramID int64) *UserProgress {
	now := time.Now()
	return &UserProgress{
		TelegramID:         telegramID,
		GoogleLinked:       false,
		State:              StateStart,
		CreatedAt:          now,
		LastActivity:       now,
		OnboardingComplete: false,
		VocabularyTest: VocabTestData{
			CurrentGroup: 1,
		},
		CurrentLesson: LessonProgress{
			CurrentWords:  make([]LessonWord, 0),
			WordsShown:    make([]WordShown, 0),
			ExerciseQueue: make([]Exercise, 0),
		},
		DailyStats: DailyStats{
			Date: now.Format("2006-01-02"),
		},
		TotalStats: TotalStats{
			LastLessonDate: now,
		},
	}
}

// UpdateState updates the user's current state
func (up *UserProgress) UpdateState(newState UserState) error {
	if !IsValidTransition(up.State, newState) {
		return fmt.Errorf("invalid state transition from %s to %s", up.State, newState)
	}
	up.State = newState
	up.LastActivity = time.Now()
	return nil
}

// AddWordShown adds a word to the list of shown words
func (up *UserProgress) AddWordShown(wordID uuid.UUID, word, translation string) {
	wordShown := WordShown{
		WordID:        wordID,
		Word:          word,
		Translation:   translation,
		ShownAt:       time.Now(),
		TimesRepeated: 1,
	}
	up.CurrentLesson.WordsShown = append(up.CurrentLesson.WordsShown, wordShown)
}

// UpdateWordRepetition updates the repetition count for a word
func (up *UserProgress) UpdateWordRepetition(wordID uuid.UUID, correct bool) {
	for i := range up.CurrentLesson.WordsShown {
		if up.CurrentLesson.WordsShown[i].WordID == wordID {
			up.CurrentLesson.WordsShown[i].TimesRepeated++
			up.CurrentLesson.WordsShown[i].LastCorrect = correct
			break
		}
	}
}

// AddExerciseToQueue adds an exercise to the exercise queue
func (up *UserProgress) AddExerciseToQueue(exercise Exercise) {
	up.CurrentLesson.ExerciseQueue = append(up.CurrentLesson.ExerciseQueue, exercise)
}

// GetNextExercise retrieves the next exercise from the queue
func (up *UserProgress) GetNextExercise() (*Exercise, bool) {
	if len(up.CurrentLesson.ExerciseQueue) == 0 {
		return nil, false
	}

	exercise := up.CurrentLesson.ExerciseQueue[0]
	up.CurrentLesson.ExerciseQueue = up.CurrentLesson.ExerciseQueue[1:]
	up.CurrentLesson.CurrentExercise = exercise

	return &exercise, true
}

// CompleteExercise marks the current exercise as completed
func (up *UserProgress) CompleteExercise(correct bool) {
	up.CurrentLesson.ExercisesCompleted++
	up.DailyStats.ExercisesCompleted++
	up.TotalStats.TotalExercises++
	up.DailyStats.TotalAttempts++

	if correct {
		up.DailyStats.CorrectAnswers++
	}

	// Update word repetition
	if up.CurrentLesson.CurrentExercise.WordID != uuid.Nil {
		up.UpdateWordRepetition(up.CurrentLesson.CurrentExercise.WordID, correct)
	}
}

// StartNewLesson starts a new lesson
func (up *UserProgress) StartNewLesson(lessonID uuid.UUID) {
	up.CurrentLesson = LessonProgress{
		LessonID:           lessonID,
		WordsShown:         make([]WordShown, 0),
		CurrentWordIndex:   0,
		ExercisesCompleted: 0,
		LessonStarted:      time.Now(),
		LessonComplete:     false,
		CurrentWords:       make([]LessonWord, 0),
		ExerciseQueue:      make([]Exercise, 0),
		FirstBlockShown:    false,
	}
}

// CompleteLesson marks the current lesson as complete
func (up *UserProgress) CompleteLesson() {
	up.CurrentLesson.LessonComplete = true
	up.DailyStats.LessonsCompleted++
	up.TotalStats.TotalLessons++
	up.TotalStats.LastLessonDate = time.Now()

	// Calculate time spent (approximate)
	timeSpent := int(time.Since(up.CurrentLesson.LessonStarted).Minutes())
	up.DailyStats.TimeSpent += timeSpent
	up.TotalStats.TotalTimeSpent += timeSpent
}

// UpdateDailyStats updates daily statistics
func (up *UserProgress) UpdateDailyStats() {
	today := time.Now().Format("2006-01-02")
	if up.DailyStats.Date != today {
		// New day - reset daily stats
		up.DailyStats = DailyStats{
			Date: today,
		}
	}
}

// SetVocabTestScore sets the score for a vocabulary test group
func (up *UserProgress) SetVocabTestScore(group int, score int) {
	switch group {
	case 1:
		up.VocabularyTest.Group1Score = score
	case 2:
		up.VocabularyTest.Group2Score = score
	case 3:
		up.VocabularyTest.Group3Score = score
	case 4:
		up.VocabularyTest.Group4Score = score
	case 5:
		up.VocabularyTest.Group5Score = score
	}

	if group == 5 {
		up.VocabularyTest.TestCompleted = true
	}
}

// GetAccuracyRate returns the user's accuracy rate for exercises
func (up *UserProgress) GetAccuracyRate() float64 {
	if up.DailyStats.TotalAttempts == 0 {
		return 0.0
	}
	return float64(up.DailyStats.CorrectAnswers) / float64(up.DailyStats.TotalAttempts) * 100
}

// ToJSON converts UserProgress to JSON
func (up *UserProgress) ToJSON() ([]byte, error) {
	return json.Marshal(up)
}

// FromJSON creates UserProgress from JSON
func FromJSON(data []byte) (*UserProgress, error) {
	var up UserProgress
	err := json.Unmarshal(data, &up)
	return &up, err
}

// CreateSessionData creates new session data
func CreateSessionData(userID int64, expiration time.Duration) *SessionData {
	return &SessionData{
		UserID:    userID,
		TempData:  make(map[string]interface{}),
		ExpiresAt: time.Now().Add(expiration),
	}
}

// SetTempData stores temporary data in the session
func (sd *SessionData) SetTempData(key string, value interface{}) {
	sd.TempData[key] = value
}

// GetTempData retrieves temporary data from the session
func (sd *SessionData) GetTempData(key string) (interface{}, bool) {
	value, exists := sd.TempData[key]
	return value, exists
}

// IsExpired checks if the session has expired
func (sd *SessionData) IsExpired() bool {
	return time.Now().After(sd.ExpiresAt)
}

// ExtendExpiration extends the session expiration time
func (sd *SessionData) ExtendExpiration(duration time.Duration) {
	sd.ExpiresAt = time.Now().Add(duration)
}
