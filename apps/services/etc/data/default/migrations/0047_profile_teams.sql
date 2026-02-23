-- +goose Up

CREATE TABLE IF NOT EXISTS "profile_team" (
  "id"          CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id"  CHAR(26) NOT NULL REFERENCES "profile" ("id"),
  "name"        TEXT NOT NULL,
  "description" TEXT,
  "created_at"  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "deleted_at"  TIMESTAMP WITH TIME ZONE
);

CREATE INDEX "profile_team_profile_id_idx"
  ON "profile_team" ("profile_id")
  WHERE "deleted_at" IS NULL;

CREATE TABLE IF NOT EXISTS "profile_membership_team" (
  "id"                    CHAR(26) NOT NULL PRIMARY KEY,
  "profile_membership_id" CHAR(26) NOT NULL REFERENCES "profile_membership" ("id"),
  "profile_team_id"       CHAR(26) NOT NULL REFERENCES "profile_team" ("id"),
  "created_at"            TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "deleted_at"            TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX "profile_membership_team_unique"
  ON "profile_membership_team" ("profile_membership_id", "profile_team_id")
  WHERE "deleted_at" IS NULL;

CREATE INDEX "profile_membership_team_team_idx"
  ON "profile_membership_team" ("profile_team_id")
  WHERE "deleted_at" IS NULL;

-- +goose Down
DROP TABLE IF EXISTS "profile_membership_team";
DROP TABLE IF EXISTS "profile_team";
