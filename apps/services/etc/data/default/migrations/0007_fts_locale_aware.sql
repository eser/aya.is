-- +goose Up
-- Replace 'simple' text search config with locale-aware configurations
-- for better stemming and linguistic analysis.

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION locale_to_regconfig(locale TEXT) RETURNS regconfig AS $$
BEGIN
  RETURN CASE locale
    WHEN 'de' THEN 'german'::regconfig
    WHEN 'en' THEN 'english'::regconfig
    WHEN 'es' THEN 'spanish'::regconfig
    WHEN 'fr' THEN 'french'::regconfig
    WHEN 'it' THEN 'italian'::regconfig
    WHEN 'nl' THEN 'dutch'::regconfig
    WHEN 'pt-PT' THEN 'portuguese'::regconfig
    WHEN 'ru' THEN 'russian'::regconfig
    WHEN 'tr' THEN 'turkish'::regconfig
    ELSE 'simple'::regconfig
  END;
END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +goose StatementEnd

-- Rebuild search_vector columns with locale-aware configs

ALTER TABLE "profile_tx" DROP COLUMN "search_vector";
ALTER TABLE "profile_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector(locale_to_regconfig(locale_code), COALESCE(title, '')), 'A') ||
  setweight(to_tsvector(locale_to_regconfig(locale_code), COALESCE(description, '')), 'B')
) STORED;

ALTER TABLE "story_tx" DROP COLUMN "search_vector";
ALTER TABLE "story_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector(locale_to_regconfig(locale_code), COALESCE(title, '')), 'A') ||
  setweight(to_tsvector(locale_to_regconfig(locale_code), COALESCE(summary, '')), 'B') ||
  setweight(to_tsvector(locale_to_regconfig(locale_code), COALESCE(content, '')), 'C')
) STORED;

ALTER TABLE "profile_page_tx" DROP COLUMN "search_vector";
ALTER TABLE "profile_page_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector(locale_to_regconfig(locale_code), COALESCE(title, '')), 'A') ||
  setweight(to_tsvector(locale_to_regconfig(locale_code), COALESCE(summary, '')), 'B') ||
  setweight(to_tsvector(locale_to_regconfig(locale_code), COALESCE(content, '')), 'C')
) STORED;

-- Recreate GIN indexes
CREATE INDEX "idx_profile_tx_search" ON "profile_tx" USING GIN ("search_vector");
CREATE INDEX "idx_story_tx_search" ON "story_tx" USING GIN ("search_vector");
CREATE INDEX "idx_profile_page_tx_search" ON "profile_page_tx" USING GIN ("search_vector");

-- +goose Down
DROP INDEX IF EXISTS "idx_profile_page_tx_search";
DROP INDEX IF EXISTS "idx_story_tx_search";
DROP INDEX IF EXISTS "idx_profile_tx_search";

ALTER TABLE "profile_page_tx" DROP COLUMN "search_vector";
ALTER TABLE "profile_page_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
  setweight(to_tsvector('simple', COALESCE(summary, '')), 'B') ||
  setweight(to_tsvector('simple', COALESCE(content, '')), 'C')
) STORED;

ALTER TABLE "story_tx" DROP COLUMN "search_vector";
ALTER TABLE "story_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
  setweight(to_tsvector('simple', COALESCE(summary, '')), 'B') ||
  setweight(to_tsvector('simple', COALESCE(content, '')), 'C')
) STORED;

ALTER TABLE "profile_tx" DROP COLUMN "search_vector";
ALTER TABLE "profile_tx" ADD COLUMN "search_vector" tsvector GENERATED ALWAYS AS (
  setweight(to_tsvector('simple', COALESCE(title, '')), 'A') ||
  setweight(to_tsvector('simple', COALESCE(description, '')), 'B')
) STORED;

CREATE INDEX "idx_profile_tx_search" ON "profile_tx" USING GIN ("search_vector");
CREATE INDEX "idx_story_tx_search" ON "story_tx" USING GIN ("search_vector");
CREATE INDEX "idx_profile_page_tx_search" ON "profile_page_tx" USING GIN ("search_vector");

DROP FUNCTION IF EXISTS locale_to_regconfig(TEXT);
