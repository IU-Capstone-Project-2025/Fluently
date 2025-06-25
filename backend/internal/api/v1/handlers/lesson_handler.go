package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"

	"github.com/google/uuid"
)

type LessonHandler struct {
	Repo *postgres.LessonRepository
}

func buildLessonResponse(lesson *models.Lesson) schemas.LessonResponse {
	cards := make([]schemas.CardSchema, 0, len(lesson.Cards))
	for _, card := range lesson.Cards {
		topicTitle := ""
		if card.Word.Topic != nil {
			topicTitle = card.Word.Topic.Title
		}

		// Преобразуем предложения слова в схемы
		sentences := make([]schemas.SentenceSchema, 0, len(card.Word.Sentences))
		for _, s := range card.Word.Sentences {
			sentences = append(sentences, schemas.SentenceSchema{
				SentenceID:  s.ID,
				Text:        s.Sentence,
				Translation: s.Translation,
			})
		}

		cardSchema := schemas.CardSchema{
			WordID:        card.Word.ID,
			Word:          card.Word.Word,
			Translation:   card.Word.Translation,
			Transcription: "",
			CEFRLevel:     "",
			IsNew:         false,
			Topic:         topicTitle,
			Subtopic:      "",
			Sentences:     sentences,
			Exercise:      schemas.ExerciseSchema{},
		}

		cards = append(cards, cardSchema)
	}

	lessonSchema := schemas.LessonSchema{
		LessonID:       lesson.ID,
		UserID:         lesson.UserID,
		StartedAt:      lesson.StartedAt.Format(time.RFC3339),
		WordsPerLesson: lesson.WordsPerLesson,
		TotalWords:     lesson.TotalWords,
	}

	sync := schemas.SyncSchema{
		Dirty:        false,
		LastSyncedAt: time.Now().Format(time.RFC3339),
	}

	return schemas.LessonResponse{
		Lesson: lessonSchema,
		Cards:  cards,
		Sync:   sync,
	}
}

// CreateLesson godoc
// @Summary      Create a lesson
// @Description  Create a new lesson with list of word IDs
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
//
// @Success      201  {object}  schemas.LessonResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/lessons/ [post]
func (h *LessonHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID         uuid.UUID   `json:"user_id"`
		WordsPerLesson int         `json:"words_per_lesson"`
		TotalWords     int         `json:"total_words"`
		WordIDs        []uuid.UUID `json:"word_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	lesson := models.Lesson{
		ID:             uuid.New(),
		UserID:         req.UserID,
		StartedAt:      time.Now(),
		WordsPerLesson: req.WordsPerLesson,
		TotalWords:     req.TotalWords,
	}

	if err := h.Repo.Create(r.Context(), &lesson, req.WordIDs); err != nil {
		http.Error(w, "failed to create lesson", http.StatusInternalServerError)
		return
	}

	// После создания грузим полностью для ответа
	createdLesson, err := h.Repo.GetByID(r.Context(), lesson.ID)
	if err != nil {
		http.Error(w, "failed to load created lesson", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildLessonResponse(createdLesson))
}

// GetLesson godoc
// @Summary      Get lesson by ID
// @Description  Returns a lesson by ID with cards and exercises
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Lesson ID"
// @Success      200  {object}  schemas.LessonResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Router       /api/v1/lessons/{id} [get]
func (h *LessonHandler) GetLesson(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid lesson id", http.StatusBadRequest)
		return
	}

	lesson, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "lesson not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildLessonResponse(lesson))
}

// GetLastLessonByUser godoc
// @Summary      Get last lesson by user ID
// @Description  Returns last lesson for the specified user
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id   query      string  true  "User ID"
// @Success      200  {object}  schemas.LessonResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Router       /api/v1/lessons/last [get]
func (h *LessonHandler) GetLastLessonByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query param required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	lesson, err := h.Repo.GetLastByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "lesson not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildLessonResponse(lesson))
}

// DeleteLesson godoc
// @Summary      Delete a lesson by ID
// @Description  Deletes a lesson and its cards
// @Tags         lessons
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Lesson ID"
// @Success      204  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/lessons/{id} [delete]
func (h *LessonHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid lesson id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete lesson", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
