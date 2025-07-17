package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ChatHistory stores completed dialogue between user and AI.
// Messages field holds full array JSON (see API contract).
type ChatHistory struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Messages   datatypes.JSON `gorm:"type:jsonb;not null"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	FinishedAt time.Time      `gorm:"autoUpdateTime"`
}

func (ChatHistory) TableName() string { return "chat_histories" }
