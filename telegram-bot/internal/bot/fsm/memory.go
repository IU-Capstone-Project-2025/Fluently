// memory.go
package fsm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"telegram-bot/internal/domain"
)

// UserProgress represents the user's learning progress
type UserProgress struct {
	UserID       int64     `json:"user_id"`
	TelegramID   int64     `json:"telegram_id"`
	GoogleLinked bool      `json:"google_linked"`
	State        UserState `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`

	// Learning preferences
	WordsPerDay    int    `json:"words_per_day"`   // Number of words to learn per day
	CEFRLevel      string `json:"cefr_level"`      // User's determined CEFR level (A1-C2)
	Notifications  bool   `json:"notifications"`   // Whether notifications are enabled
	NotificationAt string `json:"notification_at"` // Time for daily notifications (HH:MM)

	// Onboarding data
	OnboardingComplete bool `json:"onboarding_complete"`

	// Questionnaire data
	Goal       string `json:"goal"`
	Confidence int    `json:"confidence"`
	Experience string `json:"experience"`

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

// UserStateManager handles FSM state for users
type UserStateManager struct {
	redisClient *redis.Client
}

// Temporary data types
type TempDataType string

const (
	// Temp data types for different flows
	TempDataCEFRTest   TempDataType = "cefr_test"
	TempDataLesson     TempDataType = "lesson"
	TempDataSettings   TempDataType = "settings"
	TempDataExercise   TempDataType = "exercise"
	TempDataOnboarding TempDataType = "onboarding"
)

// CEFRTestData holds temporary data for CEFR test flow
type CEFRTestData struct {
	Questions      []map[string]interface{} `json:"questions"`
	CurrentGroup   int                      `json:"current_group"`
	Answers        map[string]string        `json:"answers"`
	CorrectAnswers int                      `json:"correct_answers"`
	StartTime      time.Time                `json:"start_time"`
	EndTime        time.Time                `json:"end_time"`
}

// LessonData holds temporary data for a lesson flow
type LessonData struct {
	Words              []map[string]interface{} `json:"words"`
	CurrentWordIndex   int                      `json:"current_word_index"`
	CurrentBlockIndex  int                      `json:"current_block_index"`
	CompletedExercises int                      `json:"completed_exercises"`
	Progress           float64                  `json:"progress"`
	StartTime          time.Time                `json:"start_time"`
}

// SettingsData holds temporary data for settings flow
type SettingsData struct {
	SettingType   string      `json:"setting_type"`   // which setting is being modified
	CurrentValue  interface{} `json:"current_value"`  // current value being edited
	ProposedValue interface{} `json:"proposed_value"` // new value being considered
	TimeFormat    string      `json:"time_format"`    // for time settings
}

// ExerciseData holds temporary data for exercise flow
type ExerciseData struct {
	ExerciseType  string                   `json:"exercise_type"`
	Word          map[string]interface{}   `json:"word"`
	Options       []map[string]interface{} `json:"options"`
	CorrectAnswer string                   `json:"correct_answer"`
	UserAnswer    string                   `json:"user_answer"`
	IsCorrect     bool                     `json:"is_correct"`
	Attempts      int                      `json:"attempts"`
}

// OnboardingData holds temporary data for onboarding flow
type OnboardingData struct {
	Goal       string `json:"goal"`
	Confidence string `json:"confidence"`
	Experience string `json:"experience"`
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
		WordsPerDay:        10,      // Default: 10 words per day
		CEFRLevel:          "A1",    // Default: A1 level
		Notifications:      false,   // Default: notifications disabled
		NotificationAt:     "10:00", // Default: 10:00 AM notification
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

// NewUserStateManager creates a new state manager with Redis
func NewUserStateManager(redisClient *redis.Client) *UserStateManager {
	return &UserStateManager{
		redisClient: redisClient,
	}
}

// key generation helpers
func userStateKey(userID int64) string {
	return fmt.Sprintf("user:%d:state", userID)
}

func userTempDataKey(userID int64, dataType TempDataType) string {
	return fmt.Sprintf("user:%d:temp:%s", userID, dataType)
}

// GetState retrieves the current state for a user
func (m *UserStateManager) GetState(ctx context.Context, userID int64) (UserState, error) {
	state, err := m.redisClient.Get(ctx, userStateKey(userID)).Result()
	if err == redis.Nil {
		// No state found, use initial state
		return GetInitialState(), nil
	} else if err != nil {
		return "", fmt.Errorf("failed to get state: %w", err)
	}

	return UserState(state), nil
}

// SetState sets the current state for a user
func (m *UserStateManager) SetState(ctx context.Context, userID int64, state UserState) error {
	currentState, err := m.GetState(ctx, userID)
	if err != nil {
		return err
	}

	// Validate the state transition
	if !IsValidTransition(currentState, state) {
		return fmt.Errorf("invalid state transition from %s to %s", currentState, state)
	}

	// Set the new state with an expiration time (30 days)
	err = m.redisClient.Set(ctx, userStateKey(userID), string(state), 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set state: %w", err)
	}

	return nil
}

// ForceState sets the state without transition validation (for error recovery)
func (m *UserStateManager) ForceState(ctx context.Context, userID int64, state UserState) error {
	err := m.redisClient.Set(ctx, userStateKey(userID), string(state), 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to force state: %w", err)
	}

	return nil
}

// ClearState removes state information for a user
func (m *UserStateManager) ClearState(ctx context.Context, userID int64) error {
	err := m.redisClient.Del(ctx, userStateKey(userID)).Err()
	if err != nil {
		return fmt.Errorf("failed to clear state: %w", err)
	}

	return nil
}

// ResetUserToInitial resets user to initial state and clears all temp data
func (m *UserStateManager) ResetUserToInitial(ctx context.Context, userID int64) error {
	// Set initial state
	err := m.ForceState(ctx, userID, GetInitialState())
	if err != nil {
		return err
	}

	// Clear all temp data
	dataTypes := []TempDataType{
		TempDataCEFRTest,
		TempDataLesson,
		TempDataSettings,
		TempDataExercise,
		TempDataOnboarding,
	}

	for _, dt := range dataTypes {
		key := userTempDataKey(userID, dt)
		m.redisClient.Del(ctx, key)
	}

	return nil
}

// StoreTempData stores temporary data for a specific flow
func (m *UserStateManager) StoreTempData(ctx context.Context, userID int64, dataType TempDataType, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal temp data: %w", err)
	}

	// Store with 24-hour expiration
	err = m.redisClient.Set(ctx, userTempDataKey(userID, dataType), jsonData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store temp data: %w", err)
	}

	return nil
}

// GetCEFRTestData retrieves CEFR test data
func (m *UserStateManager) GetCEFRTestData(ctx context.Context, userID int64) (*CEFRTestData, error) {
	data := &CEFRTestData{}
	jsonData, err := m.redisClient.Get(ctx, userTempDataKey(userID, TempDataCEFRTest)).Result()
	if err == redis.Nil {
		return data, nil // Return empty struct if no data found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get CEFR test data: %w", err)
	}

	err = json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal CEFR test data: %w", err)
	}

	return data, nil
}

