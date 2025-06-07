package models

import (
	"github.com/google/uuid"
)

type Word struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Word         string    `gorm:"type:varchar(30);not null"`
	CEFR         string    `gorm:"type:varchar(2);not null"`
	Translation  string    `gorm:"type:varchar(30)"`
	PartOfSpeech string    `gorm:"type:varchar(30);not null"`
	Context      string    `gorm:"type:varchar(100)"`
}

func (Word) TableName() string {
	return "words"
}
