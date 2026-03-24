package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
)

// RunDailyVendorAggregate rolls vendor transaction totals into daily_vendor_metrics.
func RunDailyVendorAggregate(ctx context.Context, store *repo.Store, daily *repo.DailyVendorStore, dayUTC time.Time) error {
	if store == nil || daily == nil {
		return fmt.Errorf("nil store")
	}
	ids, err := store.ListVendorIDs(ctx)
	if err != nil {
		return err
	}
	start := time.Date(dayUTC.Year(), dayUTC.Month(), dayUTC.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	for _, vid := range ids {
		meals, revenue, err := store.AggregateTransactionsRange(ctx, vid, start, end)
		if err != nil {
			return err
		}
		v, err := store.GetVendor(ctx, vid)
		if err != nil {
			return err
		}
		eff := 0.7
		if v.FSSAIScore > 85 {
			eff = 0.9
		}
		cs := v.FSSAIScore
		payload, _ := json.Marshal(map[string]any{
			"vendor_name": v.Name,
			"location":    v.Location,
		})
		m := &repo.DailyVendorMetric{
			VendorID:              vid,
			Day:                   start,
			MealsServed:           meals,
			RevenueTotal:          revenue,
			EnergyEfficiencyScore: eff,
			ComplianceScore:       &cs,
			Payload:               payload,
		}
		if err := daily.Upsert(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

// ParseVendorID validates UUID path param.
func ParseVendorID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
