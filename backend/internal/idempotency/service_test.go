package idempotency

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var idempotencyRowColumns = []string{
	"id", "user_id", "idempotency_key",
	"method", "path", "request_hash", "response_status", "response_body", "created_at",
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

// Test reserving a brand new idempotency key
func Test_Reserve_New(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	mock.ExpectQuery(`INSERT INTO idempotency_keys`).
		WithArgs("user-1", "ikey-1", "POST", "/wishlists", "hash-1").
		WillReturnRows(pgxmock.NewRows(idempotencyRowColumns).
			AddRow("idem-1", "user-1", "ikey-1", "POST", "/wishlists", "hash-1", nil, nil, now))

	got, reserved, err := svc.Reserve(context.Background(), "user-1", "ikey-1", "POST", "/wishlists", "hash-1")

	require.NoError(t, err)
	assert.True(t, reserved)
	assert.Equal(t, "idem-1", got.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test that reserving an already-used key returns the existing record
// instead of creating a new one (the INSERT ... ON CONFLICT DO NOTHING
// returns no rows, so Reserve falls back to calling `service.get`, which is just a SELECT)
func Test_Reserve_Conflict(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	now := time.Now()
	status := 201
	body := []byte(`{"id":"wl-1"}`)

	// this insert should not return any rows in the expected result set
	mock.ExpectQuery(`INSERT INTO idempotency_keys`).
		WithArgs("user-1", "ikey-1", "POST", "/wishlists", "hash-1").
		WillReturnRows(pgxmock.NewRows(idempotencyRowColumns))

	mock.ExpectQuery(`SELECT (.+) FROM idempotency_keys`).
		WithArgs("user-1", "ikey-1", "POST", "/wishlists").
		WillReturnRows(pgxmock.NewRows(idempotencyRowColumns).
			AddRow("idem-1", "user-1", "ikey-1", "POST", "/wishlists", "hash-1", &status, body, now))

	got, reserved, err := svc.Reserve(context.Background(), "user-1", "ikey-1", "POST", "/wishlists", "hash-1")

	require.NoError(t, err)
	assert.False(t, reserved)
	assert.Equal(t, "idem-1", got.ID)
	require.NotNil(t, got.ResponseStatus)
	assert.Equal(t, status, *got.ResponseStatus)
	assert.Equal(t, body, got.ResponseBody)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Test completing a previously reserved request
func Test_Complete(t *testing.T) {
	t.Parallel()

	mock := newMockPool(t)
	svc := NewService(mock)

	body := []byte(`{"ok":true}`)
	mock.ExpectExec(`UPDATE idempotency_keys`).
		WithArgs(200, body, "idem-1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := svc.Complete(context.Background(), "idem-1", 200, body)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}




