-- +goose Up
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

-- +goose Down
DROP INDEX IF EXISTS "telegram_group_invite_code_code_idx";
DROP TABLE IF EXISTS "telegram_group_invite_code";
