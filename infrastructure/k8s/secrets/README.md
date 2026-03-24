# Kubernetes secrets (Sealed Secrets)

Use [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets) so only the cluster controller can decrypt `SealedSecret` resources committed to Git.

## What to seal

- Stripe: `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET` (payments deployment)
- Auth0: M2M / client secrets where backends call Auth0 (if any)
- Internal tokens: `INTERNAL_PAYMENTS_TOKEN`, `INTERNAL_DEVICE_TOKEN`, `INTERNAL_AGGREGATE_TOKEN`
- Database DSNs: `POSTGRES_*_DSN` per service
- SMTP: compliance service mail credentials

## Flow

1. Install the controller and `kubeseal` CLI (see `sealed-secrets/kustomization.yaml` in this repo).
2. Create a normal `Secret` manifest locally, then `kubeseal -f secret.yaml -w sealed-secret.yaml`.
3. Commit only `sealed-secret.yaml`; delete the raw secret from disk.
4. Apply the sealed manifest; the controller materializes a `Secret` in the target namespace.

## Rotation

1. Generate new secret material in the provider (Stripe, Auth0, etc.).
2. Update the sealed manifest with `kubeseal` and bump the deployment (restart pods) so env vars reload.
3. Revoke old keys where the provider supports it.

## References

- Example: `infrastructure/k8s/sealed-secrets/sample-sealedsecret.yaml.example`
