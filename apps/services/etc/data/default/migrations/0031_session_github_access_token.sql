-- +goose Up
-- Store OAuth access token, scope, and provider on session record.
-- These are set during login and used to auto-create managed profile links
-- (e.g., GitHub) when an individual profile is created within the same session.

ALTER TABLE "session" ADD COLUMN IF NOT EXISTS "oauth_provider" TEXT;
ALTER TABLE "session" ADD COLUMN IF NOT EXISTS "oauth_access_token" TEXT;
ALTER TABLE "session" ADD COLUMN IF NOT EXISTS "oauth_token_scope" TEXT;

-- +goose Down
ALTER TABLE "session" DROP COLUMN IF EXISTS "oauth_token_scope";
ALTER TABLE "session" DROP COLUMN IF EXISTS "oauth_access_token";
ALTER TABLE "session" DROP COLUMN IF EXISTS "oauth_provider";
