// Package apperr defines sentinel errors shared by service implementations
// across domain packages and translated to HTTP status codes by the
// transport layer (see internal/httpx).
package apperr

import "errors"

// Common service-layer errors.
var (
	ErrNotFound  = errors.New("resource not found")
	ErrForbidden = errors.New("you do not have access to this resource")
	ErrConflict  = errors.New("resource already exists")
)

// Sentinel errors specific to embed scraping.
var (
	ErrInvalidURL          = errors.New("invalid or disallowed url")
	ErrUpstreamUnavailable = errors.New("failed to fetch the target url")
	ErrNoMetadata          = errors.New("no open graph metadata found for this url")
)
