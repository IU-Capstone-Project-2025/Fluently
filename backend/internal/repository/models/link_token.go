package models

import (
	"time"

	"github.com/google/uuid"
)

// LinkToken is a model for link tokens
type LinkToken struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Token      string    `gorm:"type:text;not null;unique"`
	TelegramID int64     `gorm:"type:bigint;not null"`
	Used       bool      `gorm:"default:false"`
	ExpiresAt  time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// TableName returns the table name for LinkToken
func (LinkToken) TableName() string {
	return "link_tokens"
}
