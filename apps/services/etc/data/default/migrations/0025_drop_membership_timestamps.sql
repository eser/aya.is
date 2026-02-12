-- +goose Up
-- Drop created_at and updated_at from profile_membership.
-- These are now tracked via event_audit (profile_membership_created, profile_membership_updated).

ALTER TABLE "profile_membership" DROP COLUMN "created_at";
ALTER TABLE "profile_membership" DROP COLUMN "updated_at";

-- +goose Down
ALTER TABLE "profile_membership" ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL;
ALTER TABLE "profile_membership" ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE;
