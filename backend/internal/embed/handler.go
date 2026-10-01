package embed

import (
	"encoding/json"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/httpx"
)

// Handler exposes the HTTP handler for scraping Open Graph metadata.
type Handler struct {
	service *Service
}

// NewHandler constructs a Handler backed by the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /api/v1/embed.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req EmbedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.URL == "" {
		httpx.WriteError(w, http.StatusBadRequest, "url is required")
		return
	}

	embed, err := h.service.Fetch(r.Context(), req.URL)
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, embed)
}
