package middleware

import (
	"context"
	"log"
)

type contextKey string

const userIDContextKey contextKey = "userID"

// UserID extracts the authenticated Supabase user id (the JWT "sub" claim)
// from the request context. ok is false if no user is authenticated.
func UserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDContextKey).(string)
	log.Println("userId:", id)
	return id, ok
}
