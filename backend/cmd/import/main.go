package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"fluently/go-backend/internal/config"
	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/utils"
	"fluently/go-backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	csvFilePath string
	batchSize   = 100
)

type ImportStats struct {
	WordsImported     int
	TopicsCreated     int
	SentencesAdded    int
	ErrorsEncountered int
	StartTime         time.Time
	TotalRows         int
}

type ClearStats struct {
	LearnedWordsDeleted int
	SentencesDeleted    int
	WordsDeleted        int
	TopicsDeleted       int
	ErrorsEncountered   int
	StartTime           time.Time
}

type CSVRecord struct {
	Index       int
	Total       float64
	Word        string
	Topic       string
	Subtopic    string
	Subsubtopic string
	CEFRLevel   string
	Translation string
	Sentences   string
}

type SentencePair []string

var rootCmd = &cobra.Command{
	Use:   "import",
	Short: "A CLI tool for importing data",
	Long:  `A Command Line Interface (CLI) tool to import data into the backend system.`,
}

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Import data from CSV file",
	Long:  `Import words, topics, and sentences from a CSV file into the database.`,
	RunE:  runCSVImport,
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all words, sentences, topics and learned words from database",
	Long:  `Clear all words, sentences, topics and learned words from the database. Requires confirmation password from CLEAR_PASSWORD environment variable.`,
	RunE:  runClearData,
}

var enrichCmd = &cobra.Command{
	Use:   "enrich",
	Short: "Enrich existing words with dictionary data",
	Long:  `Enrich existing words in the database with phonetic transcription, audio URLs, and part of speech from the dictionary API.`,
	RunE:  runEnrichWords,
}

var enrichSentencesCmd = &cobra.Command{
	Use:   "enrich-sentences",
	Short: "Enrich existing sentences with distractor options",
	Long:  `Enrich existing sentences in the database with distractor options using the ML distractor service.`,
	RunE:  runEnrichSentences,
}

var resetEnrichmentCmd = &cobra.Command{
	Use:   "reset-enrichment",
	Short: "Reset failed word enrichments to allow retry",
	Long:  `Reset words marked as "unknown_processed" back to "unknown" to allow them to be retried for enrichment.`,
	RunE:  runResetEnrichment,
}

var (
	enrichLimit   int
	enrichDelay   int
	sentenceLimit int
	sentenceDelay int
)

func init() {
	csvCmd.Flags().StringVarP(&csvFilePath, "file", "f", "", "Path to CSV file (required)")
	csvCmd.MarkFlagRequired("file")

	enrichCmd.Flags().IntVarP(&enrichLimit, "limit", "l", 100000, "Maximum number of words to enrich in one run")
	enrichCmd.Flags().IntVarP(&enrichDelay, "delay", "d", 500, "Delay between API calls in milliseconds")

	enrichSentencesCmd.Flags().IntVarP(&sentenceLimit, "limit", "l", 1000000, "Maximum number of sentences to enrich in one run")
	enrichSentencesCmd.Flags().IntVarP(&sentenceDelay, "delay", "d", 10, "Delay between API calls in milliseconds")

	rootCmd.AddCommand(csvCmd)
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(enrichCmd)
	rootCmd.AddCommand(enrichSentencesCmd)
	rootCmd.AddCommand(resetEnrichmentCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Error("Failed to execute root command", zap.Error(err))
		os.Exit(1)
	}
}

func runCSVImport(cmd *cobra.Command, args []string) error {
	// Initialize config and logger
	config.Init()
	logger.Init(true) // Enable debug logging
	defer logger.Log.Sync()

	// Validate file exists
	if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
		return fmt.Errorf("CSV file does not exist: %s", csvFilePath)
	}

	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	logger.Log.Info("Starting CSV import", zap.String("file", csvFilePath))

	// Count total rows for progress bar
	totalRows, err := countCSVRows(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to count CSV rows: %v", err)
	}

	stats := &ImportStats{
		StartTime: time.Now(),
		TotalRows: totalRows,
	}

	// Create progress bar
	bar := progressbar.NewOptions(totalRows,
		progressbar.OptionSetDescription("Importing CSV data..."),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("rows"),
		progressbar.OptionOnCompletion(func() {
			fmt.Println("\nImport completed!")
		}),
	)

	// Process CSV file
	err = processCSVFile(db, csvFilePath, bar, stats)
	if err != nil {
		return fmt.Errorf("failed to process CSV: %v", err)
	}

	// Print final statistics
	printStats(stats)
	return nil
}

