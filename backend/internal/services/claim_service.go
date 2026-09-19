package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/blueetfm/wishlist-app/backend/internal/models"
)

// ClaimService implemenets the business logic for claiming an item in wishlist

// ClaimService constructs a ClaimService backed by the given pool.
type ClaimService struct {
	pool *pgxpool.Pool
}

func NewClaimService(pool *pgxpool.Pool) *ClaimService {
	return &ClaimService{pool: pool}
}

const claimColumns = `id, item_id, user_id, message, claimed_at`

func scanClaim(row pgx.Row) (models.Claim, error) {
	var c models.Claim
	err := row.Scan(&c.ID, &c.ItemID, &c.UserID, &c.Message, &c.ClaimedAt)
	return c, err
}

// Create checks whether this item has been claimed already.
// Fail with ErrConflict is already claimed by someone else
// If the item has not been claimed, it creates a new claim record.
func (s *ClaimService) Create(ctx context.Context, itemID, userID string, message string) (*models.Claim, error) {
	const query = `
		INSERT INTO claims (item_id, user_id, message)
		VALUES ($1, $2, $3)
		RETURNING ` + claimColumns

	row := s.pool.QueryRow(ctx, query, itemID, userID, message)
	c, err := scanClaim(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrConflict
		}
		return nil, err
	}
	return &c, nil
}

// Delete unclaims an item by removing the corresponding claim record.
// Fail with ErrNotFound if the claim does not exist.
func (s *ClaimService) Delete(ctx context.Context, itemID, userID string) error {
	const query = `
		DELETE FROM claims
		WHERE item_id = $1 AND user_id = $2
		RETURNING ` + claimColumns

	row := s.pool.QueryRow(ctx, query, itemID, userID)
	var c models.Claim
	err := row.Scan(&c.ID, &c.ItemID, &c.UserID, &c.Message, &c.ClaimedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
