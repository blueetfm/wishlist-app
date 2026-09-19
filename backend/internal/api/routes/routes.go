package routes

import (
	"net/http"

	"github.com/MicahParks/keyfunc/v3"

	"github.com/blueetfm/wishlist-app/backend/internal/api/handlers"
	"github.com/blueetfm/wishlist-app/backend/internal/api/middleware"
	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// Dependencies holds everything the router needs to construct handlers and middleware
type Dependencies struct {
	WishlistService    *services.WishlistService
	ItemService        *services.ItemService
	ClaimService       *services.ClaimService
	SplitService       *services.SplitService
	CommentService     *services.CommentService
	EmbedService       *services.EmbedService
	ShareService       *services.ShareService
	IdempotencyService *services.IdempotencyService
	JWKS               keyfunc.Keyfunc
}

// NewRouter builds the HTTP mux with all API routes mounted.
func NewRouter(deps Dependencies) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)

	auth := middleware.RequireAuth(deps.JWKS) // middleware function that takes in a http.Handler and returns a http.Handler
	optionalAuth := middleware.OptionalAuth(deps.JWKS)
	idempotent := middleware.RequireIdempotencyKey(deps.IdempotencyService)
	wishlistHandler := handlers.NewWishlistHandler(deps.WishlistService)
	itemHandler := handlers.NewItemHandler(deps.ItemService, deps.WishlistService)
	claimHandler := handlers.NewClaimHandler(deps.ClaimService, deps.ItemService)
	splitHandler := handlers.NewSplitHandler(deps.SplitService)
	commentHandler := handlers.NewCommentHandler(deps.CommentService, deps.ItemService)
	embedHandler := handlers.NewEmbedHandler(deps.EmbedService)
	shareHandler := handlers.NewShareHandler(deps.ShareService)

	// 1. Wishlist services
	mux.Handle("POST /api/v1/wishlists", auth(idempotent(http.HandlerFunc(wishlistHandler.Create))))
	mux.Handle("GET /api/v1/wishlists/me", auth(http.HandlerFunc(wishlistHandler.ListMine)))
	mux.Handle("GET /api/v1/wishlists/{wishlistId}", auth(http.HandlerFunc(wishlistHandler.Get)))
	mux.Handle("PUT /api/v1/wishlists/{wishlistId}", auth(http.HandlerFunc(wishlistHandler.Update)))
	mux.Handle("DELETE /api/v1/wishlists/{wishlistId}", auth(http.HandlerFunc(wishlistHandler.Delete)))

	// 2. Wishlist item services
	mux.Handle("POST /api/v1/wishlists/{wishlistId}/items", auth(idempotent(http.HandlerFunc(itemHandler.Create))))
	mux.Handle("GET /api/v1/wishlists/{wishlistId}/items/{itemId}", auth(http.HandlerFunc(itemHandler.Get)))
	mux.Handle("PUT /api/v1/wishlists/{wishlistId}/items/{itemId}", auth(http.HandlerFunc(itemHandler.Update)))
	mux.Handle("DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}", auth(http.HandlerFunc(itemHandler.Delete)))

	// 3. Claim services
	mux.Handle("POST /api/v1/wishlists/{wishlistId}/items/{itemId}/claim", auth(idempotent(http.HandlerFunc(claimHandler.Post))))
	mux.Handle("DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}/claim", auth(http.HandlerFunc(claimHandler.Delete)))

	// 4. Split interest services
	mux.Handle("POST /api/v1/wishlists/{wishlistId}/items/{itemId}/split", auth(idempotent(http.HandlerFunc(splitHandler.Create))))
	mux.Handle("DELETE /api/v1/wishlists/{wishlistId}/items/{itemId}/split", auth(http.HandlerFunc(splitHandler.Delete)))

	// 5. Comment services
	mux.Handle("POST /api/v1/comments", auth(idempotent(http.HandlerFunc(commentHandler.Create))))
	mux.Handle("GET /api/v1/wishlists/{wishlistId}/items/{itemId}/comments", auth(http.HandlerFunc(commentHandler.GetItemComments)))
	mux.Handle("GET /api/v1/wishlists/{wishlistId}/comments", auth(http.HandlerFunc(commentHandler.GetWishlistComments)))
	mux.Handle("PUT /api/v1/comments/{commentId}", auth(http.HandlerFunc(commentHandler.Update)))

	// 6. Embed service
	mux.Handle("POST /api/v1/embed", auth(idempotent(http.HandlerFunc(embedHandler.Create))))

	// 7. Share token service
	// a. owner-only
	mux.Handle("GET /api/v1/wishlists/{wishlistId}/share", auth(http.HandlerFunc(shareHandler.Get)))
	mux.Handle("POST /api/v1/wishlists/{wishlistId}/share", auth(idempotent(http.HandlerFunc(shareHandler.Create))))
	mux.Handle("DELETE /api/v1/wishlists/{wishlistId}/share", auth(http.HandlerFunc(shareHandler.Delete)))

	// b. guest: optional auth
	// // public share link resolves
	// // but an authenticated user is recognized so they can keep claiming normally.
	mux.Handle("GET /api/v1/shared/{shareToken}", optionalAuth(http.HandlerFunc(shareHandler.Resolve)))
	return mux
}
