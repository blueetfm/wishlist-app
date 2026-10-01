package wishlist

import (
	"encoding/json"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/httpx"
	"github.com/blueetfm/wishlist-app/backend/internal/transport/middleware"
)

// Handler exposes HTTP handlers for wishlist CRUD operations.
type Handler struct {
	service *Service
}

// NewHandler constructs a Handler backed by the given service.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles POST /api/v1/wishlists.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		httpx.WriteError(w, http.StatusBadRequest, "title is required")
		return
	}

	wishlist, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create wishlist")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, wishlist)
}

// ListMine handles GET /api/v1/wishlists/me.
func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	wishlists, err := h.service.ListByOwner(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list wishlists")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, wishlists)
}

// Get handles GET /api/v1/wishlists/{wishlistId}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	wishlist, err := h.service.GetOwned(r.Context(), userID, r.PathValue("wishlistId"))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, wishlist)
}

// Update handles PUT /api/v1/wishlists/{wishlistId}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req UpdateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wishlist, err := h.service.Update(r.Context(), userID, r.PathValue("wishlistId"), req)
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, wishlist)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.service.Delete(r.Context(), userID, r.PathValue("wishlistId")); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
