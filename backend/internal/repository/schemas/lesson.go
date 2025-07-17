package schemas

import (
	"github.com/google/uuid"
)

// Main Response
type LessonResponse struct {
	Lesson LessonInfo `json:"lesson"`
	Cards  []Card     `json:"cards"`
}

// Lesson information
type LessonInfo struct {
	StartedAt      string `json:"started_at"`
	WordsPerLesson int    `json:"words_per_lesson"`
	TotalWords     int    `json:"total_words"`
	CEFRLevel      string `json:"cefr_level"`
}

// Card with word and sentences
type Card struct {
	WordID        uuid.UUID  `json:"word_id"`
	Word          string     `json:"word"`
	Translation   string     `json:"translation"`
	Transcription string     `json:"transcription,omitempty"`
	CEFRLevel     string     `json:"cefr_level,omitempty"`
	IsLearned     *bool      `json:"is_learned,omitempty"`
	Topic         string     `json:"topic"`
	Subtopic      string     `json:"subtopic"`
	Sentences     []Sentence `json:"sentences"`
	Exercise      Exercise   `json:"exercise"`
}

// Sentence for examples
type Sentence struct {
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

// Exercise with type
type Exercise struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// ------------ Structures for different exercise types ------------------

// ExercisesTranslateRuToEn is a struct for translate exercise
type ExerciseTranslateRuToEn struct {
	Text          string   `json:"text"`
	CorrectAnswer string   `json:"correct_answer"`
	PickOptions   []string `json:"pick_options"`
}

// ExerciseWriteWordFromTranslation is a struct for write exercise
type ExerciseWriteWordFromTranslation struct {
	Translation   string `json:"translation"`
	CorrectAnswer string `json:"correct_answer"`
}

// ExercisePickOptionSentence is a struct for pick option exercise
type ExercisePickOptionSentence struct {
	Template      string   `json:"template"`
	CorrectAnswer string   `json:"correct_answer"`
	PickOptions   []string `json:"pick_options"`
}
