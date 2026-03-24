package repo

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestStore_InsertConsent_Mock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	mock.ExpectExec(`INSERT INTO consent_records`).WithArgs("auth0|x", "2025-03-01", true).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	require.NoError(t, s.InsertConsent(context.Background(), "auth0|x", "2025-03-01", true))
	require.NoError(t, mock.ExpectationsWereMet())
}
