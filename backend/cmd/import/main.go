package main

import (
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
	"fluently/go-backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
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

func init() {
	csvCmd.Flags().StringVarP(&csvFilePath, "file", "f", "", "Path to CSV file (required)")
	csvCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(csvCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCSVImport(cmd *cobra.Command, args []string) error {
	// Initialize config and logger
	config.Init()
	logger.Init(false)
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
	dsn := config.GetPostgresDSN()
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
			Sentences:   strings.TrimSpace(record[8]),
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
	return db.Transaction(func(tx *gorm.DB) error {
		ctx := context.Background()

		// Group records by word for merging translations
		wordGroups := make(map[string][]CSVRecord)
		for _, record := range batch {
			key := fmt.Sprintf("%s|%s|%s|%s|%s", record.Word, record.Topic, record.Subtopic, record.Subsubtopic, record.CEFRLevel)
			wordGroups[key] = append(wordGroups[key], record)
		}

		for _, group := range wordGroups {
			err := processWordGroup(tx, ctx, group, stats)
			if err != nil {
				logger.Log.Error("Failed to process word group", zap.Error(err), zap.String("word", group[0].Word))
				stats.ErrorsEncountered++
				continue
			}
		}

		return nil
	})
}

func processWordGroup(tx *gorm.DB, ctx context.Context, group []CSVRecord, stats *ImportStats) error {
	if len(group) == 0 {
		return nil
	}

	baseRecord := group[0]

	// Create topic hierarchy
	topicID, err := createTopicHierarchy(tx, ctx, baseRecord.Topic, baseRecord.Subtopic, baseRecord.Subsubtopic, stats)
	if err != nil {
		return fmt.Errorf("failed to create topic hierarchy: %v", err)
	}

	// Merge translations with deduplication
	translations := make(map[string]bool)
	for _, record := range group {
		if record.Translation != "" {
			translations[record.Translation] = true
		}
	}

	var translationList []string
	for translation := range translations {
		translationList = append(translationList, translation)
	}
	mergedTranslation := strings.Join(translationList, ",")

	// Check if word already exists
	var existingWord models.Word
	result := tx.WithContext(ctx).Where("word = ? AND topic_id = ?", baseRecord.Word, topicID).First(&existingWord)

	var wordID uuid.UUID
	if result.Error == gorm.ErrRecordNotFound {
		// Create new word
		newWord := models.Word{
			Word:         baseRecord.Word,
			Translation:  mergedTranslation,
			PartOfSpeech: "unknown",
			Context:      "",
			CEFRLevel:    baseRecord.CEFRLevel,
			TopicID:      topicID,
		}

		if err := tx.WithContext(ctx).Create(&newWord).Error; err != nil {
			return fmt.Errorf("failed to create word: %v", err)
		}
		wordID = newWord.ID
		stats.WordsImported++
	} else if result.Error != nil {
		return fmt.Errorf("failed to check existing word: %v", result.Error)
	} else {
		// Update existing word with merged translations
		existingTranslations := make(map[string]bool)
		if existingWord.Translation != "" {
			for _, t := range strings.Split(existingWord.Translation, ",") {
				existingTranslations[strings.TrimSpace(t)] = true
			}
		}

		// Add new translations
		for translation := range translations {
			existingTranslations[translation] = true
		}

		var updatedTranslationList []string
		for translation := range existingTranslations {
			updatedTranslationList = append(updatedTranslationList, translation)
		}

		existingWord.Translation = strings.Join(updatedTranslationList, ",")
		if baseRecord.CEFRLevel != "" {
			existingWord.CEFRLevel = baseRecord.CEFRLevel
		}

		if err := tx.WithContext(ctx).Save(&existingWord).Error; err != nil {
			return fmt.Errorf("failed to update word: %v", err)
		}
		wordID = existingWord.ID
		logger.Log.Info("Updated existing word", zap.String("word", baseRecord.Word))
	}

	// Process sentences for all records in group
	for _, record := range group {
		if record.Sentences != "" {
			sentencesAdded, err := processSentences(tx, ctx, wordID, record.Sentences)
			if err != nil {
				logger.Log.Error("Failed to process sentences", zap.Error(err), zap.String("word", record.Word))
				continue
			}
			stats.SentencesAdded += sentencesAdded
		}
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

		if err := tx.WithContext(ctx).Create(&newTopic).Error; err != nil {
			return uuid.Nil, false, fmt.Errorf("failed to create topic: %v", err)
		}

		return newTopic.ID, true, nil
	} else if result.Error != nil {
		return uuid.Nil, false, fmt.Errorf("failed to check existing topic: %v", result.Error)
	}

	return topic.ID, false, nil
}

func processSentences(tx *gorm.DB, ctx context.Context, wordID uuid.UUID, sentencesJSON string) (int, error) {
	if sentencesJSON == "" {
		return 0, nil
	}

	// Parse JSON sentences
	var sentencePairs []SentencePair
	if err := json.Unmarshal([]byte(sentencesJSON), &sentencePairs); err != nil {
		return 0, fmt.Errorf("failed to parse sentences JSON: %v", err)
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
