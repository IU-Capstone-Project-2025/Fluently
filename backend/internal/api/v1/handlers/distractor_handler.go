package handlers

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
)

// swagger:ignore – искусственное использование, чтобы пакет schemas считался используемым
var _ schemas.ErrorResponse

type DistractorRequest struct {
	Sentence string `json:"sentence"`
	Word     string `json:"word"`
}

type DistractorResponse struct {
	PickOptions []string `json:"pick_options"`
}

type DistractorHandler struct {
	Client *utils.DistractorClient
}

// GenerateDistractors godoc
// @Summary Генерация дистракторов (вариантов ответа)
// @Description Отправляет предложение и слово во внутренний ML-сервис и получает список подходящих Distractor-вариантов.
// @Tags Distractor
// @Accept json
// @Produce json
// @Param request body DistractorRequest true "Данные для генерации"
// @Success 200 {object} DistractorResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /api/v1/distractors [post]
func (h *DistractorHandler) Generate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req DistractorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	picks, err := h.Client.GenerateDistractors(ctx, req.Sentence, req.Word)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(DistractorResponse{PickOptions: picks})
}
