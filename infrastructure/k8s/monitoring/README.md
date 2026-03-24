# Monitoring and logging

## kube-prometheus-stack (Prometheus, Alertmanager, Grafana)

After Argo CD is installed and can sync Helm charts:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade --install kube-prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring --create-namespace \
  -f helm/values-kube-prometheus.yaml
```

Access Grafana (initial setup — prefer SSO and external secrets in production):

```bash
kubectl -n monitoring port-forward svc/kube-prometheus-grafana 3000:80
```

## Loki + Promtail

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm upgrade --install loki grafana/loki \
  --namespace monitoring \
  -f helm/values-loki.yaml
```

Validate chart field names against the chart version you install; upgrade `values-loki.yaml` when bumping majors.

## Argo CD (optional)

You can add `Application` resources that reference these Helm charts once your Argo CD instance has the Helm repositories configured. Inline `spec.source.helm.values` in those Applications is often simpler than multi-source `valueFiles` for a first pass.

## Alerts

`kube-prometheus-stack` ships default alert rules. Add custom `PrometheusRule` objects under `rules/` and apply them after the operator CRDs exist, or extend `values-kube-prometheus.yaml` with `additionalPrometheusRulesMap`.
