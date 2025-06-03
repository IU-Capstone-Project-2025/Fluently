package handler

import (
	"encoding/json"
	"net/http"
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
}
