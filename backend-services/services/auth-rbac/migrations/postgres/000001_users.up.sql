CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    auth0_id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    role TEXT NOT NULL,
    client_id TEXT NULL,
    vendor_id TEXT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
