-- +goose Up
CREATE TABLE IF NOT EXISTS "event_queue" (
  "id"                      CHAR(26) NOT NULL PRIMARY KEY,
  "type"                    TEXT NOT NULL,
  "payload"                 JSONB NOT NULL DEFAULT '{}',
  "status"                  TEXT NOT NULL DEFAULT 'pending',
  "retry_count"             INTEGER NOT NULL DEFAULT 0,
  "max_retries"             INTEGER NOT NULL DEFAULT 3,
  "visible_at"              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "visibility_timeout_secs" INTEGER NOT NULL DEFAULT 300,
  "started_at"              TIMESTAMP WITH TIME ZONE,
  "completed_at"            TIMESTAMP WITH TIME ZONE,
  "failed_at"               TIMESTAMP WITH TIME ZONE,
  "created_at"              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "updated_at"              TIMESTAMP WITH TIME ZONE,
  "error_message"           TEXT,
  "worker_id"               TEXT
);

-- Polling index: pending events ready to run, OR stale processing events past visibility
CREATE INDEX IF NOT EXISTS "event_queue_poll_idx"
  ON "event_queue" ("status", "visible_at")
  WHERE "status" IN ('pending', 'processing');

-- Audit/history index
CREATE INDEX IF NOT EXISTS "event_queue_type_idx"
  ON "event_queue" ("type", "created_at");

-- +goose Down
DROP TABLE IF EXISTS "event_queue";
