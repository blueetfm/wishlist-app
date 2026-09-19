package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/api/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/models"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// ItemHandler exposes HTTP handlers for wishlist CRUD operations.
type ItemHandler struct {
	service         *services.ItemService
	wishlistService *services.WishlistService
}

// NewItemHandler constructs a ItemHandler backed by the given services.
func NewItemHandler(service *services.ItemService, wishlistService *services.WishlistService) *ItemHandler {
	return &ItemHandler{service: service, wishlistService: wishlistService}
}

// Get handles GET /api/v1/wishlists/{wishlistId}/items/{itemId}.
func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	item, err := h.service.GetOwned(r.Context(), userID, r.PathValue("wishlistId"), r.PathValue("itemId"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// Create handles POST /api/v1/wishlists/{wishlistId}/items.
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	wishlistID := r.PathValue("wishlistId")
	if _, err := h.wishlistService.GetOwned(r.Context(), userID, wishlistID); err != nil {
		writeServiceError(w, err)
		return
	}

	var req models.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	item, err := h.service.Create(r.Context(), wishlistID, req)
	if err != nil {
		log.Printf("items.Create: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create item")
		return
	}
	writeJSON(w, http.StatusCreated, item)
}


// Update handles PUT /api/v1/wishlists/{wishlistId}/items/{itemId}.
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req models.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.service.Update(r.Context(), userID, r.PathValue("wishlistId"), r.PathValue("itemId"), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// Delete handles DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}.
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.service.Delete(r.Context(), userID, r.PathValue("wishlistId"), r.PathValue("itemId")); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
