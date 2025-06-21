package schemas

import "github.com/google/uuid"

type CreateSentenceRequest struct {
	WordID      uuid.UUID `json:"word_id" binding:"required"`
	Sentence    string    `json:"sentence" binding:"required"`
	Translation string    `json:"translation"`
}

type SentenceResponse struct {
	ID          string `json:"id"`
	WordID      string `json:"word_id"`
	Sentence    string `json:"sentence"`
	Translation string `json:"translation"`
}
