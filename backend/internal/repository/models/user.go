package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GoogleID     string     `gorm:"type:varchar(100);not null"`
	Provider     string     `gorm:"type:varchar(50)"`
	Name         string     `gorm:"type:varchar(100);not null"`
	Role         string     `gorm:"type:varchar(10);default:'user'"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex"`
	PasswordHash string     `gorm:"type:text"`
	PrefID       *uuid.UUID `gorm:"type:uuid"`
	IsActive     bool       `gorm:"default:true"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`

	Pref *Preference `gorm:"foreignKey:PrefID;constraint:OnUpdate:CASCADE,OnDelete: SET NULL"`
}

func (User) TableName() string {
	return "users"
}
