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

// Sentinel errors specific to embed scraping, translated to HTTP status
// codes by writeServiceError.
var (
	ErrInvalidURL          = errors.New("invalid or disallowed url")
	ErrUpstreamUnavailable = errors.New("failed to fetch the target url")
	ErrNoMetadata          = errors.New("no open graph metadata found for this url")
)