// GetLessonData retrieves lesson data (legacy)
func (m *UserStateManager) GetLessonData(ctx context.Context, userID int64) (*LessonData, error) {
	data := &LessonData{}
	jsonData, err := m.redisClient.Get(ctx, userTempDataKey(userID, TempDataLesson)).Result()
	if err == redis.Nil {
		return data, nil // Return empty struct if no data found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get lesson data: %w", err)
	}

	err = json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal lesson data: %w", err)
	}

	return data, nil
}

// GetSettingsData retrieves settings data
func (m *UserStateManager) GetSettingsData(ctx context.Context, userID int64) (*SettingsData, error) {
	data := &SettingsData{}
	jsonData, err := m.redisClient.Get(ctx, userTempDataKey(userID, TempDataSettings)).Result()
	if err == redis.Nil {
		return data, nil // Return empty struct if no data found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get settings data: %w", err)
	}

	err = json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings data: %w", err)
	}

	return data, nil
}

// GetExerciseData retrieves exercise data
func (m *UserStateManager) GetExerciseData(ctx context.Context, userID int64) (*ExerciseData, error) {
	data := &ExerciseData{}
	jsonData, err := m.redisClient.Get(ctx, userTempDataKey(userID, TempDataExercise)).Result()
	if err == redis.Nil {
		return data, nil // Return empty struct if no data found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get exercise data: %w", err)
	}

	err = json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal exercise data: %w", err)
	}

	return data, nil
}

