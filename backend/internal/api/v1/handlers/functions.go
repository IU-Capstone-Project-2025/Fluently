package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// parseUUIDParam parses UUID from URL param
func parseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	idStr := chi.URLParam(r, param)
	return uuid.Parse(idStr)
}
