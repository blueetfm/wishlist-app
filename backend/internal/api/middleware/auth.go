package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth returns middleware that verifies a Supabase Auth JWT sent as
// "Authorization: Bearer <token>" against Supabase's JWKS endpoint, and
// injects the authenticated user's id (the "sub" claim) into the request
// context. Requests without a valid token receive 401 Unauthorized.
func RequireAuth(k keyfunc.Keyfunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
			if !found || tokenString == "" {
				writeUnauthorized(w, "missing bearer token")
				return
			}

			token, err := jwt.Parse(tokenString, k.Keyfunc, jwt.WithValidMethods([]string{"ES256", "RS256"}))
			if err != nil || !token.Valid {
				writeUnauthorized(w, "invalid or expired token")
				return
			}

			sub, err := token.Claims.GetSubject()
			if err != nil || sub == "" {
				writeUnauthorized(w, "token missing subject claim")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
