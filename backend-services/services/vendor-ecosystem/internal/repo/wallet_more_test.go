package repo

import (
	"context"
	"testing"

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
