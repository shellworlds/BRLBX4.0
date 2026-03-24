# k6 load tests (staging)

Install [k6](https://k6.io/docs/get-started/installation/). Point `BASE_URL` at the **public** API host (through ingress), e.g. `https://api.borelsigma.com` or a staging URL.

```bash
# Public read / snapshot style traffic (tune VUs/duration in script)
BASE_URL=https://api.borelsigma.com k6 run perf/k6/portal-read.js

# Energy ingest style POSTs (requires valid paths/auth in script)
BASE_URL=https://api.borelsigma.com k6 run perf/k6/energy-ingest.js

# Vendor transaction posts
BASE_URL=https://api.borelsigma.com k6 run perf/k6/vendor-transactions.js
```

Run each scenario for **≥10 minutes** before launch; capture p95 latency, error rate, and (if available) HPA/KEDA events. Archive summaries under `docs/` or attach to the release.
