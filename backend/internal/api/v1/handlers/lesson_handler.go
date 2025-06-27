package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
	"strings"
	"unicode/utf8"
)

var exerciseTypes = []string{
	"translate_ru_to_en",
	"write_word_from_translation",
	"pick_option_sentence",
}

type LessonHandler struct {
	PreferenceRepo *postgres.PreferenceRepository
	TopicRepo      *postgres.TopicRepository
	SentenceRepo   *postgres.SentenceRepository
	Repo           *postgres.LessonRepository
}

func replaceWordWithUnderscores(text, word string) string {
	replacement := strings.Repeat("_", utf8.RuneCountInString(word))
	return strings.ReplaceAll(text, word, replacement)
}

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

	words, err := h.Repo.GetWordsForLesson(
		r.Context(),
		userID,
		userPref.Goal,
		userPref.CEFRLevel,
		lessonInfo.TotalWords,
	)

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
		sentence, err := h.SentenceRepo.GetByWordID(r.Context(), word.ID)
		if err != nil {
			http.Error(w, "failed to get sentence", http.StatusBadRequest)
			return
		}

		card.Sentences = append(card.Sentences, schemas.Sentence{
			Text:        sentence.Sentence,
			Translation: "Violets are Blue", //TODO: Remove moke in future
		})

		// Exercise process
		option := rand.Intn(3)

		var exercise schemas.Exercise
		exercise.Type = exerciseTypes[option]

		switch option {
		case 0: //translate_ru_to_en
			var translateRuToEn schemas.ExerciseTranslateRuToEn

			translateRuToEn.Text = word.Translation
			translateRuToEn.CorrectAnswer = word.Word
			translateRuToEn.PickOptions = []string{"biba", "boba", "bimba"}
			exercise.Data = translateRuToEn
		case 1: //write_word_from_translation
			var writeWordFromTranslation schemas.ExerciseWriteWordFromTranslation

			writeWordFromTranslation.Translation = word.Translation
			writeWordFromTranslation.CorrectAnswer = word.Word
			exercise.Data = writeWordFromTranslation
		case 2: //pick_option_sentence
			var pickOptionSentence schemas.ExercisePickOptionSentence

			pickOptionSentence.Template = replaceWordWithUnderscores(
				sentence.Sentence,
				word.Word,
			)
			pickOptionSentence.CorrectAnswer = word.Word
			pickOptionSentence.PickOptions = []string{"biba", "boba", "bimba"}
			exercise.Data = pickOptionSentence
		default:
		}

		cards = append(cards, card)
	}

	var lesson schemas.LessonResponse
	lesson.Lesson = lessonInfo
	lesson.Cards = cards
}
