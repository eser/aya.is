-- +goose Up

CREATE TABLE IF NOT EXISTS "profile_resource_team" (
  "id"                  CHAR(26) NOT NULL PRIMARY KEY,
  "profile_resource_id" CHAR(26) NOT NULL REFERENCES "profile_resource" ("id"),
  "profile_team_id"     CHAR(26) NOT NULL REFERENCES "profile_team" ("id"),
  "created_at"          TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "deleted_at"          TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX "profile_resource_team_unique"
  ON "profile_resource_team" ("profile_resource_id", "profile_team_id")
  WHERE "deleted_at" IS NULL;

CREATE INDEX "profile_resource_team_team_idx"
  ON "profile_resource_team" ("profile_team_id")
  WHERE "deleted_at" IS NULL;

-- +goose Down
DROP TABLE IF EXISTS "profile_resource_team";
