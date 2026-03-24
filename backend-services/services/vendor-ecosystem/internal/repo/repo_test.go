package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestStore_CreateVendor_GetVendor(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	ctx := context.Background()

	id := uuid.New()
	mock.ExpectQuery(`INSERT INTO vendors`).WithArgs("n", 90, "L", "c").WillReturnRows(
		pgxmock.NewRows([]string{"id", "onboarding_date"}).AddRow(id, time.Now()))

	v := &Vendor{Name: "n", FSSAIScore: 90, Location: "L", Contact: "c"}
	require.NoError(t, s.CreateVendor(ctx, v))
	require.Equal(t, id, v.ID)

	mock.ExpectQuery(`SELECT id, name`).WithArgs(id).WillReturnRows(
		pgxmock.NewRows([]string{"id", "name", "fssai_score", "location", "contact", "onboarding_date"}).
			AddRow(id, "n", 90, "L", "c", time.Now()))

	got, err := s.GetVendor(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "n", got.Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_AggregateTransactionsRange(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	vid := uuid.New()
	mock.ExpectQuery(`SELECT COALESCE`).WithArgs(vid, pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnRows(pgxmock.NewRows([]string{"meals", "rev"}).AddRow(12, 34.5))
	meals, rev, err := s.AggregateTransactionsRange(context.Background(), vid, time.Now(), time.Now().Add(time.Hour))
	require.NoError(t, err)
	require.Equal(t, 12, meals)
	require.InDelta(t, 34.5, rev, 1e-6)
	require.NoError(t, mock.ExpectationsWereMet())
}
