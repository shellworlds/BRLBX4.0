package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

type DailyVendorMetric struct {
	VendorID              uuid.UUID       `json:"vendor_id"`
	Day                   time.Time       `json:"day"`
	MealsServed           int             `json:"meals_served"`
	RevenueTotal          float64         `json:"revenue_total"`
	EnergyEfficiencyScore float64         `json:"energy_efficiency_score"`
	ComplianceScore       *int            `json:"compliance_score,omitempty"`
	Payload               json.RawMessage `json:"payload,omitempty"`
}

type DailyVendorStore struct {
	DB pgxutil.Querier
}

func (s *DailyVendorStore) Upsert(ctx context.Context, m *DailyVendorMetric) error {
	if s.DB == nil {
		return fmt.Errorf("nil db")
	}
	if len(m.Payload) == 0 {
		m.Payload = []byte(`{}`)
	}
	day := time.Date(m.Day.Year(), m.Day.Month(), m.Day.Day(), 0, 0, 0, 0, time.UTC)
	var comp any
	if m.ComplianceScore != nil {
		comp = *m.ComplianceScore
	}
	const q = `
INSERT INTO daily_vendor_metrics (vendor_id, day, meals_served, revenue_total, energy_efficiency_score, compliance_score, payload)
VALUES ($1, $2::date, $3, $4, $5, $6, $7)
ON CONFLICT (vendor_id, day) DO UPDATE SET
  meals_served = EXCLUDED.meals_served,
  revenue_total = EXCLUDED.revenue_total,
  energy_efficiency_score = EXCLUDED.energy_efficiency_score,
  compliance_score = EXCLUDED.compliance_score,
  payload = EXCLUDED.payload,
  created_at = now()`
	_, err := s.DB.Exec(ctx, q, m.VendorID, day.Format("2006-01-02"), m.MealsServed, m.RevenueTotal, m.EnergyEfficiencyScore, comp, m.Payload)
	return err
}

func (s *DailyVendorStore) List(ctx context.Context, vendor uuid.UUID, from, to time.Time) ([]DailyVendorMetric, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("nil db")
	}
	const q = `
SELECT vendor_id, day, meals_served, revenue_total, energy_efficiency_score, compliance_score, payload
FROM daily_vendor_metrics
WHERE vendor_id=$1 AND day >= $2::date AND day <= $3::date
ORDER BY day ASC`
	rows, err := s.DB.Query(ctx, q, vendor, from.Format("2006-01-02"), to.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DailyVendorMetric
	for rows.Next() {
		var m DailyVendorMetric
		var comp sql.NullInt64
		if err := rows.Scan(&m.VendorID, &m.Day, &m.MealsServed, &m.RevenueTotal, &m.EnergyEfficiencyScore, &comp, &m.Payload); err != nil {
			return nil, err
		}
		if comp.Valid {
			v := int(comp.Int64)
			m.ComplianceScore = &v
		}
		out = append(out, m)
	}
	return out, rows.Err()
}
