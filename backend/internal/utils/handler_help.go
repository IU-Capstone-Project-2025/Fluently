package utils

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ParseUUIDParam parses UUID from URL param
func ParseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, param)
	return uuid.Parse(idStr)
}
