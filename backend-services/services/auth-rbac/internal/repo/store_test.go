package repo

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestStore_Upsert_GetByAuth0(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	ctx := context.Background()

	mail := "a@b.c"
	mock.ExpectQuery(`INSERT INTO users`).WithArgs("sub1", mail, "vendor", pgxmock.AnyArg(), pgxmock.AnyArg(), "global").
		WillReturnRows(pgxmock.NewRows([]string{"updated_at"}).AddRow(time.Now()))

	u := &User{Auth0ID: "sub1", Email: mail, Role: "vendor"}
	require.NoError(t, s.Upsert(ctx, u))

	mock.ExpectQuery(`SELECT auth0_id`).WithArgs("sub1").WillReturnRows(
		pgxmock.NewRows([]string{"auth0_id", "email", "role", "client_id", "vendor_id", "region", "updated_at"}).
			AddRow("sub1", mail, "vendor", nil, nil, "global", time.Now()))

	got, err := s.GetByAuth0(ctx, "sub1")
	require.NoError(t, err)
	require.Equal(t, mail, got.Email)
	require.NoError(t, mock.ExpectationsWereMet())
}
