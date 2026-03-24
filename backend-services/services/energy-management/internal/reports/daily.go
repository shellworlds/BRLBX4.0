package reports

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shellworlds/BRLBX4.0/backend-services/services/energy-management/internal/repo"
)

const (
	gridEmissionsKgPerKWh = 0.42 // grid intensity stub
	solarAvoidedKgPerKWh  = 0.48 // avoided grid + marginal renewables stub
)

// RunDailyClientAggregate rolls the previous UTC calendar day's telemetry into daily_client_reports per vendor client_id.
func RunDailyClientAggregate(ctx context.Context, kitchens *repo.KitchenStore, readings *repo.ReadingStore, reports *repo.DailyReportStore, factors *repo.EmissionFactorStore, dayUTC time.Time) error {
	if kitchens == nil || readings == nil || reports == nil {
		return fmt.Errorf("nil store")
	}
	vendors, err := kitchens.DistinctVendorIDs(ctx)
	if err != nil {
		return err
	}
	start := time.Date(dayUTC.Year(), dayUTC.Month(), dayUTC.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	for _, vendor := range vendors {
		kitch, err := kitchens.ListByVendor(ctx, vendor)
		if err != nil {
			return err
		}
		var solarSum, gridSum, battSum, uptimeSum float64
		var sampleCount int
		var scope3 float64
		for _, k := range kitch {
			gridKg := gridEmissionsKgPerKWh
			if factors != nil {
				if g, err := factors.GridGPerKWh(ctx, k.Region); err == nil {
					gridKg = g / 1000.0
				}
			}
			rows, err := readings.List(ctx, k.ID, start, end)
			if err != nil {
				return err
			}
			for _, r := range rows {
				solarSum += r.SolarPower
				gridSum += r.GridPower
				battSum += r.BatteryPower
				uptimeSum += r.UptimePct
				sampleCount++

				d := r.SolarPower + r.GridPower + r.BatteryPower
				if d <= 0 {
					continue
				}
				sS := r.SolarPower / d
				gS := r.GridPower / d
				assumeKWh := d * 1.0 // kW sample as rough energy proxy (stub)
				rawKg := assumeKWh * (sS*solarAvoidedKgPerKWh - gS*gridKg)
				if rawKg > 0 {
					scope3 += rawKg / 1000.0
				}
			}
		}
		denomPL := solarSum + gridSum + battSum
		var solarShare, gridShare, battShare float64
		if denomPL > 0 {
			solarShare = solarSum / denomPL
			gridShare = gridSum / denomPL
			battShare = battSum / denomPL
		}
		var uptimeAvg float64
		if sampleCount > 0 {
			uptimeAvg = uptimeSum / float64(sampleCount)
		}

		payload, _ := json.Marshal(map[string]any{
			"esg": map[string]any{
				"energy_mix": map[string]any{"solar_kw_avg": solarSum / float64(max(1, sampleCount)), "grid_kw_avg": gridSum / float64(max(1, sampleCount))},
			},
			"kitchens_aggregated": len(kitch),
			"sample_count":        sampleCount,
		})
		rep := &repo.DailyClientReport{
			ClientID:         vendor.String(),
			Day:              start,
			Scope3TCO2eAvoid: scope3,
			SolarShare:       solarShare,
			GridShare:        gridShare,
			BatteryShare:     battShare,
			UptimeAvg:        uptimeAvg,
			Payload:          payload,
		}
		if err := reports.Upsert(ctx, rep); err != nil {
			return err
		}
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ParseClientID accepts a UUID string for client routes.
func ParseClientID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
