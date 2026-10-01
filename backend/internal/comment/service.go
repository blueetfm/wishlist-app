package comment

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
	"github.com/blueetfm/wishlist-app/backend/internal/db"
)

// Service implements the business logic for creating comments.
type Service struct {
	pool db.Pool
}

// NewService constructs a Service backed by the given pool.
func NewService(pool db.Pool) *Service {
	return &Service{pool: pool}
}

const commentColumns = `id, wishlist_id, item_id, parent_id, author_id, content, created_at`

func scanComment(row pgx.Row) (Comment, error) {
	var c Comment
	err := row.Scan(&c.ID, &c.WishlistID, &c.ItemID, &c.ParentID, &c.AuthorID, &c.Content, &c.CreatedAt)
	return c, err
}

// buildCommentForest nests replies under their parent comment (to any depth),
// returning only the top-level (ParentID == nil) comments as tree roots.
// https://vikkrraant.medium.com/the-anatomy-of-reddit-style-comments-a-weekend-engineering-dive-6b52cd80139d
func buildCommentForest(comments []Comment) []CommentTree {
	childrenOf := make(map[string][]Comment) // Key: ParentID (string), Value: slice of child Comments
	var rootComments []Comment
	var forest []CommentTree

	for _, comment := range comments {
		if comment.ParentID == nil {
			// 1. add comments with no parent IDs to root comments
			rootComments = append(rootComments, comment)
		} else {
			// add comments with parent IDs to the map
			parentID := *comment.ParentID
			childrenOf[parentID] = append(childrenOf[parentID], comment)
		}
	}

	for _, root := range rootComments {
		forest = append(forest, buildCommentNode(root, childrenOf))
	}

	return forest
}

// buildCommentNode recursively attaches comment's descendants as Children.
func buildCommentNode(comment Comment, childrenOf map[string][]Comment) CommentTree {
	node := CommentTree{
		ID:       comment.ID,
		Content:  comment.Content,
		Children: []CommentTree{}, // Initialize as empty slice instead of nil for better JSON serialization
	}

	// fetch children for the current comment ID
	for _, child := range childrenOf[comment.ID] {
		node.Children = append(node.Children, buildCommentNode(child, childrenOf))
	}

	return node
}

// Create inserts a comment with authorID owned by the current user.
func (s *Service) Create(ctx context.Context, authorID string, req CreateCommentRequest) (*Comment, error) {
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

// ListByWishlistID lists comments by wishlist ID, nested into a comment forest.
// Perform recursive CTE to fetch all comments and their nested replies for a given wishlist ID
func (s *Service) ListByWishlistID(ctx context.Context, wishlistID string) ([]CommentTree, error) {
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

	var comments []Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return buildCommentForest(comments), rows.Err()
}

// ListByItemID lists comments by item ID, nested into a comment forest.
// Perform recursive CTE to fetch all comments and their nested replies for a given item ID
func (s *Service) ListByItemID(ctx context.Context, itemID string) ([]CommentTree, error) {
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

	var comments []Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return buildCommentForest(comments), rows.Err()
}

// GetOwned returns the comment by id, verifying it belongs to userID.
// Returns ErrNotFound if it doesn't exist, ErrForbidden if it belongs to
// someone else.
func (s *Service) GetOwned(ctx context.Context, userID, commentID string) (*Comment, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+commentColumns+` FROM comments WHERE id = $1`, commentID)

	c, err := scanComment(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if c.AuthorID != userID {
		return nil, apperr.ErrForbidden
	}
	return &c, nil
}

// Update applies changes to a comment owned by userID. Re-verifies
// ownership server-side via GetOwned - the frontend's is_mine flag is
// display-only and must never be trusted for authorization.
func (s *Service) Update(ctx context.Context, userID, commentID string, req UpdateCommentRequest) (*Comment, error) {
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
