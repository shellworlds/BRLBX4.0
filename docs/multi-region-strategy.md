# Multi-region data residency (roadmap)

## Data model (Day 5)

- `kitchens.region`, vendor records, and auth `users.region` carry a **region code** (e.g. `EU`, `IN`, `global`).
- `emission_factors` in `energy-management` stores grid intensity per region for GHG reporting.

## Application behavior

- **Middleware**: resolve `region` from JWT (`app_metadata` / custom claim) and compare against resource region for sensitive reads (extend per service as needed).
- **Storage**: for full residency, deploy **one Postgres cluster per region** (or row-level policies with region column) and route traffic from a regional ingress.

## Deployment sketch

1. **Regional GKE clusters** (or namespaces with network policies) with the same Helm/Kustomize bases.
2. **Global DNS** (GeoDNS or latency-based routing) to the nearest ingress.
3. **Auth0** organizations or separate tenants per region if policy requires strict isolation.
4. **Replication**: cross-region read replicas only where compliance allows; never replicate EU personal data to non-EU without legal basis.

## Current repo state

This repo encodes **region fields** and **reporting**; full multi-region routing and per-region DBs are **not** automated here—use this document as the runbook for Phase 2.
