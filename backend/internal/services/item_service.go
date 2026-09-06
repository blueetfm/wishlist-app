package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/blueetfm/wishlist-app/backend/internal/models"
)

// ItemService implemenets the business logic for creating,
// updating, and deleting wishlist items in a wishlist.

// ItemService constructs a ItemService backed by the given pool.
type ItemService struct {
	pool *pgxpool.Pool
}

func NewItemService(pool *pgxpool.Pool) *ItemService {
	return &ItemService{pool: pool}
}

const itemColumns = `id, wishlist_id, name, product_url, price, embed_data, created_at`

func scanItem(row pgx.Row) (models.Item, error) {
	var i models.Item
	err := row.Scan(&i.ID, &i.WishlistID, &i.Name, &i.ProductURL, &i.Price, &i.EmbedData, &i.CreatedAt)
	return i, err
}

// Create inserts a new wishlist item into wishlistID. The database is
// responsible for generating the id (e.g. via a column default such as
// gen_random_uuid()). Callers should verify the caller owns wishlistID
// (e.g. via WishlistService.GetOwned) before calling Create.
func (s *ItemService) Create(ctx context.Context, wishlistID string, req models.CreateItemRequest) (*models.Item, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO items (wishlist_id, name, product_url, price, embed_data)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+itemColumns, wishlistID, req.Name, req.ProductURL, req.Price, req.EmbedData)

	i, err := scanItem(row)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// ListByWishlist returns all items belonging to wishlistID, oldest first.
func (s *ItemService) ListByWishlist(ctx context.Context, wishlistID string) ([]models.Item, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+itemColumns+`
		FROM items
		WHERE wishlist_id = $1
		ORDER BY created_at ASC
	`, wishlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Item, 0)
	for rows.Next() {
		i, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// GetOwned returns the item by id, verifying its parent wishlist is owned by
// ownerID. Returns ErrNotFound if it doesn't exist, ErrForbidden if the
// parent wishlist belongs to someone else.
func (s *ItemService) GetOwned(ctx context.Context, ownerID, itemID string) (*models.Item, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT i.id, i.wishlist_id, i.name, i.product_url, i.price, i.embed_data, i.created_at, w.user_id
		FROM items i
		JOIN wishlists w ON w.id = i.wishlist_id
		WHERE i.id = $1
	`, itemID)

	var i models.Item
	var ownerCheck string
	err := row.Scan(&i.ID, &i.WishlistID, &i.Name, &i.ProductURL, &i.Price, &i.EmbedData, &i.CreatedAt, &ownerCheck)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if ownerCheck != ownerID {
		return nil, ErrForbidden
	}
	return &i, nil
}

// Update applies partial changes to an item, verifying ownerID owns the
// item's parent wishlist.
func (s *ItemService) Update(ctx context.Context, ownerID, itemID string, req models.UpdateItemRequest) (*models.Item, error) {
	existing, err := s.GetOwned(ctx, ownerID, itemID)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	if req.Name != nil {
		name = *req.Name
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

	row := s.pool.QueryRow(ctx, `
		UPDATE items
		SET name = $1, product_url = $2, price = $3, embed_data = $4
		WHERE id = $5
		RETURNING `+itemColumns, name, productURL, price, embedData, itemID)

	i, err := scanItem(row)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// Delete removes an item, verifying ownerID owns the item's parent wishlist.
// The database is expected to cascade-delete related claims, split
// interests, and comments via foreign key constraints (ON DELETE CASCADE).
func (s *ItemService) Delete(ctx context.Context, ownerID, itemID string) error {
	if _, err := s.GetOwned(ctx, ownerID, itemID); err != nil {
		return err
	}

	_, err := s.pool.Exec(ctx, `DELETE FROM items WHERE id = $1`, itemID)
	return err
}
