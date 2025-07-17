package schemas

import "github.com/google/uuid"

// CreateSentenceRequest is a request body for creating a sentence
type CreateSentenceRequest struct {
	WordID      uuid.UUID `json:"word_id" binding:"required"`
	Sentence    string    `json:"sentence" binding:"required"`
	Translation string    `json:"translation"`
	AudioURL    string    `json:"audio_url"`
}

// SentenceResponse is a response for a sentence
type SentenceResponse struct {
	ID          string `json:"id"`
	WordID      string `json:"word_id"`
	Sentence    string `json:"sentence"`
	Translation string `json:"translation"`
	AudioURL    string `json:"audio_url,omitempty"`
}
