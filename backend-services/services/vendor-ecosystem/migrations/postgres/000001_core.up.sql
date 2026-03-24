CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS vendors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    fssai_score INT NOT NULL,
    location TEXT NOT NULL,
    contact TEXT NOT NULL,
    onboarding_date TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS financing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id UUID NOT NULL REFERENCES vendors (id) ON DELETE CASCADE,
    amount DOUBLE PRECISION NOT NULL,
    status TEXT NOT NULL, -- approved|pending|rejected
    repayment_schedule TEXT NOT NULL DEFAULT 'auto_pct_10',
    remaining_balance DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS financing_vendor_id_idx ON financing (vendor_id);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id UUID NOT NULL REFERENCES vendors (id) ON DELETE CASCADE,
    kitchen_id UUID NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    meal_count INT NOT NULL DEFAULT 0,
    ts TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS transactions_vendor_ts_idx ON transactions (vendor_id, ts DESC);
