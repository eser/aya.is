-- Add configurable digest frequency to bulletin subscriptions.
-- Values: 'daily' (default), 'bidaily' (every 2 days), 'weekly'.
-- "Don't send" is represented by soft-deleting all subscriptions.
ALTER TABLE "bulletin_subscription"
  ADD COLUMN "frequency" TEXT NOT NULL DEFAULT 'daily'
    CHECK (frequency IN ('daily', 'bidaily', 'weekly'));
