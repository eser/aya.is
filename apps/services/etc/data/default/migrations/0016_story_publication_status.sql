-- +goose Up
-- Move status, is_featured, published_at from story to story_publication.
-- Each publication now independently tracks its own featured and publish state.
-- A story with no active publications is considered a draft.

-- Step 1: Add new columns to story_publication
ALTER TABLE "story_publication"
  ADD COLUMN IF NOT EXISTS "is_featured" BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS "published_at" TIMESTAMP WITH TIME ZONE;

-- Step 2: Migrate data from story to story_publication (for published stories)
UPDATE "story_publication" sp
SET
  is_featured = s.is_featured,
  published_at = COALESCE(s.published_at, s.created_at)
FROM "story" s
WHERE sp.story_id = s.id
  AND sp.deleted_at IS NULL
  AND s.status = 'published';

-- Step 3: Soft-delete publications for draft stories (shouldn't normally exist, but handle edge case)
UPDATE "story_publication" sp
SET deleted_at = NOW()
FROM "story" s
WHERE sp.story_id = s.id
  AND sp.deleted_at IS NULL
  AND s.status != 'published';

-- Step 4: Drop old columns from story
ALTER TABLE "story" DROP COLUMN IF EXISTS "status";
ALTER TABLE "story" DROP COLUMN IF EXISTS "is_featured";
ALTER TABLE "story" DROP COLUMN IF EXISTS "published_at";

-- +goose Down
-- Re-add columns to story
ALTER TABLE "story"
  ADD COLUMN IF NOT EXISTS "status" TEXT NOT NULL DEFAULT 'draft',
  ADD COLUMN IF NOT EXISTS "is_featured" BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS "published_at" TIMESTAMP WITH TIME ZONE;

-- Backfill story columns from story_publication
UPDATE "story" s
SET
  status = 'published',
  is_featured = sp.is_featured,
  published_at = sp.published_at
FROM "story_publication" sp
WHERE sp.story_id = s.id
  AND sp.deleted_at IS NULL;

-- Drop columns from story_publication
ALTER TABLE "story_publication" DROP COLUMN IF EXISTS "is_featured";
ALTER TABLE "story_publication" DROP COLUMN IF EXISTS "published_at";
