-- +goose Up
-- Add icon column to profile_link_tx for custom emoticon or initials
ALTER TABLE "profile_link_tx"
ADD COLUMN "icon" TEXT;

-- +goose Down
ALTER TABLE "profile_link_tx" DROP COLUMN IF EXISTS "icon";
