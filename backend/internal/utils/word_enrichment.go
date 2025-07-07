package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// WordEnrichmentService handles enriching words with dictionary data and sentences with distractor options
type WordEnrichmentService struct {
	dictionaryClient *DictionaryClient
	distractorClient *DistractorClient
	db               *gorm.DB
}

// NewWordEnrichmentService creates a new word enrichment service
func NewWordEnrichmentService(db *gorm.DB) *WordEnrichmentService {
	return &WordEnrichmentService{
		dictionaryClient: NewDictionaryClient(DictionaryClientConfig{
			Timeout: 15 * time.Second, // Longer timeout for batch operations
		}),
		distractorClient: NewDistractorClient(DistractorClientConfig{
			Timeout: 15 * time.Second, // Longer timeout for batch operations
		}),
		db: db,
	}
}

// EnrichWordWithDictionary fetches dictionary data and updates the word model
func (s *WordEnrichmentService) EnrichWordWithDictionary(ctx context.Context, word *models.Word) error {
	if word.Word == "" {
		return fmt.Errorf("word text cannot be empty")
	}

	logger.Log.Debug("Enriching word with dictionary data",
		zap.String("word", word.Word),
		zap.String("word_id", word.ID.String()))

	// Fetch dictionary information
	wordInfo, err := s.dictionaryClient.GetWordInfo(ctx, word.Word)
	if err != nil {
		// Log the error but don't fail the operation
		logger.Log.Warn("Failed to fetch dictionary data",
			zap.String("word", word.Word),
			zap.Error(err))
		return err
	}

	// Update word fields with dictionary data
	updated := false

	// Update part of speech if not already set or if it's "unknown"
	if word.PartOfSpeech == "" || word.PartOfSpeech == "unknown" {
		if wordInfo.PartOfSpeech != "" {
			word.PartOfSpeech = wordInfo.PartOfSpeech
			updated = true
			logger.Log.Debug("Updated part of speech",
				zap.String("word", word.Word),
				zap.String("part_of_speech", wordInfo.PartOfSpeech))
		}
	}

	// Update audio URL if not already set
	if word.AudioURL == "" && wordInfo.AudioURL != "" {
		word.AudioURL = wordInfo.AudioURL
		updated = true
		logger.Log.Debug("Updated audio URL",
			zap.String("word", word.Word),
			zap.String("audio_url", wordInfo.AudioURL))
	}

	// Update phonetic transcription if not already set
	if word.Phonetic == "" && wordInfo.Phonetic != "" {
		word.Phonetic = wordInfo.Phonetic
		updated = true
		logger.Log.Debug("Updated phonetic transcription",
			zap.String("word", word.Word),
			zap.String("phonetic", wordInfo.Phonetic))
	}

	// Update context with definition if context is empty and we have a definition
	if word.Context == "" && wordInfo.Definition != "" {
		// Truncate definition to fit in context field (varchar(100))
		definition := wordInfo.Definition
		if len(definition) > 100 {
			definition = definition[:97] + "..."
		}
		word.Context = definition
		updated = true
		logger.Log.Debug("Updated context with definition",
			zap.String("word", word.Word),
			zap.String("context", definition))
	}

	if updated {
		logger.Log.Info("Word enriched with dictionary data",
			zap.String("word", word.Word),
			zap.String("part_of_speech", word.PartOfSpeech),
			zap.String("phonetic", word.Phonetic),
			zap.Bool("has_audio", word.AudioURL != ""))
	}

	return nil
}

// EnrichWordInDatabase fetches a word from database, enriches it, and saves it back
func (s *WordEnrichmentService) EnrichWordInDatabase(ctx context.Context, wordID string) error {
	var word models.Word
	if err := s.db.WithContext(ctx).First(&word, "id = ?", wordID).Error; err != nil {
		return fmt.Errorf("failed to find word: %w", err)
	}

	if err := s.EnrichWordWithDictionary(ctx, &word); err != nil {
		return fmt.Errorf("failed to enrich word: %w", err)
	}

	if err := s.db.WithContext(ctx).Save(&word).Error; err != nil {
		return fmt.Errorf("failed to save enriched word: %w", err)
	}

	return nil
}

