package share

import (
	"github.com/blueetfm/wishlist-app/backend/internal/item"
	"github.com/blueetfm/wishlist-app/backend/internal/wishlist"
)

// Share reports whether a wishlist's share_token currently resolves via the
// public /shared/{shareToken} endpoint.
type Share struct {
	Token    string `json:"token"`
	IsActive bool   `json:"is_active"`
}

// WishlistDetail is a wishlist plus its items, returned by the public
// /shared/{shareToken} endpoint.
type WishlistDetail struct {
	wishlist.Wishlist
	Items []item.Item `json:"items"`
}