// GetOnboardingData retrieves onboarding data
func (m *UserStateManager) GetOnboardingData(ctx context.Context, userID int64) (*OnboardingData, error) {
	data := &OnboardingData{}
	jsonData, err := m.redisClient.Get(ctx, userTempDataKey(userID, TempDataOnboarding)).Result()
	if err == redis.Nil {
		return data, nil // Return empty struct if no data found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get onboarding data: %w", err)
	}

	err = json.Unmarshal([]byte(jsonData), data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal onboarding data: %w", err)
	}

	return data, nil
}

// ClearTempData removes temporary data for a specific flow
func (m *UserStateManager) ClearTempData(ctx context.Context, userID int64, dataType TempDataType) error {
	err := m.redisClient.Del(ctx, userTempDataKey(userID, dataType)).Err()
	if err != nil {
		return fmt.Errorf("failed to clear temp data: %w", err)
	}

	return nil
}

// IsInState checks if the user is in a specific state
func (m *UserStateManager) IsInState(ctx context.Context, userID int64, state UserState) (bool, error) {
	currentState, err := m.GetState(ctx, userID)
	if err != nil {
		return false, err
	}

	return currentState == state, nil
}

// IsInStateGroup checks if the user is in any of the specified states
func (m *UserStateManager) IsInStateGroup(ctx context.Context, userID int64, stateGroup func(UserState) bool) (bool, error) {
	currentState, err := m.GetState(ctx, userID)
	if err != nil {
		return false, err
	}

	return stateGroup(currentState), nil
}

// IsInLessonState checks if the user is in a lesson state
func (m *UserStateManager) IsInLessonState(ctx context.Context, userID int64) (bool, error) {
	return m.IsInStateGroup(ctx, userID, IsLessonState)
}

// IsInExerciseState checks if the user is in an exercise state
func (m *UserStateManager) IsInExerciseState(ctx context.Context, userID int64) (bool, error) {
	return m.IsInStateGroup(ctx, userID, IsExerciseState)
}

// IsInSettingsState checks if the user is in a settings state
func (m *UserStateManager) IsInSettingsState(ctx context.Context, userID int64) (bool, error) {
	return m.IsInStateGroup(ctx, userID, IsSettingsState)
}

// IsInCEFRTestState checks if the user is in a CEFR test state
func (m *UserStateManager) IsInCEFRTestState(ctx context.Context, userID int64) (bool, error) {
	return m.IsInStateGroup(ctx, userID, IsCEFRTestState)
}

// GetUserChatID extracts the chat ID from a user ID for Telegram
// In most cases, chat ID = user ID for private chats
func GetUserChatID(userID int64) int64 {
	return userID
}

// GetUserIDFromChatID extracts the user ID from a chat ID
// For group chats, this would be different, but for now we're focusing on private chats
func GetUserIDFromChatID(chatID int64) (int64, error) {
	return chatID, nil
}

// GetUserIDFromString converts a string user ID to int64
func GetUserIDFromString(userIDStr string) (int64, error) {
	return strconv.ParseInt(userIDStr, 10, 64)
}

// WrongStateError is returned when a handler is called with the wrong state
type WrongStateError struct {
	Expected UserState
	Actual   UserState
}

func (e *WrongStateError) Error() string {
	return fmt.Sprintf("wrong state: expected %s, got %s", e.Expected, e.Actual)
}

// NewWrongStateError creates a new WrongStateError
func NewWrongStateError(expected, actual UserState) *WrongStateError {
	return &WrongStateError{
		Expected: expected,
		Actual:   actual,
	}
}

