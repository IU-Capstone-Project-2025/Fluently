package models

import (
	"github.com/google/uuid"
)

// Topic is a model for topics
type Topic struct {
	ID       uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title    string     `gorm:"type:varchar(100);not null;unique"`
	ParentID *uuid.UUID `gorm:"type:uuid"`
}

// TableName returns the table name for Topic
func (Topic) TableName() string {
	return "topics"
}
