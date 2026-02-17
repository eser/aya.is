-- +goose Up
-- Fix unique constraint on profile_link to exclude soft-deleted rows.
-- The old constraint includes deleted rows, preventing re-creation of links
-- with the same (profile_id, kind, remote_id) after soft-delete.

DROP INDEX IF EXISTS "profile_link_profile_id_kind_remote_id_unique";

CREATE UNIQUE INDEX "profile_link_profile_id_kind_remote_id_unique"
  ON "profile_link" ("profile_id", "kind", "remote_id")
  WHERE "remote_id" IS NOT NULL AND "deleted_at" IS NULL;

-- +goose Down
DROP INDEX IF EXISTS "profile_link_profile_id_kind_remote_id_unique";

CREATE UNIQUE INDEX "profile_link_profile_id_kind_remote_id_unique"
  ON "profile_link" ("profile_id", "kind", "remote_id")
  WHERE "remote_id" IS NOT NULL;
