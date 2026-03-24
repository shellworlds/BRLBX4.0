CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS ingestion_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kitchen_id UUID NOT NULL,
    level TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    acknowledged_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS ingestion_alerts_kitchen_created_idx ON ingestion_alerts (kitchen_id, created_at DESC);
