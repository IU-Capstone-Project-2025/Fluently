package service

import (
	"context"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"
)

type WordRepository interface {
	Create(ctx context.Context, word *models.Word) error
	GetByID(ctx context.Context, id string) (*models.Word, error)
	List(ctx context.Context) ([]*models.Word, error)
	Update(ctx context.Context, id string, updates *schemas.WordUpdateRequest) error
	Delete(ctx context.Context, id string) error 
}

type WordService struct {
	repo WordRepository
}

func NewWordService(repo WordRepository) *WordService {
	return &WordService{repo: repo}
}

func (s *WordService) Create(ctx context.Context, req *schemas.WordCreateRequest) error {
	word := &models.Word{
		Word: req.Word,
		CEFR: req.CEFR,
		Translation: req.Translation,
		PartOfSpeech: req.PartOfSpeech,
		Context: req.Context,
	}

	return s.repo.Create(ctx, word)
}

func (s *WordService) GetByID(ctx context.Context, id string) (*models.Word, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WordService) List(ctx context.Context) ([]*models.Word, error) {
	return s.repo.List(ctx)
}

func (s *WordService) Update(ctx context.Context, id string, req *schemas.WordUpdateRequest) error {
	if req.Word == nil && req.CEFR == nil && req.Translation == nil &&
		req.PartOfSpeech == nil && req.Context == nil {
		return nil 
	}

	return s.repo.Update(ctx, id, req)
}


func (s *WordService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}