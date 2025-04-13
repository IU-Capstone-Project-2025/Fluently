package service

import (
	"context"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/schemas"
)

type WordRepository interface {
	Create(ctx context.Context, word *models.Word) error
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
		CEFR: (*models.CEFRLevel)(req.CEFR),
		Translation: req.Translation,
		PartOfSpeech: req.PartOfSpeech,
		Context: req.Context,
	}

	return s.repo.Create(ctx, word)
}
