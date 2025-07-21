package schemas

import (
	"time"

	"github.com/google/uuid"
)

// AddNotLearnedWordRequest is a request body for adding a word to not learned words
type AddNotLearnedWordRequest struct {
	Word string `json:"word" binding:"required"` // The word text to add
}

// NotLearnedWordResponse is a response body for a not learned word
type NotLearnedWordResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	WordID    uuid.UUID `json:"word_id"`
	Word      string    `json:"word"`
	CreatedAt time.Time `json:"created_at"`
}

// ListNotLearnedWordsResponse is a response for listing not learned words
type ListNotLearnedWordsResponse struct {
	Words []NotLearnedWordResponse `json:"words"`
	Total int                      `json:"total"`
}
