package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpsertKitchenDevice inserts or updates a device row by (kitchen_id, label).
func (s *Store) UpsertKitchenDevice(ctx context.Context, kitchenID uuid.UUID, label, csrPEM, certPEM, serial, status string) (uuid.UUID, error) {
	if s.PG == nil {
		return uuid.Nil, fmt.Errorf("nil pg")
	}
	if label == "" {
		label = "default"
	}
	const q = `
INSERT INTO kitchen_devices (kitchen_id, label, csr_pem, cert_pem, serial_number, status)
VALUES ($1,$2,$3,$4,$5,$6)
ON CONFLICT (kitchen_id, label) DO UPDATE SET
  csr_pem = EXCLUDED.csr_pem,
  cert_pem = EXCLUDED.cert_pem,
  serial_number = EXCLUDED.serial_number,
  status = EXCLUDED.status
RETURNING id`
	var id uuid.UUID
	err := s.PG.QueryRow(ctx, q, kitchenID, label, csrPEM, strOrNil(certPEM), strOrNil(serial), status).Scan(&id)
	return id, err
}

func strOrNil(s string) any {
	if s == "" {
		return nil
	}
	return s
}