func connectToDatabase() (*gorm.DB, error) {
	dsn := config.GetPostgresDSNForImport()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate to ensure CEFR_Level column exists
	err = db.AutoMigrate(&models.Word{})
	if err != nil {
		logger.Log.Warn("Failed to auto-migrate Word model", zap.Error(err))
	}

	return db, nil
}

func countCSVRows(filepath string) (int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	count := 0
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}

	// Subtract 1 for header row
	if count > 0 {
		count--
	}
	return count, nil
}

func processCSVFile(db *gorm.DB, filepath string, bar *progressbar.ProgressBar, stats *ImportStats) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %v", err)
	}

	// Validate headers
	expectedHeaders := []string{"", "Total", "word", "topic", "subtopic", "subsubtopic", "CEFR_level", "translation", "sentences"}
	if len(headers) != len(expectedHeaders) {
		return fmt.Errorf("invalid CSV format: expected %d columns, got %d", len(expectedHeaders), len(headers))
	}

	var batch []CSVRecord
	rowNum := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Error("Failed to read CSV row", zap.Error(err), zap.Int("row", rowNum+2))
			stats.ErrorsEncountered++
			continue
		}

		rowNum++

		if len(record) != len(expectedHeaders) {
			logger.Log.Error("Invalid CSV row format", zap.Int("row", rowNum+1), zap.Int("expected_cols", len(expectedHeaders)), zap.Int("actual_cols", len(record)))
			stats.ErrorsEncountered++
			bar.Add(1)
			continue
		}

		// Parse total field
		total := 0.0
		if record[1] != "" {
			if t, err := strconv.ParseFloat(record[1], 64); err == nil {
				total = t
			}
		}

		csvRecord := CSVRecord{
			Index:       rowNum,
			Total:       total,
			Word:        strings.TrimSpace(record[2]),
			Topic:       strings.TrimSpace(record[3]),
			Subtopic:    strings.TrimSpace(record[4]),
			Subsubtopic: strings.TrimSpace(record[5]),
			CEFRLevel:   strings.TrimSpace(record[6]),
			Translation: strings.TrimSpace(record[7]),
			Sentences:   normalizeJSONQuotes(strings.TrimSpace(record[8])),
		}

		// Skip empty words
		if csvRecord.Word == "" {
			logger.Log.Warn("Skipping row with empty word", zap.Int("row", rowNum+1))
			bar.Add(1)
			continue
		}

		batch = append(batch, csvRecord)

		// Process batch when it reaches batchSize
		if len(batch) >= batchSize {
			err := processBatch(db, batch, stats)
			if err != nil {
				logger.Log.Error("Failed to process batch", zap.Error(err))
				stats.ErrorsEncountered += len(batch)
			}
			bar.Add(len(batch))
			batch = nil
		}
	}

	// Process remaining records
	if len(batch) > 0 {
		err := processBatch(db, batch, stats)
		if err != nil {
			logger.Log.Error("Failed to process final batch", zap.Error(err))
			stats.ErrorsEncountered += len(batch)
		}
		bar.Add(len(batch))
	}

	return nil
}

func processBatch(db *gorm.DB, batch []CSVRecord, stats *ImportStats) error {
	ctx := context.Background()

	// Process each record individually instead of grouping by word
	for _, record := range batch {
		err := db.Transaction(func(tx *gorm.DB) error {
			return processWordRecord(tx, ctx, record, stats)
		})

		if err != nil {
			logger.Log.Error("Failed to process word record",
				zap.Error(err),
				zap.String("word", record.Word),
				zap.String("translation", record.Translation),
				zap.String("topic", record.Topic),
				zap.String("subtopic", record.Subtopic),
				zap.String("subsubtopic", record.Subsubtopic))
			stats.ErrorsEncountered++
			continue
		}
	}

	return nil
}

