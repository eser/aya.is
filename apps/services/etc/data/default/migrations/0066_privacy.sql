-- +goose Up

-- Self-service privacy: hiding a profile until the next login, exporting personal data, and account deletion.

-- A profile hidden by its owner. Hiding reuses "deleted_at" (every read path already filters on it) and records
-- hidden_at so that the next login can tell a self-hide apart from an admin deletion and restore only the former.
ALTER TABLE "profile" ADD COLUMN "hidden_at" TIMESTAMP WITH TIME ZONE;

-- Set once an account is deleted and its personal data has been redacted. Never reverted.
ALTER TABLE "profile" ADD COLUMN "redacted_at" TIMESTAMP WITH TIME ZONE;
ALTER TABLE "user" ADD COLUMN "redacted_at" TIMESTAMP WITH TIME ZONE;

-- Redacted users all carry the literal '**redacted**' as their e-mail, so a *unique* e-mail index must ignore them.
-- Production runs this index as a plain, non-unique one (it has drifted from migration 0001, and 62 users share the
-- empty e-mail), so only a genuinely unique index is rebuilt; anything else is left exactly as it is.
-- +goose StatementBegin
DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM pg_indexes
    WHERE schemaname = current_schema()
      AND indexname = 'user_email_unique'
      AND indexdef LIKE 'CREATE UNIQUE INDEX%'
  ) THEN
    DROP INDEX "user_email_unique";
    CREATE UNIQUE INDEX "user_email_unique" ON "user" ("email") WHERE "email" IS NOT NULL AND "redacted_at" IS NULL;
  END IF;
END
$$;
-- +goose StatementEnd

-- Personal data export requests ("download my data"). The archive itself lives in a private object storage bucket;
-- this table only tracks the request lifecycle.
CREATE TABLE IF NOT EXISTS "profile_data_export" (
  "id"            CHAR(26) NOT NULL PRIMARY KEY,
  "user_id"       CHAR(26) NOT NULL CONSTRAINT "profile_data_export_user_id_fk" REFERENCES "user",
  "profile_id"    CHAR(26) NOT NULL CONSTRAINT "profile_data_export_profile_id_fk" REFERENCES "profile",
  "status"        TEXT NOT NULL DEFAULT 'pending',
  "requested_at"  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "started_at"    TIMESTAMP WITH TIME ZONE,
  "completed_at"  TIMESTAMP WITH TIME ZONE,
  "expires_at"    TIMESTAMP WITH TIME ZONE,
  "notified_at"   TIMESTAMP WITH TIME ZONE,
  "storage_key"   TEXT,
  "size_bytes"    BIGINT,
  "error_message" TEXT,
  "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "updated_at"    TIMESTAMP WITH TIME ZONE
);

CREATE INDEX "idx_profile_data_export_user_requested" ON "profile_data_export" ("user_id", "requested_at" DESC);
CREATE INDEX "idx_profile_data_export_status" ON "profile_data_export" ("status", "expires_at");

-- +goose Down

DROP TABLE IF EXISTS "profile_data_export";

-- Dropping the column also drops an index whose predicate mentions it; restore the original definition then.
-- A plain production index never depended on the column, so IF NOT EXISTS leaves it alone.
ALTER TABLE "user" DROP COLUMN IF EXISTS "redacted_at";
CREATE UNIQUE INDEX IF NOT EXISTS "user_email_unique" ON "user" ("email") WHERE "email" IS NOT NULL;
ALTER TABLE "profile" DROP COLUMN IF EXISTS "redacted_at";
ALTER TABLE "profile" DROP COLUMN IF EXISTS "hidden_at";
