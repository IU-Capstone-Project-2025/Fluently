package models

import (
	"github.com/google/uuid"
)

// LessonCard is a model for lesson cards
type LessonCard struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	LessonID uuid.UUID `gorm:"type:uuid;not null;index"`
	WordID   uuid.UUID `gorm:"type:uuid;not null"`
	Order    int       `gorm:"not null"`

	Lesson Lesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE"` // lesson cards are part of the lesson
	Word   Word   `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"`   // lesson cards can't be without a word
}

// TableName returns the table name for LessonCard
func (LessonCard) TableName() string {
	return "lesson_cards"
}
