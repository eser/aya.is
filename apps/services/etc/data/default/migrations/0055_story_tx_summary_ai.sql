-- +goose Up

ALTER TABLE "story_tx" ADD COLUMN "summary_ai" TEXT;

-- +goose Down

ALTER TABLE "story_tx" DROP COLUMN IF EXISTS "summary_ai";
