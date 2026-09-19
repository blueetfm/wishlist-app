package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/api/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/models"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// SplitHandler exposes HTTP handlers for managing split interest 
type SplitHandler struct {
	service *services.SplitService
}

// NewSplitHandler constructs a new SplitHandler backed by the SplitService
func NewSplitHandler(service *services.SplitService) *SplitHandler {
	return &SplitHandler{service: service}
}

// Create handles POST /api/v1/wishlists/{wishlistId}/items/{itemId}/split.
func (h *SplitHandler) Create (w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req models.SplitInterestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.SocialHandle == "" {
		writeError(w, http.StatusBadRequest, "social handle is required")
		return
	}

	split, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		log.Printf("split.Create: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create split interest")
		return
	}
	split.IsMine = true
	writeJSON(w, http.StatusCreated, split)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}/split.
func (h *SplitHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.service.Delete(r.Context(), r.PathValue("itemId"), userID); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

