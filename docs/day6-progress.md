# Day 6 — E2E, performance, hardening

## Implemented in repo

- **Auth0 sync**: `auth-rbac` accepts `Authorization: Bearer` or `X-Webhook-Secret`; body accepts `sub` or `auth0_id`. Action script and staging/prod secrets documented in `docs/auth0-post-login-sync.md`.
- **Stripe Connect**: vendor-ecosystem creates **Express** accounts, Account Links, **Transfer** payout processing (`POST /internal/payouts/process`), Connect **webhook** handler, migration `failure_reason` on `payout_requests`, CronJob every 10m, vendor dashboard wallet + payout list + onboarding button.
- **Ingress**: `borel-web-ingress.yaml` for `www` + `staging` frontends (port 3000); `ingress/kustomization.yaml` for CI validation.
- **k6**: scripts under `perf/k6/` (energy ingest, public snapshot, vendor transactions).
- **Playwright**: `e2e/full-flows.spec.ts` (auth optional via env); workflow `.github/workflows/e2e-playwright.yml` (push + nightly).
- **Monitoring**: example `BorelAPI5xxRateHigh` rule when `http_requests_total` exists; existing CPU/restart/target rules retained.
- **Runbooks**: `docs/runbooks/auth0-action.md`, `docs/runbooks/stripe-connect-vendors.md`, `docs/runbooks/disaster-recovery.md`, `docs/day7-launch-checklist.md`.
- **Tests / coverage**: vendor-ecosystem router + payout repo mocks; blended coverage gate ≥35% (`stripepayout` excluded from blend as Stripe SDK glue).
- **npm audit**: 8 issues (4 high) as of Day 6; fixes require `next@16` / breaking `eslint-config-next` — track in Security tab; patch when Next major is scheduled.

## Manual / cluster steps (not automated here)

- Install **cert-manager** + `letsencrypt-prod` issuer in each cluster.
- Point **Stripe Connect** webhook to the public vendor URL (through `api` ingress `/api/vendor/...`).
- Run **k6** for 10+ minutes against staging and archive results (Grafana + k6 cloud optional).
- **Dependabot**: review GitHub Security tab; Next.js major upgrade tracked separately from patch-only fixes.

## Open gaps

- Full client/admin Playwright flows need stable test users + seeded kitchens (extend `full-flows.spec.ts`).
- Redis-backed caching for hot report endpoints (optional tuning after k6).
- Multi-region active-active remains future work (see `docs/multi-region-strategy.md`).
