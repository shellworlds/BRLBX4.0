package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

type Kitchen struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	VendorID   uuid.UUID `json:"vendor_id"`
	CapacityKW float64   `json:"capacity_kw"`
}

type KitchenStore struct {
	DB pgxutil.Querier
}

func (s *KitchenStore) Create(ctx context.Context, k *Kitchen) error {
	if s.DB == nil {
		return fmt.Errorf("nil db")
	}
	const q = `
INSERT INTO kitchens (name, location, vendor_id, capacity_kw)
VALUES ($1, $2, $3, $4)
RETURNING id`
	return s.DB.QueryRow(ctx, q, k.Name, k.Location, k.VendorID, k.CapacityKW).Scan(&k.ID)
}

func (s *KitchenStore) Get(ctx context.Context, id uuid.UUID) (*Kitchen, error) {
	const q = `SELECT id, name, location, vendor_id, capacity_kw FROM kitchens WHERE id=$1`
	var k Kitchen
	err := s.DB.QueryRow(ctx, q, id).Scan(&k.ID, &k.Name, &k.Location, &k.VendorID, &k.CapacityKW)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

// ListByVendor returns kitchens owned by the vendor.
func (s *KitchenStore) ListByVendor(ctx context.Context, vendor uuid.UUID) ([]Kitchen, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("nil db")
	}
	const q = `SELECT id, name, location, vendor_id, capacity_kw FROM kitchens WHERE vendor_id=$1 ORDER BY created_at`
	rows, err := s.DB.Query(ctx, q, vendor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Kitchen
	for rows.Next() {
		var k Kitchen
		if err := rows.Scan(&k.ID, &k.Name, &k.Location, &k.VendorID, &k.CapacityKW); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

// DistinctVendorIDs lists vendors that have at least one kitchen.
func (s *KitchenStore) DistinctVendorIDs(ctx context.Context) ([]uuid.UUID, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("nil db")
	}
	const q = `SELECT DISTINCT vendor_id FROM kitchens`
	rows, err := s.DB.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}
