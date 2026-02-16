-- +goose Up
-- Change answered_by from referencing user(id) to referencing profile(id).
-- Migrate existing user IDs to their individual profile IDs.

-- First, update existing answered_by values from user IDs to profile IDs
UPDATE "profile_question" pq
SET answered_by = u.individual_profile_id
FROM "user" u
WHERE pq.answered_by = u.id
  AND u.individual_profile_id IS NOT NULL;

-- Drop the old FK referencing user
ALTER TABLE "profile_question" DROP CONSTRAINT IF EXISTS "profile_question_answered_by_fk";

-- Add the new FK referencing profile
ALTER TABLE "profile_question"
  ADD CONSTRAINT "profile_question_answered_by_fk" FOREIGN KEY ("answered_by") REFERENCES "profile" ("id");

-- +goose Down
-- Revert: change answered_by back from profile(id) to user(id)

-- Update profile IDs back to user IDs
UPDATE "profile_question" pq
SET answered_by = u.id
FROM "user" u
WHERE pq.answered_by = u.individual_profile_id
  AND u.individual_profile_id IS NOT NULL;

-- Drop the profile FK
ALTER TABLE "profile_question" DROP CONSTRAINT IF EXISTS "profile_question_answered_by_fk";

-- Re-add the user FK
ALTER TABLE "profile_question"
  ADD CONSTRAINT "profile_question_answered_by_fk" FOREIGN KEY ("answered_by") REFERENCES "user" ("id");