func processWordRecord(tx *gorm.DB, ctx context.Context, record CSVRecord, stats *ImportStats) error {
	// Log detailed information about what we're processing
	logger.Log.Debug("Processing word record",
		zap.String("word", record.Word),
		zap.String("translation", record.Translation),
		zap.String("topic", record.Topic),
		zap.String("subtopic", record.Subtopic),
		zap.String("subsubtopic", record.Subsubtopic),
		zap.String("cefr_level", record.CEFRLevel))

	// Create topic hierarchy
	topicID, err := createTopicHierarchy(tx, ctx, record.Topic, record.Subtopic, record.Subsubtopic, stats)
	if err != nil {
		logger.Log.Error("Failed to create topic hierarchy",
			zap.Error(err),
			zap.String("word", record.Word),
			zap.String("translation", record.Translation),
			zap.String("topic", record.Topic),
			zap.String("subtopic", record.Subtopic),
			zap.String("subsubtopic", record.Subsubtopic))
		return fmt.Errorf("failed to create topic hierarchy: %v", err)
	}

	if topicID != nil {
		logger.Log.Debug("Topic hierarchy created",
			zap.String("topicID", topicID.String()),
			zap.String("word", record.Word),
			zap.String("translation", record.Translation))
	} else {
		logger.Log.Debug("No topic hierarchy created (all empty)",
			zap.String("word", record.Word),
			zap.String("translation", record.Translation))
	}

	// Use the translation from this specific record
	translation := record.Translation

	// Truncate translation if too long (varchar(255) limit)
	if len(translation) > 255 {
		translation = translation[:255]
		logger.Log.Warn("Translation truncated due to length limit",
			zap.String("word", record.Word),
			zap.String("original_translation", record.Translation),
			zap.String("truncated_translation", translation))
	}

	// Check word length (varchar(255) limit)
	wordText := record.Word
	if len(wordText) > 255 {
		wordText = wordText[:255]
		logger.Log.Warn("Word truncated due to length limit",
			zap.String("original_word", record.Word),
			zap.String("truncated_word", wordText))
	}

	// Check if word with this translation already exists
	// Uniqueness is now based on (word, translation, topic_id) combination
	var existingWord models.Word
	query := tx.WithContext(ctx).Where("word = ? AND translation = ?", wordText, translation)
	if topicID != nil {
		query = query.Where("topic_id = ?", *topicID)
	} else {
		query = query.Where("topic_id IS NULL")
	}
	result := query.First(&existingWord)

	var wordID uuid.UUID
	if result.Error == gorm.ErrRecordNotFound {
		// Create new word - each (word, translation) pair gets its own entity
		newWord := models.Word{
			Word:         wordText,
			Translation:  translation,
			PartOfSpeech: "unknown",
			Context:      "",
			CEFRLevel:    record.CEFRLevel,
			TopicID:      topicID,
		}

		logger.Log.Debug("Creating new word",
			zap.String("word", wordText),
			zap.String("translation", translation),
			zap.String("cefr_level", record.CEFRLevel))

		if err := tx.WithContext(ctx).Create(&newWord).Error; err != nil {
			topicIDStr := "NULL"
			if topicID != nil {
				topicIDStr = topicID.String()
			}
			logger.Log.Error("Failed to create word",
				zap.Error(err),
				zap.String("word", wordText),
				zap.String("translation", translation),
				zap.String("topic_id", topicIDStr))
			return fmt.Errorf("failed to create word: %v", err)
		}
		wordID = newWord.ID
		stats.WordsImported++

		logger.Log.Debug("Successfully created word",
			zap.String("id", newWord.ID.String()),
			zap.String("word", wordText),
			zap.String("translation", translation))
	} else if result.Error != nil {
		return fmt.Errorf("failed to check existing word: %v", result.Error)
	} else {
		// Update existing word with CEFR level if provided
		if record.CEFRLevel != "" {
			existingWord.CEFRLevel = record.CEFRLevel
		}

		if err := tx.WithContext(ctx).Save(&existingWord).Error; err != nil {
			return fmt.Errorf("failed to update word: %v", err)
		}
		wordID = existingWord.ID
		logger.Log.Info("Updated existing word",
			zap.String("word", record.Word),
			zap.String("translation", record.Translation))
	}

	// Process sentences for this record
	if record.Sentences != "" {
		sentencesAdded, err := processSentences(tx, ctx, wordID, record.Sentences)
		if err != nil {
			logger.Log.Error("Failed to process sentences",
				zap.Error(err),
				zap.String("word", record.Word),
				zap.String("translation", record.Translation))
			return err
		}
		stats.SentencesAdded += sentencesAdded
	}

	return nil
}

