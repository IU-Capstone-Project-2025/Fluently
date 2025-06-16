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

type UserHandler struct {
	Repo *postgres.UserRepository
}

// buildUserResponse returns user response
func buildUserResponse(user *models.User) schemas.UserResponse {
	return schemas.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		PrefID:    *user.PrefID,
		CreatedAt: user.CreatedAt,
	}
}

// CreateUser godoc
// @Summary      Create a user
// @Description  Registers a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      schemas.CreateUserRequest  true  "User data"
// @Success      201  {object}  schemas.UserResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/ [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Name,
		PasswordHash: req.PasswordHash,
		Provider:     req.Provider,
		GoogleID:     req.GoogleID,
		Role:         req.Role,
		IsActive:     req.IsActive,
		PrefID:       &req.PrefID,
	}

	if err := h.Repo.Create(r.Context(), &user); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(buildUserResponse(&user))
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Returns a user by their unique identifier
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  schemas.UserResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(buildUserResponse(user))
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Updates user data by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "User ID"
// @Param        user  body      schemas.CreateUserRequest  true  "User data"
// @Success      200  {object}  schemas.UserResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
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
	user.Email = req.Email
	user.Role = req.Role
	user.IsActive = req.IsActive
	user.Provider = req.Provider
	user.GoogleID = req.GoogleID
	user.PasswordHash = req.PasswordHash
	user.PrefID = &req.PrefID

	if err := h.Repo.Update(r.Context(), user); err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(buildUserResponse(user))
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Deletes a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      204  ""
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      404  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUIDParam(r, "id")
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delte user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GoogleAuthHandler godoc
// @Summary      Authenticates with Google
// @Description  Authenticates with Google using the authorization code flow
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        code query     string  true  "Authorization code"
// @Success      200  {object}  schemas.JwtResponse
// @Failure      400  {object}  schemas.ErrorResponse
// @Failure      401  {object}  schemas.ErrorResponse
// @Failure      500  {object}  schemas.ErrorResponse
// @Router       /auth/google [get]
// func (h *UserHandler) GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
// Getting idtoken and platform (for audience) from request
// 	var req struct{
// 		IDToken  string `json:"id_token"`
// 		Platform string `json:"platform"`
// 	}

// 	googleToken := req.IDToken
// 	platform := req.Platform

// 	var googleClientIDs = map[string]string{
// 		"ios":     fmt.Sprintf("%s.apps.googleusercontent.com", appConfig.GetIosGoogleClientID()),
// 		"android": fmt.Sprintf("%s.apps.googleusercontent.com", appConfig.GetAndroidGoogleClientID()),
// 		"web":     appConfig.GetWebGoogleClientID(),
// 	}

// 	// В момент логина — клиент передаёт "platform": "ios" | "android" | "web"
// 	audience := googleClientIDs[platform]

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		logger.Log.Error("Invalid request", zap.Error(err))
// 		http.Error(w, "invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	payload, err := idtoken.Validate(r.Context(), googleToken, audience)

// 	if err != nil {
// 		logger.Log.Error("Invalid token", zap.Error(err))
// 		http.Error(w, "invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	claims := payload.Claims
// 	sub := claims["sub"].(string)
// 	email := claims["email"].(string)
// 	name := claims["name"].(string)
// 	picture := claims["picture"].(string)

// 	// Check if user exists
// 	var user models.User
// 	var user_preferences models.Preference

// 	if err := h.Repo.GetByEmail(r.Context(), email); err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			logger.Log.Info("Creating new user with email: ", zap.String("email", email))
// 			user = models.User{
// 				Email:    email,
// 				Name:     name,
// 				SubLevel: 0,
// 				PrefID:   1,
// 			}
// 			h.Repo.Create(r.Context(), )
// 		}
// 		logger.Log.Error("Failed to get user", zap.Error(err))
// 		http.Error(w, "failed to get user", http.StatusInternalServerError)
// 		return
// 	}
// 		logger.Log.Error("User not found", zap.Error(err))
// 		http.Error(w, "user not found", http.StatusNotFound)
// 		return
// 	}

// }
