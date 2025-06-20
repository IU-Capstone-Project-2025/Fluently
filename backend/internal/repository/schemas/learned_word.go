package schemas

import (
	"time"

	"github.com/google/uuid"
)

type CreateLearnedWordRequest struct {
	UserID          uuid.UUID `json:"user_id" binding:"required"`
	WordID          uuid.UUID `json:"word_id" binding:"required"`
	LearnedAt       time.Time `json:"learned_at"`
	LastReviewed    time.Time `json:"last_reviewed"`
	CntReviewed     int       `json:"cnt_reviewed"`
	ConfidenceScore int       `json:"confidence_score"`
}

type LearnedWordResponse struct {
	UserID          uuid.UUID `json:"user_id"`
	WordID          uuid.UUID `json:"word_id"`
	LearnedAt       time.Time `json:"learned_at"`
	LastReviewed    time.Time `json:"last_reviewed"`
	CntReviewedAt   int       `json:"cnt_reviewed"`
	ConfidenceScore int       `json:"confidence_score"`
}
