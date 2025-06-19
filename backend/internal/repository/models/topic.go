package models

import (
	"github.com/google/uuid"
)

type Topic struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title string    `gorm:"type:varchar(100);not null"`

	Words []Word `gorm:"foreignKey:TopicID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
