-- +goose Up

-- Bulletin subscription: one per profile per channel
CREATE TABLE IF NOT EXISTS "bulletin_subscription" (
  "id"               CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id"       CHAR(26) NOT NULL
    CONSTRAINT "bulletin_subscription_profile_id_fk" REFERENCES "profile" ("id") ON DELETE CASCADE,
  "channel"          TEXT NOT NULL,
  "preferred_time"   SMALLINT NOT NULL DEFAULT 8
    CHECK (preferred_time >= 0 AND preferred_time <= 23),
  "last_bulletin_at" TIMESTAMP WITH TIME ZONE,
  "created_at"       TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"       TIMESTAMP WITH TIME ZONE,
  "deleted_at"       TIMESTAMP WITH TIME ZONE
);

-- One active subscription per profile per channel
CREATE UNIQUE INDEX IF NOT EXISTS "bulletin_subscription_profile_channel_uniq"
  ON "bulletin_subscription" ("profile_id", "channel")
  WHERE "deleted_at" IS NULL;

CREATE INDEX IF NOT EXISTS "bulletin_subscription_active_time_idx"
  ON "bulletin_subscription" ("preferred_time", "last_bulletin_at")
  WHERE "deleted_at" IS NULL;

-- Bulletin log: one row per sent digest attempt
CREATE TABLE IF NOT EXISTS "bulletin_log" (
  "id"               CHAR(26) NOT NULL PRIMARY KEY,
  "subscription_id"  CHAR(26) NOT NULL
    CONSTRAINT "bulletin_log_subscription_id_fk" REFERENCES "bulletin_subscription" ("id"),
  "story_count"      INT NOT NULL,
  "status"           TEXT NOT NULL,
  "error_message"    TEXT,
  "created_at"       TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS "bulletin_log_subscription_sent_idx"
  ON "bulletin_log" ("subscription_id", "created_at" DESC);

-- +goose Down

DROP TABLE IF EXISTS "bulletin_log";
DROP TABLE IF EXISTS "bulletin_subscription";
