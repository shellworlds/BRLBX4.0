# Runbook: Stripe Connect (vendor payouts)

## Configure

1. Platform Stripe account: enable **Connect**.
2. Set on **vendor-ecosystem** deployment:
   - `STRIPE_SECRET_KEY`
   - `STRIPE_CONNECT_WEBHOOK_SECRET` (signing secret for endpoint `https://api.../api/vendor/api/v1/webhooks/stripe/connect`)
   - `STRIPE_CONNECT_REFRESH_URL` / `STRIPE_CONNECT_RETURN_URL` (vendor portal URLs after onboarding)
3. In Stripe Dashboard, add Connect webhook URL and select events: `transfer.created`, `transfer.reversed`.

## Flow

1. Vendor clicks **Stripe Connect onboarding** in portal → Express account + Account Link.
2. Vendor requests **wallet withdraw** → row in `payout_requests` (`pending`).
3. CronJob `vendor-process-payouts` or `POST /api/v1/internal/payouts/process` creates **Transfers** to the connected account.
4. Webhooks reconcile idempotently if the worker partially fails.

## Troubleshooting

- **no_connect_not_onboarded** in `failure_reason`: complete Connect onboarding.
- **Insufficient balance** on Stripe: use test mode balances or fund platform test account.
- Check vendor-ecosystem logs for Stripe API errors after withdraw.
