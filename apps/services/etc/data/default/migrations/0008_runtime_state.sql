-- +goose Up
CREATE TABLE IF NOT EXISTS "runtime_state" (
  "key" TEXT NOT NULL PRIMARY KEY,
  "value" TEXT NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "runtime_state";
