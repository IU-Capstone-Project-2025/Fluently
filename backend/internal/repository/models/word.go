package models

import (
	"github.com/google/uuid"
)

type Word struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CEFRLevel    string    `gorm:"varchar(2)"`
	Word         string    `gorm:"type:varchar(30);not null"`
	Translation  string    `gorm:"type:varchar(30)"`
	PartOfSpeech string    `gorm:"type:varchar(30);not null"`
	Context      string    `gorm:"type:varchar(100)"`
	AudioURL     string    `gorm:"type:text"`

	TopicID *uuid.UUID `gorm:"type:uuid"` // foreign key to Topic
	Topic   *Topic     `gorm:"foreignKey:TopicID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	Sentences []Sentence `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"`
}

func (Word) TableName() string {
	return "words"
}
