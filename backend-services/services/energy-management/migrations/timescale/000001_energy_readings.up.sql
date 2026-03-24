CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS energy_readings (
    kitchen_id UUID NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    grid_power DOUBLE PRECISION NOT NULL DEFAULT 0,
    battery_power DOUBLE PRECISION NOT NULL DEFAULT 0,
    solar_power DOUBLE PRECISION NOT NULL DEFAULT 0,
    lpg_status TEXT NOT NULL DEFAULT 'unknown',
    uptime_percent DOUBLE PRECISION NOT NULL DEFAULT 0
);

SELECT public.create_hypertable('energy_readings', 'ts', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS energy_readings_kitchen_ts_idx ON energy_readings (kitchen_id, ts DESC);
