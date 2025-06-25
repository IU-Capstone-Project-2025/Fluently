package models

import (
	"time"

	"github.com/google/uuid"
)

type Preference struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid"`
	CEFRLevel       string    `gorm:"type:varchar(2);not null"`
	FactEveryday    bool      `gorm:"default:false"`
	Notifications   bool      `gorm:"default:true"`
	NotificationsAt *time.Time
	WordsPerDay     int    `gorm:"default:10"`
	Goal            string `gorm:"type:varchar(255)"`
	Subscribed      bool   `gorm:"default:false"`
	AvatarImage     []byte `gorm:"type:blob"`
}

func (Preference) TableName() string {
	return "user_preferences"
}
