package comment

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
)

var commentRowColumns = []string{
	"id", "wishlist_id", "item_id", "parent_id", "author_id", "content", "created_at",
}

func newMockPool(t *testing.T) pgxmock.PgxPoolIface {
	// 1. mark this function as test helper
	// when a test fails, the file name and line number of the calling test func
	// will be called
	t.Helper()

	// 2. initialize mock db pool
	mock, err := pgxmock.NewPool()

	// 3. test executions will stop immediately when errors occur
	require.NoError(t, err)

	// 4. each test will schedule its own cleanup, instead of having to call
	// `defer mock.Close()` in each test
	t.Cleanup(func() {
        assert.NoError(t, mock.ExpectationsWereMet())
    })

	return mock
}

// Test creating a new comment on a wishlist
func Test_Create_Wishlist(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	wishlistID := "wl-1"
	req := CreateCommentRequest{
		WishlistID: &wishlistID,
		Content:    "test",
	}

	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(req.WishlistID, req.ItemID, req.ParentID, "u-1", req.Content).
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-2", &wishlistID, (*string)(nil), (*string)(nil), "u-1", req.Content, now))

	got, err := svc.Create(context.Background(), "u-1", req)

	require.NoError(t, err)
	assert.Equal(t, "c-2", got.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test creating a new comment on an item
func Test_Create_Item(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	itemID := "i-1"
	req := CreateCommentRequest{
		ItemID:  &itemID,
		Content: "nice gift",
	}

	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(req.WishlistID, req.ItemID, req.ParentID, "u-1", req.Content).
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-3", (*string)(nil), &itemID, (*string)(nil), "u-1", req.Content, now))

	got, err := svc.Create(context.Background(), "u-1", req)

	require.NoError(t, err)
	assert.Equal(t, "c-3", got.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}


// Test listing comments for a wishlist nests replies under their parent
func Test_ListByWishlistID(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	wishlistID := "wl-1"
	parentID := "c-1"
	mock.ExpectQuery(`FROM comments WHERE wishlist_id = \$1`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-1", &wishlistID, (*string)(nil), (*string)(nil), "u-1", "hi", now).
			AddRow("c-2", &wishlistID, (*string)(nil), &parentID, "u-2", "reply", now))

	got, err := svc.ListByWishlistID(context.Background(), "wl-1")

	require.NoError(t, err)
	require.Len(t, got, 1, "the reply should be nested, not returned as a top-level comment")
	assert.Equal(t, "c-1", got[0].ID)
	require.Len(t, got[0].Children, 1)
	assert.Equal(t, "c-2", got[0].Children[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test listing comments for an item nests replies under their parent
func Test_ListByItemID(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	itemID := "i-1"
	parentID := "c-1"
	mock.ExpectQuery(`FROM comments WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-1", (*string)(nil), &itemID, (*string)(nil), "u-1", "hi", now).
			AddRow("c-2", (*string)(nil), &itemID, &parentID, "u-2", "reply", now))

	got, err := svc.ListByItemID(context.Background(), "i-1")

	require.NoError(t, err)
	require.Len(t, got, 1, "the reply should be nested, not returned as a top-level comment")
	assert.Equal(t, "c-1", got[0].ID)
	require.Len(t, got[0].Children, 1)
	assert.Equal(t, "c-2", got[0].Children[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test that buildCommentForest nests replies to arbitrary depth and leaves
// comments with no replies untouched
func Test_BuildCommentForest(t *testing.T) {
	t.Parallel()

	parent := "root"
	child := "child"
	flat := []Comment{
		{ID: "root"},
		{ID: "child", ParentID: &parent},
		{ID: "grandchild", ParentID: &child},
		{ID: "orphan-root"},
	}

	got := buildCommentForest(flat)

	require.Len(t, got, 2)
	require.Len(t, got[0].Children, 1)
	assert.Equal(t, "child", got[0].Children[0].ID)
	require.Len(t, got[0].Children[0].Children, 1)
	assert.Equal(t, "grandchild", got[0].Children[0].Children[0].ID)
	assert.Empty(t, got[1].Children)
}

// Test returning a comment by id, verifying it belongs to userID
func Test_GetOwned(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	wishlistID := "wl-1"
	mock.ExpectQuery(`SELECT (.+) FROM comments WHERE id = \$1`).
		WithArgs("c-1").
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-1", &wishlistID, (*string)(nil), (*string)(nil), "u-1", "hi", now))

	got, err := svc.GetOwned(context.Background(), "u-1", "c-1")

	require.NoError(t, err)
	assert.Equal(t, "c-1", got.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrNotFound if the comment does not exist
func Test_GetOwned_NotFound(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	mock.ExpectQuery(`SELECT (.+) FROM comments WHERE id = \$1`).
		WithArgs("missing").
		WillReturnRows(pgxmock.NewRows(commentRowColumns))

	_, err := svc.GetOwned(context.Background(), "u-1", "missing")

	assert.ErrorIs(t, err, apperr.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrForbidden if the comment belongs to someone else
func Test_GetOwned_Forbidden(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	wishlistID := "wl-1"
	mock.ExpectQuery(`SELECT (.+) FROM comments WHERE id = \$1`).
		WithArgs("c-1").
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-1", &wishlistID, (*string)(nil), (*string)(nil), "someone-else", "hi", now))

	_, err := svc.GetOwned(context.Background(), "u-1", "c-1")

	assert.ErrorIs(t, err, apperr.ErrForbidden)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test updating a comment owned by userID
func Test_Update(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	wishlistID := "wl-1"
	mock.ExpectQuery(`SELECT (.+) FROM comments WHERE id = \$1`).
		WithArgs("c-1").
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-1", &wishlistID, (*string)(nil), (*string)(nil), "u-1", "old", now))

	newContent := "edited"
	mock.ExpectQuery(`UPDATE comments`).
		WithArgs(newContent, "c-1").
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-1", &wishlistID, (*string)(nil), (*string)(nil), "u-1", newContent, now))

	got, err := svc.Update(context.Background(), "u-1", "c-1", UpdateCommentRequest{Content: &newContent})

	require.NoError(t, err)
	assert.Equal(t, newContent, got.Content)
	assert.True(t, got.IsMine)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrForbidden when users somehow attempt to update a comment
// that does not belong to them
func Test_Update_Forbidden(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	wishlistID := "wl-1"
	mock.ExpectQuery(`SELECT (.+) FROM comments WHERE id = \$1`).
		WithArgs("c-1").
		WillReturnRows(pgxmock.NewRows(commentRowColumns).
			AddRow("c-1", &wishlistID, (*string)(nil), (*string)(nil), "someone-else", "old", now))

	_, err := svc.Update(context.Background(), "u-1", "c-1", UpdateCommentRequest{})

	assert.ErrorIs(t, err, apperr.ErrForbidden)
	assert.NoError(t, mock.ExpectationsWereMet())
}
