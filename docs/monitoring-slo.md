# Monitoring & SLO notes

## Grafana dashboards (recommended panels)

- **Business**: daily meals (vendor reports), energy snapshot (`/api/v1/public/snapshot` scrape), emissions (GHG report job).
- **Technical**: request rate, p95 latency, 5xx ratio per service, pod restarts, MQTT publish rate (EMQX exporter if used).

## Prometheus

- Standardize on `http_requests_total` with labels `handler`, `method`, `status` (add middleware in Go services if missing).
- Ingest lag: compare `max(time()) - max(reading_ts)` via recording rule on Timescale / custom exporter.

## Slack

- Route `iot-ingestion` anomaly webhooks (existing `SLACK_WEBHOOK_URL`) to `#ops-alerts`.
- Route Alertmanager critical to on-call; warnings to Slack.

## Testing alerts

- `kubectl delete pod -l app.kubernetes.io/name=energy-management -n staging` — expect restart alert if thresholds tuned.
