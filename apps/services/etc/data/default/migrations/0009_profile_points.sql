-- +goose Up
ALTER TABLE "profile" ADD COLUMN "points" INTEGER NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS "profile_point_transaction" (
  "id"                CHAR(26) NOT NULL PRIMARY KEY,
  "target_profile_id" CHAR(26) NOT NULL CONSTRAINT "profile_point_transaction_target_profile_id_fk" REFERENCES "profile",
  "origin_profile_id" CHAR(26) CONSTRAINT "profile_point_transaction_origin_profile_id_fk" REFERENCES "profile",
  "transaction_type"  TEXT NOT NULL,
  "triggering_event"  TEXT,
  "description"       TEXT NOT NULL,
  "amount"            INTEGER NOT NULL,
  "balance_after"     INTEGER NOT NULL,
  "created_at"        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "profile_point_transaction_target_profile_id_idx"
  ON "profile_point_transaction" ("target_profile_id", "created_at" DESC);

CREATE INDEX IF NOT EXISTS "profile_point_transaction_origin_profile_id_idx"
  ON "profile_point_transaction" ("origin_profile_id", "created_at" DESC)
  WHERE "origin_profile_id" IS NOT NULL;

ALTER TABLE "profile_point_transaction"
  ADD CONSTRAINT "profile_point_transaction_origin_check"
  CHECK (
    (transaction_type = 'TRANSFER' AND origin_profile_id IS NOT NULL) OR
    (transaction_type IN ('GAIN', 'SPEND') AND origin_profile_id IS NULL)
  );

ALTER TABLE "profile_point_transaction"
  ADD CONSTRAINT "profile_point_transaction_amount_positive"
  CHECK (amount > 0);

-- +goose Down
ALTER TABLE "profile_point_transaction" DROP CONSTRAINT IF EXISTS "profile_point_transaction_amount_positive";
ALTER TABLE "profile_point_transaction" DROP CONSTRAINT IF EXISTS "profile_point_transaction_origin_check";
DROP INDEX IF EXISTS "profile_point_transaction_origin_profile_id_idx";
DROP INDEX IF EXISTS "profile_point_transaction_target_profile_id_idx";
DROP TABLE IF EXISTS "profile_point_transaction";
ALTER TABLE "profile" DROP COLUMN IF EXISTS "points";
