package handlers

import (
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/api/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// ShareHandler exposes HTTP handlers for managing and resolving wishlist share links.
type ShareHandler struct {
	service *services.ShareService
}

// NewShareHandler constructs a ShareHandler backed by the given service.
func NewShareHandler(service *services.ShareService) *ShareHandler {
	return &ShareHandler{service: service}
}

// Get handles GET /api/v1/wishlists/{wishlistId}/share.
func (h *ShareHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	token, err := h.service.GetActive(r.Context(), userID, r.PathValue("wishlistId"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, token)
}

// Create handles POST /api/v1/wishlists/{wishlistId}/share.
func (h *ShareHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	token, err := h.service.CreateOrGet(r.Context(), userID, r.PathValue("wishlistId"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, token)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}/share.
func (h *ShareHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.service.Revoke(r.Context(), userID, r.PathValue("wishlistId")); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Resolve handles GET /api/v1/shared/{shareToken}. Optional auth as required for ShareService's ResolvePublic
// An authenticated caller (if any) is passed through so their claims still show up as "mine",
// but an anonymous guest won't be able to claim
func (h *ShareHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	viewerID, _ := middleware.UserID(r.Context())
	detail, err := h.service.ResolvePublic(r.Context(), r.PathValue("shareToken"), viewerID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, detail)
}
