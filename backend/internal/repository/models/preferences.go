package models

import (
	"time"

	"github.com/google/uuid"
)

type Preference struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CEFRLevel       string    `gorm:"type:varchar(2);not null"`
	Points          int       `gorm:"default:0"`
	FactEveryday    bool      `gorm:"default:false"`
	Notifications   bool      `gorm:"default:true"`
	NotificationsAt *time.Time
	WordsPerDay     int    `gorm:"default:10"`
	Goal            string `gorm:"type:varchar(255)"`
}

func (Preference) TableName() string {
	return "user_preferences"
}