// IsWrongStateError checks if an error is a WrongStateError
func IsWrongStateError(err error) bool {
	var wse *WrongStateError
	return errors.As(err, &wse)
}

// GetLessonProgress retrieves the new lesson progress structure
func (m *UserStateManager) GetLessonProgress(ctx context.Context, userID int64) (*domain.LessonProgress, error) {
	jsonData, err := m.redisClient.Get(ctx, userTempDataKey(userID, TempDataLesson)).Result()
	if err == redis.Nil {
		return nil, nil // Return nil if no data found
	} else if err != nil {
		return nil, fmt.Errorf("failed to get lesson progress: %w", err)
	}

	var progress domain.LessonProgress
	err = json.Unmarshal([]byte(jsonData), &progress)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal lesson progress: %w", err)
	}

	return &progress, nil
}

// StoreLessonProgress stores the new lesson progress structure
func (m *UserStateManager) StoreLessonProgress(ctx context.Context, userID int64, progress *domain.LessonProgress) error {
	jsonData, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("failed to marshal lesson progress: %w", err)
	}

	// Store with 24-hour expiration
	err = m.redisClient.Set(ctx, userTempDataKey(userID, TempDataLesson), jsonData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store lesson progress: %w", err)
	}

	return nil
}

// UpdateLessonProgress updates specific fields in lesson progress
func (m *UserStateManager) UpdateLessonProgress(ctx context.Context, userID int64, updateFunc func(*domain.LessonProgress) error) error {
	progress, err := m.GetLessonProgress(ctx, userID)
	if err != nil {
		return err
	}

	if progress == nil {
		return fmt.Errorf("no lesson progress found for user %d", userID)
	}

	err = updateFunc(progress)
	if err != nil {
		return err
	}

	return m.StoreLessonProgress(ctx, userID, progress)
}

// ClearLessonProgress removes lesson progress data
func (m *UserStateManager) ClearLessonProgress(ctx context.Context, userID int64) error {
	return m.ClearTempData(ctx, userID, TempDataLesson)
}

// HasActiveLessonProgress checks if user has an active lesson in progress
func (m *UserStateManager) HasActiveLessonProgress(ctx context.Context, userID int64) (bool, error) {
	progress, err := m.GetLessonProgress(ctx, userID)
	if err != nil {
		return false, err
	}

	return progress != nil && progress.LessonData != nil, nil
}

// GetCurrentWordSet gets the current set of 3 words being studied
func (m *UserStateManager) GetCurrentWordSet(ctx context.Context, userID int64) ([]domain.Card, error) {
	progress, err := m.GetLessonProgress(ctx, userID)
	if err != nil {
		return nil, err
	}

	if progress == nil || len(progress.WordsInCurrentSet) == 0 {
		return nil, fmt.Errorf("no current word set found")
	}

	return progress.WordsInCurrentSet, nil
}

// AddWordProgress adds a word to the learned words list
func (m *UserStateManager) AddWordProgress(ctx context.Context, userID int64, wordProgress domain.WordProgress) error {
	return m.UpdateLessonProgress(ctx, userID, func(progress *domain.LessonProgress) error {
		progress.WordsLearned = append(progress.WordsLearned, wordProgress)
		progress.LearnedCount++
		progress.LastActivity = time.Now()
		return nil
	})
}

// GetJWTToken retrieves user's JWT token for API calls
func (m *UserStateManager) GetJWTToken(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("user:%d:jwt_token", userID)
	token, err := m.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("no JWT token found for user %d", userID)
	} else if err != nil {
		return "", fmt.Errorf("failed to get JWT token: %w", err)
	}

	return token, nil
}

// StoreJWTToken stores user's JWT token with expiration
func (m *UserStateManager) StoreJWTToken(ctx context.Context, userID int64, token string, expiration time.Duration) error {
	key := fmt.Sprintf("user:%d:jwt_token", userID)
	err := m.redisClient.Set(ctx, key, token, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store JWT token: %w", err)
	}

	return nil
}
