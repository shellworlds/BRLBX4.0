package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestGetWallet_Mock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	vid := uuid.New()
	mock.ExpectQuery(`SELECT vendor_id, balance::float8, pending_payout::float8`).
		WithArgs(vid).WillReturnRows(
		pgxmock.NewRows([]string{"vendor_id", "balance", "pending_payout"}).AddRow(vid, 12.5, 1.0),
	)
	w, err := s.GetWallet(context.Background(), vid)
	require.NoError(t, err)
	require.InDelta(t, 12.5, w.Balance, 1e-9)
	require.InDelta(t, 1.0, w.PendingPayout, 1e-9)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListPendingPayouts_Mock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	pid := uuid.New()
	vid := uuid.New()
	ts := time.Now().UTC().Truncate(time.Second)
	mock.ExpectQuery(`SELECT id, vendor_id, amount::float8, status, stripe_transfer_id, failure_reason, created_at`).
		WithArgs(50).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "vendor_id", "amount", "status", "stripe_transfer_id", "failure_reason", "created_at"}).
				AddRow(pid, vid, 10.0, "pending", nil, nil, ts),
		)
	rows, err := s.ListPendingPayouts(context.Background(), 50)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, pid, rows[0].ID)
	require.Equal(t, "pending", rows[0].Status)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListPayoutsForVendor_Mock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	vid := uuid.New()
	pid := uuid.New()
	ts := time.Now().UTC().Truncate(time.Second)
	mock.ExpectQuery(`SELECT id, vendor_id, amount::float8, status, stripe_transfer_id, failure_reason, created_at`).
		WithArgs(vid, 50).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "vendor_id", "amount", "status", "stripe_transfer_id", "failure_reason", "created_at"}).
				AddRow(pid, vid, 5.0, "paid", strPtr("tr_1"), nil, ts),
		)
	rows, err := s.ListPayoutsForVendor(context.Background(), vid, 50)
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "paid", rows[0].Status)
	require.NotNil(t, rows[0].StripeTransferID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func strPtr(s string) *string { return &s }
