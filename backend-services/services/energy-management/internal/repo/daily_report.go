package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

// DailyClientReport stores Scope 3–style ESG hints plus operational mix for a logical client (vendor UUID string).
type DailyClientReport struct {
	ClientID         string          `json:"client_id"`
	Day              time.Time       `json:"day"`
	Scope3TCO2eAvoid float64         `json:"scope3_tco2e_avoided"`
	SolarShare       float64         `json:"solar_share"`
	GridShare        float64         `json:"grid_share"`
	BatteryShare     float64         `json:"battery_share"`
	UptimeAvg        float64         `json:"uptime_avg"`
	Payload          json.RawMessage `json:"payload,omitempty"`
}

type DailyReportStore struct {
	DB pgxutil.Querier
}

func (s *DailyReportStore) Upsert(ctx context.Context, r *DailyClientReport) error {
	if s.DB == nil {
		return fmt.Errorf("nil db")
	}
	if len(r.Payload) == 0 {
		r.Payload = []byte(`{}`)
	}
	const q = `
INSERT INTO daily_client_reports (client_id, day, scope3_tco2e_avoided, solar_share, grid_share, battery_share, uptime_avg, payload)
VALUES ($1, $2::date, $3, $4, $5, $6, $7, $8)
ON CONFLICT (client_id, day) DO UPDATE SET
  scope3_tco2e_avoided = EXCLUDED.scope3_tco2e_avoided,
  solar_share = EXCLUDED.solar_share,
  grid_share = EXCLUDED.grid_share,
  battery_share = EXCLUDED.battery_share,
  uptime_avg = EXCLUDED.uptime_avg,
  payload = EXCLUDED.payload,
  created_at = now()`
	day := time.Date(r.Day.Year(), r.Day.Month(), r.Day.Day(), 0, 0, 0, 0, time.UTC)
	_, err := s.DB.Exec(ctx, q, r.ClientID, day.Format("2006-01-02"), r.Scope3TCO2eAvoid, r.SolarShare, r.GridShare, r.BatteryShare, r.UptimeAvg, r.Payload)
	return err
}

func (s *DailyReportStore) List(ctx context.Context, clientID string, from, to time.Time) ([]DailyClientReport, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("nil db")
	}
	const q = `
SELECT client_id, day, scope3_tco2e_avoided, solar_share, grid_share, battery_share, uptime_avg, payload
FROM daily_client_reports
WHERE client_id=$1 AND day >= $2::date AND day <= $3::date
ORDER BY day ASC`
	rows, err := s.DB.Query(ctx, q, clientID, from.Format("2006-01-02"), to.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DailyClientReport
	for rows.Next() {
		var r DailyClientReport
		if err := rows.Scan(&r.ClientID, &r.Day, &r.Scope3TCO2eAvoid, &r.SolarShare, &r.GridShare, &r.BatteryShare, &r.UptimeAvg, &r.Payload); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
