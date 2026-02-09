-- +goose Up
ALTER TABLE "story" ADD COLUMN "is_managed" BOOLEAN DEFAULT FALSE NOT NULL;
CREATE INDEX "story_is_managed_idx"
  ON "story" ("author_profile_id", "is_managed")
  WHERE "is_managed" = TRUE AND "deleted_at" IS NULL;

-- +goose Down
DROP INDEX IF EXISTS "story_is_managed_idx";
ALTER TABLE "story" DROP COLUMN IF EXISTS "is_managed";
