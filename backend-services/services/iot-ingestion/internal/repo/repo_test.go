package repo

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestStore_InsertRaw_InsertAlert(t *testing.T) {
	mockTS, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockTS.Close()
	mockPG, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPG.Close()
	s := &Store{PG: mockPG, TS: mockTS}
	k := uuid.New()
	mockTS.ExpectExec(`INSERT INTO raw_telemetry`).WithArgs(pgxmock.AnyArg(), k, "t", json.RawMessage(`{}`)).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	require.NoError(t, s.InsertRaw(context.Background(), k, "t", json.RawMessage(`{}`)))

	mockPG.ExpectExec(`INSERT INTO ingestion_alerts`).WithArgs(k, "offline", "msg").WillReturnResult(pgxmock.NewResult("INSERT", 1))
	require.NoError(t, s.InsertAlert(context.Background(), k, "offline", "msg"))
	require.NoError(t, mockTS.ExpectationsWereMet())
	require.NoError(t, mockPG.ExpectationsWereMet())
}

func TestStore_ListAlerts_Ack(t *testing.T) {
	mockPG, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mockPG.Close()
	s := &Store{PG: mockPG}
	id := uuid.New()
	ts := time.Now().UTC()
	mockPG.ExpectQuery(`SELECT id, kitchen_id`).WithArgs(50).WillReturnRows(
		pgxmock.NewRows([]string{"id", "kitchen_id", "level", "message", "created_at", "acknowledged_at"}).
			AddRow(id, uuid.New(), "offline", "x", ts, nil))
	items, err := s.ListAlerts(context.Background(), 50)
	require.NoError(t, err)
	require.Len(t, items, 1)

	mockPG.ExpectExec(`UPDATE ingestion_alerts`).WithArgs(id).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	require.NoError(t, s.AckAlert(context.Background(), id))
	require.NoError(t, mockPG.ExpectationsWereMet())
}
