-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "profile_point_pending_award" (
  "id"                CHAR(26) NOT NULL PRIMARY KEY,
  "target_profile_id" CHAR(26) NOT NULL CONSTRAINT "profile_point_pending_award_target_profile_id_fk" REFERENCES "profile",
  "triggering_event"  TEXT NOT NULL,
  "description"       TEXT NOT NULL,
  "amount"            INTEGER NOT NULL,
  "status"            TEXT NOT NULL DEFAULT 'pending',
  "reviewed_by"       CHAR(26) CONSTRAINT "profile_point_pending_award_reviewed_by_fk" REFERENCES "user",
  "reviewed_at"       TIMESTAMP WITH TIME ZONE,
  "rejection_reason"  TEXT,
  "metadata"          JSONB,
  "created_at"        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX "profile_point_pending_award_status_idx"
  ON "profile_point_pending_award" ("status", "created_at" DESC);

CREATE INDEX "profile_point_pending_award_profile_idx"
  ON "profile_point_pending_award" ("target_profile_id", "created_at" DESC);

ALTER TABLE "profile_point_pending_award"
  ADD CONSTRAINT "profile_point_pending_award_status_check"
  CHECK (status IN ('pending', 'approved', 'rejected'));

ALTER TABLE "profile_point_pending_award"
  ADD CONSTRAINT "profile_point_pending_award_amount_positive"
  CHECK (amount > 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "profile_point_pending_award";
-- +goose StatementEnd
