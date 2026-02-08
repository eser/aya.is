-- +goose Up
ALTER TABLE "profile" ADD COLUMN "hide_qa" BOOLEAN DEFAULT TRUE NOT NULL;

CREATE TABLE IF NOT EXISTS "profile_question" (
  "id"              CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id"      CHAR(26) NOT NULL CONSTRAINT "profile_question_profile_id_fk" REFERENCES "profile",
  "author_user_id"  CHAR(26) NOT NULL CONSTRAINT "profile_question_author_user_id_fk" REFERENCES "user",
  "content"         TEXT NOT NULL,
  "answer_content"  TEXT,
  "answered_at"     TIMESTAMP WITH TIME ZONE,
  "answered_by"     CHAR(26) CONSTRAINT "profile_question_answered_by_fk" REFERENCES "user",
  "is_anonymous"    BOOLEAN NOT NULL DEFAULT FALSE,
  "is_hidden"       BOOLEAN NOT NULL DEFAULT FALSE,
  "vote_count"      INTEGER NOT NULL DEFAULT 0,
  "created_at"      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "updated_at"      TIMESTAMP WITH TIME ZONE,
  "deleted_at"      TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS "profile_question_profile_id_idx"
  ON "profile_question" ("profile_id", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

CREATE INDEX IF NOT EXISTS "profile_question_author_user_id_idx"
  ON "profile_question" ("author_user_id", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

CREATE TABLE IF NOT EXISTS "profile_question_vote" (
  "id"          CHAR(26) NOT NULL PRIMARY KEY,
  "question_id" CHAR(26) NOT NULL CONSTRAINT "profile_question_vote_question_id_fk" REFERENCES "profile_question",
  "user_id"     CHAR(26) NOT NULL CONSTRAINT "profile_question_vote_user_id_fk" REFERENCES "user",
  "score"       INTEGER NOT NULL DEFAULT 1,
  "created_at"  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "profile_question_vote_question_user_idx"
  ON "profile_question_vote" ("question_id", "user_id");

-- +goose Down
DROP INDEX IF EXISTS "profile_question_vote_question_user_idx";
DROP TABLE IF EXISTS "profile_question_vote";
DROP INDEX IF EXISTS "profile_question_author_user_id_idx";
DROP INDEX IF EXISTS "profile_question_profile_id_idx";
DROP TABLE IF EXISTS "profile_question";
ALTER TABLE "profile" DROP COLUMN IF EXISTS "hide_qa";
