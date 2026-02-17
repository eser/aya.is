-- +goose Up
-- Profile Resources table: extensible external resource references (GitHub repos, etc.)

CREATE TABLE "profile_resource" (
  "id"                   CHAR(26) NOT NULL,
  "profile_id"           CHAR(26) NOT NULL,
  "kind"                 TEXT NOT NULL,
  "is_managed"           BOOLEAN NOT NULL DEFAULT false,
  "remote_id"            TEXT,
  "public_id"            TEXT,
  "url"                  TEXT,
  "title"                TEXT NOT NULL DEFAULT '',
  "description"          TEXT,
  "properties"           JSONB,
  "added_by_profile_id"  CHAR(26) NOT NULL,
  "created_at"           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at"           TIMESTAMPTZ,
  "deleted_at"           TIMESTAMPTZ,
  CONSTRAINT "profile_resource_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "profile_resource_profile_id_fk" FOREIGN KEY ("profile_id") REFERENCES "profile" ("id"),
  CONSTRAINT "profile_resource_added_by_profile_id_fk" FOREIGN KEY ("added_by_profile_id") REFERENCES "profile" ("id")
);

-- Prevent duplicate resources per profile+kind+remote_id
CREATE UNIQUE INDEX "profile_resource_profile_kind_remote_id_unique"
  ON "profile_resource" ("profile_id", "kind", "remote_id")
  WHERE "remote_id" IS NOT NULL AND "deleted_at" IS NULL;

-- For listing resources by profile
CREATE INDEX "profile_resource_profile_id_idx"
  ON "profile_resource" ("profile_id")
  WHERE "deleted_at" IS NULL;

-- +goose Down
DROP TABLE IF EXISTS "profile_resource";
