package repo

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestStore_GetActiveSubscription_NoRows(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	mock.ExpectQuery(`SELECT id, client_id, plan, status, start_date, next_billing, stripe_customer_id, stripe_subscription_id`).
		WithArgs("client-1").WillReturnError(pgx.ErrNoRows)
	sub, err := s.GetActiveSubscription(context.Background(), "client-1")
	require.NoError(t, err)
	require.Nil(t, sub)
	require.NoError(t, mock.ExpectationsWereMet())
}