// BatchEnrichWords enriches multiple words with dictionary data
func (s *WordEnrichmentService) BatchEnrichWords(ctx context.Context, words []*models.Word, rateLimitDelay time.Duration) []error {
	var errors []error

	for i, word := range words {
		// Add rate limiting to avoid overwhelming the API
		if i > 0 && rateLimitDelay > 0 {
			select {
			case <-ctx.Done():
				errors = append(errors, ctx.Err())
				return errors
			case <-time.After(rateLimitDelay):
				// Continue after delay
			}
		}

		if err := s.EnrichWordWithDictionary(ctx, word); err != nil {
			errors = append(errors, fmt.Errorf("failed to enrich word '%s': %w", word.Word, err))
			logger.Log.Error("Failed to enrich word in batch",
				zap.String("word", word.Word),
				zap.Error(err))
		}
	}

	return errors
}

// EnrichWordsInDatabase finds words that need enrichment and processes them
func (s *WordEnrichmentService) EnrichWordsInDatabase(ctx context.Context, limit int, rateLimitDelay time.Duration) error {
	// Find words that need enrichment (missing phonetic, audio, or have "unknown" part of speech)
	var words []models.Word
	err := s.db.WithContext(ctx).
		Where("phonetic = ? OR audio_url = ? OR part_of_speech = ?", "", "", "unknown").
		Limit(limit).
		Find(&words).Error

	if err != nil {
		return fmt.Errorf("failed to find words for enrichment: %w", err)
	}

	if len(words) == 0 {
		logger.Log.Info("No words found that need enrichment")
		return nil
	}

	logger.Log.Info("Starting word enrichment process",
		zap.Int("word_count", len(words)),
		zap.Duration("rate_limit_delay", rateLimitDelay))

	// Convert to pointer slice for BatchEnrichWords
	wordPtrs := make([]*models.Word, len(words))
	for i := range words {
		wordPtrs[i] = &words[i]
	}

	// Enrich words
	errors := s.BatchEnrichWords(ctx, wordPtrs, rateLimitDelay)

	// Save all enriched words in a transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, word := range words {
			if err := tx.Save(word).Error; err != nil {
				return fmt.Errorf("failed to save word %s: %w", word.Word, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to save enriched words: %w", err)
	}

	logger.Log.Info("Word enrichment process completed",
		zap.Int("total_words", len(words)),
		zap.Int("errors", len(errors)))

	// Log errors but don't fail the operation
	for _, enrichErr := range errors {
		logger.Log.Warn("Word enrichment error", zap.Error(enrichErr))
	}

	return nil
}

// ValidatePartOfSpeech checks if a part of speech is valid according to the dictionary
func (s *WordEnrichmentService) ValidatePartOfSpeech(ctx context.Context, word, partOfSpeech string) (bool, error) {
	partsOfSpeech, err := s.dictionaryClient.GetAllPartsOfSpeech(ctx, word)
	if err != nil {
		return false, err
	}

	for _, pos := range partsOfSpeech {
		if strings.EqualFold(pos, partOfSpeech) {
			return true, nil
		}
	}

	return false, nil
}

// GetSuggestedPartOfSpeech returns the most common part of speech for a word
func (s *WordEnrichmentService) GetSuggestedPartOfSpeech(ctx context.Context, word string) (string, error) {
	wordInfo, err := s.dictionaryClient.GetWordInfo(ctx, word)
	if err != nil {
		return "", err
	}

	return wordInfo.PartOfSpeech, nil
}

// EnrichSentenceWithDistractors generates distractor options for a sentence using the distractor API
func (s *WordEnrichmentService) EnrichSentenceWithDistractors(ctx context.Context, sentence *models.Sentence) error {
	if sentence.WordID == uuid.Nil || sentence.Sentence == "" {
		return fmt.Errorf("sentence must have valid WordID and text")
	}

	// Get the word associated with this sentence
	var word models.Word
	if err := s.db.WithContext(ctx).First(&word, "id = ?", sentence.WordID).Error; err != nil {
		return fmt.Errorf("failed to find word for sentence: %w", err)
	}

	logger.Log.Debug("Enriching sentence with distractor options",
		zap.String("sentence_id", sentence.ID.String()),
		zap.String("word", word.Word),
		zap.String("sentence", sentence.Sentence))

	// Check if pick options already exist for this sentence
	var existingOption models.PickOption
	err := s.db.WithContext(ctx).Where("word_id = ? AND sentence_id = ?", sentence.WordID, sentence.ID).First(&existingOption).Error
	if err == nil {
		logger.Log.Debug("Pick options already exist for sentence", zap.String("sentence_id", sentence.ID.String()))
		return nil // Already has options
	} else if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing pick options: %w", err)
	}

	// Generate distractor options using the distractor API
	distractorOptions, err := s.distractorClient.GenerateDistractors(ctx, sentence.Sentence, word.Word)
	if err != nil {
		logger.Log.Warn("Failed to generate distractor options",
			zap.String("word", word.Word),
			zap.String("sentence", sentence.Sentence),
			zap.Error(err))
		return err
	}

	if len(distractorOptions) == 0 {
		logger.Log.Warn("No distractor options generated",
			zap.String("word", word.Word),
			zap.String("sentence", sentence.Sentence))
		return fmt.Errorf("no distractor options generated for word '%s' in sentence", word.Word)
	}

	// Create pick option record
	pickOption := models.PickOption{
		ID:         uuid.New(),
		WordID:     sentence.WordID,
		SentenceID: sentence.ID,
		Option:     models.StringArray(distractorOptions),
	}

	if err := s.db.WithContext(ctx).Create(&pickOption).Error; err != nil {
		return fmt.Errorf("failed to save pick options: %w", err)
	}

	logger.Log.Info("Successfully created distractor options",
		zap.String("word", word.Word),
		zap.String("sentence", sentence.Sentence),
		zap.Int("option_count", len(distractorOptions)),
		zap.Strings("options", distractorOptions))

	return nil
}

