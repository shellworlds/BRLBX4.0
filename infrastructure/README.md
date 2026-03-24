# infrastructure

Terraform and Kubernetes foundation for **Borel Sigma** on GCP (Day 1 baseline).  
This monorepo uses `shellworlds/BRLBX4.0`; Argo CD `repoURL` values point at that repository.

## Repository bootstrap (optional split)

If you later split into three GitHub repositories, recreate the same folder roots under each repo and update Argo `repoURL` fields accordingly. For a single monorepo, you can ignore this block.

```bash
export GH_OWNER="shellworlds"
for repo in infrastructure backend-services frontend; do
  gh repo create "${GH_OWNER}/${repo}" --public --clone=false --description "Borel Sigma ${repo}"
done
```

## GCP prerequisites

```bash
gcloud config set project borel-sigma-prod
gcloud services enable \
  container.googleapis.com \
  sqladmin.googleapis.com \
  artifactregistry.googleapis.com \
  compute.googleapis.com \
  servicenetworking.googleapis.com \
  iam.googleapis.com \
  storage.googleapis.com

gsutil mb -l us-central1 "gs://tfstate-borel-sigma" || true
gsutil versioning set on "gs://tfstate-borel-sigma"
```

## Terraform apply

State backend: **GCS** bucket `tfstate-borel-sigma` (see `terraform/versions.tf`).

```bash
cd terraform
terraform init
terraform plan -out=tfplan
terraform apply tfplan
```

### After apply — tighten GKE API access

The GKE module may include a **temporary** wide authorized-networks CIDR for bootstrap. Replace `0.0.0.0/0` with your workstation / bastion CIDRs before production traffic.

## GitOps bootstrap (Argo CD)

1. Fetch GKE credentials:

```bash
gcloud container clusters get-credentials borel-sigma-cluster --region us-central1 --project borel-sigma-prod
```

2. Install namespaces (if not already applied by Argo later):

```bash
kubectl apply -k k8s/namespaces
```

3. Install Argo CD:

```bash
kubectl apply -k k8s/argocd
kubectl wait --for=condition=available deployment/argocd-server -n argocd --timeout=300s
```

4. **Initial admin password** (default install):

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d; echo
```

5. Port-forward the UI (ingress can be added later):

```bash
kubectl -n argocd port-forward svc/argocd-server 8080:443
```

Open `https://localhost:8080` (accept the self-signed cert in dev), user `admin`, password from the command above.

6. Register this Git repo in Argo CD (if not using the default in-cluster credentials / public HTTPS):

```bash
argocd repo add https://github.com/shellworlds/BRLBX4.0.git
```

7. Apply the **App of Apps** root `Application` once (or let CI apply it):

```bash
kubectl apply -f k8s/bootstrap/root-application.yaml
```

The root app syncs `infrastructure/k8s/apps/app-of-apps` with **automated sync + prune** for child apps (namespaces, `hello-world` dev, sealed-secrets controller).

### Prometheus CRDs and optional overlays

- Default `hello-world` Dev overlay avoids `ServiceMonitor` so a fresh cluster without Prometheus Operator stays green.
- After installing `kube-prometheus-stack` (see `k8s/monitoring/README.md`), switch the Argo Application path to `infrastructure/k8s/apps/hello-world/overlays/dev-with-prometheus` or duplicate an Application for that overlay.

## Manual post-provision database steps

1. Verify point-in-time recovery is enabled with 7-day retention on all Cloud SQL instances (Terraform enables this; confirm in console).
2. Create databases:
   - `borelsigma_main` on each transactional instance
   - `borelsigma_timeseries` on the timeseries instance
3. Create restricted app user `backend_user` with minimal privileges.

Example:

```bash
gcloud sql users create backend_user --instance=borel-sigma-tx-1 --password='REPLACE_ME'
gcloud sql databases create borelsigma_main --instance=borel-sigma-tx-1
gcloud sql databases create borelsigma_main --instance=borel-sigma-tx-2
gcloud sql databases create borelsigma_timeseries --instance=borel-sigma-timeseries-1
```

4. **Cloud SQL Proxy** sidecars: see `k8s/jobs/cloudsql-proxy-test-pod.yaml` as a template (replace connection name, secrets).

5. **Logical backups**: Terraform provisions `${project_id}-db-backups` (e.g. `borel-sigma-prod-db-backups`). Apply `k8s/jobs/db-backup-cronjob.yaml` after creating `db-backup-credentials` and Workload Identity bindings for `gsutil`.

### Sealed Secrets

Controller installs from `k8s/sealed-secrets`. Generate real secrets with `kubeseal`; see `k8s/sealed-secrets/sample-sealedsecret.yaml.example`.

## Validation checklist (Day 1)

- [ ] Terraform apply succeeds; outputs list cluster, Artifact Registry, Cloud SQL connection names, backup bucket.
- [ ] Argo CD UI reachable; `borelsigma-root` and child apps **Healthy/Synced**.
- [ ] `hello-world` Deployment in `dev` serves `/` and `/metrics`.
- [ ] GitHub Actions green on `main` (tests; image push after WIF secrets).
- [ ] After Prometheus Operator is installed, optional ServiceMonitor overlay; confirm Prometheus **Targets** include `hello-world`.
- [ ] Cloud SQL proxy test pod connects after DB/user/password exist.

## Provisioned resources (summary)

- Dedicated VPC, private subnet, private services access range, Cloud NAT
- Firewall posture (internal allow + deny untagged public paths for node tags used here)
- GKE regional cluster **borel-sigma-cluster** (`REGULAR` channel), node pools **general-pool** and **high-mem-pool**
- Cloud SQL: 2× HA PostgreSQL (transactional), 1× HA PostgreSQL with Timescale preload flag
- Artifact Registry Docker repository, GCS bucket for DB backups, least-privilege service accounts for GKE / DB identity / CI
