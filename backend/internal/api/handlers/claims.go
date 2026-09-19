package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/api/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/models"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// ClaimHandler exposes HTTP handlers for claim CRUD operations.
type ClaimHandler struct {
	service     *services.ClaimService
	itemService *services.ItemService
}

// NewClaimHandler constructs a ClaimHandler backed by the given services.
func NewClaimHandler(service *services.ClaimService, itemService *services.ItemService) *ClaimHandler {
	return &ClaimHandler{service: service, itemService: itemService}
}

// Post handles POST /api/v1/wishlists/{wishlistId}/items/{itemId}/claim.
func (h *ClaimHandler) Post(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	itemID := r.PathValue("itemId")
	if err := h.itemService.ExistsInWishlist(r.Context(), r.PathValue("wishlistId"), itemID); err != nil {
		writeServiceError(w, err)
		return
	}

	var req models.ClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claim, err := h.service.Create(r.Context(), itemID, userID, req.Message)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	claim.IsMine = true
	writeJSON(w, http.StatusOK, claim)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}/claim.
func (h *ClaimHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	itemID := r.PathValue("itemId")
	if err := h.itemService.ExistsInWishlist(r.Context(), r.PathValue("wishlistId"), itemID); err != nil {
		writeServiceError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), itemID, userID); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
