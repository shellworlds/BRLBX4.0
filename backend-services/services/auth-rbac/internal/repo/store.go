package repo

import (
	"context"
	"time"

	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

type User struct {
	Auth0ID   string  `json:"auth0_id"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	ClientID  *string `json:"client_id,omitempty"`
	VendorID  *string `json:"vendor_id,omitempty"`
	Region    string  `json:"region"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store struct {
	DB pgxutil.Querier
}

func (s *Store) Upsert(ctx context.Context, u *User) error {
	if u.Region == "" {
		u.Region = "global"
	}
	const q = `
INSERT INTO users (auth0_id, email, role, client_id, vendor_id, region, updated_at)
VALUES ($1,$2,$3,$4,$5,$6,now())
ON CONFLICT (auth0_id) DO UPDATE
SET email=excluded.email,
    role=excluded.role,
    client_id=excluded.client_id,
    vendor_id=excluded.vendor_id,
    region=excluded.region,
    updated_at=now()
RETURNING updated_at`
	return s.DB.QueryRow(ctx, q, u.Auth0ID, u.Email, u.Role, u.ClientID, u.VendorID, u.Region).Scan(&u.UpdatedAt)
}

func (s *Store) GetByAuth0(ctx context.Context, auth0ID string) (*User, error) {
	const q = `SELECT auth0_id, email, role, client_id, vendor_id, region, updated_at FROM users WHERE auth0_id=$1`
	var u User
	if err := s.DB.QueryRow(ctx, q, auth0ID).Scan(&u.Auth0ID, &u.Email, &u.Role, &u.ClientID, &u.VendorID, &u.Region, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
