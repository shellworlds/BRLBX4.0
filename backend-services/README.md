# backend-services

Go microservices for **Borel Sigma** with shared packages under `pkg/`.

## Services

| Service | Responsibility | Primary data stores |
|--------|------------------|------------------------|
| `services/energy-management` | Kitchens, energy readings, tri-modal controller status | PostgreSQL (kitchens) + TimescaleDB (readings hypertable) |
| `services/vendor-ecosystem` | Vendors, financing, transactions | PostgreSQL |
| `services/iot-ingestion` | MQTT telemetry ingest, raw archive, offline alerts | TimescaleDB (raw) + PostgreSQL (alerts) |
| `services/auth-rbac` | Auth0 JWT validation, user sync webhook, RBAC surface | PostgreSQL |
| `services/ml-predictor` | Stub ML inference (energy load curve, meal demand), optional Redis | optional Redis |
| `services/payments` | Stripe subscriptions, webhooks, carbon purchases, internal meal records | PostgreSQL |
| `services/compliance` | Consent, audit log, GDPR deletion queue, contact CRM + SMTP | PostgreSQL |
| `services/hello-world` | Tiny health/metrics demo | none |

## Shared packages (`pkg/`)

- `pkg/config` — environment config via Viper
- `pkg/logger` — zap log setup
- `pkg/db` — pgxpool helper
- `pkg/metrics` — `/metrics` Gin handler (Prometheus)
- `pkg/migrate` — `golang-migrate` runner against `migrations/` trees
- `pkg/auth` — Auth0 JWKS validation + Gin middleware + admin header helper
- `pkg/pgxutil` — small `Querier` interface shared with `pgxmock` tests

## Running locally

From `backend-services/`:

```bash
export POSTGRES_KITCHEN_DSN='postgres://user:pass@localhost:5432/kitchens?sslmode=disable'
export POSTGRES_TIMESCALE_DSN='postgres://user:pass@localhost:5433/ts?sslmode=disable'
go run ./services/energy-management
```

Each service reads `SERVICE_ROOT` (defaults to `.` locally; Docker sets `/app` with migrations baked in).

### Common environment variables

**energy-management**

- `POSTGRES_KITCHEN_DSN`, `POSTGRES_TIMESCALE_DSN` (required)
- `ADMIN_API_KEY` **or** `AUTH0_DOMAIN` + `AUTH0_AUDIENCE` for admin-only kitchen registration
- `INGEST_BEARER_TOKEN` (optional) protects `POST /api/v1/kitchens/:id/readings`
- `INTERNAL_AGGREGATE_TOKEN` enables `POST /api/v1/internal/aggregate/daily` (used by CronJob)
- `ENABLE_SWAGGER` (default true) serves `/swagger/*`

**vendor-ecosystem**

- `POSTGRES_VENDOR_DSN`
- `INTERNAL_AGGREGATE_TOKEN` for vendor daily rollup + payout processor CronJob
- `STRIPE_SECRET_KEY` — Connect transfers and Express onboarding
- `STRIPE_CONNECT_WEBHOOK_SECRET` — `POST /api/v1/webhooks/stripe/connect`
- `STRIPE_CONNECT_REFRESH_URL`, `STRIPE_CONNECT_RETURN_URL` — Account Link URLs for vendors
- `STRIPE_CONNECT_DEFAULT_COUNTRY` (e.g. `US`)
- `ENABLE_SWAGGER`

**iot-ingestion**

- `POSTGRES_IOT_DSN`, `POSTGRES_IOT_TIMESCALE_DSN`
- `ENERGY_SERVICE_URL` (e.g. `http://localhost:8080` or in-cluster `http://energy-management`)
- `MQTT_BROKER_URL` (`tcp://...`)
- `MQTT_USERNAME`, `MQTT_PASSWORD` (optional)
- `INGEST_BEARER_TOKEN` (must match energy service)
- `SLACK_WEBHOOK_URL` (optional)
- `ANOMALY_WINDOW` (default 20), `ANOMALY_SIGMA` (default 3), `ANOMALY_DISABLED` to turn off z-score alerts
- `INTERNAL_DEVICE_TOKEN` — enables `POST /api/v1/internal/devices/register`
- `IOT_DEVICE_CA_CERT_FILE`, `IOT_DEVICE_CA_KEY_FILE` — optional PEM paths for signing device CSRs (both required if either set)
- `ENABLE_SWAGGER`

