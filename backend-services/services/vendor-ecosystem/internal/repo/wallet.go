package repo

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Wallet struct {
	VendorID      uuid.UUID `json:"vendor_id"`
	Balance       float64   `json:"balance"`
	PendingPayout float64   `json:"pending_payout"`
}

type PayoutRequest struct {
	ID                uuid.UUID `json:"id"`
	VendorID          uuid.UUID `json:"vendor_id"`
	Amount            float64   `json:"amount"`
	Status            string    `json:"status"`
	StripeTransferID  *string   `json:"stripe_transfer_id,omitempty"`
}

// ApplyMealNetCredits vendor wallet with gross meal revenue minus platform fee (fee in basis points).
func (s *Store) ApplyMealNetCredits(ctx context.Context, vendorID uuid.UUID, gross float64, feeBps int, mealTxID uuid.UUID) error {
	if feeBps < 0 || feeBps > 10000 {
		feeBps = 0
	}
	fee := math.Round(gross*float64(feeBps)) / 10000.0
	net := gross - fee
	if net < 0 {
		net = 0
	}
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const ensure = `
INSERT INTO vendor_wallets (vendor_id) VALUES ($1)
ON CONFLICT (vendor_id) DO NOTHING`
	if _, err := tx.Exec(ctx, ensure, vendorID); err != nil {
		return err
	}
	const up = `UPDATE vendor_wallets SET balance = balance + $2, updated_at = now() WHERE vendor_id = $1`
	if _, err := tx.Exec(ctx, up, vendorID, net); err != nil {
		return err
	}
	const led = `
INSERT INTO wallet_ledger (vendor_id, delta, reason, ref_type, ref_id)
VALUES ($1, $2, 'meal_net', 'transaction', $3)`
	if _, err := tx.Exec(ctx, led, vendorID, net, mealTxID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) GetWallet(ctx context.Context, vendorID uuid.UUID) (*Wallet, error) {
	const q = `
SELECT vendor_id, balance::float8, pending_payout::float8
FROM vendor_wallets WHERE vendor_id=$1`
	var w Wallet
	err := s.DB.QueryRow(ctx, q, vendorID).Scan(&w.VendorID, &w.Balance, &w.PendingPayout)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &Wallet{VendorID: vendorID, Balance: 0, PendingPayout: 0}, nil
		}
		return nil, err
	}
	return &w, nil
}

func (s *Store) RequestPayout(ctx context.Context, vendorID uuid.UUID, amount float64) (*PayoutRequest, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount")
	}
	w, err := s.GetWallet(ctx, vendorID)
	if err != nil {
		return nil, err
	}
	if w.Balance < amount {
		return nil, fmt.Errorf("insufficient balance")
	}
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const dec = `UPDATE vendor_wallets SET balance = balance - $2, pending_payout = pending_payout + $2, updated_at=now() WHERE vendor_id=$1 AND balance >= $2`
	ct, err := tx.Exec(ctx, dec, vendorID, amount)
	if err != nil {
		return nil, err
	}
	if ct.RowsAffected() == 0 {
		return nil, fmt.Errorf("insufficient balance")
	}
	pr := &PayoutRequest{VendorID: vendorID, Amount: amount, Status: "pending"}
	const ins = `
INSERT INTO payout_requests (vendor_id, amount, status) VALUES ($1,$2,'pending') RETURNING id`
	if err := tx.QueryRow(ctx, ins, vendorID, amount).Scan(&pr.ID); err != nil {
		return nil, err
	}
	const led = `
INSERT INTO wallet_ledger (vendor_id, delta, reason, ref_type, ref_id)
VALUES ($1, $2, 'payout_reserved', 'payout_request', $3)`
	if _, err := tx.Exec(ctx, led, vendorID, -amount, pr.ID); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return pr, nil
}
