package models

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	StartedAt      time.Time `gorm:"autoCreateTime"`
	WordsPerLesson int       `gorm:"not null"`
	TotalWords     int       `gorm:"not null"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	Cards []LessonCard `gorm:"foreignKey:LessonID"`
}

func (Lesson) TableName() string {
	return "lessons"
}
