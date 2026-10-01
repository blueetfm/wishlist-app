package share

import (
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/httpx"
	"github.com/blueetfm/wishlist-app/backend/internal/transport/middleware"
)

// Handler exposes HTTP handlers for managing and resolving wishlist share links.
type Handler struct {
	service *Service
}

// NewHandler constructs a Handler backed by the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Get handles GET /api/v1/wishlists/{wishlistId}/share.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	token, err := h.service.GetActive(r.Context(), userID, r.PathValue("wishlistId"))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, token)
}

// Create handles POST /api/v1/wishlists/{wishlistId}/share.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	token, err := h.service.CreateOrGet(r.Context(), userID, r.PathValue("wishlistId"))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, token)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}/share.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.service.Revoke(r.Context(), userID, r.PathValue("wishlistId")); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Resolve handles GET /api/v1/shared/{shareToken}. Optional auth as required for Service's ResolvePublic
// An authenticated caller (if any) is passed through so their claims still show up as "mine",
// but an anonymous guest won't be able to claim
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	viewerID, _ := middleware.UserID(r.Context())
	detail, err := h.service.ResolvePublic(r.Context(), r.PathValue("shareToken"), viewerID)
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, detail)
}