**auth-rbac**

- `POSTGRES_AUTH_DSN`
- `AUTH0_DOMAIN`, `AUTH0_AUDIENCE`
- `AUTH0_SYNC_WEBHOOK_SECRET` — required for `/api/v1/users/sync`; send as `X-Webhook-Secret` or `Authorization: Bearer`
- `ENABLE_SWAGGER`

**ml-predictor**

- Optional `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`, `REDIS_PREFIX` for prediction memoization
- `ENABLE_SWAGGER`

**payments**

- `POSTGRES_PAYMENTS_DSN`
- `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`, `STRIPE_DEFAULT_PRICE_ID` (subscription price)
- `INTERNAL_PAYMENTS_TOKEN` — protects internal meal recording from vendor-ecosystem
- `AUTH0_DOMAIN`, `AUTH0_AUDIENCE` (JWT routes)

**compliance**

- `POSTGRES_COMPLIANCE_DSN`
- `AUTH0_DOMAIN`, `AUTH0_AUDIENCE`
- Optional SMTP: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`, `SALES_NOTIFY_EMAIL` (defaults to sales@borelsigma.com)

## Migrations

Each service keeps SQL migrations under `services/<name>/migrations/...` and runs `migrate` automatically on startup.

## Docker images

Build context must be **`backend-services/`** (the module root):

```bash
docker build -f services/energy-management/Dockerfile -t energy-management:dev .
```

## OpenAPI / Swagger

Handlers are annotated with `swag` conventions in selected `main.go` files. To generate docs locally:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
# Example (per service): swag init -g main.go -d services/energy-management -o services/energy-management/docs
```

## Tests

```bash
go test ./...
go test -cover ./...
bash scripts/check_coverage.sh   # default min 40% blended; override with COVERAGE_THRESHOLD=70
```

Long-running **Docker integration** flows (full MQTT + multi-DB) are documented for a future `docker compose` profile; `testcontainers-go` was not pinned here due to Go 1.22 / dependency drift—see `docs/day3-progress.md`.

## Kubernetes secrets (expected keys)

Create SealedSecrets / ExternalSecrets with the same names referenced by manifests:

- `energy-management-secrets`: `POSTGRES_KITCHEN_DSN`, `POSTGRES_TIMESCALE_DSN`, optional `ADMIN_API_KEY`, `INGEST_BEARER_TOKEN`, `AUTH0_*`, `INTERNAL_AGGREGATE_TOKEN`
- `vendor-ecosystem-secrets`: `POSTGRES_VENDOR_DSN`, `INTERNAL_AGGREGATE_TOKEN`, Stripe Connect keys (`STRIPE_SECRET_KEY`, `STRIPE_CONNECT_WEBHOOK_SECRET`, refresh/return URLs)
- `iot-ingestion-secrets`: `POSTGRES_IOT_DSN`, `POSTGRES_IOT_TIMESCALE_DSN`, `ENERGY_SERVICE_URL`, `MQTT_BROKER_URL`, optional MQTT creds, `INGEST_BEARER_TOKEN`, `SLACK_WEBHOOK_URL`
- `auth-rbac-secrets`: `POSTGRES_AUTH_DSN`, `AUTH0_DOMAIN`, `AUTH0_AUDIENCE`, `AUTH0_SYNC_WEBHOOK_SECRET`
- `ml-predictor-secrets`: optional Redis secrets if not using env literals

Cloud SQL Auth Proxy sidecars should target these DSNs with `127.0.0.1` hosts (advanced overlays).
