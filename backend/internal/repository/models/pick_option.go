package models

import (
	"github.com/google/uuid"
)

// PickOption is a model for pick options
type PickOption struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey"`
	WordID     uuid.UUID   `gorm:"type:uuid"`
	SentenceID uuid.UUID   `gorm:"type:uuid"`
	Option     StringArray `gorm:"type:text[]"` // for storing []string like text[] in PostgreSQL
}

// TableName returns the table name for PickOption
func (PickOption) TableName() string {
	return "pick_options"
}
