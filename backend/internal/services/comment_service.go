package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/blueetfm/wishlist-app/backend/internal/models"
)

// CommentService implements the business logic for creating comments
type CommentService struct {
	pool *pgxpool.Pool
}

// NewCommentService constructs a CommentService backed by the given pool.
func NewCommentService(pool *pgxpool.Pool) *CommentService {
	return &CommentService{pool: pool}
}

const commentColumns = `id, wishlist_id, item_id, parent_id, author_id, content, created_at`

func scanComment(row pgx.Row) (models.Comment, error) {
	var c models.Comment
	err := row.Scan(&c.ID, &c.WishlistID, &c.ItemID, &c.ParentID, &c.AuthorID, &c.Content, &c.CreatedAt)
	return c, err
}

// Create inserts a comment with authorID owned by the current user
func (s *CommentService) Create(ctx context.Context, authorID string, req models.CreateCommentRequest) (*models.Comment, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO comments (wishlist_id, item_id, parent_id, author_id, content)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+commentColumns, req.WishlistID, req.ItemID, req.ParentID, authorID, req.Content)

	c, err := scanComment(row)
	if err != nil {
		return nil, err

	}

	return &c, nil
}

// List comments by wishlist ID
// Perform recursive CTE to fetch all comments and their nested replies for a given wishlist ID
func (s *CommentService) ListByWishlistID(ctx context.Context, wishlistID string) ([]models.Comment, error) {
	rows, err := s.pool.Query(ctx, `
		WITH RECURSIVE comment_tree AS (
			SELECT `+commentColumns+`
			FROM comments
			WHERE wishlist_id = $1
			UNION ALL
			SELECT c.`+commentColumns+`
			FROM comments c
			WHERE c.parent_id = comment_tree.id
			)
			SELECT * FROM comment_tree
			`, wishlistID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// List comments by item ID
// Perform recursive CTE to fetch all comments and their nested replies for a given item ID
func (s *CommentService) ListByItemID(ctx context.Context, itemID string) ([]models.Comment, error) {
	rows, err := s.pool.Query(ctx, `
		WITH RECURSIVE comment_tree AS (
			SELECT `+commentColumns+`
			FROM comments
			WHERE item_id = $1
			UNION ALL
			SELECT c.`+commentColumns+`
			FROM comments c
			WHERE c.parent_id = comment_tree.id
			)
			SELECT * FROM comment_tree
			`, itemID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// GetOwned returns the comment by id, verifying it belongs to userID.
// Returns ErrNotFound if it doesn't exist, ErrForbidden if it belongs to
// someone else.
func (s *CommentService) GetOwned(ctx context.Context, userID, commentID string) (*models.Comment, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+commentColumns+` FROM comments WHERE id = $1`, commentID)

	c, err := scanComment(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if c.AuthorID != userID {
		return nil, ErrForbidden
	}
	return &c, nil
}

// Update applies changes to a comment owned by userID. Re-verifies
// ownership server-side via GetOwned - the frontend's is_mine flag is
// display-only and must never be trusted for authorization.
func (s *CommentService) Update(ctx context.Context, userID, commentID string, req models.UpdateCommentRequest) (*models.Comment, error) {
	existing, err := s.GetOwned(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}

	content := existing.Content
	if req.Content != nil {
		content = *req.Content
	}

	row := s.pool.QueryRow(ctx, `
		UPDATE comments
		SET content = $1
		WHERE id = $2
		RETURNING `+commentColumns, content, commentID)

	c, err := scanComment(row)
	if err != nil {
		return nil, err
	}
	c.IsMine = true
	return &c, nil
}