func createTopicHierarchy(tx *gorm.DB, ctx context.Context, topicName, subtopicName, subsubtopicName string, stats *ImportStats) (*uuid.UUID, error) {
	var parentID *uuid.UUID

	// Create main topic
	if topicName != "" {
		topicID, created, err := getOrCreateTopic(tx, ctx, topicName, nil)
		if err != nil {
			return nil, err
		}
		if created {
			stats.TopicsCreated++
		}
		parentID = &topicID
	}

	// Create subtopic
	if subtopicName != "" {
		subtopicID, created, err := getOrCreateTopic(tx, ctx, subtopicName, parentID)
		if err != nil {
			return nil, err
		}
		if created {
			stats.TopicsCreated++
		}
		parentID = &subtopicID
	}

	// Create subsubtopic
	if subsubtopicName != "" {
		subsubtopicID, created, err := getOrCreateTopic(tx, ctx, subsubtopicName, parentID)
		if err != nil {
			return nil, err
		}
		if created {
			stats.TopicsCreated++
		}
		parentID = &subsubtopicID
	}

	return parentID, nil
}

func getOrCreateTopic(tx *gorm.DB, ctx context.Context, title string, parentID *uuid.UUID) (uuid.UUID, bool, error) {
	var topic models.Topic

	// Log what we're trying to find/create
	parentIDStr := "NULL"
	if parentID != nil {
		parentIDStr = parentID.String()
	}
	logger.Log.Debug("Getting or creating topic",
		zap.String("title", title),
		zap.String("parent_id", parentIDStr))

	query := tx.WithContext(ctx).Where("title = ?", title)
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	result := query.First(&topic)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new topic
		newTopic := models.Topic{
			Title:    title,
			ParentID: parentID,
		}

		logger.Log.Debug("Creating new topic",
			zap.String("title", title),
			zap.String("parent_id", parentIDStr))

		if err := tx.WithContext(ctx).Create(&newTopic).Error; err != nil {
			logger.Log.Error("Failed to create topic",
				zap.Error(err),
				zap.String("title", title),
				zap.String("parent_id", parentIDStr))
			return uuid.Nil, false, fmt.Errorf("failed to create topic: %v", err)
		}

		logger.Log.Debug("Successfully created topic",
			zap.String("id", newTopic.ID.String()),
			zap.String("title", title))

		return newTopic.ID, true, nil
	} else if result.Error != nil {
		logger.Log.Error("Database error when checking topic",
			zap.Error(result.Error),
			zap.String("title", title),
			zap.String("parent_id", parentIDStr))
		return uuid.Nil, false, fmt.Errorf("failed to check existing topic: %v", result.Error)
	}

	logger.Log.Debug("Found existing topic",
		zap.String("id", topic.ID.String()),
		zap.String("title", title))

	return topic.ID, false, nil
}

func processSentences(tx *gorm.DB, ctx context.Context, wordID uuid.UUID, sentencesJSON string) (int, error) {
	if sentencesJSON == "" {
		return 0, nil
	}

	// Log the JSON we're trying to parse for debugging
	logger.Log.Debug("Processing sentences JSON",
		zap.String("word_id", wordID.String()),
		zap.String("sentences_json", sentencesJSON))

	// Parse JSON sentences
	var sentencePairs []SentencePair
	if err := json.Unmarshal([]byte(sentencesJSON), &sentencePairs); err != nil {
		logger.Log.Error("Failed to parse sentences JSON",
			zap.Error(err),
			zap.String("word_id", wordID.String()),
			zap.String("sentences_json", sentencesJSON))
		return 0, fmt.Errorf("failed to parse sentences JSON: %v (JSON: %s)", err, sentencesJSON)
	}

	sentencesAdded := 0
	for _, pair := range sentencePairs {
		if len(pair) >= 2 {
			sentence := models.Sentence{
				WordID:      wordID,
				Sentence:    pair[0],
				Translation: pair[1],
			}

			// Check if sentence already exists
			var existingSentence models.Sentence
			result := tx.WithContext(ctx).Where("word_id = ? AND sentence = ?", wordID, pair[0]).First(&existingSentence)

			if result.Error == gorm.ErrRecordNotFound {
				if err := tx.WithContext(ctx).Create(&sentence).Error; err != nil {
					logger.Log.Error("Failed to create sentence", zap.Error(err))
					continue
				}
				sentencesAdded++
			} else if result.Error != nil {
				logger.Log.Error("Failed to check existing sentence", zap.Error(result.Error))
				continue
			}
			// If sentence exists, skip it (no update needed)
		}
	}

	return sentencesAdded, nil
}

