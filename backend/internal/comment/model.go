package comment

import "time"

// Comment matches the Comment schema in docs/swagger.yaml. It is polymorphic
// (attached to either a wishlist or an item) and can be threaded via ParentID.
type Comment struct {
	ID         string    `json:"id"`
	AuthorID   string    `json:"author_id"`
	WishlistID *string   `json:"wishlist_id"`
	ItemID     *string   `json:"item_id"`
	ParentID   *string   `json:"parent_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	IsMine     bool      `json:"is_mine"`
}

// CreateCommentRequest is the payload for POST /api/v1/comments. Exactly one
// of WishlistID or ItemID must be set.
type CreateCommentRequest struct {
	WishlistID *string `json:"wishlistId"`
	ItemID     *string `json:"itemId"`
	ParentID   *string `json:"parentId"`
	Content    string  `json:"content"`
}

// UpdateCommentRequest is the payload for PUT /api/v1/comments/{commentId}
type UpdateCommentRequest struct {
	Content *string `json:"content"`
}

// CommentTree is a single node in a nested reply structure: a comment
// together with its threaded replies, to arbitrary depth. Returned by
// ListByWishlistID/ListByItemID instead of the flat Comment rows they scan.
type CommentTree struct {
	ID       string        `json:"id"`
	Content  string        `json:"content"`
	Children []CommentTree `json:"children"`
}


