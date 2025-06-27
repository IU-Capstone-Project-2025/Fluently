package utils

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"fluently/go-backend/internal/middleware"
	"fluently/go-backend/internal/repository/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

// ParseUUIDParam parses UUID from URL param
func ParseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, param)
	return uuid.Parse(idStr)
}

// ParseIntParam extracts an integer from URL parameters
func ParseIntParam(r *http.Request, param string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, param))
}

// GetCurrentUser retrieves the current authenticated user from context
// This works with the new go-chi/jwtauth system
func GetCurrentUser(ctx context.Context) (*models.User, error) {
	// Try the custom middleware approach first (includes full user info)
	if user := middleware.GetUserFromContext(ctx); user != nil {
		return user, nil
	}

	// Fallback to JWT claims (minimal user info)
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, errors.New("no authentication context")
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid UUID in token")
	}

	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return &models.User{
		ID:    userID,
		Email: email,
		Role:  role,
	}, nil
}

// GetUserClaims gets the JWT claims for the current user
func GetUserClaims(ctx context.Context) (map[string]interface{}, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	return claims, err
}