// normalizeJSONQuotes converts mixed quote JSON to valid JSON format
// Handles cases like: [["text", 'текст'], ['text2', 'текст2']]
// Converts to: [["text", "текст"], ["text2", "текст2"]]
func normalizeJSONQuotes(jsonStr string) string {
	if jsonStr == "" {
		return ""
	}

	// Log original for debugging
	logger.Log.Debug("Normalizing JSON quotes", zap.String("original", jsonStr))

	result := make([]rune, 0, len(jsonStr))
	runes := []rune(jsonStr)
	inString := false
	currentQuote := rune(0)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Handle escaped characters first
		if r == '\\' && inString && i+1 < len(runes) {
			result = append(result, r)
			i++
			nextChar := runes[i]
			result = append(result, nextChar)
			continue
		}

		switch r {
		case '"', '\'':
			if !inString {
				// Starting a string - normalize to double quote
				result = append(result, '"')
				inString = true
				currentQuote = r
			} else if r == currentQuote {
				// Ending the current string - use double quote
				result = append(result, '"')
				inString = false
				currentQuote = 0
			} else {
				// Different quote inside string - escape if it's a double quote
				if r == '"' {
					result = append(result, '\\', '"')
				} else {
					// Single quote inside double-quoted string - keep as is
					result = append(result, r)
				}
			}
		default:
			result = append(result, r)
		}
	}

	normalized := string(result)
	logger.Log.Debug("Normalized JSON quotes", zap.String("normalized", normalized))

	return normalized
}

func printStats(stats *ImportStats) {
	duration := time.Since(stats.StartTime)
	speed := float64(stats.TotalRows) / duration.Seconds()

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 IMPORT STATISTICS")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("⏱️  Duration: %v\n", duration.Round(time.Second))
	fmt.Printf("🚀 Speed: %.1f rows/second\n", speed)
	fmt.Printf("📄 Total rows processed: %d\n", stats.TotalRows)
	fmt.Printf("📝 Words imported: %d\n", stats.WordsImported)
	fmt.Printf("🏷️  Topics created: %d\n", stats.TopicsCreated)
	fmt.Printf("💬 Sentences added: %d\n", stats.SentencesAdded)
	fmt.Printf("❌ Errors encountered: %d\n", stats.ErrorsEncountered)
	fmt.Println(strings.Repeat("=", 50))
}

func runClearData(cmd *cobra.Command, args []string) error {
	// Initialize config and logger first
	config.Init()
	logger.Init(true) // Enable debug logging
	defer logger.Log.Sync()

	// Load .env from project root (one level up from backend)
	err := godotenv.Load("../.env")
	if err != nil {
		// Try loading from current directory as fallback
		err = godotenv.Load()
		if err != nil {
			logger.Log.Warn("No .env file found, reading from environment variables", zap.Error(err))
		}
	} else {
		logger.Log.Info("Loaded .env file from project root")
	}

	// Re-initialize config after loading .env to pick up new variables
	config.Init()

	// Check for confirmation password
	clearPassword := viper.GetString("CLEAR_PASSWORD")
	if clearPassword == "" {
		return fmt.Errorf("CLEAR_PASSWORD environment variable is not set. This is required for safety")
	}

	// Prompt for password confirmation
	fmt.Print("⚠️  WARNING: This will permanently delete ALL words, sentences, topics, and learned words!\n")
	fmt.Print("Enter confirmation password: ")

	reader := bufio.NewReader(os.Stdin)
	inputPassword, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read password: %v", err)
	}

	inputPassword = strings.TrimSpace(inputPassword)
	if inputPassword != clearPassword {
		return fmt.Errorf("incorrect password. Operation cancelled")
	}

	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	logger.Log.Info("Starting database cleanup")

	stats := &ClearStats{
		StartTime: time.Now(),
	}

	// Perform cleanup in transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		return clearAllData(tx, stats)
	})

	if err != nil {
		return fmt.Errorf("failed to clear data: %v", err)
	}

	// Print final statistics
	printClearStats(stats)
	return nil
}

