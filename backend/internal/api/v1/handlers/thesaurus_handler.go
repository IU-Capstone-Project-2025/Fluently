package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/utils"
)

// swagger:ignore
var _ schemas.ErrorResponse

// LearnedWordPayload represents minimal data sent by client
// We only care about the English `word` field.
type LearnedWordPayload struct {
	Word string `json:"word"`
}

type ThesaurusHandler struct {
	Client *utils.ThesaurusClient
}

// RecommendWords godoc
// @Summary Рекомендации словаря (Thesaurus)
// @Description Принимает список выученных пользователем слов и возвращает новые рекомендации для изучения.
// @Tags Thesaurus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body []LearnedWordPayload true "Массив выученных слов пользователя"
// @Success 200 {array} utils.ThesaurusRecommendation
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /api/v1/thesaurus/recommend [post]
func (h *ThesaurusHandler) Recommend(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	endpoint := "/api/v1/thesaurus/recommend"
	method := r.Method
	statusCode := 200
	defer func() {
		httpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(time.Since(start).Seconds())
	}()
	ctx := r.Context()
	var payload []LearnedWordPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		statusCode = 400
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var words []string
	for _, lw := range payload {
		if lw.Word != "" {
			words = append(words, lw.Word)
		}
	}
	if len(words) == 0 {
		statusCode = 400
		http.Error(w, "no words provided", http.StatusBadRequest)
		return
	}

	recs, err := h.Client.Recommend(ctx, words)
	if err != nil {
		statusCode = 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recs)
}
