package models

import "time"

// Claim matches the Claim schema in docs/swagger.yaml (strict 1-to-1 item claim).
type Claim struct {
	ItemID    string    `json:"item_id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message,omitempty"`
	ClaimedAt time.Time `json:"claimed_at"`
}

// ClaimRequest is the payload for POST /api/v1/items/{itemId}/claim.
type ClaimRequest struct {
	Message string `json:"message,omitempty"`
}
