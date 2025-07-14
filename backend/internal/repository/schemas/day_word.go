package schemas

import "github.com/google/uuid"

type DayWordResponse struct {
	WordID        uuid.UUID  `json:"word_id"`
	Word          string     `json:"word"`
	Translation   string     `json:"translation"`
	Transcription *string     `json:"transcription,omitempty"`
	CEFRLevel     string     `json:"cefr_level"`
	IsLearned     bool       `json:"is_learned"`
	Topic         string     `json:"topic"`
	Subtopic      string     `json:"subtopic"`
	Sentences     []Sentence `json:"sentences"`
	Exercise      Exercise   `json:"exercise"`
}
