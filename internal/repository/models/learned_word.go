package models

import (
	"time"

	"github.com/google/uuid"
)

type LearnedWords struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	WordID uuid.UUID `gorm:"type:uuid;primaryKey"`

	LearnedAt    time.Time `gorm:"not null"`
	LastReviewed time.Time

	CountOfRevisions int `gorm:"default:0"`
	ConfidenceScore  int `gorm:"default:0"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Word Word `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"`
}

func (LearnedWords) TableName() string {
	return "learned_words"
}
