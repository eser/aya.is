-- +goose Up

-- Replace token-based flow with verification-code flow.
-- Old flow: web generates token → user opens deep link → bot validates.
-- New flow: user sends /start → bot generates code → user pastes code in web UI.
DROP TABLE IF EXISTS "telegram_link_token";

CREATE TABLE IF NOT EXISTS "telegram_verification_code" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "code" VARCHAR(10) NOT NULL CONSTRAINT "telegram_verification_code_code_unique" UNIQUE,
  "telegram_user_id" BIGINT NOT NULL,
  "telegram_username" TEXT NOT NULL DEFAULT '',
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "expires_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  "consumed_at" TIMESTAMP WITH TIME ZONE
);

-- Index for code lookup (primary query path — only unconsumed codes)
CREATE INDEX "telegram_verification_code_code_idx"
  ON "telegram_verification_code" ("code")
  WHERE "consumed_at" IS NULL;

-- +goose Down

DROP INDEX IF EXISTS "telegram_verification_code_code_idx";
DROP TABLE IF EXISTS "telegram_verification_code";

-- Re-create the old table
CREATE TABLE IF NOT EXISTS "telegram_link_token" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "token" TEXT NOT NULL CONSTRAINT "telegram_link_token_token_unique" UNIQUE,
  "profile_id" CHAR(26) NOT NULL
    CONSTRAINT "telegram_link_token_profile_id_fk" REFERENCES "profile",
  "profile_slug" TEXT NOT NULL,
  "created_by_user_id" TEXT NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "expires_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  "consumed_at" TIMESTAMP WITH TIME ZONE
);

CREATE INDEX "telegram_link_token_token_idx"
  ON "telegram_link_token" ("token")
  WHERE "consumed_at" IS NULL;
