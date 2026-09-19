package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/models"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// EmbedHandler exposes the HTTP handler for scraping Open Graph metadata.
type EmbedHandler struct {
	service *services.EmbedService
}

// NewEmbedHandler constructs an EmbedHandler backed by the given service.
func NewEmbedHandler(service *services.EmbedService) *EmbedHandler {
	return &EmbedHandler{service: service}
}

// Create handles POST /api/v1/embed.
func (h *EmbedHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.EmbedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	embed, err := h.service.Fetch(r.Context(), req.URL)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, embed)
}
