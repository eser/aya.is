-- +goose Up

-- Session preferences table (extends existing session)
CREATE TABLE IF NOT EXISTS "session_preference" (
  "session_id" CHAR(26) NOT NULL CONSTRAINT "session_preference_session_id_fk" REFERENCES "session" ON DELETE CASCADE,
  "key" VARCHAR(50) NOT NULL,
  "value" TEXT NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  PRIMARY KEY ("session_id", "key")
);

-- PoW challenges for protection (short-lived, unlogged for performance)
CREATE UNLOGGED TABLE IF NOT EXISTS "protection_pow_challenge" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "prefix" CHAR(64) NOT NULL,
  "difficulty" SMALLINT NOT NULL,
  "ip_hash" CHAR(64) NOT NULL,
  "used" BOOLEAN DEFAULT FALSE NOT NULL,
  "expires_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Index for cleanup job (expired challenges)
CREATE INDEX IF NOT EXISTS "protection_pow_challenge_expires_at_idx" ON "protection_pow_challenge" ("expires_at");

-- Index for IP-based lookups
CREATE INDEX IF NOT EXISTS "protection_pow_challenge_ip_hash_idx" ON "protection_pow_challenge" ("ip_hash");

-- Rate limiting for session creation (unlogged for performance)
CREATE UNLOGGED TABLE IF NOT EXISTS "session_rate_limit" (
  "ip_hash" CHAR(64) NOT NULL PRIMARY KEY,
  "count" INTEGER DEFAULT 1 NOT NULL,
  "window_start" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "session_rate_limit";

DROP INDEX IF EXISTS "protection_pow_challenge_ip_hash_idx";

DROP INDEX IF EXISTS "protection_pow_challenge_expires_at_idx";

DROP TABLE IF EXISTS "protection_pow_challenge";

DROP TABLE IF EXISTS "session_preference";
