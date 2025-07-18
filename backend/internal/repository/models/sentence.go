package models

import (
	"github.com/google/uuid"
)

// Sentence is a model for sentences
type Sentence struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	WordID      uuid.UUID `gorm:"type:uuid;not null"`
	Sentence    string    `gorm:"type:text;not null"`
	Translation string    `gorm:"type:text"`
	AudioURL    string    `gorm:"type:text"`

	Word Word `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"` // sentence belongs to a word
}

// TableName returns the table name for Sentence
func (Sentence) TableName() string {
	return "sentences"
}
