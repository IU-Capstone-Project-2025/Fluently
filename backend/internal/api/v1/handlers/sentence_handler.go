package handlers

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
)

type SentenceHandler struct {
	Repo *postgres.SentenceRepository
}

func buildSentenceResponse(sentence *models.Sentence) schemas.SentenceResponse {
	return schemas.SentenceResponse{
		ID:          sentence.ID.String(),
		WordID:      sentence.WordID.String(),
		Sentence:    sentence.Sentence,
		Translation: sentence.Translation,
	}
}

// ListSentences godoc
// @Summary      Get sentences for a word
// @Description  Returns all sentences for the specified word
// @Tags         sentences
// @Accept       json
// @Produce      json
// @Param        word_id   path      string  true  "Word ID"
// @Success      200  {array}   schemas.SentenceResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /words/{word_id}/sentences [get]
func (h *SentenceHandler) ListSentences(w http.ResponseWriter, r *http.Request) {
	wordID, err := utils.ParseUUIDParam(r, "word_id")
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	sentences, err := h.Repo.ListByWord(r.Context(), wordID)
	if err != nil {
		http.Error(w, "failed to fetch sentences", http.StatusInternalServerError)
		return
	}

	var resp []schemas.SentenceResponse
	for _, s := range sentences {
		resp = append(resp, buildSentenceResponse(&s))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateSentence godoc
// @Summary      Create a sentence
// @Description  Adds a new sentence for a word
// @Tags         sentences
// @Accept       json
// @Produce      json
// @Param        sentence  body      schemas.CreateSentenceRequest  true  "Sentence data"
// @Success      201  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /sentences/ [post]
func (h *SentenceHandler) CreateSentence(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateSentenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	s := models.Sentence{
		WordID:      req.WordID,
		Sentence:    req.Sentence,
		Translation: req.Translation,
	}

	if err := h.Repo.Create(r.Context(), &s); err != nil {
		http.Error(w, "failed to create sentence", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(buildSentenceResponse(&s))
}

// UpdateSentence godoc
// @Summary      Update a sentence
// @Description  Updates an existing sentence by ID
// @Tags         sentences
// @Accept       json
// @Produce      json
// @Param        id        path      string                        true  "Sentence ID"
// @Param        sentence  body      schemas.CreateSentenceRequest true  "Sentence data"
// @Success      200  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /sentences/{id} [put]
func (h *SentenceHandler) UpdateSentence(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req schemas.CreateSentenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sentence, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "sentence not found", http.StatusNotFound)
		return
	}

	sentence.WordID = req.WordID
	sentence.Sentence = req.Sentence
	sentence.Translation = req.Translation

	if err := h.Repo.Update(r.Context(), sentence); err != nil {
		http.Error(w, "failed to update sentence", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildSentenceResponse(sentence))
}

// DeleteSentence godoc
// @Summary      Delete a sentence
// @Description  Deletes a sentence by ID
// @Tags         sentences
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Sentence ID"
// @Success      204  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /sentences/{id} [delete]
func (h *SentenceHandler) DeleteSentence(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete sentence", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
