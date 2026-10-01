package item

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
)

var itemRowColumns = []string{
	"id", "wishlist_id", "name", "description", "product_url",
	"price", "embed_data", "created_at", "image_url",
}

// itemOwnerRowColumns matches GetOwned's join query, which returns the
// item's parent wishlist owner as a trailing column.
var itemOwnerRowColumns = []string{
	"id", "wishlist_id", "name", "description", "product_url",
	"price", "embed_data", "created_at", "image_url", "user_id",
}

var claimRowColumns = []string{"id", "item_id", "user_id", "message", "claimed_at"}

var splitRowColumns = []string{"id", "item_id", "user_id", "social_handle", "created_at"}

// newTestService wires up a Service backed by the same mock pool as its
// ClaimService and SplitService collaborators.
func newTestService(mock pgxmock.PgxPoolIface) *Service {
	return NewService(mock, NewClaimService(mock), NewSplitService(mock))
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

// Test creating a new wishlist item
func Test_Create(t *testing.T) {
	// 1. tests will run concurrenlty, speeding up test suite
	t.Parallel()

	// 2. setting up test pools and serivce
	mock := newMockPool(t)
	svc := newTestService(mock)

	// 3. set up test req data
	// with no productURL nor embeddata
	now := time.Now()
	req := CreateItemRequest{
		Name:       "Test",
		ProductURL: "",
		Price:      40,
		EmbedData:  nil,
	}

	// 4. set up database mock
	mock.ExpectQuery(`INSERT INTO wishlist_items`).
		WithArgs("wl-1", req.Name, req.Description, req.ProductURL, req.Price, req.EmbedData, req.ImageURL).
		WillReturnRows(pgxmock.NewRows(itemRowColumns).
			AddRow("i-1", "wl-1", req.Name, req.Description, req.ProductURL,
				req.Price, req.EmbedData, now, req.ImageURL))

	// 5. calls service with a blank context
	got, err := svc.Create(context.Background(), "wl-1", req)

	// 6. asserts
	require.NoError(t, err)
	assert.Equal(t, "i-1", got.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test load claim and split interests
func Test_LoadClaimAndSplitInterests(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM claims WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(claimRowColumns).
			AddRow("c-1", "i-1", "user-1", "msg", now))

	mock.ExpectQuery(`SELECT (.+) FROM split_interests WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(splitRowColumns).
			AddRow("s-1", "i-1", "user-2", "@handle", now))

	item := Item{ID: "i-1"}
	err := svc.loadClaimAndSplitInterests(context.Background(), &item, "user-1")

	require.NoError(t, err)
	require.NotNil(t, item.Claim)
	assert.True(t, item.Claim.IsMine, "claim placed by the viewer should be marked IsMine")
	require.Len(t, item.SplitInterests, 1)
	assert.False(t, item.SplitInterests[0].IsMine, "split placed by someone else should not be marked IsMine")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test listing wishlist items by wishlist id
func Test_ListByWishlist(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlist_items WHERE wishlist_id = \$1`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(itemRowColumns).
			AddRow("i-1", "wl-1", "Mug", "", "", 10.0, nil, now, nil))

	mock.ExpectQuery(`SELECT (.+) FROM claims WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(claimRowColumns))

	mock.ExpectQuery(`SELECT (.+) FROM split_interests WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(splitRowColumns))

	got, err := svc.ListByWishlist(context.Background(), "wl-1", "user-1")

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "i-1", got[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning details of a wishlist item,
// given its item id, its parent wishlist id, and a user id
func Test_GetOwned(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlist_items i`).
		WithArgs("i-1", "wl-1").
		WillReturnRows(pgxmock.NewRows(itemOwnerRowColumns).
			AddRow("i-1", "wl-1", "Mug", "", "", 10.0, nil, now, nil, "user-1"))

	mock.ExpectQuery(`SELECT (.+) FROM claims WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(claimRowColumns))

	mock.ExpectQuery(`SELECT (.+) FROM split_interests WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(splitRowColumns))

	got, err := svc.GetOwned(context.Background(), "user-1", "wl-1", "i-1")

	require.NoError(t, err)
	assert.Equal(t, "i-1", got.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrNotFound if the item is not present in the parent wishlist
func Test_GetOwned_NotFound(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	mock.ExpectQuery(`SELECT (.+) FROM wishlist_items i`).
		WithArgs("missing", "wl-1").
		WillReturnRows(pgxmock.NewRows(itemOwnerRowColumns))

	_, err := svc.GetOwned(context.Background(), "user-1", "wl-1", "missing")

	assert.ErrorIs(t, err, apperr.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrForbidden if the item belogns to some other wishlist
func Test_GetOwned_Forbidden(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlist_items i`).
		WithArgs("i-1", "wl-1").
		WillReturnRows(pgxmock.NewRows(itemOwnerRowColumns).
			AddRow("i-1", "wl-1", "Mug", "", "", 10.0, nil, now, nil, "someone-else"))

	_, err := svc.GetOwned(context.Background(), "user-1", "wl-1", "i-1")

	assert.ErrorIs(t, err, apperr.ErrForbidden)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test verifying a wishlist item belongs to a wishlist
func Test_ExistsInWishlist(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("i-1", "wl-1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	err := svc.ExistsInWishlist(context.Background(), "wl-1", "i-1")

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrNotFound if the item is not present in the parent wishlist
func Test_ExistsInWishlist_NotFound(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs("missing", "wl-1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	err := svc.ExistsInWishlist(context.Background(), "wl-1", "missing")

	assert.ErrorIs(t, err, apperr.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test updating a wishlist item
func Test_Update(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlist_items i`).
		WithArgs("i-1", "wl-1").
		WillReturnRows(pgxmock.NewRows(itemOwnerRowColumns).
			AddRow("i-1", "wl-1", "Mug", "", "", 10.0, nil, now, nil, "user-1"))

	mock.ExpectQuery(`SELECT (.+) FROM claims WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(claimRowColumns))

	mock.ExpectQuery(`SELECT (.+) FROM split_interests WHERE item_id = \$1`).
		WithArgs("i-1").
		WillReturnRows(pgxmock.NewRows(splitRowColumns))

	newName := "New Mug"
	mock.ExpectQuery(`UPDATE wishlist_items`).
		WithArgs(newName, pgxmock.AnyArg(), "", 10.0, pgxmock.AnyArg(), pgxmock.AnyArg(), "i-1").
		WillReturnRows(pgxmock.NewRows(itemRowColumns).
			AddRow("i-1", "wl-1", newName, "", "", 10.0, nil, now, nil))

	got, err := svc.Update(context.Background(), "user-1", "wl-1", "i-1", UpdateItemRequest{Name: &newName})

	require.NoError(t, err)
	assert.Equal(t, newName, got.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrForbidden when users somehow attempt to update a wishlist item
// that does not belong to them
func Test_Update_Forbidden(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlist_items i`).
		WithArgs("i-1", "wl-1").
		WillReturnRows(pgxmock.NewRows(itemOwnerRowColumns).
			AddRow("i-1", "wl-1", "Mug", "", "", 10.0, nil, now, nil, "someone-else"))

	_, err := svc.Update(context.Background(), "user-1", "wl-1", "i-1", UpdateItemRequest{})

	assert.ErrorIs(t, err, apperr.ErrForbidden)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test hiding claim info from owner
func Test_HideClaimsForOwner(t *testing.T) {
	t.Parallel()

	items := []Item{
		{
			ID:             "i-1",
			Claim:          &Claim{ID: "c-1"},
			SplitInterests: []SplitInterestEntry{{ID: "s-1"}},
		},
	}

	got := HideClaimsForOwner(items)

	require.Len(t, got, 1)
	assert.Nil(t, got[0].Claim)
	assert.Nil(t, got[0].SplitInterests)
}
