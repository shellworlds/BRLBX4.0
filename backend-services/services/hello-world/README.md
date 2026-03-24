# hello-world

Minimal HTTP service exposing `/`, `/healthz`, and Prometheus metrics at `/metrics`.

```bash
go test ./...
docker build -t hello-world:dev .
```

Image published by GitHub Actions to Artifact Registry; Argo CD syncs tag from `infrastructure/k8s/apps/hello-world/overlays/dev/kustomization.yaml`.
