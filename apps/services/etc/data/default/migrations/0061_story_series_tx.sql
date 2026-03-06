-- +goose Up

-- 1. Create story_series_tx for locale-aware title/description
CREATE TABLE IF NOT EXISTS "story_series_tx" (
  "story_series_id" CHAR(26) NOT NULL
    CONSTRAINT "story_series_tx_story_series_id_fk" REFERENCES "story_series",
  "locale_code" CHAR(12) NOT NULL,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  PRIMARY KEY ("story_series_id", "locale_code")
);

-- 2. Migrate existing data into 'en' locale
INSERT INTO "story_series_tx" (story_series_id, locale_code, title, description)
SELECT id, 'en', title, description
FROM "story_series"
WHERE deleted_at IS NULL;

-- 3. Drop inline text columns from story_series
ALTER TABLE "story_series" DROP COLUMN "title";
ALTER TABLE "story_series" DROP COLUMN "description";

-- 4. Add sort_order to story for series ordering (NULL = use published_at)
ALTER TABLE "story" ADD COLUMN "sort_order" INTEGER;

-- 5. Index for efficient series story listing
CREATE INDEX "story_series_id_sort_idx"
  ON "story" ("series_id", "sort_order" ASC NULLS LAST)
  WHERE "deleted_at" IS NULL AND "series_id" IS NOT NULL;

-- +goose Down

DROP INDEX IF EXISTS "story_series_id_sort_idx";
ALTER TABLE "story" DROP COLUMN IF EXISTS "sort_order";

ALTER TABLE "story_series" ADD COLUMN "title" TEXT NOT NULL DEFAULT '';
ALTER TABLE "story_series" ADD COLUMN "description" TEXT NOT NULL DEFAULT '';
UPDATE "story_series" ss SET
  title = COALESCE((SELECT sst.title FROM "story_series_tx" sst WHERE sst.story_series_id = ss.id AND RTRIM(sst.locale_code) = 'en' LIMIT 1), ''),
  description = COALESCE((SELECT sst.description FROM "story_series_tx" sst WHERE sst.story_series_id = ss.id AND RTRIM(sst.locale_code) = 'en' LIMIT 1), '');

DROP TABLE IF EXISTS "story_series_tx";
