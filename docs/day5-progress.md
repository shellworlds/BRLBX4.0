# Day 5 — Payments, compliance, security hardening

## Delivered in repo

### Payments (`backend-services/services/payments`)

- Stripe subscriptions, webhooks, carbon purchase, internal meal transaction recording.
- Docker image + K8s base/overlays (dev/staging/prod) + ArgoCD App of Apps entries.
- CI: build/push in `backend-ci-cd.yml`, Kustomize validation in `k8s-manifests.yml`.

### Compliance (`backend-services/services/compliance`)

- GDPR/CCPA consent, audit trail, deletion requests, contact form with optional SMTP.
- Same Docker/K8s/CI wiring as payments.

### Vendor ecosystem

- Wallet ledger, payouts, region fields (see prior commits + migrations).

### Energy management

- `kitchens.region`, `emission_factors` table, `GET /api/v1/reports/ghg` (JSON + CSV), daily aggregate uses regional grid intensity.

### IoT ingestion

- `kitchen_devices` table, internal device registration with optional CA signing (`IOT_DEVICE_CA_*`, `INTERNAL_DEVICE_TOKEN`).
- Documentation: `docs/emqx-mtls-devices.md`.

### Frontend

- BFF `payments` and `compliance` upstream routes; public contact via `COMPLIANCE_SERVICE_URL` in `/api/contact`.
- Dev proxy ports for payments/compliance; helpers for GHG report, wallet, subscriptions, consent.

### Ingress & secrets

- `infrastructure/k8s/ingress/borel-api-ingress.yaml`: TLS (cert-manager), `/api/payments`, `/api/compliance`, legacy paths retained.
- `infrastructure/k8s/secrets/README.md`: Sealed Secrets and rotation.

### Docs

- `docs/payments-runbook.md`, `docs/multi-region-strategy.md`, `docs/auth0-post-login-sync.md`, `docs/emqx-mtls-devices.md`.

## CI note

- Blended Go coverage gate is **35%** (see `backend-services/scripts/check_coverage.sh`); payments/compliance Gin layers are excluded from the blend and should gain dedicated integration tests later.

## Pending / follow-up

- **Stripe Connect**: full automated transfers to connected accounts in production.
- **Auth0 Action**: deploy the Action in the tenant using `docs/auth0-post-login-sync.md`.
- **E2E**: expand Playwright flows against a test tenant (smoke tests remain minimal).
- **Dependabot**: weekly PRs; resolve alerts as they appear (`npm audit` / `go mod` updates). Frontend `npm audit` may still report Next/eslint advisories that only clear with a **major** Next upgrade—schedule for Day 6+ rather than `npm audit fix --force` on this branch.
- **GitHub Wiki**: paste a summary from this file into the wiki Home (manual step on GitHub).

## Known issues

- Ingress host `api.borelsigma.com` requires DNS + cert-manager issuer `letsencrypt-prod` in cluster.
- Contact email requires SMTP env on compliance + `COMPLIANCE_SERVICE_URL` on frontend for full persistence.
