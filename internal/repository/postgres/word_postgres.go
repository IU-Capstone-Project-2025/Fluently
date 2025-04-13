package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fluently/go-backend/internal/repository/models"
)

type WordPostgres struct {
	db *gorm.DB
}

func NewWordPostgres(db *gorm.DB) *WordPostgres {
	return &WordPostgres{db: db}
}

func (r *WordPostgres) Create(ctx context.Context, word *models.Word) error {
	word.ID = uuid.New()
	return r.db.WithContext(ctx).Create(word).Error
}
