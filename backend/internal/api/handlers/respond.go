package handlers

import (
	"encoding/json"
	"net/http"
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
