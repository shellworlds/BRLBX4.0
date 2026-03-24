package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Kitchen struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	VendorID   uuid.UUID `json:"vendor_id"`
	CapacityKW float64   `json:"capacity_kw"`
}

type KitchenStore struct {
	Pool *pgxpool.Pool
}

func (s *KitchenStore) Create(ctx context.Context, k *Kitchen) error {
	if s.Pool == nil {
		return fmt.Errorf("nil pool")
	}
	const q = `
INSERT INTO kitchens (name, location, vendor_id, capacity_kw)
VALUES ($1, $2, $3, $4)
RETURNING id`
	return s.Pool.QueryRow(ctx, q, k.Name, k.Location, k.VendorID, k.CapacityKW).Scan(&k.ID)
}

func (s *KitchenStore) Get(ctx context.Context, id uuid.UUID) (*Kitchen, error) {
	const q = `SELECT id, name, location, vendor_id, capacity_kw FROM kitchens WHERE id=$1`
	var k Kitchen
	err := s.Pool.QueryRow(ctx, q, id).Scan(&k.ID, &k.Name, &k.Location, &k.VendorID, &k.CapacityKW)
	if err != nil {
		return nil, err
	}
	return &k, nil
}
