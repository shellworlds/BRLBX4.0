package stripepayout

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"

	"github.com/shellworlds/BRLBX4.0/backend-services/services/vendor-ecosystem/internal/repo"
)

// HandleConnectWebhook verifies the Stripe signature and reconciles payout / Connect state.
func HandleConnectWebhook(ctx context.Context, rawBody []byte, sigHeader, secret string, st *repo.Store) error {
	ev, err := webhook.ConstructEvent(rawBody, sigHeader, secret)
	if err != nil {
		return fmt.Errorf("signature: %w", err)
	}
	switch ev.Type {
	case "transfer.created":
		var tr stripe.Transfer
		if err := json.Unmarshal(ev.Data.Raw, &tr); err != nil {
			return err
		}
		pid := ""
		if tr.Metadata != nil {
			pid = tr.Metadata["payout_request_id"]
		}
		return CompleteFromTransferMetadata(ctx, st, pid, tr.ID)
	case "transfer.reversed":
		var tr stripe.Transfer
		if err := json.Unmarshal(ev.Data.Raw, &tr); err != nil {
			return err
		}
		pid := ""
		if tr.Metadata != nil {
			pid = tr.Metadata["payout_request_id"]
		}
		if pid == "" {
			return nil
		}
		puid, err := uuid.Parse(pid)
		if err != nil {
			return nil
		}
		row, err := st.GetPayoutByID(ctx, puid)
		if err != nil {
			return nil
		}
		if row.Status != "pending" {
			return nil
		}
		return st.FailPayoutAndRelease(ctx, row.ID, row.VendorID, row.Amount, "stripe_transfer_reversed")
	default:
		return nil
	}
}
