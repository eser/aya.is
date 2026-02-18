-- +goose Up
ALTER TABLE "story" ADD COLUMN "remote_id" TEXT;

UPDATE "story" SET remote_id = properties->>'remote_id'
  WHERE is_managed = TRUE AND properties->>'remote_id' IS NOT NULL;

CREATE UNIQUE INDEX "story_author_profile_id_remote_id_unique"
  ON "story" ("author_profile_id", "remote_id")
  WHERE "remote_id" IS NOT NULL AND "deleted_at" IS NULL;

-- +goose Down
DROP INDEX IF EXISTS "story_author_profile_id_remote_id_unique";
ALTER TABLE "story" DROP COLUMN IF EXISTS "remote_id";
