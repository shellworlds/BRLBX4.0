package stripepayout

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/transfer"

	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
)

// ProcessPendingStripeTransfers creates Stripe Transfers for pending payout rows (test/live mode per API key).
func ProcessPendingStripeTransfers(ctx context.Context, st *repo.Store, limit int) (success int, _ error) {
	pending, err := st.ListPendingPayouts(ctx, limit)
	if err != nil {
		return 0, err
	}
	var lastErr error
	for _, p := range pending {
		v, err := st.GetVendor(ctx, p.VendorID)
		if err != nil {
			lastErr = err
			continue
		}
		if v.StripeConnectAccountID == nil || strings.TrimSpace(*v.StripeConnectAccountID) == "" {
			_ = st.FailPayoutAndRelease(ctx, p.ID, p.VendorID, p.Amount, "stripe_connect_not_onboarded")
			continue
		}
		cents := int64(math.Round(p.Amount * 100))
		if cents < 1 {
			cents = 1
		}
		dest := strings.TrimSpace(*v.StripeConnectAccountID)
		tr, err := transfer.New(&stripe.TransferParams{
			Amount:      stripe.Int64(cents),
			Currency:    stripe.String("usd"),
			Destination: stripe.String(dest),
			Metadata: map[string]string{
				"payout_request_id": p.ID.String(),
			},
		})
		if err != nil {
			_ = st.FailPayoutAndRelease(ctx, p.ID, p.VendorID, p.Amount, truncateReason(err))
			lastErr = err
			continue
		}
		if err := st.CompletePayoutAsPaid(ctx, p.ID, p.VendorID, p.Amount, tr.ID); err != nil {
			lastErr = err
			continue
		}
		success++
	}
	return success, lastErr
}

func truncateReason(err error) string {
	s := err.Error()
	if len(s) > 480 {
		return s[:480]
	}
	return s
}

// CompleteFromTransferMetadata marks a payout paid when Stripe confirms a transfer (idempotent).
func CompleteFromTransferMetadata(ctx context.Context, st *repo.Store, payoutIDStr, transferID string) error {
	if payoutIDStr == "" || transferID == "" {
		return fmt.Errorf("metadata")
	}
	pid, err := parseUUID(payoutIDStr)
	if err != nil {
		return err
	}
	row, err := st.GetPayoutByID(ctx, pid)
	if err != nil {
		return err
	}
	if row.Status != "pending" {
		return nil
	}
	return st.CompletePayoutAsPaid(ctx, row.ID, row.VendorID, row.Amount, transferID)
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
