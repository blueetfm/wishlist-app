package wishlist

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
)

var wishlistRowColumns = []string{
	"id", "user_id", "title", "description", "share_token",
	"is_shared", "hide_claims_from_owner", "created_at",
}

// Create a mock db connection pool.
// Returns the full pgxmock interface (not db.Pool) since tests need
// Expect*/ExpectationsWereMet, not just the QueryRow/Query/Exec that
// db.Pool exposes to services
// pgxmock.pgxPoolIFace still satisfies db.Pool when passed
// to NewService in all the below test funcs
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

// Test creating a new wishlist
func Test_Create(t *testing.T) {
	// 1. tests will run concurrenlty, speeding up test suite
	t.Parallel()

	// 2. setting up test pools and serivce
	mock := newMockPool(t)
	svc := NewService(mock)

	// 3. set up test req data
	now := time.Now()
	req := CreateWishlistRequest{Title: "Birthday", Description: "gifts pls"}

	// 4. set up database mock
	mock.ExpectQuery(`INSERT INTO wishlists`).
		WithArgs("user-1", req.Title, req.Description).
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			// populating my fake db 
			AddRow("wl-1", "user-1", req.Title, req.Description, "tok", false, false, now))

	// 5. calls service with a blank context 
	got, err := svc.Create(context.Background(), "user-1", req)

	// 6. asserts
	require.NoError(t, err)
	assert.Equal(t, "wl-1", got.ID)
	assert.Equal(t, req.Title, got.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test listing wishists belonging to a given user_id
func Test_ListByOwner(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("user-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "A", "", "tok-1", false, false, now).
			AddRow("wl-2", "user-1", "B", "", "tok-2", true, false, now))

	got, err := svc.ListByOwner(context.Background(), "user-1")

	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "wl-1", got[0].ID)
	assert.Equal(t, "wl-2", got[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning owner status of a wishlist, given a user and wishlist id
func Test_GetOwned_Success(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "A", "", "tok", false, false, now))

	got, err := svc.GetOwned(context.Background(), "user-1", "wl-1")

	require.NoError(t, err)
	assert.Equal(t, "user-1", got.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrNotFound if a given wishlist id does not exist, given a user id
func Test_GetOwned_NotFound(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("missing").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns))

	_, err := svc.GetOwned(context.Background(), "user-1", "missing")

	assert.ErrorIs(t, err, apperr.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrForbidden if a given wishlist id exists but does not belong to a given user id
func Test_GetOwned_Forbidden(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "someone-else", "A", "", "tok", false, false, now))

	_, err := svc.GetOwned(context.Background(), "user-1", "wl-1")

	assert.ErrorIs(t, err, apperr.ErrForbidden)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning a wishlist's details when given a shar etoken
// includes sub-tests here, just another alternative to having
// 2 separate functions
func Test_GetByShareToken(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	t.Run("found", func(t *testing.T) {
		now := time.Now()
		mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
			WithArgs("tok-1").
			WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
				AddRow("wl-1", "user-1", "A", "", "tok-1", true, false, now))

		got, err := svc.GetByShareToken(context.Background(), "tok-1")

		require.NoError(t, err)
		assert.Equal(t, "wl-1", got.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
			WithArgs("bad-token").
			WillReturnRows(pgxmock.NewRows(wishlistRowColumns))

		_, err := svc.GetByShareToken(context.Background(), "bad-token")

		assert.ErrorIs(t, err, apperr.ErrNotFound)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test updating a wishlist
func Test_Update(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "Old", "old desc", "tok", false, false, now))

	newTitle := "New"
	mock.ExpectQuery(`UPDATE wishlists`).
		WithArgs(newTitle, "old desc", false, "wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", newTitle, "old desc", "tok", false, false, now))

	got, err := svc.Update(context.Background(), "user-1", "wl-1", UpdateWishlistRequest{Title: &newTitle})

	require.NoError(t, err)
	assert.Equal(t, newTitle, got.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrForbidden when users somehow attempt to update a wishlist 
// that does not belong to them
func Test_Update_NotOwned(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "someone-else", "Old", "old desc", "tok", false, false, now))

	_, err := svc.Update(context.Background(), "user-1", "wl-1", UpdateWishlistRequest{})

	assert.ErrorIs(t, err, apperr.ErrForbidden)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test deleting a wishlist
func Test_Delete(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-1", "A", "", "tok", false, false, now))

	mock.ExpectExec(`DELETE FROM wishlists`).
		WithArgs("wl-1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := svc.Delete(context.Background(), "user-1", "wl-1")

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrNotFound when users somehow try to delete a wishlist
// that doesn't exist
func Test_Delete_NotOwned(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("missing").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns))

	err := svc.Delete(context.Background(), "user-1", "missing")

	assert.True(t, errors.Is(err, apperr.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test returning ErrForbidden when users somehow try to delete a wishlist
// that doesn't belong to them
func Test_Delete_Forbidden(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM wishlists`).
		WithArgs("wl-1").
		WillReturnRows(pgxmock.NewRows(wishlistRowColumns).
			AddRow("wl-1", "user-2", "A", "", "tok", false, false, now))

	err := svc.Delete(context.Background(), "user-1", "wl-1")

	assert.ErrorIs(t, err, apperr.ErrForbidden)
	assert.NoError(t, mock.ExpectationsWereMet())
}
