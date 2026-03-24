package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

// EmissionFactorStore reads grid emission factors per region (GHG Protocol Scope 2 / market-based hints).
type EmissionFactorStore struct {
	DB pgxutil.Querier
}

func (s *EmissionFactorStore) GridGPerKWh(ctx context.Context, region string) (float64, error) {
	if s.DB == nil {
		return 450, nil
	}
	if region == "" {
		region = "global"
	}
	const q = `SELECT grid_g_co2e_per_kwh::float8 FROM emission_factors WHERE region=$1`
	var v float64
	err := s.DB.QueryRow(ctx, q, region).Scan(&v)
	if err != nil {
		if err == pgx.ErrNoRows {
			return s.GridGPerKWh(ctx, "global")
		}
		return 0, err
	}
	return v, nil
}
