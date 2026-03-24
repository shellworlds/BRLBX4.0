package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PayoutRow is a payout request with timestamps for APIs and webhooks.
type PayoutRow struct {
	ID               uuid.UUID  `json:"id"`
	VendorID         uuid.UUID  `json:"vendor_id"`
	Amount           float64    `json:"amount"`
	Status           string     `json:"status"`
	StripeTransferID *string    `json:"stripe_transfer_id,omitempty"`
	FailureReason    *string    `json:"failure_reason,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

func (s *Store) ListPendingPayouts(ctx context.Context, limit int) ([]PayoutRow, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	const q = `
SELECT id, vendor_id, amount::float8, status, stripe_transfer_id, failure_reason, created_at
FROM payout_requests WHERE status = 'pending' ORDER BY created_at ASC LIMIT $1`
	rows, err := s.DB.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPayoutRows(rows)
}

func (s *Store) ListPayoutsForVendor(ctx context.Context, vendorID uuid.UUID, limit int) ([]PayoutRow, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	const q = `
SELECT id, vendor_id, amount::float8, status, stripe_transfer_id, failure_reason, created_at
FROM payout_requests WHERE vendor_id=$1 ORDER BY created_at DESC LIMIT $2`
	rows, err := s.DB.Query(ctx, q, vendorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPayoutRows(rows)
}

func scanPayoutRows(rows pgx.Rows) ([]PayoutRow, error) {
	var out []PayoutRow
	for rows.Next() {
		var r PayoutRow
		if err := rows.Scan(&r.ID, &r.VendorID, &r.Amount, &r.Status, &r.StripeTransferID, &r.FailureReason, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// CompletePayoutAsPaid marks a pending payout paid and reduces pending_payout on the wallet.
func (s *Store) CompletePayoutAsPaid(ctx context.Context, payoutID, vendorID uuid.UUID, amount float64, transferID string) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const up = `UPDATE payout_requests SET status='paid', stripe_transfer_id=$2, failure_reason=NULL WHERE id=$1 AND status='pending'`
	ct, err := tx.Exec(ctx, up, payoutID, transferID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return tx.Commit(ctx) // idempotent: already finalized
	}
	const w = `UPDATE vendor_wallets SET pending_payout = pending_payout - $2, updated_at=now() WHERE vendor_id=$1`
	if _, err := tx.Exec(ctx, w, vendorID, amount); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// FailPayoutAndRelease returns reserved funds from pending_payout back to balance.
func (s *Store) FailPayoutAndRelease(ctx context.Context, payoutID, vendorID uuid.UUID, amount float64, reason string) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const up = `UPDATE payout_requests SET status='failed', failure_reason=$2 WHERE id=$1 AND status='pending'`
	ct, err := tx.Exec(ctx, up, payoutID, reason)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("payout not pending")
	}
	const w = `UPDATE vendor_wallets SET pending_payout = pending_payout - $2, balance = balance + $2, updated_at=now() WHERE vendor_id=$1`
	if _, err := tx.Exec(ctx, w, vendorID, amount); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// SetStripeConnectAccountID persists the Stripe Connect account id for a vendor.
func (s *Store) SetStripeConnectAccountID(ctx context.Context, vendorID uuid.UUID, accountID string) error {
	const q = `UPDATE vendors SET stripe_connect_account_id=$2 WHERE id=$1`
	_, err := s.DB.Exec(ctx, q, vendorID, accountID)
	return err
}

// GetPayoutByID loads one payout row (any status).
func (s *Store) GetPayoutByID(ctx context.Context, id uuid.UUID) (*PayoutRow, error) {
	const q = `
SELECT id, vendor_id, amount::float8, status, stripe_transfer_id, failure_reason, created_at
FROM payout_requests WHERE id=$1`
	var r PayoutRow
	err := s.DB.QueryRow(ctx, q, id).Scan(&r.ID, &r.VendorID, &r.Amount, &r.Status, &r.StripeTransferID, &r.FailureReason, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
