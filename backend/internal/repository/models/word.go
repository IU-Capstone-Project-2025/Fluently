package models

import (
	"github.com/google/uuid"
)

type Word struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Word         string    `gorm:"type:varchar(30);not null"`
	Translation  string    `gorm:"type:varchar(30)"`
	PartOfSpeech string    `gorm:"type:varchar(30);not null"`
	Context      string    `gorm:"type:varchar(100)"`
	CEFRLevel    string    `gorm:"type:varchar(2)"`

	TopicID *uuid.UUID `gorm:"type:uuid"` // foreign key to Topic
	Topic   *Topic     `gorm:"foreignKey:TopicID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Word) TableName() string {
	return "words"
}
