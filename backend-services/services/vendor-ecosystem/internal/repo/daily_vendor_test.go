package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestDailyVendorUpsert_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &DailyVendorStore{DB: mock}
	vid := uuid.New()
	day := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)
	cs := 88
	m := &DailyVendorMetric{
		VendorID: vid, Day: day, MealsServed: 10, RevenueTotal: 100, EnergyEfficiencyScore: 0.8,
		ComplianceScore: &cs, Payload: []byte(`{}`),
	}
	mock.ExpectExec(`INSERT INTO daily_vendor_metrics`).
		WithArgs(vid, "2024-01-05", 10, 100.0, 0.8, 88, pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	require.NoError(t, s.Upsert(context.Background(), m))

	mock.ExpectQuery(`SELECT vendor_id`).WithArgs(vid, "2024-01-01", "2024-01-07").WillReturnRows(
		pgxmock.NewRows([]string{"vendor_id", "day", "meals_served", "revenue_total", "energy_efficiency_score", "compliance_score", "payload"}).
			AddRow(vid, day, 10, 100.0, 0.8, int64(88), []byte(`{}`)))
	items, err := s.List(context.Background(), vid, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}
