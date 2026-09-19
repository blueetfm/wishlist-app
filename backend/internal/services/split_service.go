package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/blueetfm/wishlist-app/backend/internal/models"
)

// SplitService implements the business logic for displaying split interest for a wishlist item
type SplitService struct {
	pool *pgxpool.Pool
}

// NewSplitService creates a new instance of SplitService
func NewSplitService(pool *pgxpool.Pool) *SplitService {
	return &SplitService{pool: pool}
}

const splitColumns = `id, item_id, user_id, social_handle, created_at`

func scanSplitInterestEntry(row pgx.Row) (models.SplitInterestEntry, error) {
	var s models.SplitInterestEntry
	err := row.Scan(&s.ID, &s.ItemID, &s.UserID, &s.SocialHandle, &s.CreatedAt)
	return s, err	
}

// Create inserts a new split interest entry into the database.
func (s* SplitService) Create (ctx context.Context, userId string, req models.SplitInterestRequest) (*models.SplitInterestEntry, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO split_interests (item_id, user_id, social_handle)
		VALUES ($1, $2, $3)
		RETURNING `+splitColumns, req.ItemID, userId, req.SocialHandle)
	
	split, err := scanSplitInterestEntry(row)
	if err != nil {
		return nil, err
	}
	return &split, nil
}


// ListByItemID retrieves all split interest entries for a given item ID.
func (s* SplitService) ListByItemID(ctx context.Context, itemId string) ([]models.SplitInterestEntry, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+splitColumns+`
		FROM split_interests
		WHERE item_id = $1
		`,
		itemId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	splits := make([]models.SplitInterestEntry, 0)
	for rows.Next() {
		split, err := scanSplitInterestEntry(rows)
		if err != nil {
			return nil, err
		}
		splits = append(splits, split)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return splits, nil
}

// Delete removes a split interest entry by item ID and user ID
func (s* SplitService) Delete(ctx context.Context, itemId string, userId string) error {
	_, err := s.pool.Exec(ctx, `
		DELETE FROM split_interests
		WHERE item_id = $1 AND user_id = $2
	`, itemId, userId)
	return err
}
