package postgres

import (
	"context"
	"testing"

	"fluently/go-backend/internal/repository/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCreateAndGetWord tests the creation and retrieval of a word
func TestCreateAndGetWord(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		Word:         "apple",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
		Translation:  "яблоко",
		Context:      "I ate an apple",
		AudioURL:     "",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	found, err := wordRepo.GetByValue(ctx, word.Word)
	assert.NoError(t, err)
	assert.Equal(t, word.Word, found.Word)
	assert.Equal(t, word.Translation, found.Translation)
	assert.Equal(t, word.PartOfSpeech, found.PartOfSpeech)
	assert.Equal(t, word.Context, found.Context)
	assert.Equal(t, word.AudioURL, found.AudioURL)
}

// TestListWord tests the listing of words
func TestListWord(t *testing.T) {
	ctx := context.Background()

	words, err := wordRepo.ListWords(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, words)
}

// TestUpdateWord tests the update of a word
func TestUpdateWord(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		Word:         "banana",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
		Translation:  "банан",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	word.Translation = "Новый перевод"
	word.AudioURL = "Rickroll URL"

	err = wordRepo.Update(ctx, word)
	assert.NoError(t, err)

	updated, err := wordRepo.GetByValue(ctx, word.Word)
	assert.NoError(t, err)
	assert.Equal(t, "Новый перевод", updated.Translation)
	assert.Equal(t, "Rickroll URL", updated.AudioURL)
}

// DeleteWord tests the deletion of a word
func DeleteWord(t *testing.T) {
	ctx := context.Background()
	word := &models.Word{
		Word:         "orange",
		CEFRLevel:    "A1",
		PartOfSpeech: "noun",
	}

	err := wordRepo.Create(ctx, word)
	assert.NoError(t, err)

	// Get the word by value to get its ID
	found, err := wordRepo.GetByValue(ctx, word.Word)
	assert.NoError(t, err)

	err = wordRepo.Delete(ctx, found.ID)
	assert.NoError(t, err)

	_, err = wordRepo.GetByID(ctx, found.ID)
	assert.Error(t, err) // Not Found
}

// TestGetRandomWordsWithTopic tests the GetRandomWordsWithTopic method
func TestGetRandomWordsWithTopic(t *testing.T) {
	ctx := context.Background()

	// Create a test topic
	topic := &models.Topic{
		ID:    uuid.New(),
		Title: "Test Topic",
	}
	err := topicRepo.Create(ctx, topic)
	assert.NoError(t, err)

	// Create test words with topic
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

	// Test GetRandomWordsWithTopic
	words, err := wordRepo.GetRandomWordsWithTopic(ctx, 10)
	assert.NoError(t, err)
	assert.Len(t, words, 2)

	// Verify the words are returned with topic information
	wordMap := make(map[string]models.Word)
	for _, word := range words {
		wordMap[word.Word] = word
	}

	assert.Contains(t, wordMap, "apple")
	assert.Contains(t, wordMap, "book")
	assert.NotNil(t, wordMap["apple"].Topic)
	assert.NotNil(t, wordMap["book"].Topic)
	assert.Equal(t, "Test Topic", wordMap["apple"].Topic.Title)
	assert.Equal(t, "Test Topic", wordMap["book"].Topic.Title)
}
