package handler

import (
	"encoding/json"
	"net/http"

	"fluently/go-backend/internal/repository/models"
	"fluently/go-backend/internal/repository/postgres"
	"fluently/go-backend/internal/repository/schemas"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	Repo *postgres.UserRepository
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:       uuid.New(),
		Name:     req.Name,
		SubLevel: req.SubLevel,
		PrefID:   req.PrefID,
	}

	if err := h.Repo.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	resp := schemas.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		SubLevel: user.SubLevel,
		Pref: &schemas.PreferenceMini{
			ID:        user.Pref.ID,
			CEFRLevel: user.Pref.CEFRLevel,
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	resp := schemas.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		SubLevel: user.SubLevel,
		Pref: &schemas.PreferenceMini{
			ID:        user.Pref.ID,
			CEFRLevel: user.Pref.CEFRLevel,
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	var req schemas.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	user.Name = req.Name
	user.SubLevel = req.SubLevel
	user.PrefID = req.PrefID

	if err := h.Repo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	resp := schemas.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		SubLevel: user.SubLevel,
		Pref: &schemas.PreferenceMini{
			ID:        user.Pref.ID,
			CEFRLevel: user.Pref.CEFRLevel,
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, "failed to delte user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
