package item

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/httpx"
	"github.com/blueetfm/wishlist-app/backend/internal/transport/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/wishlist"
)

// Handler exposes HTTP handlers for wishlist item CRUD operations.
type Handler struct {
	service   *Service
	wishlists *wishlist.Service
}

// NewHandler constructs a Handler backed by the given services.
func NewHandler(service *Service, wishlists *wishlist.Service) *Handler {
	return &Handler{service: service, wishlists: wishlists}
}

// Get handles GET /api/v1/wishlists/{wishlistId}/items/{itemId}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	item, err := h.service.GetOwned(r.Context(), userID, r.PathValue("wishlistId"), r.PathValue("itemId"))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, item)
}

// Create handles POST /api/v1/wishlists/{wishlistId}/items.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	wishlistID := r.PathValue("wishlistId")
	if _, err := h.wishlists.GetOwned(r.Context(), userID, wishlistID); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		httpx.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	item, err := h.service.Create(r.Context(), wishlistID, req)
	if err != nil {
		log.Printf("items.Create: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create item")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, item)
}

// Update handles PUT /api/v1/wishlists/{wishlistId}/items/{itemId}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.service.Update(r.Context(), userID, r.PathValue("wishlistId"), r.PathValue("itemId"), req)
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, item)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.service.Delete(r.Context(), userID, r.PathValue("wishlistId"), r.PathValue("itemId")); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