func clearAllData(tx *gorm.DB, stats *ClearStats) error {
	ctx := context.Background()

	// Step 1: Delete learned_words
	logger.Log.Info("Deleting learned words...")
	result := tx.WithContext(ctx).Delete(&models.LearnedWords{}, "1=1")
	if result.Error != nil {
		logger.Log.Error("Failed to delete learned words", zap.Error(result.Error))
		stats.ErrorsEncountered++
		return fmt.Errorf("failed to delete learned words: %v", result.Error)
	}
	stats.LearnedWordsDeleted = int(result.RowsAffected)
	logger.Log.Info("Deleted learned words", zap.Int("count", stats.LearnedWordsDeleted))

	// Step 2: Delete sentences
	logger.Log.Info("Deleting sentences...")
	result = tx.WithContext(ctx).Delete(&models.Sentence{}, "1=1")
	if result.Error != nil {
		logger.Log.Error("Failed to delete sentences", zap.Error(result.Error))
		stats.ErrorsEncountered++
		return fmt.Errorf("failed to delete sentences: %v", result.Error)
	}
	stats.SentencesDeleted = int(result.RowsAffected)
	logger.Log.Info("Deleted sentences", zap.Int("count", stats.SentencesDeleted))

	// Step 3: Delete words
	logger.Log.Info("Deleting words...")
	result = tx.WithContext(ctx).Delete(&models.Word{}, "1=1")
	if result.Error != nil {
		logger.Log.Error("Failed to delete words", zap.Error(result.Error))
		stats.ErrorsEncountered++
		return fmt.Errorf("failed to delete words: %v", result.Error)
	}
	stats.WordsDeleted = int(result.RowsAffected)
	logger.Log.Info("Deleted words", zap.Int("count", stats.WordsDeleted))

	// Step 4: Delete topics (in correct order for hierarchy)
	logger.Log.Info("Deleting topics...")
	err := deleteTopicsHierarchy(tx, ctx, stats)
	if err != nil {
		logger.Log.Error("Failed to delete topics", zap.Error(err))
		stats.ErrorsEncountered++
		return fmt.Errorf("failed to delete topics: %v", err)
	}

	return nil
}

func deleteTopicsHierarchy(tx *gorm.DB, ctx context.Context, stats *ClearStats) error {
	// Delete topics in reverse hierarchy order (children first, then parents)
	// We'll do this by repeatedly deleting topics that have no children

	for {
		// Find topics that have no children
		var topicsToDelete []models.Topic
		err := tx.WithContext(ctx).Raw(`
			SELECT t1.* FROM topics t1 
			WHERE NOT EXISTS (
				SELECT 1 FROM topics t2 WHERE t2.parent_id = t1.id
			)
		`).Scan(&topicsToDelete).Error

		if err != nil {
			return fmt.Errorf("failed to find leaf topics: %v", err)
		}

		// If no topics to delete, we're done
		if len(topicsToDelete) == 0 {
			break
		}

		// Delete the leaf topics
		var topicIDs []uuid.UUID
		for _, topic := range topicsToDelete {
			topicIDs = append(topicIDs, topic.ID)
		}

		result := tx.WithContext(ctx).Delete(&models.Topic{}, "id IN ?", topicIDs)
		if result.Error != nil {
			return fmt.Errorf("failed to delete topics batch: %v", result.Error)
		}

		deletedCount := int(result.RowsAffected)
		stats.TopicsDeleted += deletedCount

		logger.Log.Debug("Deleted topics batch",
			zap.Int("count", deletedCount),
			zap.Int("total_deleted", stats.TopicsDeleted))
	}

	logger.Log.Info("Deleted all topics", zap.Int("total_count", stats.TopicsDeleted))
	return nil
}

