# Day 3 — Data intelligence, autonomy, gap closure

## Shipped in this milestone

- **ML predictor** (`services/ml-predictor`): stub APIs for energy curves and vendor demand using `gonum`, optional Redis cache, Prometheus latency histogram, Swagger UI, Dockerfile and Kubernetes overlays (dev/staging/prod).
- **Reporting / ESG-style aggregates**
  - Energy: `daily_client_reports` table, `GET /api/v1/reports/client/:client_id`, internal `POST /api/v1/internal/aggregate/daily` (protected by `INTERNAL_AGGREGATE_TOKEN` + CronJob).
  - Vendor: `daily_vendor_metrics`, `GET /api/v1/reports/vendor/:vendor_id`, internal aggregate route + CronJob.
- **IoT**: sliding-window anomaly detection (3-sigma on grid power), Prometheus counters `mqtt_ingest_messages_total` and `iot_anomaly_alerts_total`, Slack on anomaly (existing webhook).
- **Auth / repos**: `pgxutil.Querier` abstraction for pgxmock-backed tests across services.
- **Swagger**: generated `docs/` committed for energy, vendor, iot, auth, ml-predictor; `/swagger/index.html` when `ENABLE_SWAGGER` is true (default).
- **Tests & CI**: broadened unit tests (pkg auth/config/logger/metrics, repos, routers, anomaly, ML); `scripts/check_coverage.sh` enforces a **40% blended floor** on `pkg/*` (minus bootstrap) + all `services/*/internal/*` until handler coverage closes the gap toward the **70% target**. Integration tests via `testcontainers-go` were **not** pinned (toolchain/Docker API drift on Go 1.22); prefer Docker Compose or upgrade to Go 1.23+ with maintained testcontainers later.
- **Kubernetes / GitOps**
  - Cloud SQL Auth Proxy **sidecar** on energy, vendor, iot, auth deployments; ConfigMap keys for `cloudsql_instances` (placeholder `PROJECT:REGION:…`).
  - **CronJobs** for daily aggregation (energy + vendor).
  - **KEDA `ScaledObject`** manifests (CPU; iot also has a Prometheus trigger aimed at kube-prometheus-stack).
  - Argo CD: **KEDA** and **EMQX** dev apps; **ml-predictor** dev app; **energy-management** staging + prod apps; IoT ConfigMap defaults `MQTT_BROKER_URL` to in-cluster EMQX service.
  - Helm wrapper chart under `infrastructure/helm/emqx` (Argo primarily uses the upstream chart repo).

## Gaps / next steps

1. Raise **coverage toward 70%**: expand Gin handler table-tests; optional Go 1.23+ and latest `testcontainers-go` for real integration flows (vendor + kitchen + MQTT + reports + ML).
2. **Staging/prod Argo** for vendor, iot, auth, ml (only energy + patterns for others are in app-of-apps today).
3. **Cloud SQL**: replace placeholder ConfigMap values; align SealedSecrets with `INTERNAL_AGGREGATE_TOKEN`, DSNs pointing at `127.0.0.1:5432/5433`, and Workload Identity for the proxy.
4. **EMQX**: harden auth/TLS via chart values; align `MQTT_BROKER_URL` per environment.
5. **KEDA / Prometheus**: tune `ScaledObject` queries when scrape labels are known; ensure KEDA installs before app sync or allow CRD install order in Argo waves.

## Key references

- `backend-services/README.md` — env vars, Swagger, coverage script.
- `infrastructure/k8s/apps/app-of-apps/` — canonical Argo Application list.
