package models

import (
	"github.com/google/uuid"
)

type Sentence struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	WordID      uuid.UUID `gorm:"type:uuid;not null"`
	Sentence    string    `gorm:"type:text;not null"`
	Translation string    `gorm:"type:text"`

	Word Word `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"`
}

func (Sentence) TableName() string {
	return "sentences"
}
