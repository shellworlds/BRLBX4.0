package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

type Vendor struct {
	ID                     uuid.UUID `json:"id"`
	Name                   string    `json:"name"`
	FSSAIScore             int       `json:"fssai_score"`
	Location               string    `json:"location"`
	Contact                string    `json:"contact"`
	OnboardingDate         time.Time `json:"onboarding_date"`
	Region                 string    `json:"region"`
	StripeConnectAccountID *string   `json:"stripe_connect_account_id,omitempty"`
}

type Financing struct {
	ID                uuid.UUID `json:"id"`
	VendorID          uuid.UUID `json:"vendor_id"`
	Amount            float64   `json:"amount"`
	Status            string    `json:"status"`
	RepaymentSchedule string    `json:"repayment_schedule"`
	RemainingBalance  float64   `json:"remaining_balance"`
	CreatedAt         time.Time `json:"created_at"`
}

type Transaction struct {
	ID        uuid.UUID `json:"id"`
	VendorID  uuid.UUID `json:"vendor_id"`
	KitchenID uuid.UUID `json:"kitchen_id"`
	Amount    float64   `json:"amount"`
	MealCount int       `json:"meal_count"`
	TS        time.Time `json:"timestamp"`
}

type Store struct {
	DB pgxutil.Querier
}

func (s *Store) CreateVendor(ctx context.Context, v *Vendor) error {
	const q = `INSERT INTO vendors (name, fssai_score, location, contact) VALUES ($1,$2,$3,$4) RETURNING id, onboarding_date`
	return s.DB.QueryRow(ctx, q, v.Name, v.FSSAIScore, v.Location, v.Contact).Scan(&v.ID, &v.OnboardingDate)
}

func (s *Store) GetVendor(ctx context.Context, id uuid.UUID) (*Vendor, error) {
	const q = `
SELECT id, name, fssai_score, location, contact, onboarding_date, region, stripe_connect_account_id
FROM vendors WHERE id=$1`
	var v Vendor
	if err := s.DB.QueryRow(ctx, q, id).Scan(
		&v.ID, &v.Name, &v.FSSAIScore, &v.Location, &v.Contact, &v.OnboardingDate, &v.Region, &v.StripeConnectAccountID,
	); err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Store) AvgTransactionVolume(ctx context.Context, vendor uuid.UUID, from, to time.Time) (float64, error) {
	const q = `SELECT COALESCE(AVG(amount), 0) FROM transactions WHERE vendor_id=$1 AND ts >= $2 AND ts <= $3`
	var v float64
	if err := s.DB.QueryRow(ctx, q, vendor, from, to).Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (s *Store) CreateFinancing(ctx context.Context, f *Financing) error {
	const q = `
INSERT INTO financing (vendor_id, amount, status, repayment_schedule, remaining_balance)
VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`
	rb := f.Amount
	if f.Status == "rejected" {
		rb = 0
	}
	return s.DB.QueryRow(ctx, q, f.VendorID, f.Amount, f.Status, f.RepaymentSchedule, rb).Scan(&f.ID, &f.CreatedAt)
}

func (s *Store) ListFinancing(ctx context.Context, vendor uuid.UUID) ([]Financing, error) {
	const q = `
SELECT id, vendor_id, amount, status, repayment_schedule, remaining_balance, created_at
FROM financing WHERE vendor_id=$1 ORDER BY created_at DESC`
	rows, err := s.DB.Query(ctx, q, vendor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Financing
	for rows.Next() {
		var f Financing
		if err := rows.Scan(&f.ID, &f.VendorID, &f.Amount, &f.Status, &f.RepaymentSchedule, &f.RemainingBalance, &f.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (s *Store) LatestOpenFinancing(ctx context.Context, vendor uuid.UUID) (*Financing, error) {
	const q = `
SELECT id, vendor_id, amount, status, repayment_schedule, remaining_balance, created_at
FROM financing
WHERE vendor_id=$1 AND status='approved' AND remaining_balance > 0
ORDER BY created_at DESC
LIMIT 1`
	var f Financing
	err := s.DB.QueryRow(ctx, q, vendor).Scan(&f.ID, &f.VendorID, &f.Amount, &f.Status, &f.RepaymentSchedule, &f.RemainingBalance, &f.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *Store) UpdateFinancingBalance(ctx context.Context, id uuid.UUID, remaining float64) error {
	const q = `UPDATE financing SET remaining_balance=$2 WHERE id=$1`
	ct, err := s.DB.Exec(ctx, q, id, remaining)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("no financing updated")
	}
	return nil
}

func (s *Store) InsertTransaction(ctx context.Context, t *Transaction) error {
	const q = `
INSERT INTO transactions (vendor_id, kitchen_id, amount, meal_count, ts)
VALUES ($1,$2,$3,$4,$5) RETURNING id`
	if t.TS.IsZero() {
		t.TS = time.Now().UTC()
	}
	return s.DB.QueryRow(ctx, q, t.VendorID, t.KitchenID, t.Amount, t.MealCount, t.TS).Scan(&t.ID)
}

func (s *Store) ListTransactions(ctx context.Context, vendor uuid.UUID, limit int) ([]Transaction, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	const q = `
SELECT id, vendor_id, kitchen_id, amount, meal_count, ts
FROM transactions WHERE vendor_id=$1 ORDER BY ts DESC LIMIT $2`
	rows, err := s.DB.Query(ctx, q, vendor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.VendorID, &t.KitchenID, &t.Amount, &t.MealCount, &t.TS); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// AggregateTransactionsRange sums meals and revenue for vendor between start inclusive and end exclusive.
func (s *Store) AggregateTransactionsRange(ctx context.Context, vendor uuid.UUID, start, end time.Time) (meals int, revenue float64, err error) {
	const q = `
SELECT COALESCE(SUM(meal_count),0)::int, COALESCE(SUM(amount),0)::float8
FROM transactions WHERE vendor_id=$1 AND ts >= $2 AND ts < $3`
	err = s.DB.QueryRow(ctx, q, vendor, start, end).Scan(&meals, &revenue)
	return
}

// ListVendorIDs returns all vendor primary keys (for batch jobs).
func (s *Store) ListVendorIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := s.DB.Query(ctx, `SELECT id FROM vendors ORDER BY onboarding_date`)
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
