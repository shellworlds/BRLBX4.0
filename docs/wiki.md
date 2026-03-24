# GitHub Wiki setup

The repository has **Wiki** enabled. GitHub creates the wiki Git endpoint the first time content exists.

1. Open the repository on GitHub → **Wiki**.
2. Click **Create the first page** (or equivalent), save **Home**.
3. Optionally replace the default home with content that points to:
   - `README.md` (overview)
   - `infrastructure/README.md` (ops)
   - `docs/architecture.md` (diagram)
   - `docs/day5-progress.md` (Day 5 payments/compliance/security status)
   - `docs/day6-progress.md` (Day 6 E2E, load tests, Connect, ingress)
   - `docs/runbooks/` (Auth0, Stripe Connect, DR)
   - `docs/day7-launch-checklist.md`

Advanced contributors can clone the wiki repository:

`https://github.com/shellworlds/BRLBX4.0.wiki.git`

**E2E (Playwright) secrets** (optional): `E2E_BASE_URL`, `E2E_AUTH0_EMAIL` / `PASSWORD` (vendor default), plus `E2E_AUTH0_CLIENT_*`, `E2E_AUTH0_ADMIN_*`, `E2E_AUTH0_VENDOR_*` for role-specific suites. CI serves the **standalone** Node server after `next build`.
