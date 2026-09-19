package models

// Share reports whether a wishlist's share_token currently resolves via the
// public /shared/{shareToken} endpoint.
type Share struct {
	Token    string `json:"token"`
	IsActive bool   `json:"is_active"`
}
