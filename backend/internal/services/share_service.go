package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/blueetfm/wishlist-app/backend/internal/models"
)

// ShareService implements the business logic for toggling and resolving a
// wishlist's public share link, backed by the wishlists.is_shared flag and
// its existing (permanent) share_token column.
type ShareService struct {
	pool      *pgxpool.Pool
	wishlists *WishlistService
	items     *ItemService
}

// NewShareService constructs a ShareService backed by the given pool.
func NewShareService(pool *pgxpool.Pool, wishlists *WishlistService, items *ItemService) *ShareService {
	return &ShareService{pool: pool, wishlists: wishlists, items: items}
}

// GetActive returns the current share status for wishlistID, verifying
// userID owns it. Returns ErrNotFound if the wishlist isn't currently shared.
func (s *ShareService) GetActive(ctx context.Context, userID, wishlistID string) (*models.Share, error) {
	w, err := s.wishlists.GetOwned(ctx, userID, wishlistID)
	if err != nil {
		return nil, err
	}
	if !w.IsShared {
		return nil, ErrNotFound
	}
	return &models.Share{Token: w.ShareToken, IsActive: w.IsShared}, nil
}

// CreateOrGet shares wishlistID, verifying userID owns it. If it's already
// shared, the existing token is returned unchanged so retries stay
// idempotent. Otherwise sharing is turned back on with a freshly rotated
// share_token, so links revoked in the past can't be resurrected.
func (s *ShareService) CreateOrGet(ctx context.Context, userID, wishlistID string) (*models.Share, error) {
	w, err := s.wishlists.GetOwned(ctx, userID, wishlistID)
	if err != nil {
		return nil, err
	}

	if w.IsShared {
		return &models.Share{Token: w.ShareToken, IsActive: w.IsShared}, nil
	}

	row := s.pool.QueryRow(ctx, `
		UPDATE wishlists
		SET is_shared = true, share_token = gen_random_uuid()
		WHERE id = $1
		RETURNING share_token, is_shared
	`, wishlistID)

	var status models.Share
	if err := row.Scan(&status.Token, &status.IsActive); err != nil {
		return nil, err
	}
	return &status, nil
}

// Revoke turns off sharing for wishlistID, verifying userID owns it.
// Returns ErrNotFound if the wishlist isn't currently shared.
func (s *ShareService) Revoke(ctx context.Context, userID, wishlistID string) error {
	if _, err := s.wishlists.GetOwned(ctx, userID, wishlistID); err != nil {
		return err
	}

	tag, err := s.pool.Exec(ctx, `
		UPDATE wishlists
		SET is_shared = false
		WHERE id = $1 AND is_shared = true
	`, wishlistID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ResolvePublic looks up a wishlist and its items by its share_token, for
// the public /shared/{shareToken} endpoint.
// Optioanl auth.  viewerID is the optional authenticated caller (empty for guests);
func (s *ShareService) ResolvePublic(ctx context.Context, token, viewerID string) (*models.WishlistDetail, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+wishlistColumns+`
		FROM wishlists
		WHERE share_token = $1 AND is_shared = true
	`, token)

	w, err := scanWishlist(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	items, err := s.items.ListByWishlist(ctx, w.ID, viewerID)
	if err != nil {
		return nil, err
	}

	return &models.WishlistDetail{Wishlist: w, Items: items}, nil
}
