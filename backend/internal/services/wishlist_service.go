package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/blueetfm/wishlist-app/backend/internal/models"
)

// WishlistService implements the business logic for creating, reading,
// updating, and deleting wishlists.
//
// NOTE: queries assume a `wishlists` table with columns
// (id, user_id, title, description, share_token, created_at).
// Adjust the SQL below if your actual Supabase schema differs.
type WishlistService struct {
	pool *pgxpool.Pool
}

// NewWishlistService constructs a WishlistService backed by the given pool.
func NewWishlistService(pool *pgxpool.Pool) *WishlistService {
	return &WishlistService{pool: pool}
}

const wishlistColumns = `id, user_id, title, description, share_token, created_at`

func scanWishlist(row pgx.Row) (models.Wishlist, error) {
	var w models.Wishlist
	err := row.Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.ShareToken, &w.CreatedAt)
	return w, err
}

// Create inserts a new wishlist owned by userID. The database is
// responsible for generating the id and share_token (e.g. via column
// defaults such as gen_random_uuid()).
func (s *WishlistService) Create(ctx context.Context, userID string, req models.CreateWishlistRequest) (*models.Wishlist, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO wishlists (user_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING `+wishlistColumns, userID, req.Title, req.Description)

	w, err := scanWishlist(row)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// ListByOwner returns all wishlists owned by userID, most recent first.
func (s *WishlistService) ListByOwner(ctx context.Context, userID string) ([]models.Wishlist, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+wishlistColumns+`
		FROM wishlists
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wishlists := make([]models.Wishlist, 0)
	for rows.Next() {
		w, err := scanWishlist(rows)
		if err != nil {
			return nil, err
		}
		wishlists = append(wishlists, w)
	}
	return wishlists, rows.Err()
}

// GetOwned returns the wishlist by id, verifying it is owned by userID.
// Returns ErrNotFound if it doesn't exist, ErrForbidden if it's owned by
// someone else.
func (s *WishlistService) GetOwned(ctx context.Context, userID, wishlistID string) (*models.Wishlist, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+wishlistColumns+`
		FROM wishlists
		WHERE id = $1
	`, wishlistID)

	w, err := scanWishlist(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	return &w, nil
}

// Update applies partial changes to a wishlist owned by userID.
func (s *WishlistService) Update(ctx context.Context, userID, wishlistID string, req models.UpdateWishlistRequest) (*models.Wishlist, error) {
	existing, err := s.GetOwned(ctx, userID, wishlistID)
	if err != nil {
		return nil, err
	}

	title := existing.Title
	if req.Title != nil {
		title = *req.Title
	}
	description := existing.Description
	if req.Description != nil {
		description = *req.Description
	}

	row := s.pool.QueryRow(ctx, `
		UPDATE wishlists
		SET title = $1, description = $2
		WHERE id = $3
		RETURNING `+wishlistColumns, title, description, wishlistID)

	w, err := scanWishlist(row)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// Delete removes a wishlist owned by userID. The database is expected to
// cascade-delete related items, claims, split interests, and comments via
// foreign key constraints (ON DELETE CASCADE).
func (s *WishlistService) Delete(ctx context.Context, userID, wishlistID string) error {
	if _, err := s.GetOwned(ctx, userID, wishlistID); err != nil {
		return err
	}

	_, err := s.pool.Exec(ctx, `DELETE FROM wishlists WHERE id = $1`, wishlistID)
	return err
}
