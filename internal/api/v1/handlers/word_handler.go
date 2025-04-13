package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/repository/service"
)

type WordHandler struct {
    service *service.WordService
}

func NewWordHandler(service *service.WordService) *WordHandler {
    return &WordHandler{service: service}
}

var validate = validator.New()

func ListWords(w http.ResponseWriter, r *http.Request) {
    // TODO: Получить список слов с фильтрами
}

func GetWord(w http.ResponseWriter, r *http.Request) {
    // TODO: Получить слово по ID
}

func (h *WordHandler) CreateWord(w http.ResponseWriter, r *http.Request) {
    // TODO: Добавить новое слово
    var req schemas.WordCreateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if err := validate.Struct(&req); err != nil {
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.Create(r.Context(), &req); err != nil {
        http.Error(w, "Failed to create word: "+err.Error(), http.StatusInternalServerError)
    }

    w.WriteHeader(http.StatusCreated)
}

func UpdateWord(w http.ResponseWriter, r *http.Request) {
    // TODO: Обновить слово
}

func DeleteWord(w http.ResponseWriter, r *http.Request) {
    // TODO: Удалить слово
}
