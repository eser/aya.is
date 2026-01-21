-- +goose Up
-- Add tsvector columns with weighted search for full-text search

ALTER TABLE "profile_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
  setweight(to_tsvector('simple', COALESCE(description, '')), 'B')
) STORED;

ALTER TABLE "story_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
  setweight(to_tsvector('simple', COALESCE(summary, '')), 'B') ||
  setweight(to_tsvector('simple', COALESCE(content, '')), 'C')
) STORED;

ALTER TABLE "profile_page_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
  setweight(to_tsvector('simple', COALESCE(summary, '')), 'B') ||
  setweight(to_tsvector('simple', COALESCE(content, '')), 'C')
) STORED;

-- GIN indexes for fast lookups
CREATE INDEX "idx_profile_tx_search" ON "profile_tx" USING GIN ("search_vector");
CREATE INDEX "idx_story_tx_search" ON "story_tx" USING GIN ("search_vector");
CREATE INDEX "idx_profile_page_tx_search" ON "profile_page_tx" USING GIN ("search_vector");

-- +goose Down
DROP INDEX IF EXISTS "idx_profile_page_tx_search";
DROP INDEX IF EXISTS "idx_story_tx_search";
DROP INDEX IF EXISTS "idx_profile_tx_search";

ALTER TABLE "profile_page_tx" DROP COLUMN IF EXISTS "search_vector";
ALTER TABLE "story_tx" DROP COLUMN IF EXISTS "search_vector";
ALTER TABLE "profile_tx" DROP COLUMN IF EXISTS "search_vector";
