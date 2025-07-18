package handlers

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
	"strconv"
	"time"
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
// @Security BearerAuth
// @Param request body DistractorRequest true "Данные для генерации"
// @Success 200 {object} DistractorResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /api/v1/distractors [post]
func (h *DistractorHandler) Generate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/distractors"
	method := r.Method
	statusCode := 200
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()
	ctx := r.Context()
	var req DistractorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		statusCode = 400
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Sentence == "" || req.Word == "" {
		statusCode = 400
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	picks, err := h.Client.GenerateDistractors(ctx, req.Sentence, req.Word)
	if err != nil {
		statusCode = 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DistractorResponse{PickOptions: picks})
}
