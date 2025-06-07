package handler

import (
	"encoding/json"
	"net/http"
<<<<<<< HEAD

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

	if err := h.Repo.Create(r.Context(), &user); err != nil {
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

	user, err := h.Repo.GetByID(r.Context(), id)
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

	user, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	user.Name = req.Name
	user.SubLevel = req.SubLevel
	user.PrefID = req.PrefID

	if err := h.Repo.Update(r.Context(), user); err != nil {
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

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delte user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
=======
	"strconv"

	"fluently/go-backend/internal/repository/schemas"
	"fluently/go-backend/internal/repository/service"
)

type UserHandler struct {
    service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req schemas.UserCreateRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if err := validate.Struct(&req); err != nil {
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }

    user, err := h.service.Create(r.Context(), &req)
    if err != nil {
        http.Error(w, "Failled to create user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user, err := h.service.GetByID(r.Context(), uint(id))
    if err != nil {
        http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    if user == nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var req schemas.UserUpdateRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if req.Name == nil {
        http.Error(w, "At least one field must be updated", http.StatusBadRequest)
        return
    }

    if err := h.service.Update(r.Context(), uint(id), &req); err != nil {
        http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User updated successfully"))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    if err := h.service.Delete(r.Context(), uint(id)); err != nil {
        http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User deleted successfully"))
>>>>>>> d67dbcc (Add all user logic)
}
