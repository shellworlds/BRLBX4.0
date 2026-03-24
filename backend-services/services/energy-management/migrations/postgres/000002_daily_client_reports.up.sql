CREATE TABLE IF NOT EXISTS daily_client_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    client_id TEXT NOT NULL,
    day DATE NOT NULL,
    scope3_tco2e_avoided DOUBLE PRECISION NOT NULL DEFAULT 0,
    solar_share DOUBLE PRECISION NOT NULL DEFAULT 0,
    grid_share DOUBLE PRECISION NOT NULL DEFAULT 0,
    battery_share DOUBLE PRECISION NOT NULL DEFAULT 0,
    uptime_avg DOUBLE PRECISION NOT NULL DEFAULT 0,
    payload JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (client_id, day)
);

CREATE INDEX IF NOT EXISTS daily_client_reports_client_day_idx ON daily_client_reports (client_id, day DESC);
