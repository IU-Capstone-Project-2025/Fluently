package models

import (
	"github.com/google/uuid"
)

type LessonCard struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	LessonID uuid.UUID `gorm:"type:uuid;not null;index"`
	WordID   uuid.UUID `gorm:"type:uuid;not null"`
	Order    int       `gorm:"not null"`

	Lesson Lesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE"`
	Word   Word   `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE"`
}

func (LessonCard) TableName() string {
	return "lesson_cards"
}
