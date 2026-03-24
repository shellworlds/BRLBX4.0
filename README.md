# BRLBX4.0 — Borel Sigma platform monorepo

Cloud-native baseline for **Borel Sigma** (GCP, GKE, GitOps, GitHub Actions). Layout mirrors the Day‑1 plan while keeping a single repository for velocity; split into separate GitHub repos later if you want stricter boundaries.

| Path | Purpose |
|------|---------|
| `infrastructure/` | Terraform (GCP), Kubernetes manifests, Argo CD, monitoring docs |
| `backend-services/` | Microservices (Go `hello-world` today) |
| `frontend/` | Next.js app (standalone Docker image) |
| `docs/` | Architecture and runbooks |

## Quick links

- Terraform and GCP: `infrastructure/README.md`
- Architecture diagram (text): `docs/architecture.md`
- Argo CD bootstrap: `infrastructure/README.md` (“GitOps bootstrap”)
- Monitoring (Prometheus / Loki): `infrastructure/k8s/monitoring/README.md`

## GitHub automation

Workflows live in `.github/workflows/`:

- `backend-ci-cd.yml` — Go tests; on `main`, optional image push + Kustomize bump PR (needs GCP WIF secrets).
- `frontend-ci-cd.yml` — lint, test, `next build`; optional Docker push on `main`.
- `infrastructure-terraform.yml` — `fmt`, `validate`.
- `k8s-manifests.yml` — `kustomize build` smoke checks.

### Required repository secrets (for image push)

Configure **Workload Identity Federation** for GitHub ↔ GCP, then add:

- `GCP_WORKLOAD_IDENTITY_PROVIDER`
- `GCP_SERVICE_ACCOUNT_EMAIL`

Until these are set, CI still runs builds/tests but skips pushes (see workflow `::notice` lines).

## Local development

```bash
# Backend
cd backend-services/services/hello-world && go test ./...

# Frontend
cd frontend && npm ci && npm run dev
```

## License

See `infrastructure/LICENSE` (MIT). Copy or consolidate per your legal preference.
