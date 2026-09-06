package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/api/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/models"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// WishlistHandler exposes HTTP handlers for wishlist CRUD operations.
type WishlistHandler struct {
	service *services.WishlistService
}

// NewWishlistHandler constructs a WishlistHandler backed by the given service.
func NewWishlistHandler(service *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{service: service}
}

// Create handles POST /api/v1/wishlists.
func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req models.CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	wishlist, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		log.Printf("wishlists.Create: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create wishlist")
		return
	}
	writeJSON(w, http.StatusCreated, wishlist)
}

// ListMine handles GET /api/v1/wishlists/me.
func (h *WishlistHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	wishlists, err := h.service.ListByOwner(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list wishlists")
		return
	}
	writeJSON(w, http.StatusOK, wishlists)
}

// Get handles GET /api/v1/wishlists/{wishlistId}.
func (h *WishlistHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	wishlist, err := h.service.GetOwned(r.Context(), userID, r.PathValue("wishlistId"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, wishlist)
}

// Update handles PUT /api/v1/wishlists/{wishlistId}.
func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req models.UpdateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wishlist, err := h.service.Update(r.Context(), userID, r.PathValue("wishlistId"), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, wishlist)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}.
func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.service.Delete(r.Context(), userID, r.PathValue("wishlistId")); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeServiceError maps known service-layer sentinel errors to HTTP status codes.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		writeError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, services.ErrForbidden):
		writeError(w, http.StatusForbidden, "you do not have access to this resource")
	case errors.Is(err, services.ErrConflict):
		writeError(w, http.StatusConflict, "resource already exists")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
