-- +goose Up

-- 1. Reclassify conversations with a real sender from system → direct.
-- "system" should only be for automated platform messages with no sender.
UPDATE "mailbox_conversation"
SET "kind" = 'direct'
WHERE "kind" = 'system'
  AND "created_by_profile_id" IS NOT NULL;

-- 2. Fix participant ordering: the conversation creator gets an earlier joined_at
-- so ORDER BY joined_at deterministically returns the creator first.
UPDATE "mailbox_participant" mp
SET "joined_at" = mc."created_at" - INTERVAL '1 second'
FROM "mailbox_conversation" mc
WHERE mp.conversation_id = mc.id
  AND mp.profile_id = mc.created_by_profile_id;

-- +goose Down

UPDATE "mailbox_conversation"
SET "kind" = 'system'
WHERE "kind" = 'direct'
  AND "created_by_profile_id" IS NOT NULL;
