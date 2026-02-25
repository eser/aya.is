-- +goose Up
ALTER TABLE "story_tx" ADD COLUMN "is_managed" BOOLEAN DEFAULT FALSE NOT NULL;

-- Mark all existing translations of managed stories as managed
UPDATE "story_tx" SET is_managed = TRUE
WHERE story_id IN (SELECT id FROM "story" WHERE is_managed = TRUE);

-- +goose Down
ALTER TABLE "story_tx" DROP COLUMN IF EXISTS "is_managed";
