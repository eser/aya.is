-- +goose Up

-- Profile Envelopes: inbox items sent to profiles (invitations, messages, badges, passes).
CREATE TABLE IF NOT EXISTS "profile_envelope" (
  "id"                 CHAR(26) NOT NULL PRIMARY KEY,
  "target_profile_id"  CHAR(26) NOT NULL REFERENCES "profile" ("id"),
  "sender_profile_id"  CHAR(26) REFERENCES "profile" ("id"),
  "sender_user_id"     CHAR(26) REFERENCES "user" ("id"),
  "kind"               TEXT NOT NULL,
  "status"             TEXT NOT NULL DEFAULT 'pending',
  "title"              TEXT NOT NULL,
  "description"        TEXT,
  "properties"         JSONB,
  "accepted_at"        TIMESTAMP WITH TIME ZONE,
  "rejected_at"        TIMESTAMP WITH TIME ZONE,
  "revoked_at"         TIMESTAMP WITH TIME ZONE,
  "redeemed_at"        TIMESTAMP WITH TIME ZONE,
  "created_at"         TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"         TIMESTAMP WITH TIME ZONE,
  "deleted_at"         TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS "profile_envelope_target_profile_id_idx"
  ON "profile_envelope" ("target_profile_id", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

CREATE INDEX IF NOT EXISTS "profile_envelope_status_idx"
  ON "profile_envelope" ("status", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

-- +goose Down

DROP INDEX IF EXISTS "profile_envelope_status_idx";
DROP INDEX IF EXISTS "profile_envelope_target_profile_id_idx";
DROP TABLE IF EXISTS "profile_envelope";
