-- +goose Up
-- Add last_activity_at and user_agent columns to session table
ALTER TABLE "session" ADD COLUMN "last_activity_at" TIMESTAMPTZ;
ALTER TABLE "session" ADD COLUMN "user_agent" TEXT;

-- Create index for querying sessions by user (for session management feature)
CREATE INDEX "session_logged_in_user_id_idx" ON "session" ("logged_in_user_id")
WHERE "logged_in_user_id" IS NOT NULL;

-- Create index for querying active sessions
CREATE INDEX "session_status_idx" ON "session" ("status")
WHERE "status" = 'active';

-- +goose Down
DROP INDEX IF EXISTS "session_status_idx";
DROP INDEX IF EXISTS "session_logged_in_user_id_idx";
ALTER TABLE "session" DROP COLUMN IF EXISTS "user_agent";
ALTER TABLE "session" DROP COLUMN IF EXISTS "last_activity_at";
