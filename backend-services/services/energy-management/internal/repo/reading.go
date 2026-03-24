package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EnergyReading struct {
	KitchenID     uuid.UUID `json:"kitchen_id"`
	TS            time.Time `json:"timestamp"`
	GridPower     float64   `json:"grid_power"`
	BatteryPower  float64   `json:"battery_power"`
	SolarPower    float64   `json:"solar_power"`
	LPGStatus     string    `json:"lpg_status"`
	UptimePct     float64   `json:"uptime_percent"`
}

type ReadingStore struct {
	Pool *pgxpool.Pool
}

func (s *ReadingStore) Insert(ctx context.Context, r *EnergyReading) error {
	if r == nil {
		return fmt.Errorf("nil reading")
	}
	if r.TS.IsZero() {
		r.TS = time.Now().UTC()
	}
	const q = `
INSERT INTO energy_readings (kitchen_id, ts, grid_power, battery_power, solar_power, lpg_status, uptime_percent)
VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.Pool.Exec(ctx, q, r.KitchenID, r.TS, r.GridPower, r.BatteryPower, r.SolarPower, r.LPGStatus, r.UptimePct)
	return err
}

func (s *ReadingStore) List(ctx context.Context, kitchen uuid.UUID, from, to time.Time) ([]EnergyReading, error) {
	const q = `
SELECT kitchen_id, ts, grid_power, battery_power, solar_power, lpg_status, uptime_percent
FROM energy_readings
WHERE kitchen_id=$1 AND ts >= $2 AND ts <= $3
ORDER BY ts ASC`
	rows, err := s.Pool.Query(ctx, q, kitchen, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EnergyReading
	for rows.Next() {
		var r EnergyReading
		if err := rows.Scan(&r.KitchenID, &r.TS, &r.GridPower, &r.BatteryPower, &r.SolarPower, &r.LPGStatus, &r.UptimePct); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// AggregateUptime returns average uptime_percent over window.
func (s *ReadingStore) AggregateUptime(ctx context.Context, kitchen uuid.UUID, from, to time.Time) (float64, error) {
	const q = `
SELECT COALESCE(AVG(uptime_percent), 0)::float8
FROM energy_readings
WHERE kitchen_id=$1 AND ts >= $2 AND ts <= $3`
	var v float64
	if err := s.Pool.QueryRow(ctx, q, kitchen, from, to).Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// AverageGridKW returns average grid power (kW) over window for LCOE stub.
func (s *ReadingStore) AverageGridKW(ctx context.Context, kitchen uuid.UUID, from, to time.Time) (float64, error) {
	const q = `
SELECT COALESCE(AVG(grid_power), 0)::float8
FROM energy_readings
WHERE kitchen_id=$1 AND ts >= $2 AND ts <= $3`
	var v float64
	if err := s.Pool.QueryRow(ctx, q, kitchen, from, to).Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}
