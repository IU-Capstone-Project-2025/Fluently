package models

import (
	"time"

	"github.com/google/uuid"
)

type LearnedWords struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	WordID uuid.UUID `gorm:"type:uuid;not null"`

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
