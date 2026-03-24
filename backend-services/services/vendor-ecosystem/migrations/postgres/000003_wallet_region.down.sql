DROP TABLE IF EXISTS payout_requests;
DROP TABLE IF EXISTS wallet_ledger;
DROP TABLE IF EXISTS vendor_wallets;
ALTER TABLE vendors DROP COLUMN IF EXISTS stripe_connect_account_id;
ALTER TABLE vendors DROP COLUMN IF EXISTS region;
