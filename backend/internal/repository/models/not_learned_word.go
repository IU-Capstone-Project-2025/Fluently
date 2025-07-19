package models

import (
	"github.com/google/uuid"
)

// NotLearnedWords is a model for not learned words
type NotLearnedWords struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	WordID uuid.UUID `gorm:"type:uuid;not null"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // user that has not learned the word
	Word Word `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"` // word that the user has not learned
}

func (NotLearnedWords) TableName() string {
	return "not_learned_words"
}
