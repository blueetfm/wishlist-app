package models

import "time"

// SplitInterestEntry matches the SplitInterestEntry schema in
// docs/swagger.yaml: one user's plain-text handle in the "split the cost"
// pool for an item. No fractional amounts are tracked.
type SplitInterestEntry struct {
	ItemID       string    `json:"item_id"`
	UserID       string    `json:"user_id"`
	SocialHandle string    `json:"social_handle"`
	CreatedAt    time.Time `json:"created_at"`
}

// SplitInterestRequest is the payload for POST /api/v1/items/{itemId}/split-interest.
type SplitInterestRequest struct {
	SocialHandle string `json:"social_handle"`
}
