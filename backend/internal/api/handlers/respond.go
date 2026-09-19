package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// writeJSON writes payload as a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

type errorResponse struct {
	Error string `json:"error"`
}

// writeError writes a JSON {"error": message} response with the given status code.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

// writeServiceError maps known service-layer sentinel errors to HTTP status codes.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		writeError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, services.ErrForbidden):
		writeError(w, http.StatusForbidden, "you do not have access to this resource")
	case errors.Is(err, services.ErrConflict):
		writeError(w, http.StatusConflict, "resource already exists")
	case errors.Is(err, services.ErrInvalidURL):
		writeError(w, http.StatusBadRequest, "invalid or disallowed url")
	case errors.Is(err, services.ErrUpstreamUnavailable):
		writeError(w, http.StatusBadGateway, "failed to fetch the target url")
	case errors.Is(err, services.ErrNoMetadata):
		writeError(w, http.StatusUnprocessableEntity, "no open graph metadata found for this url")
	default:
		log.Printf("unhandled service error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
