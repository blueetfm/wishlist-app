package models

import "time"

// Item matches the Item schema in docs/swagger.yaml.
type Item struct {
	ID             string               `json:"id"`
	WishlistID     string               `json:"wishlist_id"`
	Name           string               `json:"name"`
	ProductURL     string               `json:"product_url"`
	Price          float64              `json:"price"`
	EmbedData      map[string]any       `json:"embed_data,omitempty"`
	Claim          *Claim               `json:"claim"`
	SplitInterests []SplitInterestEntry `json:"split_interests"`
	CreatedAt      time.Time            `json:"created_at"`
}

// CreateItemRequest is the payload for POST /api/v1/wishlists/{wishlistId}/items.
type CreateItemRequest struct {
	Name       string         `json:"name"`
	ProductURL string         `json:"product_url"`
	Price      float64        `json:"price"`
	EmbedData  map[string]any `json:"embed_data,omitempty"`
}

// UpdateItemRequest is the payload for PUT /api/v1/items/{itemId}.
type UpdateItemRequest struct {
	Name       *string        `json:"name"`
	ProductURL *string        `json:"product_url"`
	Price      *float64       `json:"price"`
	EmbedData  map[string]any `json:"embed_data,omitempty"`
}
