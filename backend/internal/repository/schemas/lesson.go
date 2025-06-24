package schemas

import "github.com/google/uuid"

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

type SyncSchema struct {
	Dirty        bool   `json:"dirty"`
	LastSyncedAt string `json:"last_synced_at"`
}

type LessonResponse struct {
	Lesson LessonSchema `json:"lesson"`
	Cards  []CardSchema `json:"cards"`
	Sync   SyncSchema   `json:"sync"`
}
