package models

import (
	"time"

	"github.com/google/uuid"
)

// LearnedWords is a model for learned words
type LearnedWords struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	WordID uuid.UUID `gorm:"type:uuid;not null"`

	LearnedAt    time.Time `gorm:"not null"`
	LastReviewed time.Time

	CountOfRevisions int `gorm:"default:0"` // count of revisions
	ConfidenceScore  int `gorm:"default:0"` // confidence score

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // user who learned the word
	Word Word `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"` // word that was learned
}

// TableName returns the table name for LearnedWords
func (LearnedWords) TableName() string {
	return "learned_words"
}