// BatchEnrichSentencesWithDistractors enriches multiple sentences with distractor options
func (s *WordEnrichmentService) BatchEnrichSentencesWithDistractors(ctx context.Context, sentences []*models.Sentence, rateLimitDelay time.Duration) []error {
	var errors []error

	for i, sentence := range sentences {
		// Add rate limiting to avoid overwhelming the API
		if i > 0 && rateLimitDelay > 0 {
			select {
			case <-ctx.Done():
				errors = append(errors, ctx.Err())
				return errors
			case <-time.After(rateLimitDelay):
				// Continue after delay
			}
		}

		if err := s.EnrichSentenceWithDistractors(ctx, sentence); err != nil {
			errors = append(errors, fmt.Errorf("failed to enrich sentence '%s': %w", sentence.Sentence, err))
			logger.Log.Error("Failed to enrich sentence in batch",
				zap.String("sentence_id", sentence.ID.String()),
				zap.String("sentence", sentence.Sentence),
				zap.Error(err))
		}
	}

	return errors
}

// EnrichSentencesInDatabase finds sentences without distractor options and processes them
func (s *WordEnrichmentService) EnrichSentencesInDatabase(ctx context.Context, limit int, rateLimitDelay time.Duration) error {
	// Find sentences that don't have pick options yet
	var sentences []models.Sentence
	err := s.db.WithContext(ctx).Raw(`
		SELECT s.* FROM sentences s 
		LEFT JOIN pick_options po ON s.id = po.sentence_id 
		WHERE po.sentence_id IS NULL 
		AND s.sentence != '' 
		LIMIT ?
	`, limit).Scan(&sentences).Error

	if err != nil {
		return fmt.Errorf("failed to find sentences for enrichment: %w", err)
	}

	if len(sentences) == 0 {
		logger.Log.Info("No sentences found that need distractor options")
		return nil
	}

	logger.Log.Info("Starting sentence enrichment process",
		zap.Int("sentence_count", len(sentences)),
		zap.Duration("rate_limit_delay", rateLimitDelay))

	// Convert to pointer slice for BatchEnrichSentencesWithDistractors
	sentencePtrs := make([]*models.Sentence, len(sentences))
	for i := range sentences {
		sentencePtrs[i] = &sentences[i]
	}

	// Enrich sentences
	errors := s.BatchEnrichSentencesWithDistractors(ctx, sentencePtrs, rateLimitDelay)

	logger.Log.Info("Sentence enrichment process completed",
		zap.Int("total_sentences", len(sentences)),
		zap.Int("errors", len(errors)))

	// Log errors but don't fail the operation
	for _, enrichErr := range errors {
		logger.Log.Warn("Sentence enrichment error", zap.Error(enrichErr))
	}

	return nil
}

// EnrichSentenceInDatabase fetches a sentence from database, enriches it, and saves the options
func (s *WordEnrichmentService) EnrichSentenceInDatabase(ctx context.Context, sentenceID string) error {
	var sentence models.Sentence
	if err := s.db.WithContext(ctx).First(&sentence, "id = ?", sentenceID).Error; err != nil {
		return fmt.Errorf("failed to find sentence: %w", err)
	}

	return s.EnrichSentenceWithDistractors(ctx, &sentence)
}
