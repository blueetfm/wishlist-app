package routes

import (
	"net/http"

	"github.com/MicahParks/keyfunc/v3"

	"github.com/blueetfm/wishlist-app/backend/internal/api/handlers"
	"github.com/blueetfm/wishlist-app/backend/internal/api/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// Dependencies holds everything the router needs to construct handlers and
// middleware.
type Dependencies struct {
	WishlistService *services.WishlistService
	JWKS            keyfunc.Keyfunc
}

// NewRouter builds the HTTP mux with all API routes mounted.
func NewRouter(deps Dependencies) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)

	auth := middleware.RequireAuth(deps.JWKS)
	wishlistHandler := handlers.NewWishlistHandler(deps.WishlistService)

	mux.Handle("POST /api/v1/wishlists", auth(http.HandlerFunc(wishlistHandler.Create)))
	mux.Handle("GET /api/v1/wishlists/me", auth(http.HandlerFunc(wishlistHandler.ListMine)))
	mux.Handle("GET /api/v1/wishlists/{wishlistId}", auth(http.HandlerFunc(wishlistHandler.Get)))
	mux.Handle("PUT /api/v1/wishlists/{wishlistId}", auth(http.HandlerFunc(wishlistHandler.Update)))
	mux.Handle("DELETE /api/v1/wishlists/{wishlistId}", auth(http.HandlerFunc(wishlistHandler.Delete)))

	// TODO: mount items, claims, split-interest, comments, embed, and shared
	// (public, unauthenticated) routes here as they are implemented.

	return mux
}
