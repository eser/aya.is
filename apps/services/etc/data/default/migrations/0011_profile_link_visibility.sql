-- +goose Up
-- Add visibility and featured flag columns
ALTER TABLE "profile_link"
ADD COLUMN "visibility" TEXT DEFAULT 'public' NOT NULL;

ALTER TABLE "profile_link"
ADD COLUMN "is_featured" BOOLEAN DEFAULT TRUE NOT NULL;

-- Create translation table for profile links
CREATE TABLE IF NOT EXISTS "profile_link_tx" (
  "profile_link_id" CHAR(26) NOT NULL,
  "locale_code" CHAR(12) NOT NULL,
  "title" TEXT NOT NULL,
  "group" TEXT,
  "description" TEXT,
  PRIMARY KEY ("profile_link_id", "locale_code")
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS "profile_link_visibility_idx"
ON "profile_link" ("profile_id", "visibility")
WHERE "deleted_at" IS NULL;

CREATE INDEX IF NOT EXISTS "profile_link_featured_idx"
ON "profile_link" ("profile_id", "is_featured")
WHERE "deleted_at" IS NULL;

-- +goose Down
DROP INDEX IF EXISTS "profile_link_featured_idx";
DROP INDEX IF EXISTS "profile_link_visibility_idx";
DROP TABLE IF EXISTS "profile_link_tx";
ALTER TABLE "profile_link" DROP COLUMN IF EXISTS "is_featured";
ALTER TABLE "profile_link" DROP COLUMN IF EXISTS "visibility";
