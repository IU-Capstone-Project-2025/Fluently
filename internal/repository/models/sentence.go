package models

import (
	"github.com/google/uuid"
)

type Sentence struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Sentence    string    `gorm:"type:text;not null"`
	Translation string    `gorm:"type:text;not null"`
	WordID      uuid.UUID `gorm:"type:uuid;not null;column:word_id"`

	Word Word `gorm:"foreignKey:WordID;references:ID"`
}

func (Sentence) TableName() string {
	return "sentences"
}
