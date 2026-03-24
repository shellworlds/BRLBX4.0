package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestReadingInsert_List_Uptime(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &ReadingStore{DB: mock}
	k := uuid.New()
	ts := time.Now().UTC()
	mock.ExpectExec(`INSERT INTO energy_readings`).WithArgs(k, ts, 1.0, 2.0, 3.0, "ok", 99.0).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	require.NoError(t, s.Insert(context.Background(), &EnergyReading{
		KitchenID: k, TS: ts, GridPower: 1, BatteryPower: 2, SolarPower: 3, LPGStatus: "ok", UptimePct: 99,
	}))

	from, to := ts.Add(-time.Hour), ts.Add(time.Hour)
	mock.ExpectQuery(`SELECT kitchen_id`).WithArgs(k, from, to).WillReturnRows(
		pgxmock.NewRows([]string{"kitchen_id", "ts", "grid_power", "battery_power", "solar_power", "lpg_status", "uptime_percent"}).
			AddRow(k, ts, 1.0, 2.0, 3.0, "ok", 99.0))
	rows, err := s.List(context.Background(), k, from, to)
	require.NoError(t, err)
	require.Len(t, rows, 1)

	mock.ExpectQuery(`SELECT COALESCE`).WithArgs(k, from, to).
		WillReturnRows(pgxmock.NewRows([]string{"avg"}).AddRow(99.0))
	avgUp, err := s.AggregateUptime(context.Background(), k, from, to)
	require.NoError(t, err)
	require.InDelta(t, 99.0, avgUp, 1e-6)
	require.NoError(t, mock.ExpectationsWereMet())
}
