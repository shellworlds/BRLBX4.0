package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestKitchenCreate_Mock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	k := &Kitchen{Name: "k1", Location: "L", VendorID: uuid.New(), CapacityKW: 10}
	rows := pgxmock.NewRows([]string{"id"}).AddRow(uuid.New())
	mock.ExpectQuery(`INSERT INTO kitchens`).WithArgs(k.Name, k.Location, k.VendorID, k.CapacityKW).WillReturnRows(rows)

	s := &KitchenStore{DB: mock}
	require.NoError(t, s.Create(context.Background(), k))
	require.NotEqual(t, uuid.Nil, k.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}
