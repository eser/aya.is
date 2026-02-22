-- +goose Up

-- 1. Backfill conversation titles from the first envelope's title.
UPDATE "mailbox_conversation" mc
SET "title" = sub.first_title
FROM (
  SELECT DISTINCT ON (me.conversation_id)
    me.conversation_id, me.title AS first_title
  FROM "mailbox_envelope" me
  WHERE me.deleted_at IS NULL AND me.title != ''
  ORDER BY me.conversation_id, me.created_at ASC
) sub
WHERE mc.id = sub.conversation_id
  AND (mc.title IS NULL OR mc.title = '');

-- 2. Rename description → message.
ALTER TABLE "mailbox_envelope" RENAME COLUMN "description" TO "message";

-- 3. Drop title from envelopes.
ALTER TABLE "mailbox_envelope" DROP COLUMN "title";

-- +goose Down
ALTER TABLE "mailbox_envelope" ADD COLUMN "title" TEXT NOT NULL DEFAULT '';
ALTER TABLE "mailbox_envelope" RENAME COLUMN "message" TO "description";
