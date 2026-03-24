ALTER TABLE vendors
ADD COLUMN IF NOT EXISTS region TEXT NOT NULL DEFAULT 'global';

ALTER TABLE vendors
ADD COLUMN IF NOT EXISTS stripe_connect_account_id TEXT NULL;

CREATE TABLE IF NOT EXISTS vendor_wallets (
    vendor_id UUID PRIMARY KEY REFERENCES vendors (id) ON DELETE CASCADE,
    balance NUMERIC NOT NULL DEFAULT 0,
    pending_payout NUMERIC NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS wallet_ledger (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id UUID NOT NULL REFERENCES vendors (id) ON DELETE CASCADE,
    delta NUMERIC NOT NULL,
    reason TEXT NOT NULL,
    ref_type TEXT,
    ref_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS wallet_ledger_vendor_ts_idx ON wallet_ledger (vendor_id, created_at DESC);

CREATE TABLE IF NOT EXISTS payout_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id UUID NOT NULL REFERENCES vendors (id) ON DELETE CASCADE,
    amount NUMERIC NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    stripe_transfer_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS payout_requests_vendor_idx ON payout_requests (vendor_id, created_at DESC);
