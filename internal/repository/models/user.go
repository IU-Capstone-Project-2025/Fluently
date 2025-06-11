package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string    `gorm:"type:varchar(100);not null"`
	SubLevel bool      `gorm:"default:false"`
	PrefID   uuid.UUID `gorm:"type:uuid"`

	Pref Preference `gorm:"foreignKey:PrefID;constraint:OnUpdate:CASCADE,OnDelete: SET NULL"`
}

func (User) TableName() string {
	return "users"
}
