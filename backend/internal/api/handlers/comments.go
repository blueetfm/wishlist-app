package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/api/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/models"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// CommentHandler exposes HTTP handlers for comment CRUD operations.
type CommentHandler struct {
	service     *services.CommentService
	itemService *services.ItemService
}

// NewCommentHandler construct a new handler backed by the corresponding comment service
func NewCommentHandler(service *services.CommentService, itemService *services.ItemService) *CommentHandler {
	return &CommentHandler{service: service, itemService: itemService}
}

// Create handles POST/comments
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	comment, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		log.Printf("comments.Create: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create comment")
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

// GetWishlistComments handles GET /api/v1/wishlists/{wishlistId}/comments.
func (h *CommentHandler) GetWishlistComments(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	comments, err := h.service.ListByWishlistID(r.Context(), r.PathValue("wishlistId"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

// GetItemComments handles GET /api/v1/wishlists/{wishlistId}/items/{itemId}/comments
func (h *CommentHandler) GetItemComments(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.itemService.ExistsInWishlist(r.Context(), r.PathValue("wishlistId"), r.PathValue("itemId")); err != nil {
		writeServiceError(w, err)
		return
	}

	comments, err := h.service.ListByItemID(r.Context(), r.PathValue("itemId"))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, comments)

}

// Update handles PUT/comments/{commentId}
func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing unauthenticated user")
		return
	}

	var req models.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	comment, err := h.service.Update(r.Context(), userID, r.PathValue("commentId"), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, comment)
}
