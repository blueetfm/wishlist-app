package idempotency

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/blueetfm/wishlist-app/backend/internal/db"
)

// Service persists Idempotency-Key reservations so retried requests can
// replay the original response instead of re-running the handler.
type Service struct {
	pool db.Pool
}

// NewService constructs a Service backed by the given pool.
func NewService(pool db.Pool) *Service {
	return &Service{pool: pool}
}

const idempotencyColumns = `id, user_id, idempotency_key, method, path, request_hash, response_status, response_body, created_at`

func scanRecord(row pgx.Row) (Record, error) {
	var rec Record
	err := row.Scan(&rec.ID, &rec.UserID, &rec.IdempotencyKey, &rec.Method, &rec.Path,
		&rec.RequestHash, &rec.ResponseStatus, &rec.ResponseBody, &rec.CreatedAt)
	return rec, err
}

// Reserve attempts to claim (userID, key, method, path) for a new request.
// reserved is true if this call created the row (the caller should execute
// the handler and call Complete). If the key was already used, reserved is
// false and the existing record is returned instead - callers should check
// RequestHash to detect key reuse with a different payload, and check
// whether ResponseStatus is nil (still in flight) before replaying it.
func (s *Service) Reserve(ctx context.Context, userID, key, method, path, requestHash string) (*Record, bool, error) {
	// no rows returned if the tuple key already exists
	row := s.pool.QueryRow(ctx, `
		INSERT INTO idempotency_keys (user_id, idempotency_key, method, path, request_hash)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, idempotency_key, method, path) DO NOTHING
		RETURNING `+idempotencyColumns, userID, key, method, path, requestHash)

	rec, err := scanRecord(row)

	// essentially, if there was a previous conflict, get the already existing one
	if errors.Is(err, pgx.ErrNoRows) {
		existing, getErr := s.get(ctx, userID, key, method, path)
		if getErr != nil {
			return nil, false, getErr
		}
		return existing, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	// else return the newly created one
	return &rec, true, nil
}

func (s *Service) get(ctx context.Context, userID, key, method, path string) (*Record, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+idempotencyColumns+`
		FROM idempotency_keys
		WHERE user_id = $1 AND idempotency_key = $2 AND method = $3 AND path = $4
	`, userID, key, method, path)

	rec, err := scanRecord(row)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// Complete records the response produced for a previously reserved request.
func (s *Service) Complete(ctx context.Context, id string, status int, body []byte) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE idempotency_keys
		SET response_status = $1, response_body = $2
		WHERE id = $3
	`, status, body, id)
	return err
}
