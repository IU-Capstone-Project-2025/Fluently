package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/utils"
	"fluently/go-backend/pkg/logger"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProgressHandler handles the progress endpoint
type ProgressHandler struct {
	LearnedWordRepo    *postgres.LearnedWordRepository
	WordRepo           *postgres.WordRepository
	NotLearnedWordRepo *postgres.NotLearnedWordRepository
	LLMClient          *utils.LLMClient
	Redis              *goredis.Client
}

// ProgressRequest is a request body for updating user progress
type ProgressRequest struct {
	WordID          uuid.UUID  `json:"word_id"`
	LearnedAt       *time.Time `json:"learned_at,omitempty"`
	ConfidenceScore *int       `json:"confidence_score,omitempty"`
	CntReviewed     *int       `json:"cnt_reviewed,omitempty"`
}

// UpdateUserProgress godoc
// @Summary      Update user progress
// @Description  Updates the user's learned words progress. Accepts an array of word progress objects with word-translation pairs.
// @Tags         progress
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        progress body []ProgressRequest true "List of progress updates with word-translation pairs"
// @Success      200  {string}  string  "ok"
// @Router       /api/v1/progress [post]
func (h *ProgressHandler) UpdateUserProgress(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/progress"
	method := r.Method
	statusCode := 200
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()

	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		statusCode = 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := user.ID

	var progress []ProgressRequest
	err = json.NewDecoder(r.Body).Decode(&progress)
	if err != nil {
		statusCode = 400
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(progress) == 0 {
		statusCode = 400
		http.Error(w, "no progress provided", http.StatusBadRequest)
		return
	}

	for _, p := range progress {
		word, err := h.WordRepo.GetByID(r.Context(), p.WordID)
		if err != nil {
			statusCode = 404
			http.Error(w, "word not found: "+err.Error(), http.StatusNotFound)
			return
		}

		existing, err := h.LearnedWordRepo.GetByUserWordID(r.Context(), userID, word.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			statusCode = 500
			http.Error(w, "failed to get learned word", http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC()

		if isNotLearnedWord(p) {
			notLearnedExists, err := h.NotLearnedWordRepo.Exists(r.Context(), userID, word.ID)
			if err != nil {
				statusCode = 500
				http.Error(w, "failed to check if not learned word exists", http.StatusInternalServerError)
				return
			}

			if !notLearnedExists {
				nlw := &models.NotLearnedWords{
					ID:     uuid.New(),
					UserID: userID,
					WordID: word.ID,
				}
				if err := h.NotLearnedWordRepo.Create(r.Context(), nlw); err != nil {
					statusCode = 500
					http.Error(w, "failed to create not learned word", http.StatusInternalServerError)
					return
				}
			}

			continue
		}

		if existing == nil {
			lw := &models.LearnedWords{
				ID:           uuid.New(),
				UserID:       userID,
				WordID:       word.ID,
				LastReviewed: now,
			}

			if p.LearnedAt != nil {
				lw.LearnedAt = *p.LearnedAt
			} else {
				lw.LearnedAt = now
			}

			if p.ConfidenceScore != nil {
				lw.ConfidenceScore = *p.ConfidenceScore
			} else {
				lw.ConfidenceScore = 0
			}

			if p.CntReviewed != nil {
				lw.CountOfRevisions = *p.CntReviewed
			} else {
				lw.CountOfRevisions = 0
			}

			if err := h.LearnedWordRepo.Create(r.Context(), lw); err != nil {
				statusCode = 500
				http.Error(w, "failed to create learned word", http.StatusInternalServerError)
				return
			}

			if err := h.NotLearnedWordRepo.DeleteIfExists(r.Context(), userID, word.ID); err != nil {
				statusCode = 500
				http.Error(w, "failed to delete not learned word", http.StatusInternalServerError)
				return
			}
		} else {
			existing.LastReviewed = now

			if p.LearnedAt != nil {
				existing.LearnedAt = *p.LearnedAt
			}

			if p.ConfidenceScore != nil {
				existing.ConfidenceScore = *p.ConfidenceScore
			}

			if p.CntReviewed != nil {
				existing.CountOfRevisions = *p.CntReviewed
			}

			if err := h.LearnedWordRepo.Update(r.Context(), existing); err != nil {
				statusCode = 500
				http.Error(w, "failed to update learned word", http.StatusInternalServerError)
				return
			}
		}
	}

	// Generate conversation topic for newly learned words
	if err := h.generateAndStoreConversationTopic(r.Context(), userID, progress); err != nil {
		logger.Log.Warn("failed to generate conversation topic", zap.Error(err))
		// Don't fail the request, just log the warning
	}

	// Return ok
	w.WriteHeader(http.StatusOK)
}

func isNotLearnedWord(p ProgressRequest) bool {
	return p.LearnedAt == nil && p.ConfidenceScore == nil && p.CntReviewed == nil
}

// generateAndStoreConversationTopic generates a conversation topic based on learned words and stores it in Redis
func (h *ProgressHandler) generateAndStoreConversationTopic(ctx context.Context, userID uuid.UUID, progress []ProgressRequest) error {
	// Collect learned words from the progress
	var learnedWords []models.Word
	for _, p := range progress {
		if !isNotLearnedWord(p) {
			word, err := h.WordRepo.GetByID(ctx, p.WordID)
			if err != nil {
				logger.Log.Warn("failed to get word for topic generation", zap.Error(err), zap.String("word_id", p.WordID.String()))
				continue
			}
			learnedWords = append(learnedWords, *word)
		}
	}

	if len(learnedWords) == 0 {
		logger.Log.Info("no learned words found for topic generation")
		return nil
	}

	// Generate topic using LLM
	topic, err := h.generateTopicFromWords(ctx, learnedWords)
	if err != nil {
		return fmt.Errorf("failed to generate topic: %w", err)
	}

	// Store topic and words in Redis
	if err := h.storeConversationTopic(ctx, userID, topic, learnedWords); err != nil {
		return fmt.Errorf("failed to store conversation topic: %w", err)
	}

	logger.Log.Info("generated and stored conversation topic",
		zap.String("topic", topic),
		zap.Int("word_count", len(learnedWords)),
		zap.String("user_id", userID.String()))

	return nil
}

// generateTopicFromWords generates a conversation topic based on the learned words
func (h *ProgressHandler) generateTopicFromWords(ctx context.Context, words []models.Word) (string, error) {
	// Build words list for the prompt
	var wordsList strings.Builder
	for i, word := range words {
		wordsList.WriteString(fmt.Sprintf("- %s (%s): %s", word.Word, word.PartOfSpeech, word.Translation))
		if i < len(words)-1 {
			wordsList.WriteString("\n")
		}
	}

	// Create prompt for topic generation
	prompt := fmt.Sprintf(`Based on these English words that a user just learned, generate a natural conversation topic that would allow them to practice these words in context. The topic should be something people would naturally talk about.

Words learned:
%s

Generate a short, specific conversation topic (2-4 words) that relates to these words. Examples: "buying a ticket in airport", "cooking dinner", "planning a vacation", "shopping for clothes".

Return only the topic, nothing else.`, wordsList.String())

	// Call LLM to generate topic
	llmMsgs := []utils.LLMMessage{
		{Role: "user", Content: prompt},
	}

	response, err := h.LLMClient.Chat(ctx, llmMsgs, "balanced", nil, nil)

	if err != nil {
		return "", fmt.Errorf("LLM error: %w", err)
	}

	// Clean up the response (remove quotes, extra whitespace, etc.)
	topic := strings.TrimSpace(response)
	topic = strings.Trim(topic, `"'`)

	logger.Log.Info("generated topic for conversation", zap.String("topic", topic))
	if topic == "" {
		return "general conversation", nil // fallback topic
	}

	return topic, nil
}

// storeConversationTopic stores the generated topic and words in Redis
func (h *ProgressHandler) storeConversationTopic(ctx context.Context, userID uuid.UUID, topic string, words []models.Word) error {
	// Convert words to ChatWord format for consistency with chat handler
	var chatWords []ChatWord
	for _, word := range words {
		chatWord := ChatWord{
			Word:         word.Word,
			Context:      word.Context,
			PartOfSpeech: word.PartOfSpeech,
		}
		chatWords = append(chatWords, chatWord)
	}

	// Create data structure for Redis
	topicData := map[string]interface{}{
		"topic":      topic,
		"words":      chatWords,
		"created_at": time.Now().UTC(),
	}

	// Marshal to JSON
	data, err := json.Marshal(topicData)
	if err != nil {
		return fmt.Errorf("failed to marshal topic data: %w", err)
	}

	// Store in Redis with key "chat_topic:{userID}"
	key := "chat_topic:" + userID.String()
	err = h.Redis.Set(ctx, key, data, 24*time.Hour).Err() // expire after a day
	if err != nil {
		return fmt.Errorf("failed to store topic in Redis: %w", err)
	}

	return nil
}
