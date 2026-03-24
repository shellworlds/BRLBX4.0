# Borel Sigma frontend

Next.js 14 (App Router) marketing site plus authenticated client, vendor, and admin portals. APIs are reached through a server-side BFF (`/api/upstream/...`) so the browser stays same-origin and Auth0 access tokens attach automatically.

## Prerequisites

- Node 20+
- Auth0 application (Regular Web App) with callback/logout URLs for this app
- Backend services locally or reachable from the BFF

## Environment

Copy `.env.example` to `.env.local`. Key variables: `AUTH0_SECRET`, `AUTH0_BASE_URL`, `AUTH0_ISSUER_BASE_URL`, `AUTH0_CLIENT_ID`, `AUTH0_CLIENT_SECRET`, optional `AUTH0_AUDIENCE`, service URLs `ENERGY_SERVICE_URL`, `VENDOR_SERVICE_URL`, `IOT_SERVICE_URL`, `AUTH_SERVICE_URL`, `ML_SERVICE_URL`, `PAYMENTS_SERVICE_URL`, `COMPLIANCE_SERVICE_URL` (also used by `/api/contact`), and optional `DEV_*_URL` for local rewrites.

## Local development

```bash
npm ci
npm run dev
```

Git hooks (from repo root): `git config core.hooksPath .husky`

## Scripts

- `npm run lint` — ESLint
- `npm test` — Jest
- `npm run test:e2e` — Playwright (`npx playwright install` first)
- `npm run build` / `start` — production (needs Auth0 env; CI uses placeholders)

## Docker / Kubernetes

See `Dockerfile` and `infrastructure/k8s/apps/frontend/`. Use Sealed Secrets or a secret manager for `AUTH0_SECRET` and `AUTH0_CLIENT_SECRET`.

## Separate repo

To publish as `shellworlds/borelsigma-frontend`, copy `frontend/` or use `git subtree split` on this directory.
