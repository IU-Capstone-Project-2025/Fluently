package models

import (
	"time"

	"github.com/google/uuid"
)

type LearnedWords struct {
	UserID uint      `gorm:"primaryKey;not null;column:user_id"`
	WordID uuid.UUID `gorm:"type:uuid;primaryKey;not null;column:word_id"`

	LearnedAt        time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	LastRevisionAt   *time.Time `gorm:"type:timestamptz;default:now()"`
	CountOfRevisions int        `gorm:"default:0"`
	ConfidenceScore  int        `gorm:"default:0"`

	User User `gorm:"foreignKey:UserID;references:UserID"`
	Word Word `gorm:"foreignKey:WordID;references:ID"`
}

func (LearnedWords) TableName() string {
	return "learned_words"
}
