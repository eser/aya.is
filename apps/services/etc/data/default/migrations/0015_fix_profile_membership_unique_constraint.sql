-- +goose Up
-- Fix unique constraint to only apply to non-deleted rows
-- This allows re-adding previously deleted members

DROP INDEX IF EXISTS "profile_membership_profile_id_member_profile_id_unique";

CREATE UNIQUE INDEX IF NOT EXISTS "profile_membership_profile_id_member_profile_id_unique"
ON "profile_membership" ("profile_id", "member_profile_id")
WHERE "member_profile_id" IS NOT NULL AND "deleted_at" IS NULL;

-- +goose Down
DROP INDEX IF EXISTS "profile_membership_profile_id_member_profile_id_unique";

CREATE UNIQUE INDEX IF NOT EXISTS "profile_membership_profile_id_member_profile_id_unique"
ON "profile_membership" ("profile_id", "member_profile_id")
WHERE "member_profile_id" IS NOT NULL;
