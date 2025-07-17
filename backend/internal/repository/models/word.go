package models

import (
	"github.com/google/uuid"
)

// Word is a model for words
type Word struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Word         string    `gorm:"type:varchar(30);not null"`
	Translation  string    `gorm:"type:varchar(255)"`
	PartOfSpeech string    `gorm:"type:varchar(30);not null"`
	Context      string    `gorm:"type:varchar(100)"`
	CEFRLevel    string    `gorm:"type:varchar(2)"`
	AudioURL     string    `gorm:"type:text"`
	Phonetic     string    `gorm:"type:varchar(100)"` // phonetic transcription

	TopicID *uuid.UUID `gorm:"type:uuid"`                                                        // foreign key to Topic
	Topic   *Topic     `gorm:"foreignKey:TopicID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // Topic has many words

	Sentences []Sentence `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"` // Word has many sentences
}

func (Word) TableName() string {
	return "words"
}
