-- +goose Up
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS "apple_remote_id" TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS "user_apple_remote_id_unique"
  ON "user" ("apple_remote_id") WHERE "apple_remote_id" IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS "user_apple_remote_id_unique";
ALTER TABLE "user" DROP COLUMN IF EXISTS "apple_remote_id";
