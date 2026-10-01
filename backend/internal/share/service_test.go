package share

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
	"github.com/blueetfm/wishlist-app/backend/internal/item"
	"github.com/blueetfm/wishlist-app/backend/internal/wishlist"
)

var wishlistRowColumns = []string{
	"id", "user_id", "title", "description", "share_token",
	"is_shared", "hide_claims_from_owner", "created_at",
}

var itemRowColumns = []string{
	"id", "wishlist_id", "name", "description", "product_url",
	"price", "embed_data", "created_at", "image_url",
}

// newTestService wires up a Service backed by the same mock pool as its
// WishlistService and ItemService collaborators.
func newTestService(mock pgxmock.PgxPoolIface) *Service {
	items := item.NewService(mock, item.NewClaimService(mock), item.NewSplitService(mock))
	return NewService(mock, wishlist.NewService(mock), items)
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

// Test getting current share status for a wishlist
func Test_GetActive(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "A", "", "tok-1", true, false, now))

	got, err := svc.GetActive(context.Background(), "user-1", "wl-1")

	require.NoError(t, err)
	assert.Equal(t, "tok-1", got.Token)
	assert.True(t, got.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrNotFound if wishlist does not belong
func Test_GetActive_ErrNotFound(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("missing").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns))

	_, err := svc.GetActive(context.Background(), "user-1", "missing")

	assert.ErrorIs(t, err, apperr.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test that CreateOrGet will return an existing token if it's already shared
func Test_CreateOrGet_Get(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "A", "", "tok-1", true, false, now))

	got, err := svc.CreateOrGet(context.Background(), "user-1", "wl-1")

	require.NoError(t, err)
	assert.Equal(t, "tok-1", got.Token)
	assert.Equal(t, true, got.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test that CreateOrGet will return a new token if it was not shared
func Test_CreateOrGet_Create(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "A", "", "tok-1", false, false, now))

	// fake the newly rotated token the UPDATE ... RETURNING would produce
	mock.ExpectQuery(`UPDATE wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows([]string{"share_token", "is_shared"}).
			AddRow("tok-2", true))

	got, err := svc.CreateOrGet(context.Background(), "user-1", "wl-1")

	require.NoError(t, err)
	assert.Equal(t, "tok-2", got.Token)
	assert.Equal(t, true, got.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test that revoking share status for a wishlist works
func Test_Revoke(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "A", "", "tok-1", true, false, now))

	mock.ExpectExec(`UPDATE wishlists`).
		WithArgs("wl-1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := svc.Revoke(context.Background(), "user-1", "wl-1")

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test that guests can retrieve wishlist details when using publix share token
func Test_ResolvePublic(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := newTestService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("tok-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "A", "", "tok-1", true, false, now))

	mock.ExpectQuery(`SELECT (.+) FROM wishlist_items WHERE wishlist_id = \$1`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(itemRowColumns))

	got, err := svc.ResolvePublic(context.Background(), "tok-1", "user-1")

	require.NoError(t, err)
	assert.Equal(t, "wl-1", got.ID)
	assert.Empty(t, got.Items)
	assert.NoError(t, mock.ExpectationsWereMet())
}
