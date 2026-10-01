package wishlist

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
	"github.com/blueetfm/wishlist-app/backend/internal/db"
)

// Service implements the business logic for creating, reading, updating,
// and deleting wishlists.
//
// NOTE: queries assume a `wishlists` table with columns
// (id, user_id, title, description, share_token, created_at).
type Service struct {
	pool db.Pool
}

// NewService constructs a Service backed by the given pool.
func NewService(pool db.Pool) *Service {
	return &Service{pool: pool}
}

const wishlistColumns = `id, user_id, title, description, share_token, is_shared, hide_claims_from_owner, created_at`

func scanWishlist(row pgx.Row) (Wishlist, error) {
	var w Wishlist
	err := row.Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.ShareToken, &w.IsShared, &w.HideClaimsFromOwner, &w.CreatedAt)
	return w, err
}

// Create inserts a new wishlist owned by userID. The database is
// responsible for generating the id and share_token (e.g. via column
// defaults such as gen_random_uuid()).
func (s *Service) Create(ctx context.Context, userID string, req CreateWishlistRequest) (*Wishlist, error) {
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
func (s *Service) ListByOwner(ctx context.Context, userID string) ([]Wishlist, error) {
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

	wishlists := make([]Wishlist, 0)
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
func (s *Service) GetOwned(ctx context.Context, userID, wishlistID string) (*Wishlist, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+wishlistColumns+`
		FROM wishlists
		WHERE id = $1
	`, wishlistID)

	w, err := scanWishlist(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, apperr.ErrForbidden
	}
	return &w, nil
}

// GetByShareToken returns the wishlist whose share_token is token and that
// currently has sharing enabled. Returns ErrNotFound otherwise.
func (s *Service) GetByShareToken(ctx context.Context, token string) (*Wishlist, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+wishlistColumns+`
		FROM wishlists
		WHERE share_token = $1 AND is_shared = true
	`, token)

	w, err := scanWishlist(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// Update applies partial changes to a wishlist owned by userID.
func (s *Service) Update(ctx context.Context, userID, wishlistID string, req UpdateWishlistRequest) (*Wishlist, error) {
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
	hideClaimsFromOwner := existing.HideClaimsFromOwner
	if req.HideClaimsFromOwner != nil {
		hideClaimsFromOwner = *req.HideClaimsFromOwner
	}

	row := s.pool.QueryRow(ctx, `
		UPDATE wishlists
		SET title = $1, description = $2, hide_claims_from_owner = $3
		WHERE id = $4
		RETURNING `+wishlistColumns, title, description, hideClaimsFromOwner, wishlistID)

	w, err := scanWishlist(row)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// Delete removes a wishlist owned by userID. The database is expected to
// cascade-delete related items, claims, split interests, and comments via
// foreign key constraints (ON DELETE CASCADE).
func (s *Service) Delete(ctx context.Context, userID, wishlistID string) error {
	if _, err := s.GetOwned(ctx, userID, wishlistID); err != nil {
		return err
	}

	_, err := s.pool.Exec(ctx, `DELETE FROM wishlists WHERE id = $1`, wishlistID)
	return err
}
