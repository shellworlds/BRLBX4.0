CREATE TABLE IF NOT EXISTS daily_vendor_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    vendor_id UUID NOT NULL REFERENCES vendors (id) ON DELETE CASCADE,
    day DATE NOT NULL,
    meals_served INT NOT NULL DEFAULT 0,
    revenue_total DOUBLE PRECISION NOT NULL DEFAULT 0,
    energy_efficiency_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    compliance_score INT,
    payload JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (vendor_id, day)
);

CREATE INDEX IF NOT EXISTS daily_vendor_metrics_vendor_day_idx ON daily_vendor_metrics (vendor_id, day DESC);
