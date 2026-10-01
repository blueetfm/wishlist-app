package item

import (
	"time"
)

// Item matches the Item schema in docs/swagger.yaml.
type Item struct {
	ID          string         `json:"id"`
	WishlistID  string         `json:"wishlist_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ProductURL  string         `json:"product_url"`
	Price       float64        `json:"price"`
	EmbedData   map[string]any `json:"embed_data,omitempty"`
	// ImageURL is a user-uploaded photo, taking precedence over any image in
	// EmbedData. Nil means the client should fall back to embed_data.image_url
	// or a placeholder.
	ImageURL       *string              `json:"image_url,omitempty"`
	Claim          *Claim               `json:"claim"`
	SplitInterests []SplitInterestEntry `json:"split_interests"`
	CreatedAt      time.Time            `json:"created_at"`
}

// CreateItemRequest is the payload for POST /api/v1/wishlists/{wishlistId}/items.
type CreateItemRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ProductURL  string         `json:"product_url"`
	Price       float64        `json:"price"`
	EmbedData   map[string]any `json:"embed_data,omitempty"`
	ImageURL    *string        `json:"image_url,omitempty"`
}

// UpdateItemRequest is the payload for PUT /api/v1/wishlists/{wishlistId}/items/{itemId}.
type UpdateItemRequest struct {
	Name        *string        `json:"name"`
	Description *string        `json:"description"`
	ProductURL  *string        `json:"product_url"`
	Price       *float64       `json:"price"`
	EmbedData   map[string]any `json:"embed_data,omitempty"`
	ImageURL    *string        `json:"image_url,omitempty"`
}
