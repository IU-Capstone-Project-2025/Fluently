package service

import (
	"context"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, id uint, updates *schemas.UserUpdateRequest) error
	Delete(ctx context.Context, id uint) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req *schemas.UserCreateRequest) (*models.User, error) {
	user := &models.User{
		Name: req.Name,
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id uint, updates *schemas.UserUpdateRequest) error {
	if updates.Name != nil {
		return s.repo.Update(ctx, id, updates)
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}