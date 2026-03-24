# Disaster recovery drill (single-region)

## Scope

This repo assumes a **single active GKE region** with multi-zone Cloud SQL where enabled.

## Node cordon drill

1. `kubectl cordon <node>` on a subset of nodes.
2. Confirm workloads reschedule within **2 minutes** (`kubectl get pods -A -w`).
3. `kubectl uncordon <node>` when finished.

## Namespace loss (tabletop)

1. Restore from Terraform / Argo CD: re-apply `infrastructure/k8s` and namespaces.
2. Restore Postgres from latest **automated backup** (Cloud SQL PITR).
3. Rotate any secrets that may have been exposed.

## RTO / RPO targets

Document actual times after each drill in `docs/day6-progress.md` (or wiki).

## Cloud SQL failover

If HA is enabled, trigger **failover** in a maintenance window and verify DSN connectivity from Cloud SQL Proxy pods.
