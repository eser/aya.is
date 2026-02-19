-- +goose Up

-- Temporary tokens for the Telegram account linking flow.
-- Generated when a user clicks "Connect Telegram" in the web UI.
-- Consumed when the user sends /start <token> to the bot.
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

-- Index for token lookup (primary query path — only unconsumed tokens)
CREATE INDEX "telegram_link_token_token_idx"
  ON "telegram_link_token" ("token")
  WHERE "consumed_at" IS NULL;

-- +goose Down

DROP INDEX IF EXISTS "telegram_link_token_token_idx";
DROP TABLE IF EXISTS "telegram_link_token";
