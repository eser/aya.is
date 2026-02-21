-- +goose Up
-- Add visibility column to profile_page and story tables.
-- Three states: 'public' (default), 'unlisted', 'private'.
-- Deprecates the old publication-based visibility rule for stories.

ALTER TABLE "profile_page"
ADD COLUMN "visibility" TEXT DEFAULT 'public' NOT NULL;

ALTER TABLE "story"
ADD COLUMN "visibility" TEXT DEFAULT 'public' NOT NULL;

-- +goose Down
ALTER TABLE "story"
DROP COLUMN IF EXISTS "visibility";

ALTER TABLE "profile_page"
DROP COLUMN IF EXISTS "visibility";
