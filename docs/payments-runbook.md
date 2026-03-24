# Payments runbook (Stripe)

## Webhook failures

**Symptoms:** Subscriptions or invoices out of sync; Stripe Dashboard shows delivery errors.

**Checks:**

1. `GET /healthz` on `payments` service; verify `payments-secrets` contains `STRIPE_WEBHOOK_SECRET` matching the endpoint in Stripe Dashboard.
2. Ingress must expose `POST /api/v1/webhooks/stripe` (path `/api/payments/api/v1/webhooks/stripe` when using `/api/payments` prefix).
3. Stripe sends **Stripe-Signature**; the handler uses `ConstructEvent` — body must be raw (no JSON re-encoding in proxies).

**Fix:** Rotate webhook secret in Stripe, update sealed secret, restart deployment.

## Payout / Connect issues

Vendor withdrawals depend on **Stripe Connect** and vendor `stripe_connect_account_id`. If `wallet/withdraw` returns errors:

1. Confirm vendor completed Connect onboarding in Stripe Dashboard.
2. Check `payments` logs for transfer errors.
3. Verify platform balance and Connect settings in Stripe.

## Idempotency

Webhook events are deduplicated by Stripe `event.id` in Postgres (`TryMarkWebhookEvent`). If you replay events manually, expect `duplicate` responses.
