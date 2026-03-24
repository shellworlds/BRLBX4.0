# Day 4 progress — public site and portals

## Summary

Delivered the Day 4 scope inside the monorepo `frontend/` package (with notes for splitting to `shellworlds/borelsigma-frontend` if desired): marketing pages, Auth0-protected `/portal/*` routes with role-based redirects, client/vendor/admin UIs backed by typed API helpers and a BFF proxy, Docker standalone build, Kubernetes manifests, Argo CD applications, ingress example, backend CORS support, and a public energy snapshot endpoint for landing metrics.

## Frontend

- **Stack:** Next.js 14, TypeScript, Tailwind, ESLint, Prettier, lint-staged, Husky hook at repo `.husky/pre-commit`, Jest + RTL, Playwright smoke tests.
- **Marketing:** `/`, `/about`, `/contact` (API + mailto), `/blog` placeholder; 20x6 diligence table with external sources; hero, metrics, patent narrative; live-ish ticker via `/api/public/energy-snapshot`.
- **Auth:** `@auth0/nextjs-auth0`, middleware on `/portal`, `/portal` resolves role via `auth-rbac` `/api/v1/users/me` and redirects to client, vendor, or admin dashboards.
- **Portals:** Client (dashboard with Recharts, reports, kitchens + time series, settings stub), vendor (dashboard, transactions, financing form, compliance + IoT alerts), admin (overview, vendors/kitchens copy, alerts with ack, health checklist).
- **API layer:** `lib/api.ts` + `/api/upstream/[service]/[[...path]]` forwards bearer tokens; `lib/service-url.ts` for server-side URLs; dev rewrites `/dev-proxy/{energy|vendor|iot|auth|ml}`.

## Backend (gap handling)

- **`pkg/cors`:** `CORS_ALLOWED_ORIGINS` comma list or allow-all when unset; wired into energy, vendor, IoT, auth-rbac, ml-predictor routers.
- **Energy:** `GET /api/v1/public/snapshot` (marketing metrics), `GET /api/v1/kitchens/vendor/:vendor_id` (optional JWT + roles vendor/admin/client when `AUTH0_DOMAIN` is set).

## GitOps and ingress

- **Kustomize:** `infrastructure/k8s/apps/frontend/` base + dev/staging/prod overlays.
- **Argo:** `frontend-dev`, `frontend-staging`, `frontend-prod` in app-of-apps.
- **Ingress example:** `infrastructure/k8s/ingress/borel-api-ingress.yaml` (nginx path routing).

## CI

- **Frontend workflow:** Auth0 placeholder env for `next build`; Docker build-args; optional PR to bump `frontend` dev image tag when GCP WIF is configured.
- **Kustomize workflow:** Builds new frontend overlays.

## Next steps

1. Create Auth0 Actions or rules to emit roles and call `auth-rbac` `/users/sync` on login.
2. Seal secrets for `AUTH0_SECRET` / client secret; set real `AUTH0_BASE_URL` per environment.
3. Point `NEXT_PUBLIC_*` and ingress hosts at production domains; rebuild images if client bundle must embed public URLs.
4. Wire `POST /api/contact` to SendGrid/SES.
5. Expand Playwright coverage (login flow) with a test tenant or mocks.
6. Optional: split `frontend/` to `borelsigma-frontend` repo for deployment isolation.
