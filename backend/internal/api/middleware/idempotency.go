package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/blueetfm/wishlist-app/backend/internal/services"
)

// RequireIdempotencyKey returns middleware that, when a client sends an
// "Idempotency-Key" header, ensures the request is only executed once per
// (user, key, method, path): a retry with the same key and body replays the
// original response instead of re-running the handler. Requests without the
// header are passed through unchanged.
func RequireIdempotencyKey(svc *services.IdempotencyService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			userID, ok := UserID(r.Context())
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				writeIdempotencyError(w, http.StatusBadRequest, "failed to read request body")
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))

			hash := sha256.Sum256(body)
			requestHash := hex.EncodeToString(hash[:])

			record, reserved, err := svc.Reserve(r.Context(), userID, key, r.Method, r.URL.Path, requestHash)
			if err != nil {
				writeIdempotencyError(w, http.StatusInternalServerError, "failed to process idempotency key")
				return
			}

			if !reserved {
				if record.RequestHash != requestHash {
					writeIdempotencyError(w, http.StatusUnprocessableEntity, "Idempotency-Key was already used with a different request")
					return
				}
				if record.ResponseStatus == nil {
					writeIdempotencyError(w, http.StatusConflict, "a request with this Idempotency-Key is already being processed")
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(*record.ResponseStatus)
				w.Write(record.ResponseBody)
				return
			}

			rec := &idempotencyResponseRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			_ = svc.Complete(r.Context(), record.ID, rec.status, rec.body.Bytes())
		})
	}
}

// idempotencyResponseRecorder captures the response written by the wrapped
// handler so it can be persisted for future replay, while still forwarding
// it to the real client.
type idempotencyResponseRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (r *idempotencyResponseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *idempotencyResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func writeIdempotencyError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
