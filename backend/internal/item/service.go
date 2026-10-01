package item

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
	"github.com/blueetfm/wishlist-app/backend/internal/db"
)

// Service implements the business logic for creating, updating, and
// deleting wishlist items.
type Service struct {
	pool   db.Pool
	claims *ClaimService
	splits *SplitService
}

// NewService constructs a Service backed by the given pool, claims, and splits.
func NewService(pool db.Pool, claims *ClaimService, splits *SplitService) *Service {
	return &Service{pool: pool, claims: claims, splits: splits}
}

const itemColumns = `id, wishlist_id, name, description, product_url, price, embed_data, created_at, image_url`

func scanItem(row pgx.Row) (Item, error) {
	var i Item
	err := row.Scan(&i.ID, &i.WishlistID, &i.Name, &i.Description, &i.ProductURL, &i.Price, &i.EmbedData, &i.CreatedAt, &i.ImageURL)
	i.SplitInterests = make([]SplitInterestEntry, 0)
	return i, err
}

// loadClaimAndSplitInterests populates item's Claim and SplitInterests fields, marking
// any that belong to viewerID as IsMine.
func (s *Service) loadClaimAndSplitInterests(ctx context.Context, item *Item, viewerID string) error {
	c, err := s.claims.GetByItemID(ctx, item.ID)
	if err != nil {
		return err
	}
	if c != nil {
		c.IsMine = c.UserID == viewerID
		item.Claim = c
	}

	splits, err := s.splits.ListByItemID(ctx, item.ID)
	if err != nil {
		return err
	}
	for i := range splits {
		splits[i].IsMine = splits[i].UserID == viewerID
	}
	item.SplitInterests = splits
	return nil
}

// Create inserts a new wishlist item into wishlistID. The database is
// responsible for generating the id (e.g. via a column default such as
// gen_random_uuid()). Callers should verify the caller owns wishlistID
// (e.g. via wishlist.Service.GetOwned) before calling Create.
func (s *Service) Create(ctx context.Context, wishlistID string, req CreateItemRequest) (*Item, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO wishlist_items (wishlist_id, name, description, product_url, price, embed_data, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+itemColumns, wishlistID, req.Name, req.Description, req.ProductURL, req.Price, req.EmbedData, req.ImageURL)

	i, err := scanItem(row)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// ListByWishlist returns all items belonging to wishlistID, oldest first, marking
// any claims/split interests owned by viewerID as IsMine.
func (s *Service) ListByWishlist(ctx context.Context, wishlistID, viewerID string) ([]Item, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+itemColumns+`
		FROM wishlist_items
		WHERE wishlist_id = $1
		ORDER BY created_at ASC
	`, wishlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	for rows.Next() {
		i, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		if err := s.loadClaimAndSplitInterests(ctx, &i, viewerID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// GetOwned returns the item by id within wishlistID, verifying its parent
// wishlist is owned by ownerID. Returns ErrNotFound if it doesn't exist or
// doesn't belong to wishlistID, ErrForbidden if the parent wishlist belongs
// to someone else.
func (s *Service) GetOwned(ctx context.Context, ownerID, wishlistID, itemID string) (*Item, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT i.id, i.wishlist_id, i.name, i.description, i.product_url, i.price, i.embed_data, i.created_at, i.image_url, w.user_id
		FROM wishlist_items i
		JOIN wishlists w ON w.id = i.wishlist_id
		WHERE i.id = $1 AND i.wishlist_id = $2
	`, itemID, wishlistID)

	var i Item
	var ownerCheck string
	err := row.Scan(&i.ID, &i.WishlistID, &i.Name, &i.Description, &i.ProductURL, &i.Price, &i.EmbedData, &i.CreatedAt, &i.ImageURL, &ownerCheck)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if ownerCheck != ownerID {
		return nil, apperr.ErrForbidden
	}
	if err := s.loadClaimAndSplitInterests(ctx, &i, ownerID); err != nil {
		return nil, err
	}
	return &i, nil
}

// ExistsInWishlist verifies itemID belongs to wishlistID, without requiring
// ownership. Used by handlers acting on someone else's wishlist (e.g.
// claiming an item), where the caller is never the wishlist owner.
func (s *Service) ExistsInWishlist(ctx context.Context, wishlistID, itemID string) error {
	var exists bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM wishlist_items WHERE id = $1 AND wishlist_id = $2)
	`, itemID, wishlistID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return apperr.ErrNotFound
	}
	return nil
}

// Update applies partial changes to an item within wishlistID, verifying
// ownerID owns the item's parent wishlist.
func (s *Service) Update(ctx context.Context, ownerID, wishlistID, itemID string, req UpdateItemRequest) (*Item, error) {
	existing, err := s.GetOwned(ctx, ownerID, wishlistID, itemID)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	if req.Name != nil {
		name = *req.Name
	}
	description := existing.Description
	if req.Description != nil {
		description = *req.Description
	}
	productURL := existing.ProductURL
	if req.ProductURL != nil {
		productURL = *req.ProductURL
	}
	price := existing.Price
	if req.Price != nil {
		price = *req.Price
	}
	embedData := existing.EmbedData
	if req.EmbedData != nil {
		embedData = req.EmbedData
	}
	imageURL := existing.ImageURL
	if req.ImageURL != nil {
		imageURL = req.ImageURL
	}

	row := s.pool.QueryRow(ctx, `
		UPDATE wishlist_items
		SET name = $1, description = $2, product_url = $3, price = $4, embed_data = $5, image_url = $6
		WHERE id = $7
		RETURNING `+itemColumns, name, description, productURL, price, embedData, imageURL, itemID)

	i, err := scanItem(row)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// Delete removes an item within wishlistID, verifying ownerID owns the
// item's parent wishlist.
// The database is expected to cascade-delete related claims, split
// interests, and comments via foreign key constraints (ON DELETE CASCADE).
func (s *Service) Delete(ctx context.Context, ownerID, wishlistID, itemID string) error {
	if _, err := s.GetOwned(ctx, ownerID, wishlistID, itemID); err != nil {
		return err
	}

	_, err := s.pool.Exec(ctx, `DELETE FROM wishlist_items WHERE id = $1`, itemID)
	return err
}

// HideClaimsForOwner strips claim and split-interest data from items so an
// owner who has opted into a spoiler-free view (Wishlist.HideClaimsFromOwner)
// doesn't see who claimed or split-interested their own gifts. Guests never
// call this - the shared/public view always returns full claim data.
func HideClaimsForOwner(items []Item) []Item {
	sanitized := make([]Item, len(items))
	for i, item := range items {
		item.Claim = nil
		item.SplitInterests = nil
		sanitized[i] = item
	}
	return sanitized
}
