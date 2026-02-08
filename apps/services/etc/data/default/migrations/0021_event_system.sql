-- +goose Up

-- Rename queue → event_queue (reverse of migration 0018)
ALTER TABLE "queue" RENAME TO "event_queue";
ALTER INDEX "queue_poll_idx" RENAME TO "event_queue_poll_idx";
ALTER INDEX "queue_type_idx" RENAME TO "event_queue_type_idx";

-- Create event_audit table
CREATE TABLE IF NOT EXISTS "event_audit" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "event_type" TEXT NOT NULL,
  "entity_type" TEXT NOT NULL,
  "entity_id" TEXT,
  "actor_id" TEXT,
  "actor_kind" TEXT NOT NULL DEFAULT 'user',
  "session_id" TEXT,
  "payload" JSONB,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE INDEX "event_audit_type_idx" ON "event_audit" ("event_type", "created_at");
CREATE INDEX "event_audit_entity_idx" ON "event_audit" ("entity_type", "entity_id", "created_at");
CREATE INDEX "event_audit_actor_idx" ON "event_audit" ("actor_id", "created_at") WHERE "actor_id" IS NOT NULL;

-- +goose Down

DROP INDEX IF EXISTS "event_audit_actor_idx";
DROP INDEX IF EXISTS "event_audit_entity_idx";
DROP INDEX IF EXISTS "event_audit_type_idx";
DROP TABLE IF EXISTS "event_audit";

ALTER TABLE "event_queue" RENAME TO "queue";
ALTER INDEX "event_queue_poll_idx" RENAME TO "queue_poll_idx";
ALTER INDEX "event_queue_type_idx" RENAME TO "queue_type_idx";
