// Package services contains business logic for the wishlist API, invoked
// by HTTP handlers.
package services

import "errors"

// Common service-layer errors, translated to HTTP status codes by handlers.
var (
	ErrNotFound  = errors.New("resource not found")
	ErrForbidden = errors.New("you do not have access to this resource")
	ErrConflict  = errors.New("resource already exists")
)
