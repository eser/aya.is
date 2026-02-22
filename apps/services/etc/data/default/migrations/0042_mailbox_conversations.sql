-- +goose Up

-- 1. Create conversation container table.
CREATE TABLE IF NOT EXISTS "mailbox_conversation" (
  "id"                    CHAR(26) NOT NULL PRIMARY KEY,
  "kind"                  TEXT NOT NULL DEFAULT 'direct',
  "title"                 TEXT,
  "created_by_profile_id" CHAR(26) REFERENCES "profile" ("id"),
  "created_at"            TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"            TIMESTAMP WITH TIME ZONE
);

-- 2. Create participant table (links profiles to conversations).
CREATE TABLE IF NOT EXISTS "mailbox_participant" (
  "id"              CHAR(26) NOT NULL PRIMARY KEY,
  "conversation_id" CHAR(26) NOT NULL REFERENCES "mailbox_conversation" ("id"),
  "profile_id"      CHAR(26) NOT NULL REFERENCES "profile" ("id"),
  "last_read_at"    TIMESTAMP WITH TIME ZONE,
  "is_archived"     BOOLEAN NOT NULL DEFAULT FALSE,
  "joined_at"       TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "left_at"         TIMESTAMP WITH TIME ZONE,
  UNIQUE ("conversation_id", "profile_id")
);

-- 3. Rename envelope table.
ALTER TABLE "profile_envelope" RENAME TO "mailbox_envelope";

-- 4. Add new columns to envelope table.
ALTER TABLE "mailbox_envelope"
  ADD COLUMN "conversation_id" CHAR(26) REFERENCES "mailbox_conversation" ("id"),
  ADD COLUMN "reply_to_id"     CHAR(26) REFERENCES "mailbox_envelope" ("id");

-- 5. Create reaction table (must come after rename so FK references mailbox_envelope).
CREATE TABLE IF NOT EXISTS "mailbox_reaction" (
  "id"          CHAR(26) NOT NULL PRIMARY KEY,
  "envelope_id" CHAR(26) NOT NULL REFERENCES "mailbox_envelope" ("id"),
  "profile_id"  CHAR(26) NOT NULL REFERENCES "profile" ("id"),
  "emoji"       TEXT NOT NULL,
  "created_at"  TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  UNIQUE ("envelope_id", "profile_id", "emoji")
);

-- 6. Backfill: wrap each existing envelope in a system conversation.
DO $$
DECLARE
  rec RECORD;
  conv_id CHAR(26);
  participant_id_1 CHAR(26);
  participant_id_2 CHAR(26);
BEGIN
  FOR rec IN
    SELECT id, target_profile_id, sender_profile_id, created_at
    FROM "mailbox_envelope"
    WHERE conversation_id IS NULL
  LOOP
    -- Generate IDs (26-char from uuid)
    conv_id := substring(replace(gen_random_uuid()::text, '-', ''), 1, 26);
    participant_id_1 := substring(replace(gen_random_uuid()::text, '-', ''), 1, 26);

    -- Create system conversation
    INSERT INTO "mailbox_conversation" (id, kind, created_by_profile_id, created_at)
    VALUES (conv_id, 'system', rec.sender_profile_id, rec.created_at);

    -- Add target as participant
    INSERT INTO "mailbox_participant" (id, conversation_id, profile_id, last_read_at, joined_at)
    VALUES (participant_id_1, conv_id, rec.target_profile_id, rec.created_at, rec.created_at);

    -- Add sender as participant (if present)
    IF rec.sender_profile_id IS NOT NULL THEN
      participant_id_2 := substring(replace(gen_random_uuid()::text, '-', ''), 1, 26);

      INSERT INTO "mailbox_participant" (id, conversation_id, profile_id, last_read_at, joined_at)
      VALUES (participant_id_2, conv_id, rec.sender_profile_id, rec.created_at, rec.created_at);
    END IF;

    -- Link envelope to conversation
    UPDATE "mailbox_envelope" SET conversation_id = conv_id WHERE id = rec.id;
  END LOOP;
END $$;

-- 7. Make conversation_id NOT NULL after backfill.
ALTER TABLE "mailbox_envelope" ALTER COLUMN "conversation_id" SET NOT NULL;

-- 8. Rename existing indexes.
ALTER INDEX IF EXISTS "profile_envelope_target_profile_id_idx"
  RENAME TO "mailbox_envelope_target_profile_id_idx";
ALTER INDEX IF EXISTS "profile_envelope_status_idx"
  RENAME TO "mailbox_envelope_status_idx";

-- 9. Create new indexes.
CREATE INDEX IF NOT EXISTS "mailbox_participant_profile_id_idx"
  ON "mailbox_participant" ("profile_id")
  WHERE "left_at" IS NULL;

CREATE INDEX IF NOT EXISTS "mailbox_participant_conversation_idx"
  ON "mailbox_participant" ("conversation_id", "profile_id")
  WHERE "left_at" IS NULL;

CREATE INDEX IF NOT EXISTS "mailbox_envelope_conversation_id_idx"
  ON "mailbox_envelope" ("conversation_id", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

CREATE INDEX IF NOT EXISTS "mailbox_reaction_envelope_id_idx"
  ON "mailbox_reaction" ("envelope_id");

-- +goose Down

-- Drop new indexes
DROP INDEX IF EXISTS "mailbox_reaction_envelope_id_idx";
DROP INDEX IF EXISTS "mailbox_envelope_conversation_id_idx";
DROP INDEX IF EXISTS "mailbox_participant_conversation_idx";
DROP INDEX IF EXISTS "mailbox_participant_profile_id_idx";

-- Rename indexes back
ALTER INDEX IF EXISTS "mailbox_envelope_target_profile_id_idx"
  RENAME TO "profile_envelope_target_profile_id_idx";
ALTER INDEX IF EXISTS "mailbox_envelope_status_idx"
  RENAME TO "profile_envelope_status_idx";

-- Drop new columns from envelope
ALTER TABLE "mailbox_envelope" DROP COLUMN IF EXISTS "reply_to_id";
ALTER TABLE "mailbox_envelope" DROP COLUMN IF EXISTS "conversation_id";

-- Rename table back
ALTER TABLE "mailbox_envelope" RENAME TO "profile_envelope";

-- Drop new tables
DROP TABLE IF EXISTS "mailbox_reaction";
DROP TABLE IF EXISTS "mailbox_participant";
DROP TABLE IF EXISTS "mailbox_conversation";
