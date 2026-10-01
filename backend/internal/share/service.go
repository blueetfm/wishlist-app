package share

import (
	"context"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
	"github.com/blueetfm/wishlist-app/backend/internal/db"
	"github.com/blueetfm/wishlist-app/backend/internal/item"
	"github.com/blueetfm/wishlist-app/backend/internal/wishlist"
)

// Service implements the business logic for toggling and resolving a
// wishlist's public share link, backed by the wishlists.is_shared flag and
// its existing (permanent) share_token column.
type Service struct {
	pool      db.Pool
	wishlists *wishlist.Service
	items     *item.Service
}

// NewService constructs a Service backed by the given pool.
func NewService(pool db.Pool, wishlists *wishlist.Service, items *item.Service) *Service {
	return &Service{pool: pool, wishlists: wishlists, items: items}
}

// GetActive returns the current share status for wishlistID, verifying
// userID owns it. Returns ErrNotFound if the wishlist isn't currently shared.
func (s *Service) GetActive(ctx context.Context, userID, wishlistID string) (*Share, error) {
	w, err := s.wishlists.GetOwned(ctx, userID, wishlistID)
	if err != nil {
		return nil, err
	}
	if !w.IsShared {
		return nil, apperr.ErrNotFound
	}
	return &Share{Token: w.ShareToken, IsActive: w.IsShared}, nil
}

// CreateOrGet shares wishlistID, verifying userID owns it. If it's already
// shared, the existing token is returned unchanged so retries stay
// idempotent. Otherwise sharing is turned back on with a freshly rotated
// share_token, so links revoked in the past can't be resurrected.
func (s *Service) CreateOrGet(ctx context.Context, userID, wishlistID string) (*Share, error) {
	w, err := s.wishlists.GetOwned(ctx, userID, wishlistID)
	if err != nil {
		return nil, err
	}

	if w.IsShared {
		return &Share{Token: w.ShareToken, IsActive: w.IsShared}, nil
	}

	row := s.pool.QueryRow(ctx, `
		UPDATE wishlists
		SET is_shared = true, share_token = gen_random_uuid()
		WHERE id = $1
		RETURNING share_token, is_shared
	`, wishlistID)

	var status Share
	if err := row.Scan(&status.Token, &status.IsActive); err != nil {
		return nil, err
	}
	return &status, nil
}

// Revoke turns off sharing for wishlistID, verifying userID owns it.
// Returns ErrNotFound if the wishlist isn't currently shared.
func (s *Service) Revoke(ctx context.Context, userID, wishlistID string) error {
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
		return apperr.ErrNotFound
	}
	return nil
}

// ResolvePublic looks up a wishlist and its items by its share_token, for
// the public /shared/{shareToken} endpoint.
// Optional auth. viewerID is the optional authenticated caller (empty for guests).
func (s *Service) ResolvePublic(ctx context.Context, token, viewerID string) (*WishlistDetail, error) {
	w, err := s.wishlists.GetByShareToken(ctx, token)
	if err != nil {
		return nil, err
	}

	items, err := s.items.ListByWishlist(ctx, w.ID, viewerID)
	if err != nil {
		return nil, err
	}

	return &WishlistDetail{Wishlist: *w, Items: items}, nil
}
