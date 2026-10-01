// Package httpx provides shared HTTP response helpers used by every domain
// handler, keeping JSON encoding and sentinel-error-to-status mapping in one
// place instead of duplicated per package.
package httpx

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
)

// WriteJSON writes payload as a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

type errorResponse struct {
	Error string `json:"error"`
}

// WriteError writes a JSON {"error": message} response with the given status code.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, errorResponse{Error: message})
}

// WriteServiceError maps known service-layer sentinel errors to HTTP status codes.
func WriteServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperr.ErrNotFound):
		WriteError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, apperr.ErrForbidden):
		WriteError(w, http.StatusForbidden, "you do not have access to this resource")
	case errors.Is(err, apperr.ErrConflict):
		WriteError(w, http.StatusConflict, "resource already exists")
	case errors.Is(err, apperr.ErrInvalidURL):
		WriteError(w, http.StatusBadRequest, "invalid or disallowed url")
	case errors.Is(err, apperr.ErrUpstreamUnavailable):
		WriteError(w, http.StatusBadGateway, "failed to fetch the target url")
	case errors.Is(err, apperr.ErrNoMetadata):
		WriteError(w, http.StatusUnprocessableEntity, "no open graph metadata found for this url")
	default:
		log.Printf("unhandled service error: %v", err)
		WriteError(w, http.StatusInternalServerError, "internal server error")
	}
}
