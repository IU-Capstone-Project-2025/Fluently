package postgres

import (
	"context"
	"errors"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"
	"gorm.io/gorm"
)

type UserPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserPostgres) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *UserPostgres) Update(ctx context.Context, id uint, updates *schemas.UserUpdateRequest) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("user_id = ?", id).Updates(updates).Error
}

func (r *UserPostgres) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}