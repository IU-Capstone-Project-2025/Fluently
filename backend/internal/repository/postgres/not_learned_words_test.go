package postgres

import (
	"context"
	"testing"
	"time"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGetRecentlyNotLearnedWords tests the GetRecentlyNotLearnedWords method
func TestGetRecentlyNotLearnedWords(t *testing.T) {
	ctx := context.Background()

	// Create a test user
	user := &models.User{
		ID:        uuid.New(),
		Name:      "Test User",
		Email:     "test@example.com",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	err := userRepo.Create(ctx, user)
	assert.NoError(t, err)

	// Create a test topic
	topic := &models.Topic{
		ID:    uuid.New(),
		Title: "Test Topic",
	}
	err = topicRepo.Create(ctx, topic)
	assert.NoError(t, err)

	// Create test words
	word1 := &models.Word{
		ID:           uuid.New(),
		Word:         "apple",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
		Translation:  "яблоко",
		Context:      "I ate an apple",
		TopicID:      &topic.ID,
	}
	err = wordRepo.Create(ctx, word1)
	assert.NoError(t, err)

	word2 := &models.Word{
		ID:           uuid.New(),
		Word:         "book",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
		Translation:  "книга",
		Context:      "I read a book",
		TopicID:      &topic.ID,
	}
	err = wordRepo.Create(ctx, word2)
	assert.NoError(t, err)

	// Create not learned words for the user
	nlw1 := &models.NotLearnedWords{
		ID:     uuid.New(),
		UserID: user.ID,
		WordID: word1.ID,
	}
	err = notLearnedWordRepo.Create(ctx, nlw1)
	assert.NoError(t, err)

	nlw2 := &models.NotLearnedWords{
		ID:     uuid.New(),
		UserID: user.ID,
		WordID: word2.ID,
	}
	err = notLearnedWordRepo.Create(ctx, nlw2)
	assert.NoError(t, err)

	// Test GetRecentlyNotLearnedWords
	words, err := notLearnedWordRepo.GetRecentlyNotLearnedWords(ctx, user.ID, 10)
	assert.NoError(t, err)
	assert.Len(t, words, 2)

	// Verify the words are returned correctly
	wordMap := make(map[string]models.Word)
	for _, word := range words {
		wordMap[word.Word] = word
	}

	assert.Contains(t, wordMap, "apple")
	assert.Contains(t, wordMap, "book")
	assert.Equal(t, "яблоко", wordMap["apple"].Translation)
	assert.Equal(t, "книга", wordMap["book"].Translation)
}
