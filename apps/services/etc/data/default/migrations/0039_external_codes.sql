-- +goose Up

-- Merge telegram_verification_code and telegram_group_invite_code into a single
-- platform-agnostic external_code table. All platform-specific data lives in
-- the JSONB properties column.
DROP TABLE IF EXISTS "telegram_group_invite_code";
DROP TABLE IF EXISTS "telegram_verification_code";

CREATE TABLE IF NOT EXISTS "external_code" (
  "id"              CHAR(26) NOT NULL PRIMARY KEY,
  "code"            VARCHAR(10) NOT NULL CONSTRAINT "external_code_code_unique" UNIQUE,
  "external_system" TEXT NOT NULL,
  "properties"      JSONB NOT NULL DEFAULT '{}',
  "created_at"      TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "expires_at"      TIMESTAMP WITH TIME ZONE NOT NULL,
  "consumed_at"     TIMESTAMP WITH TIME ZONE
);

-- Index for code lookup (primary query path — only unconsumed codes)
CREATE INDEX "external_code_code_idx"
  ON "external_code" ("code")
  WHERE "consumed_at" IS NULL;

-- +goose Down

DROP INDEX IF EXISTS "external_code_code_idx";
DROP TABLE IF EXISTS "external_code";

-- Restore old tables
CREATE TABLE IF NOT EXISTS "telegram_verification_code" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "code" VARCHAR(10) NOT NULL CONSTRAINT "telegram_verification_code_code_unique" UNIQUE,
  "telegram_user_id" BIGINT NOT NULL,
  "telegram_username" TEXT NOT NULL DEFAULT '',
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "expires_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  "consumed_at" TIMESTAMP WITH TIME ZONE
);

CREATE INDEX "telegram_verification_code_code_idx"
  ON "telegram_verification_code" ("code")
  WHERE "consumed_at" IS NULL;

CREATE TABLE IF NOT EXISTS "telegram_group_invite_code" (
  "id"                          CHAR(26) NOT NULL PRIMARY KEY,
  "code"                        VARCHAR(10) NOT NULL UNIQUE,
  "telegram_chat_id"            BIGINT NOT NULL,
  "telegram_chat_title"         TEXT NOT NULL DEFAULT '',
  "created_by_telegram_user_id" BIGINT NOT NULL,
  "created_at"                  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "expires_at"                  TIMESTAMP WITH TIME ZONE NOT NULL,
  "consumed_at"                 TIMESTAMP WITH TIME ZONE
);

CREATE INDEX "telegram_group_invite_code_code_idx"
  ON "telegram_group_invite_code" ("code")
  WHERE "consumed_at" IS NULL;
