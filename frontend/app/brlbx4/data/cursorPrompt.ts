/** Plain-text body for clipboard copy (mirrors terminal view). */
export const cursorPromptPlain = `# ═══════════════════════════════════════════════════════════════
# BRLBX4.0 — BOREL SIGMA — DAY 1 CURSOR PROMPT
# Role: Senior DevOps Engineer
# GitHub: shellworlds | GCP Project: borel-sigma-prod
# Platform: Ubuntu/WSL2 | 32GB RAM | 1TB SSD
# Sprint: Day 1 of 7 (10 hours)
# ═══════════════════════════════════════════════════════════════

── 1. REPOSITORY SETUP ──────────────────────────────────────
Action: Create 3 GitHub repos under shellworlds
  infrastructure   # Terraform, K8s manifests, ArgoCD configs
  backend-services # Go/Node.js microservices
  frontend         # Next.js portals + public site
Each: README.md + .gitignore + MIT LICENSE

── 2. TERRAFORM MODULES ─────────────────────────────────────
Provider:   GCP (google)
State:      GCS bucket tfstate-borel-sigma
Create:
  gke_cluster   # regional, asia-south1, 3 zones, REGULAR channel
    general-pool  e2-standard-4  (autoscale 2–10)
    high-mem-pool c2-standard-8  (autoscale 1–5)
  vpc           # private nodes, Cloud NAT
  cloud_sql     # PostgreSQL HA + TimescaleDB instances
  artifact_registry
  service_accounts  # least-privilege for GKE, CloudSQL, CI/CD

── 3. KUBERNETES + GITOPS ───────────────────────────────────
Namespaces:  prod, staging, dev, argocd, monitoring
ArgoCD:      latest stable via Kustomize into argocd namespace
Pattern:     App-of-Apps — apps/ folder, auto-sync + prune
Secrets:     Sealed-Secrets controller + SealedSecret example

── 4. CI/CD — GITHUB ACTIONS ────────────────────────────────
Trigger:     push to main + pull_request
Backend:     test → build Docker → push Artifact Registry
             → PR to infrastructure/ with new image tag
Frontend:    lint → test → build Next.js → push or CDN deploy
Auth:        Workload Identity Federation (no service account keys)

── 5. DATABASES ────────────────────────────────────────────
CloudSQL setup:
  PITR: 7-day retention
  Databases: borelsigma_main, borelsigma_timeseries
  User: backend_user (least-privilege)
  Connection: Cloud SQL Proxy sidecar pattern
CronJob K8s:  daily pg_dump → GCS bucket db-backups

── 6. MONITORING ───────────────────────────────────────────
Helm:        kube-prometheus-stack → monitoring namespace
Grafana:     cluster resources + custom EMS metrics dashboards
Loki:        log aggregation + Promtail DaemonSet
Alerts:      high CPU, pod restarts, service down, anomaly spike

── 7. VALIDATION ───────────────────────────────────────────
Verify:
  ArgoCD healthy + synced
  hello-world microservice builds via GH Actions → dev namespace
  DB connection test from pod (SQL Proxy sidecar)
  Prometheus scrapes metrics from test pod
Output: Directory tree + key file contents for all modules
# ═══════════════════════════════════════════════════════════════
# Start with repository creation, then Terraform modules.
# Generate production-ready, well-structured code throughout.
# ═══════════════════════════════════════════════════════════════`;
