package models

import (
	"time"

	"github.com/google/uuid"
)

type Preference struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CEFRLevel       float64   `gorm:"type:double precision;not null"`
	FactEveryday    bool      `gorm:"default:false"`
	Notifications   bool      `gorm:"default:true"`
	NotificationsAt *time.Time
	WordsPerDay     int    `gorm:"default:10"`
	Goal            string `gorm:"type:varchar(255)"`
	SubLevel        bool   `gorm:"default:false"`
	AvatarImage     []byte `gorm:"type:bytea"`
}

func (Preference) TableName() string {
	return "user_preferences"
}
