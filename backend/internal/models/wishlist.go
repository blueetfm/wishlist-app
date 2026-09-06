package models

import "time"

// Wishlist matches the Wishlist schema in docs/swagger.yaml.
type Wishlist struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ShareToken  string    `json:"share_token"`
	CreatedAt   time.Time `json:"created_at"`
}

// WishlistDetail matches the WishlistDetail schema: a Wishlist plus its items.
type WishlistDetail struct {
	Wishlist
	Items []Item `json:"items"`
}

// CreateWishlistRequest is the payload for POST /api/v1/wishlists.
type CreateWishlistRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateWishlistRequest is the payload for PUT /api/v1/wishlists/{wishlistId}.
// Pointer fields are optional; nil means "leave unchanged".
type UpdateWishlistRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
