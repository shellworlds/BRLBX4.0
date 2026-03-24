# Day 6 — E2E, performance, hardening

## Implemented in repo

- **Auth0 sync**: `auth-rbac` accepts `Authorization: Bearer` or `X-Webhook-Secret`; body accepts `sub` or `auth0_id`. Action script and staging/prod secrets documented in `docs/auth0-post-login-sync.md`.
- **Stripe Connect**: vendor-ecosystem creates **Express** accounts, Account Links, **Transfer** payout processing (`POST /internal/payouts/process`), Connect **webhook** handler, migration `failure_reason` on `payout_requests`, CronJob every 10m, vendor dashboard wallet + payout list + onboarding button.
- **Ingress**: `borel-web-ingress.yaml` for `www` + `staging` frontends (port 3000); `ingress/kustomization.yaml` for CI validation.
- **k6**: scripts under `perf/k6/` (energy ingest, public snapshot, vendor transactions).
- **Playwright**: `e2e/smoke.spec.ts` (public marketing + contact); `e2e/full-flows.spec.ts` with `e2e/helpers.ts` — vendor (`E2E_AUTH0_*` or `E2E_AUTH0_VENDOR_*`), optional client (`E2E_AUTH0_CLIENT_*`), optional admin (`E2E_AUTH0_ADMIN_*`). Workflow `.github/workflows/e2e-playwright.yml` (push + nightly).
- **Monitoring**: example `BorelAPI5xxRateHigh` rule when `http_requests_total` exists; existing CPU/restart/target rules retained.
- **Runbooks**: `docs/runbooks/auth0-action.md`, `docs/runbooks/stripe-connect-vendors.md`, `docs/runbooks/disaster-recovery.md`, `docs/day7-launch-checklist.md`.
- **Tests / coverage**: vendor-ecosystem router + payout repo mocks; blended coverage gate ≥35% (`stripepayout` excluded from blend as Stripe SDK glue).
- **Frontend stack**: `next@15.5.14`, `eslint-config-next@15.5.14`, `jest@30` + `jest-environment-jsdom@30` — **`npm audit` reports 0 vulnerabilities** (2026-03-24). App Router dynamic route handlers use async `params` (Next 15).

## Manual / cluster steps (not automated here)

- Install **cert-manager** + `letsencrypt-prod` issuer in each cluster.
- Point **Stripe Connect** webhook to the public vendor URL (through `api` ingress `/api/vendor/...`).
- Run **k6** for 10+ minutes against staging and archive results (Grafana + k6 cloud optional).
- **Dependabot**: re-scan after this push; GitHub Security tab for any residual transitive alerts.

## Open gaps (post–Day 6)

- **Playwright**: add repository secrets for `E2E_AUTH0_CLIENT_*` and `E2E_AUTH0_ADMIN_*` when test tenants exist; seed kitchens/vendors via API for stricter data assertions.
- Redis-backed caching for hot report endpoints (optional tuning after k6).
- Multi-region active-active remains future work (see `docs/multi-region-strategy.md`).