func printClearStats(stats *ClearStats) {
	duration := time.Since(stats.StartTime)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🧹 CLEAR STATISTICS")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("⏱️  Duration: %v\n", duration.Round(time.Second))
	fmt.Printf("📚 Learned words deleted: %d\n", stats.LearnedWordsDeleted)
	fmt.Printf("💬 Sentences deleted: %d\n", stats.SentencesDeleted)
	fmt.Printf("📝 Words deleted: %d\n", stats.WordsDeleted)
	fmt.Printf("🏷️  Topics deleted: %d\n", stats.TopicsDeleted)
	fmt.Printf("❌ Errors encountered: %d\n", stats.ErrorsEncountered)
	fmt.Println(strings.Repeat("=", 50))

	if stats.ErrorsEncountered == 0 {
		fmt.Println("✅ Database cleanup completed successfully!")
	} else {
		fmt.Println("⚠️  Database cleanup completed with errors!")
	}
}

func runEnrichWords(cmd *cobra.Command, args []string) error {
	// Initialize config and logger
	config.Init()
	logger.Init(true) // Enable debug logging
	defer logger.Log.Sync()

	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	logger.Log.Info("Starting word enrichment",
		zap.Int("limit", enrichLimit),
		zap.Int("delay_ms", enrichDelay))

	// Create enrichment service
	enrichmentService := utils.NewWordEnrichmentService(db)

	// Convert delay from milliseconds to duration
	delayDuration := time.Duration(enrichDelay) * time.Millisecond

	ctx := context.Background()

	// Run enrichment
	startTime := time.Now()
	err = enrichmentService.EnrichWordsInDatabase(ctx, enrichLimit, delayDuration)
	if err != nil {
		return fmt.Errorf("failed to enrich words: %v", err)
	}

	duration := time.Since(startTime)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🔍 WORD ENRICHMENT STATISTICS")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("⏱️  Duration: %v\n", duration.Round(time.Second))
	fmt.Printf("📝 Words processed: up to %d\n", enrichLimit)
	fmt.Printf("⏲️  Rate limit delay: %v\n", delayDuration)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✅ Word enrichment completed successfully!")

	return nil
}

func runEnrichSentences(cmd *cobra.Command, args []string) error {
	// Initialize config and logger
	config.Init()
	logger.Init(true) // Enable debug logging
	defer logger.Log.Sync()

	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	logger.Log.Info("Starting sentence enrichment",
		zap.Int("limit", sentenceLimit),
		zap.Int("delay_ms", sentenceDelay))

	// Create enrichment service
	enrichmentService := utils.NewWordEnrichmentService(db)

	// Convert delay from milliseconds to duration
	delayDuration := time.Duration(sentenceDelay) * time.Millisecond

	ctx := context.Background()

	// Run sentence enrichment
	startTime := time.Now()
	err = enrichmentService.EnrichSentencesInDatabase(ctx, sentenceLimit, delayDuration)
	if err != nil {
		return fmt.Errorf("failed to enrich sentences: %v", err)
	}

	duration := time.Since(startTime)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🎯 SENTENCE ENRICHMENT STATISTICS")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("⏱️  Duration: %v\n", duration.Round(time.Second))
	fmt.Printf("💬 Sentences processed: up to %d\n", sentenceLimit)
	fmt.Printf("⏲️  Rate limit delay: %v\n", delayDuration)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✅ Sentence enrichment completed successfully!")

	return nil
}

func runResetEnrichment(cmd *cobra.Command, args []string) error {
	// Initialize config and logger
	config.Init()
	logger.Init(true) // Enable debug logging
	defer logger.Log.Sync()

	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	logger.Log.Info("Starting reset enrichment")

	// Create enrichment service
	enrichmentService := utils.NewWordEnrichmentService(db)

	ctx := context.Background()

	// Run reset enrichment
	startTime := time.Now()
	count, err := enrichmentService.ResetFailedEnrichments(ctx)
	if err != nil {
		return fmt.Errorf("failed to reset enrichment: %v", err)
	}

	duration := time.Since(startTime)

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🔄 RESET ENRICHMENT STATISTICS")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("⏱️  Duration: %v\n", duration.Round(time.Second))
	fmt.Printf("🔄 Words reset: %d\n", count)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✅ Reset enrichment completed successfully!")

	return nil
}
