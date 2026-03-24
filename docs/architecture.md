# Borel Sigma — baseline architecture (Day 1)

Text diagram of the Day‑1 footprint on Google Cloud and GitHub.

```
                    ┌─────────────────────────────────────────────────────────┐
                    │                     GitHub                               │
                    │  shellworlds/BRLBX4.0 (monorepo)                        │
                    │   ├─ backend-services/   (Go microservices)             │
                    │   ├─ frontend/           (Next.js)                       │
                    │   └─ infrastructure/     (Terraform + K8s + GitOps)     │
                    │        GitHub Actions → build/test → Artifact Registry   │
                    └──────────────────────────┬──────────────────────────────┘
                                               │ WIF (preferred) or SA key
                                               ▼
┌──────────────────────────────────────────────────────────────────────────────────┐
│                         GCP project: borel-sigma-prod                             │
│  ┌────────────────────────────────────────────────────────────────────────────┐  │
│  │ VPC (private GKE nodes)  +  Cloud NAT  +  internal firewall model           │  │
│  │        │                                                                     │  │
│  │        ▼                                                                     │  │
│  │  GKE regional cluster: borel-sigma-cluster                                   │  │
│  │   ├─ namespaces: dev / staging / prod / argocd / monitoring                  │  │
│  │   ├─ Argo CD (GitOps) → App of Apps → per-service Applications               │  │
│  │   ├─ hello-world (dev) + Prometheus metrics (/metrics)                       │  │
│  │   ├─ payments + compliance (Stripe, GDPR/audit, contact)                     │  │
│  │   ├─ Sealed Secrets controller                                               │  │
│  │   └─ Monitoring: Prometheus/Grafana (Helm) + Loki/Promtail (Helm)           │  │
│  └────────────────────────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────┐   ┌────────────────────────────────────────┐   │
│  │ Cloud SQL (HA Postgres x2)    │   │ Cloud SQL (Postgres + Timescale flag)   │   │
│  │ transactional workloads       │   │ time-series workloads                   │   │
│  └──────────────────────────────┘   └────────────────────────────────────────┘   │
│        ▲ private IP (PSC)                 ▲                                         │
│        └──────── Cloud SQL Proxy (pods / CronJobs) ────────────────┘                │
│                                                                                     │
│  Artifact Registry (images)           GCS: Terraform state + DB backups           │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

## Data flow (CI/CD)

1. Engineers push to `main` or open PRs; Actions run tests and Terraform/Kustomize checks.
2. On `main`, authenticated workflows build container images and push to Artifact Registry.
3. `backend-ci-cd` opens a PR that bumps the image digest/tag in Kustomize under `infrastructure/k8s/apps/`.
4. Argo CD reconciles the cluster to match Git (auto-sync + prune where enabled).

## Operational notes

- Lock down GKE control-plane authorized networks after bootstrap (replace any temporary `0.0.0.0/0` entries).
- Apply PrometheusRule and ServiceMonitor overlays only after Prometheus Operator CRDs exist (`dev-with-prometheus` overlay, or manual Helm install per `infrastructure/k8s/monitoring/README.md`).
- Database credentials belong in Sealed Secrets or Secret Manager; never commit raw secrets.
