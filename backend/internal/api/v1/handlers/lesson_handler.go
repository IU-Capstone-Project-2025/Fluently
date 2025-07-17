package handlers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
	"fluently/go-backend/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// List of exercise types
var exerciseTypes = []string{
	"translate_ru_to_en",
	"write_word_from_translation",
	"pick_option_sentence",
}

// LessonHandler handles the lesson endpoint
type LessonHandler struct {
	PreferenceRepo  *postgres.PreferenceRepository
	TopicRepo       *postgres.TopicRepository
	SentenceRepo    *postgres.SentenceRepository
	PickOptionRepo  *postgres.PickOptionRepository
	WordRepo        *postgres.WordRepository
	Repo            *postgres.LessonRepository
	LearnedWordRepo *postgres.LearnedWordRepository
	ThesaurusClient *utils.ThesaurusClient
}

// replaceWordWithUnderscores replaces a word in a text with underscores
func replaceWordWithUnderscores(text, word string) string {
	lowerText := strings.ToLower(text)
	lowerWord := strings.ToLower(word)

	wordIndex := strings.Index(lowerText, lowerWord)
	if wordIndex == -1 {
		return text
	}

	return text[:wordIndex] + strings.Repeat("_", len(word)) + text[wordIndex+len(word):]
}

// GenerateLesson godoc
// @Summary Generate a new lesson for the user
// @Description Creates a personalized lesson based on user preferences, including words, exercises, and topics
// @Tags lessons
// @Produce json
// @Security BearerAuth
// @Success 200 {object} schemas.LessonResponse "Successfully generated lesson"
// @Failure 400 {string} string "Bad request - invalid user or preferences"
// @Failure 401 {string} string "Unauthorized - invalid or missing token"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/lesson [get]
func (h *LessonHandler) GenerateLesson(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := user.ID

	// Lesson info block
	var lessonInfo schemas.LessonInfo

	userPref, err := h.PreferenceRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get preference", http.StatusBadRequest)
		return
	}

	lessonInfo.WordsPerLesson = userPref.WordsPerDay
	lessonInfo.TotalWords = userPref.WordsPerDay * 2
	lessonInfo.CEFRLevel = userPref.CEFRLevel
	lessonInfo.StartedAt = time.Now().UTC().Format(time.RFC3339)

	// Lesson card block
	var cards []schemas.Card

	// -------------------- Fetch recommendations from Thesaurus --------------------
	// 1. Gather list of words the user already knows
	learnedWords, err := h.LearnedWordRepo.ListByUserID(r.Context(), userID)
	if err != nil {
		logger.Log.Error("Failed to list learned words", zap.Error(err))
		http.Error(w, "failed to list learned words", http.StatusInternalServerError)
		return
	}

	var known []string
	for _, lw := range learnedWords {
		wm, err := h.WordRepo.GetByID(r.Context(), lw.WordID)
		if err == nil {
			known = append(known, wm.Word)
		}
	}

	// 2. Request recommendations until we collect enough valid words from DB
	var words []models.Word
	seen := make(map[uuid.UUID]struct{})
	attempts := 0
	for len(words) < lessonInfo.TotalWords && attempts < 5 {
		recs, err := h.ThesaurusClient.Recommend(r.Context(), known)
		if err != nil {
			logger.Log.Error("Failed to get thesaurus recommendations", zap.Error(err))
			break // fallback later
		}

		for _, rec := range recs {
			if len(words) >= lessonInfo.TotalWords {
				break
			}

			wm, err := h.WordRepo.GetByValue(r.Context(), rec.Word)
			if err != nil {
				// Recommendation not present in local DB – skip
				continue
			}
			if _, exists := seen[wm.ID]; exists {
				continue
			}
			words = append(words, *wm)
			seen[wm.ID] = struct{}{}
		}
		attempts++
	}

	// 3. Fallback: fill remaining slots with random DB words
	if len(words) < lessonInfo.TotalWords {
		remaining := lessonInfo.TotalWords - len(words)
		additional, err := h.Repo.GetWordsForLesson(
			r.Context(),
			userID,
			userPref.CEFRLevel,
			userPref.Goal,
			remaining,
		)
		if err != nil {
			logger.Log.Error("Failed to get fallback words", zap.Error(err))
			http.Error(w, "failed to get words for lesson", http.StatusBadRequest)
			return
		}
		for _, wdb := range additional {
			if _, exists := seen[wdb.ID]; exists {
				continue
			}
			words = append(words, wdb)
			seen[wdb.ID] = struct{}{}
			if len(words) >= lessonInfo.TotalWords {
				break
			}
		}
	}

	if len(words) == 0 {
		http.Error(w, "no suitable words found for lesson", http.StatusInternalServerError)
		return
	}

	// Shuffle words
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	// Process words
	for _, word := range words {
		var card schemas.Card

		card.WordID = word.ID
		card.Word = word.Word
		card.Translation = word.Translation

		topic, err := h.TopicRepo.GetByID(r.Context(), *word.TopicID)
		if err != nil {
			http.Error(w, "failed to get topic", http.StatusBadRequest)
			return
		}

		// Topic and subtopic process
		card.Subtopic = topic.Title

		for topic.ParentID != nil {
			topic, err = h.TopicRepo.GetByID(r.Context(), *topic.ParentID)
			if err != nil {
				http.Error(w, "failed to get topic", http.StatusBadRequest)
				return
			}
		}

		card.Topic = topic.Title

		// Sentence process
		sentences, err := h.SentenceRepo.GetByWordID(r.Context(), word.ID)
		if err != nil {
			http.Error(w, "failed to get sentence", http.StatusBadRequest)
			return
		}

		for _, sentence := range sentences {
			card.Sentences = append(card.Sentences, schemas.Sentence{
				Text:        sentence.Sentence,
				Translation: sentence.Translation,
			})
		}

		// Exercise process
		var exercises []schemas.Exercise

		// translate_ru_to_en
		var translateRuToEn schemas.ExerciseTranslateRuToEn

		translateRuToEn.Text = word.Translation
		translateRuToEn.CorrectAnswer = word.Word

		pickOptionTranslate, err := h.PickOptionRepo.GetOptionByWordID(r.Context(), word.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				pickOptionTranslate = &models.PickOption{
					WordID: word.ID,
					Option: []string{word.Word},
				}
			} else {
				http.Error(w, "failed to get pick option", http.StatusInternalServerError)
				return
			}
		}

		if len(pickOptionTranslate.Option) == 1 {
			const optionToAdd = 3

			randomWords, err := h.WordRepo.GetRandomWordsByCEFRLevel(r.Context(), userPref.CEFRLevel, optionToAdd)
			if err != nil {
				http.Error(w, "failed to get random words", http.StatusBadRequest)
				return
			}

			for _, word := range randomWords {
				pickOptionTranslate.Option = append(pickOptionTranslate.Option, word.Word)
			}
		}

		// Shuffle options
		rand.Shuffle(len(pickOptionTranslate.Option), func(i, j int) {
			pickOptionTranslate.Option[i], pickOptionTranslate.Option[j] = pickOptionTranslate.Option[j], pickOptionTranslate.Option[i]
		})

		translateRuToEn.PickOptions = pickOptionTranslate.Option

		exercises = append(exercises, schemas.Exercise{
			Type: "translate_ru_to_en",
			Data: translateRuToEn,
		})

		// write_word_from_translation
		var writeWordFromTranslation schemas.ExerciseWriteWordFromTranslation

		writeWordFromTranslation.Translation = word.Translation
		writeWordFromTranslation.CorrectAnswer = word.Word
		exercises = append(exercises, schemas.Exercise{
			Type: "write_word_from_translation",
			Data: writeWordFromTranslation,
		})

		// pick_option_sentence
		var pickOptionSentence schemas.ExercisePickOptionSentence

		// Shuffle options
		rand.Shuffle(len(pickOptionTranslate.Option), func(i, j int) {
			pickOptionTranslate.Option[i], pickOptionTranslate.Option[j] = pickOptionTranslate.Option[j], pickOptionTranslate.Option[i]
		})

		pickOptionSentence.Template = replaceWordWithUnderscores(
			sentences[0].Sentence,
			word.Word,
		)
		pickOptionSentence.CorrectAnswer = word.Word
		pickOptionSentence.PickOptions = pickOptionTranslate.Option

		exercises = append(exercises, schemas.Exercise{
			Type: "pick_option_sentence",
			Data: pickOptionSentence,
		})

		// Pick random exercise
		randomExercise := exercises[rand.Intn(len(exercises))]

		card.Exercise = randomExercise
		cards = append(cards, card)
		logger.Log.Info("Card generated", zap.Any("card", card))
	}

	// Generate lesson
	var lesson schemas.LessonResponse
	lesson.Lesson = lessonInfo
	lesson.Cards = cards
	logger.Log.Info("Lesson generated", zap.Any("lesson", lesson))

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(lesson); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
