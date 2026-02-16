-- +goose Up
-- Add missing answer_uri and answer_kind columns to profile_question table.
-- These were defined in the original CREATE TABLE (0022) but may not exist
-- if the table was created before that migration ran.

ALTER TABLE "profile_question" ADD COLUMN IF NOT EXISTS "answer_uri" TEXT;
ALTER TABLE "profile_question" ADD COLUMN IF NOT EXISTS "answer_kind" TEXT;

-- +goose Down
ALTER TABLE "profile_question" DROP COLUMN IF EXISTS "answer_kind";
ALTER TABLE "profile_question" DROP COLUMN IF EXISTS "answer_uri";
