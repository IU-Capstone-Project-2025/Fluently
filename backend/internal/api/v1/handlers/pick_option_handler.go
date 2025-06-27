package handlers

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"

	"github.com/google/uuid"
)

type PickOptionHandler struct {
	Repo *postgres.PickOptionRepository
}

func buildPickOptionResponse(option *models.PickOption) schemas.PickOptionResponse {
	return schemas.PickOptionResponse{
		ID:         option.ID.String(),
		WordID:     option.WordID.String(),
		SentenceID: option.SentenceID.String(),
		Options:    option.Option,
	}
}

// CreatePickOption godoc
// @Summary      Create a pick option
// @Description  Creates a new pick option for a word and sentence
// @Tags         pick-options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        option  body      schemas.CreatePickOptionRequest  true  "Pick option data"
// @Success      201  {object}  schemas.PickOptionResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/pick-options/ [post]
func (h *PickOptionHandler) CreatePickOption(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreatePickOptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Options) == 0 {
		http.Error(w, "option cannot be empty", http.StatusBadRequest)
		return
	}

	wordID, err := uuid.Parse(req.WordID)
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	sentenceID, err := uuid.Parse(req.SentenceID)
	if err != nil {
		http.Error(w, "invalid sentence_id", http.StatusBadRequest)
		return
	}

	option := models.PickOption{
		ID:         uuid.New(),
		WordID:     wordID,
		SentenceID: sentenceID,
		Option:     req.Options,
	}

	if err := h.Repo.Create(r.Context(), &option); err != nil {
		http.Error(w, "failed to create option", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildPickOptionResponse(&option))
}

// UpdatePickOption godoc
// @Summary      Update a pick option
// @Description  Updates an existing pick option by ID
// @Tags         pick-options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                          true  "Pick option ID"
// @Param        option  body      schemas.CreatePickOptionRequest true  "Pick option data"
// @Success      200  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/pick-options/{id} [put]
func (h *PickOptionHandler) UpdatePickOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req schemas.CreatePickOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	option, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	}

	option.Option = req.Options

	if err := h.Repo.Update(r.Context(), option); err != nil {
		http.Error(w, "failed to update topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeletePickOption godoc
// @Summary      Delete a pick option
// @Description  Deletes a pick option by ID
// @Tags         pick-options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pick option ID"
// @Success      204  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/pick-options/{id} [delete]
func (h *PickOptionHandler) DeletePickOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete option", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPickOption godoc
// @Summary      Get pick option by ID
// @Description  Returns a pick option by its unique identifier
// @Tags         pick-options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Pick option ID"
// @Success      200  {object}  schemas.PickOptionResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Router       /api/v1/pick-options/{id} [get]
func (h *PickOptionHandler) GetPickOption(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	option, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildPickOptionResponse(option))
}

// ListPickOptions godoc
// @Summary      Get pick options for a word
// @Description  Returns all pick options for the specified word
// @Tags         pick-options
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        word_id   path      string  true  "Word ID"
// @Success      200  {array}   schemas.PickOptionResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /api/v1/words/{word_id}/pick-options [get]
func (h *PickOptionHandler) ListPickOptions(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseUUIDParam(r, "word_id")
	if err != nil {
		http.Error(w, "invalid word_id", http.StatusBadRequest)
		return
	}

	options, err := h.Repo.ListByWordID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch options", http.StatusBadRequest)
		return
	}

	var resp []schemas.PickOptionResponse
	for _, o := range options {
		resp = append(resp, buildPickOptionResponse(&o))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
