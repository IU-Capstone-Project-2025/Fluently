package postgres

import (
	"context"
	"fluently/go-backend/internal/repository/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LessonRepository is a repository for lessons
type LessonRepository struct {
	db *gorm.DB
}

// NewLessonRepository creates a new instance of LessonRepository
func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

// Create creates a new lesson
func (r *LessonRepository) Create(ctx context.Context, lesson *models.Lesson, words []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(lesson).Error; err != nil {
			return err
		}

		for i, wordID := range words {
			card := models.LessonCard{
				LessonID: lesson.ID,
				WordID:   wordID,
				Order:    i,
			}

			if err := tx.Create(&card).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByID returns a lesson by id
func (r *LessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.WithContext(ctx).
		Preload("Cards.Word.Topic").
		Preload("Cards.Word.Sentences").
		Preload("Cards.Word.Exercises.ExerciseType").
		Preload("Cards.Word.Exercises.PickOption").
		First(&lesson, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

// GetLastByUser returns the last lesson for a user
func (r *LessonRepository) GetLastByUser(ctx context.Context, userID uuid.UUID) (*models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.WithContext(ctx).
		Preload("Cards.Word.Topic").
		Preload("Cards.Word.Sentences").
		Preload("Cards.Word.Exercises.ExerciseType").
		Preload("Cards.Word.Exercises.PickOption").
		Where("user_id = ?", userID).
		Order("started_at DESC").
		First(&lesson).Error
	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

// Delete removes a lesson and all its associated cards
func (r *LessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("lesson_id = ?", id).Delete(&models.LessonCard{}).Error; err != nil {
			return err
		}

		return tx.Delete(&models.Lesson{}, "id = ?", id).Error
	})
}

// GetWordsForLesson needs for generate cards in lesson (hierarchical fallback strategy)
func (r *LessonRepository) GetWordsForLesson(
	ctx context.Context,
	userID uuid.UUID,
	cefrLevel string,
	goalTopicTitle string,
	limit int,
) ([]models.Word, error) {
	var words []models.Word
	seen := make(map[uuid.UUID]struct{})

	cefrLevel = strings.ToLower(cefrLevel)

	// Step 1: Try to find the goal topic
	var goalTopic *models.Topic
	var err error

	if strings.ToLower(goalTopicTitle) == "learn new words" ||
		goalTopicTitle == "" || strings.ToLower(goalTopicTitle) == "general" ||
		strings.ToLower(goalTopicTitle) == "all" {
		// For general goals, get a random main topic
		var topics []models.Topic
		if err := r.db.WithContext(ctx).
			Where("parent_id IS NULL").
			Order("RANDOM()").
			Limit(1).
			Find(&topics).Error; err != nil {
			return nil, err
		}
		if len(topics) > 0 {
			goalTopic = &topics[0]
		}
	} else {
		// Try to find the goal topic by title
		var foundTopic models.Topic
		err = r.db.WithContext(ctx).
			Where("title = ? AND parent_id IS NULL", goalTopicTitle).
			First(&foundTopic).Error
		if err == nil {
			goalTopic = &foundTopic
		} else {
			// If goal topic not found, fallback to random main topic
			var topics []models.Topic
			if err := r.db.WithContext(ctx).
				Where("parent_id IS NULL").
				Order("RANDOM()").
				Limit(1).
				Find(&topics).Error; err != nil {
				return nil, err
			}
			if len(topics) > 0 {
				goalTopic = &topics[0]
			}
		}
	}

	if goalTopic == nil {
		// Ultimate fallback: get random words by CEFR level only
		return r.getRandomWordsByCEFRLevelOnly(ctx, userID, cefrLevel, limit)
	}

	// Step 2: Get words from the goal topic itself
	goalTopicWords, err := r.getWordsByTopicID(ctx, userID, cefrLevel, goalTopic.ID, limit)
	if err != nil {
		return nil, err
	}

	for _, word := range goalTopicWords {
		if _, exists := seen[word.ID]; !exists {
			words = append(words, word)
			seen[word.ID] = struct{}{}
		}
	}

	// If we have enough words, return them
	if len(words) >= limit {
		return words[:limit], nil
	}

	// Step 3: Get words from subtopics of the goal topic
	var subtopicTopics []models.Topic
	if err := r.db.WithContext(ctx).
		Where("parent_id = ?", goalTopic.ID).
		Find(&subtopicTopics).Error; err == nil {

		var subtopicIDs []uuid.UUID
		for _, subtopic := range subtopicTopics {
			subtopicIDs = append(subtopicIDs, subtopic.ID)
		}

		if len(subtopicIDs) > 0 {
			remaining := limit - len(words)
			subtopicWords, err := r.getWordsByTopicIDs(ctx, userID, cefrLevel, subtopicIDs, remaining)
			if err == nil {
				for _, word := range subtopicWords {
					if _, exists := seen[word.ID]; !exists {
						words = append(words, word)
						seen[word.ID] = struct{}{}
					}
				}
			}
		}
	}

	// If we have enough words, return them
	if len(words) >= limit {
		return words[:limit], nil
	}

	// Step 4: Get words from subsubtopics (children of subtopics)
	if len(subtopicTopics) > 0 {
		var subsubtopicIDs []uuid.UUID
		for _, subtopic := range subtopicTopics {
			var subsubtopics []models.Topic
			if err := r.db.WithContext(ctx).
				Where("parent_id = ?", subtopic.ID).
				Find(&subsubtopics).Error; err == nil {
				for _, subsubtopic := range subsubtopics {
					subsubtopicIDs = append(subsubtopicIDs, subsubtopic.ID)
				}
			}
		}

		if len(subsubtopicIDs) > 0 {
			remaining := limit - len(words)
			subsubtopicWords, err := r.getWordsByTopicIDs(ctx, userID, cefrLevel, subsubtopicIDs, remaining)
			if err == nil {
				for _, word := range subsubtopicWords {
					if _, exists := seen[word.ID]; !exists {
						words = append(words, word)
						seen[word.ID] = struct{}{}
					}
				}
			}
		}
	}

	// If we have enough words, return them
	if len(words) >= limit {
		return words[:limit], nil
	}

	// Step 5: Get words from sibling topics (topics with the same parent as goal topic)
	if goalTopic.ParentID != nil {
		var siblingTopics []models.Topic
		if err := r.db.WithContext(ctx).
			Where("parent_id = ? AND id != ?", goalTopic.ParentID, goalTopic.ID).
			Find(&siblingTopics).Error; err == nil {

			var siblingIDs []uuid.UUID
			for _, sibling := range siblingTopics {
				siblingIDs = append(siblingIDs, sibling.ID)
			}

			if len(siblingIDs) > 0 {
				remaining := limit - len(words)
				siblingWords, err := r.getWordsByTopicIDs(ctx, userID, cefrLevel, siblingIDs, remaining)
				if err == nil {
					for _, word := range siblingWords {
						if _, exists := seen[word.ID]; !exists {
							words = append(words, word)
							seen[word.ID] = struct{}{}
						}
					}
				}
			}
		}
	}

	// If we have enough words, return them
	if len(words) >= limit {
		return words[:limit], nil
	}

	// Step 6: Ultimate fallback - random words by CEFR level
	remaining := limit - len(words)
	fallbackWords, err := r.getRandomWordsByCEFRLevelOnly(ctx, userID, cefrLevel, remaining)
	if err != nil {
		return nil, err
	}

	for _, word := range fallbackWords {
		if _, exists := seen[word.ID]; !exists {
			words = append(words, word)
			seen[word.ID] = struct{}{}
		}
	}

	return words, nil
}

// Helper method to get words by a single topic ID
func (r *LessonRepository) getWordsByTopicID(ctx context.Context, userID uuid.UUID, cefrLevel string, topicID uuid.UUID, limit int) ([]models.Word, error) {
	subQuery := r.db.
		Table("learned_words").
		Select("word_id").
		Where("user_id = ?", userID)

	var words []models.Word
	err := r.db.WithContext(ctx).
		Model(&models.Word{}).
		Where("cefr_level = ?", cefrLevel).
		Where("topic_id = ?", topicID).
		Where("id NOT IN (?)", subQuery).
		Order("RANDOM()").
		Limit(limit).
		Find(&words).Error

	return words, err
}

// Helper method to get words by multiple topic IDs
func (r *LessonRepository) getWordsByTopicIDs(ctx context.Context, userID uuid.UUID, cefrLevel string, topicIDs []uuid.UUID, limit int) ([]models.Word, error) {
	subQuery := r.db.
		Table("learned_words").
		Select("word_id").
		Where("user_id = ?", userID)

	var words []models.Word
	err := r.db.WithContext(ctx).
		Model(&models.Word{}).
		Where("cefr_level = ?", cefrLevel).
		Where("topic_id IN ?", topicIDs).
		Where("id NOT IN (?)", subQuery).
		Order("RANDOM()").
		Limit(limit).
		Find(&words).Error

	return words, err
}

// Helper method to get random words by CEFR level only (ultimate fallback)
func (r *LessonRepository) getRandomWordsByCEFRLevelOnly(ctx context.Context, userID uuid.UUID, cefrLevel string, limit int) ([]models.Word, error) {
	subQuery := r.db.
		Table("learned_words").
		Select("word_id").
		Where("user_id = ?", userID)

	var words []models.Word
	err := r.db.WithContext(ctx).
		Model(&models.Word{}).
		Where("cefr_level = ?", cefrLevel).
		Where("id NOT IN (?)", subQuery).
		Order("RANDOM()").
		Limit(limit).
		Find(&words).Error

	return words, err
}
