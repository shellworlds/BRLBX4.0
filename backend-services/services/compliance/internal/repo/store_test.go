package repo

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestStore_LatestConsent_NoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	mock.ExpectQuery(`SELECT policy_version, accepted FROM consent_records`).
		WithArgs("sub1").WillReturnError(pgx.ErrNoRows)
	ver, acc, ok, err := s.LatestConsent(context.Background(), "sub1")
	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, "", ver)
	require.False(t, acc)
	require.NoError(t, mock.ExpectationsWereMet())
}
