package models

import "time"

// IdempotencyRecord tracks a single (user, Idempotency-Key, method, path)
// request so retries can replay the original response instead of
// re-executing the handler. ResponseStatus/ResponseBody are nil while the
// original request is still in flight.
type IdempotencyRecord struct {
	ID             string
	UserID         string
	IdempotencyKey string
	Method         string
	Path           string
	RequestHash    string
	ResponseStatus *int
	ResponseBody   []byte
	CreatedAt      time.Time
}
