package models

import (
	"github.com/google/uuid"
)

type PickOption struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	WordID     uuid.UUID `gorm:"type:uuid"`
	SentenceID uuid.UUID `gorm:"type:uuid"`
	Option     []string  `gorm:"type:text[]"`
}
