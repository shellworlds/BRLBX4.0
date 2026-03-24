CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS raw_telemetry (
    ts TIMESTAMPTZ NOT NULL,
    kitchen_id UUID NOT NULL,
    topic TEXT NOT NULL,
    payload JSONB NOT NULL
);

SELECT public.create_hypertable('raw_telemetry', 'ts', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS raw_telemetry_kitchen_ts_idx ON raw_telemetry (kitchen_id, ts DESC);
