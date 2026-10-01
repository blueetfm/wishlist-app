package item

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/httpx"
	"github.com/blueetfm/wishlist-app/backend/internal/transport/middleware"
)

// ClaimHandler exposes HTTP handlers for claim CRUD operations.
type ClaimHandler struct {
	service *ClaimService
	items   *Service
}

// NewClaimHandler constructs a ClaimHandler backed by the given claim and item services.
func NewClaimHandler(service *ClaimService, items *Service) *ClaimHandler {
	return &ClaimHandler{service: service, items: items}
}

// Post handles POST /api/v1/wishlists/{wishlistId}/items/{itemId}/claim.
func (h *ClaimHandler) Post(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	itemID := r.PathValue("itemId")
	if err := h.items.ExistsInWishlist(r.Context(), r.PathValue("wishlistId"), itemID); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	var req ClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	claim, err := h.service.Create(r.Context(), itemID, userID, req.Message)
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	claim.IsMine = true
	httpx.WriteJSON(w, http.StatusOK, claim)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}/claim.
func (h *ClaimHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	itemID := r.PathValue("itemId")
	if err := h.items.ExistsInWishlist(r.Context(), r.PathValue("wishlistId"), itemID); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), itemID, userID); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
