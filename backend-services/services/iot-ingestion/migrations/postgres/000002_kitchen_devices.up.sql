CREATE TABLE IF NOT EXISTS kitchen_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kitchen_id UUID NOT NULL,
    label TEXT NOT NULL DEFAULT '',
    csr_pem TEXT NOT NULL,
    cert_pem TEXT,
    serial_number TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (kitchen_id, label)
);

CREATE INDEX IF NOT EXISTS kitchen_devices_kitchen_idx ON kitchen_devices (kitchen_id);
