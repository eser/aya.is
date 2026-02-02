-- +goose Up
-- Add boolean fields to control sidebar visibility

ALTER TABLE "profile"
ADD COLUMN "hide_relations" BOOLEAN DEFAULT FALSE NOT NULL;

ALTER TABLE "profile"
ADD COLUMN "hide_links" BOOLEAN DEFAULT FALSE NOT NULL;

-- +goose Down
ALTER TABLE "profile" DROP COLUMN IF EXISTS "hide_links";
ALTER TABLE "profile" DROP COLUMN IF EXISTS "hide_relations";
