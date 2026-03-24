# Contributing

## Workflow

1. Fork or create a branch from `main`.
2. Keep changes focused; match existing style and naming.
3. Run local checks:
   - `go test ./...` under changed Go modules
   - `npm ci && npm run lint && npm test && npm run build` under `frontend/` when applicable
   - `terraform fmt -recursive` and `terraform validate` under `infrastructure/terraform/` when applicable
4. Open a pull request using the template; link related issues.

## CI

Pull requests run GitHub Actions for Terraform validation, Kustomize builds, Go, and Next.js. Fix failing jobs before requesting review.
