package item

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/httpx"
	"github.com/blueetfm/wishlist-app/backend/internal/transport/middleware"
)

// SplitHandler exposes HTTP handlers for managing split interest.
type SplitHandler struct {
	service *SplitService
}

// NewSplitHandler constructs a new SplitHandler backed by the given SplitService.
func NewSplitHandler(service *SplitService) *SplitHandler {
	return &SplitHandler{service: service}
}

// Create handles POST /api/v1/wishlists/{wishlistId}/items/{itemId}/split.
func (h *SplitHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req SplitInterestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.SocialHandle == "" {
		httpx.WriteError(w, http.StatusBadRequest, "social handle is required")
		return
	}

	split, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		log.Printf("split.Create: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create split interest")
		return
	}
	split.IsMine = true
	httpx.WriteJSON(w, http.StatusCreated, split)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}/split.
func (h *SplitHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.service.Delete(r.Context(), r.PathValue("itemId"), userID); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
