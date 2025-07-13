package handlers

import (
	"errors"
	"math/rand"
	"net/http"
	"encoding/json"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"

	"gorm.io/gorm"
)

type DayWordHandler struct {
	WordRepo        *postgres.WordRepository
	PreferenceRepo  *postgres.PreferenceRepository
	TopicRepo       *postgres.TopicRepository
	SentenceRepo    *postgres.SentenceRepository
	PickOptionRepo  *postgres.PickOptionRepository
	LearnedWordRepo *postgres.LearnedWordRepository
}

// godoc
// @Summary      Get day word
// @Description  Returns the day word for the user
// @Tags         day-word
// @Produce      json
// @Security     BearerAuth
// @Success 	 200 {object}  schemas.DayWordResponse "Successfully returned day word"
// @Failure      400  {string}  string  "Invalid request - plain text error message"
// @Failure      404  {string}  string  "Resource not found - plain text error message"
// @Failure      500  {string}  string  "Internal server error - plain text error message"
// @Router       /api/v1/day-word/ [get]
func (h *DayWordHandler) GetDayWord(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetCurrentUser(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := user.ID

	userPref, err := h.PreferenceRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get preference", http.StatusInternalServerError)
		return
	}

	var dayWordResponse schemas.DayWordResponse

	dayWord, err := h.WordRepo.GetDayWord(r.Context(), userPref.CEFRLevel, userID)
	if err != nil {
		http.Error(w, "failed to get day word", http.StatusInternalServerError)
		return
	}
	dayWordResponse.Word = dayWord.Word
	dayWordResponse.Translation = dayWord.Translation
	dayWordResponse.CEFRLevel = dayWord.CEFRLevel

	if dayWord.Phonetic != "" {
		dayWordResponse.Transcription = &dayWord.Phonetic
	}

	topic, err := h.TopicRepo.GetByID(r.Context(), *dayWord.TopicID)
	if err != nil {
		http.Error(w, "failed to get topic", http.StatusBadRequest)
		return
	}

	dayWordResponse.Subtopic = topic.Title

	for topic.ParentID != nil {
		topic, err = h.TopicRepo.GetByID(r.Context(), *topic.ParentID)
		if err != nil {
			http.Error(w, "failed to get topic", http.StatusBadRequest)
			return
		}
	}

	dayWordResponse.Topic = topic.Title

	sentences, err := h.SentenceRepo.GetByWordID(r.Context(), dayWord.ID)
	if err != nil {
		http.Error(w, "failed to get sentence", http.StatusBadRequest)
		return
	}

	for _, sentence := range sentences {
		dayWordResponse.Sentences = append(dayWordResponse.Sentences, schemas.Sentence{
			Text:        sentence.Sentence,
			Translation: sentence.Translation,
		})
	}

	// Exercise process
	var exercises []schemas.Exercise

	// translate_ru_to_en
	var translateRuToEn schemas.ExerciseTranslateRuToEn

	translateRuToEn.Text = dayWord.Translation
	translateRuToEn.CorrectAnswer = dayWord.Word

	pickOptionTranslate, err := h.PickOptionRepo.GetOptionByWordID(r.Context(), dayWord.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pickOptionTranslate = &models.PickOption{
				WordID: dayWord.ID,
				Option: []string{dayWord.Word},
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

	writeWordFromTranslation.Translation = dayWord.Translation
	writeWordFromTranslation.CorrectAnswer = dayWord.Word
	exercises = append(exercises, schemas.Exercise{
		Type: "write_word_from_translation",
		Data: writeWordFromTranslation,
	})

	// pick_option_sentence
	var pickOptionSentence schemas.ExercisePickOptionSentence

	rand.Shuffle(len(pickOptionTranslate.Option), func(i, j int) {
		pickOptionTranslate.Option[i], pickOptionTranslate.Option[j] = pickOptionTranslate.Option[j], pickOptionTranslate.Option[i]
	})

	pickOptionSentence.Template = replaceWordWithUnderscores(
		sentences[0].Sentence,
		dayWord.Word,
	)
	pickOptionSentence.CorrectAnswer = dayWord.Word
	pickOptionSentence.PickOptions = pickOptionTranslate.Option

	exercises = append(exercises, schemas.Exercise{
		Type: "pick_option_sentence",
		Data: pickOptionSentence,
	})

	randomExercise := exercises[rand.Intn(len(exercises))]

	dayWordResponse.Exercise = randomExercise

	learned, err := h.LearnedWordRepo.IsLearned(r.Context(), userID, dayWord.ID)
	if err != nil {
		http.Error(w, "failed to get learned word", http.StatusInternalServerError)
		return
	}

	if learned {
		dayWordResponse.IsLearned = true
	} else {
		dayWordResponse.IsLearned = false
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dayWordResponse)
}
