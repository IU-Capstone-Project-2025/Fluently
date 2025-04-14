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

func (h *WordHandler) ListWords(w http.ResponseWriter, r *http.Request) {
    words, err := h.service.List(r.Context())
    if err != nil {
        http.Error(w, "Failed to list words: "+err.Error(), http.StatusInternalServerError)
        return 
    }

    json.NewEncoder(w).Encode(words)
}

func (h *WordHandler) GetWord(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    word, err := h.service.GetByID(r.Context(), id)
    if err != nil {
        http.Error(w, "Failed to get word: "+err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(word)
}

func (h *WordHandler) CreateWord(w http.ResponseWriter, r *http.Request) {
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

func (h *WordHandler) UpdateWord(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    var req schemas.WordUpdateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if err := validate.Struct(&req); err != nil {
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.Update(r.Context(), id, &req); err != nil {
        http.Error(w, "Failed to update word: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *WordHandler) DeleteWord(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if err := h.service.Delete(r.Context(), id); err != nil {
        http.Error(w, "Failed to delete word: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
