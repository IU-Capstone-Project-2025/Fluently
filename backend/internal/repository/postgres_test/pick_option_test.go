package postgres_test

import (
	"context"
	"testing"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetPickOption(t *testing.T) {
	ctx := context.Background()

	pickOption := &models.PickOption{
		ID:         uuid.New(),
		WordID:     uuid.New(),
		SentenceID: uuid.New(),
		Option:     []string{"Option A", "Option B", "Option C"},
	}

	err := pickOptionRepo.Create(ctx, pickOption)
	assert.NoError(t, err)

	got, err := pickOptionRepo.GetByID(ctx, pickOption.ID)
	assert.NoError(t, err)
	assert.Equal(t, pickOption.ID, got.ID)
	assert.ElementsMatch(t, pickOption.Option, got.Option)
	assert.Equal(t, pickOption.WordID, got.WordID)
	assert.Equal(t, pickOption.SentenceID, got.SentenceID)
}

func TestListPickOptionsByWordID(t *testing.T) {
	ctx := context.Background()
	wordID := uuid.New()

	p1 := &models.PickOption{
		ID:         uuid.New(),
		WordID:     wordID,
		SentenceID: uuid.New(),
		Option:     []string{"A1", "A2", "A3"},
	}
	p2 := &models.PickOption{
		ID:         uuid.New(),
		WordID:     wordID,
		SentenceID: uuid.New(),
		Option:     []string{"B1", "B2", "B3"},
	}

	err := pickOptionRepo.Create(ctx, p1)
	assert.NoError(t, err)
	err = pickOptionRepo.Create(ctx, p2)
	assert.NoError(t, err)

	opts, err := pickOptionRepo.ListByWordID(ctx, wordID)
	assert.NoError(t, err)
	assert.Len(t, opts, 2)
}

func TestUpdatePickOption(t *testing.T) {
	ctx := context.Background()

	p := &models.PickOption{
		ID:         uuid.New(),
		WordID:     uuid.New(),
		SentenceID: uuid.New(),
		Option:     []string{"Old1", "Old2", "Old3"},
	}
	err := pickOptionRepo.Create(ctx, p)
	assert.NoError(t, err)

	p.Option = []string{"New1", "New2", "New3"}
	err = pickOptionRepo.Update(ctx, p)
	assert.NoError(t, err)

	got, err := pickOptionRepo.GetByID(ctx, p.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"New1", "New2", "New3"}, got.Option)
}

func TestDeletePickOption(t *testing.T) {
	ctx := context.Background()

	p := &models.PickOption{
		ID:         uuid.New(),
		WordID:     uuid.New(),
		SentenceID: uuid.New(),
		Option:     []string{"X", "Y", "Z"},
	}
	err := pickOptionRepo.Create(ctx, p)
	assert.NoError(t, err)

	err = pickOptionRepo.Delete(ctx, p.ID)
	assert.NoError(t, err)

	_, err = pickOptionRepo.GetByID(ctx, p.ID)
	assert.Error(t, err)
}
