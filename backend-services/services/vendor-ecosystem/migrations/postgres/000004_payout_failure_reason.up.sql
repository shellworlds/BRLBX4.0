ALTER TABLE payout_requests
ADD COLUMN IF NOT EXISTS failure_reason TEXT;
