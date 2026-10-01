package comment

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/httpx"
	"github.com/blueetfm/wishlist-app/backend/internal/transport/middleware"
)

// ItemChecker verifies an item belongs to a wishlist before item-scoped
// comments can be listed. Defined here (rather than depending on the item
// package directly) since item already depends on comment's sibling claim
// package, so keeping this local avoids any risk of an import cycle.
type ItemChecker interface {
	ExistsInWishlist(ctx context.Context, wishlistID, itemID string) error
}

// Handler exposes HTTP handlers for comment CRUD operations.
type Handler struct {
	service *Service
	items   ItemChecker
}

// NewHandler constructs a new Handler backed by the given service and item checker.
func NewHandler(service *Service, items ItemChecker) *Handler {
	return &Handler{service: service, items: items}
}

// Create handles POST /api/v1/comments.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		httpx.WriteError(w, http.StatusBadRequest, "content is required")
		return
	}

	comment, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		log.Printf("comments.Create: %v", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create comment")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, comment)
}

// GetWishlistComments handles GET /api/v1/wishlists/{wishlistId}/comments.
func (h *Handler) GetWishlistComments(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	comments, err := h.service.ListByWishlistID(r.Context(), r.PathValue("wishlistId"))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, comments)
}

// GetItemComments handles GET /api/v1/wishlists/{wishlistId}/items/{itemId}/comments
func (h *Handler) GetItemComments(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	if err := h.items.ExistsInWishlist(r.Context(), r.PathValue("wishlistId"), r.PathValue("itemId")); err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	comments, err := h.service.ListByItemID(r.Context(), r.PathValue("itemId"))
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, comments)
}

// Update handles PUT /api/v1/comments/{commentId}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "missing unauthenticated user")
		return
	}

	var req UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	comment, err := h.service.Update(r.Context(), userID, r.PathValue("commentId"), req)
	if err != nil {
		httpx.WriteServiceError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, comment)
}
