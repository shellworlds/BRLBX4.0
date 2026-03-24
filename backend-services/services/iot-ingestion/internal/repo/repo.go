package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

type Store struct {
	PG pgxutil.Querier
	TS pgxutil.Querier
}

func (s *Store) InsertRaw(ctx context.Context, kitchen uuid.UUID, topic string, payload json.RawMessage) error {
	const q = `INSERT INTO raw_telemetry (ts, kitchen_id, topic, payload) VALUES ($1,$2,$3,$4)`
	_, err := s.TS.Exec(ctx, q, time.Now().UTC(), kitchen, topic, payload)
	return err
}

func (s *Store) InsertAlert(ctx context.Context, kitchen uuid.UUID, level, message string) error {
	const q = `INSERT INTO ingestion_alerts (kitchen_id, level, message) VALUES ($1,$2,$3)`
	_, err := s.PG.Exec(ctx, q, kitchen, level, message)
	return err
}

type Alert struct {
	ID          uuid.UUID  `json:"id"`
	KitchenID   uuid.UUID  `json:"kitchen_id"`
	Level       string     `json:"level"`
	Message     string     `json:"message"`
	CreatedAt   time.Time  `json:"created_at"`
	AckedAt     *time.Time `json:"acknowledged_at,omitempty"`
}

func (s *Store) ListAlerts(ctx context.Context, limit int) ([]Alert, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	const q = `
SELECT id, kitchen_id, level, message, created_at, acknowledged_at
FROM ingestion_alerts
ORDER BY created_at DESC
LIMIT $1`
	rows, err := s.PG.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Alert
	for rows.Next() {
		var a Alert
		if err := rows.Scan(&a.ID, &a.KitchenID, &a.Level, &a.Message, &a.CreatedAt, &a.AckedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) AckAlert(ctx context.Context, id uuid.UUID) error {
	const q = `UPDATE ingestion_alerts SET acknowledged_at=now() WHERE id=$1 AND acknowledged_at IS NULL`
	_, err := s.PG.Exec(ctx, q, id)
	return err
}
