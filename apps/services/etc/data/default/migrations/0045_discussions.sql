-- +goose Up

-- 1. Feature toggle on profiles (matches feature_qa, feature_links, feature_relations).
ALTER TABLE "profile" ADD COLUMN "feature_discussions" TEXT NOT NULL DEFAULT 'disabled';

-- 2. Discussion thread: anchor table, lazily created on first comment.
-- Exactly one of story_id / profile_id must be set.
CREATE TABLE IF NOT EXISTS "discussion_thread" (
  "id"            CHAR(26) NOT NULL PRIMARY KEY,
  "story_id"      CHAR(26) CONSTRAINT "discussion_thread_story_id_fk" REFERENCES "story" ("id"),
  "profile_id"    CHAR(26) CONSTRAINT "discussion_thread_profile_id_fk" REFERENCES "profile" ("id"),
  "is_locked"     BOOLEAN NOT NULL DEFAULT FALSE,
  "comment_count" INTEGER NOT NULL DEFAULT 0,
  "created_at"    TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"    TIMESTAMP WITH TIME ZONE,
  "deleted_at"    TIMESTAMP WITH TIME ZONE,

  CONSTRAINT "discussion_thread_anchor_check"
    CHECK (
      ("story_id" IS NOT NULL AND "profile_id" IS NULL) OR
      ("story_id" IS NULL AND "profile_id" IS NOT NULL)
    )
);

-- One thread per story, one thread per profile.
CREATE UNIQUE INDEX IF NOT EXISTS "discussion_thread_story_id_uniq"
  ON "discussion_thread" ("story_id")
  WHERE "story_id" IS NOT NULL AND "deleted_at" IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS "discussion_thread_profile_id_uniq"
  ON "discussion_thread" ("profile_id")
  WHERE "profile_id" IS NOT NULL AND "deleted_at" IS NULL;

-- 3. Discussion comment: tree-structured via parent_id self-reference.
CREATE TABLE IF NOT EXISTS "discussion_comment" (
  "id"              CHAR(26) NOT NULL PRIMARY KEY,
  "thread_id"       CHAR(26) NOT NULL CONSTRAINT "discussion_comment_thread_id_fk" REFERENCES "discussion_thread" ("id"),
  "parent_id"       CHAR(26) CONSTRAINT "discussion_comment_parent_id_fk" REFERENCES "discussion_comment" ("id"),
  "author_user_id"  CHAR(26) NOT NULL CONSTRAINT "discussion_comment_author_user_id_fk" REFERENCES "user" ("id"),
  "content"         TEXT NOT NULL,
  "depth"           INTEGER NOT NULL DEFAULT 0,
  "vote_score"      INTEGER NOT NULL DEFAULT 0,
  "upvote_count"    INTEGER NOT NULL DEFAULT 0,
  "downvote_count"  INTEGER NOT NULL DEFAULT 0,
  "reply_count"     INTEGER NOT NULL DEFAULT 0,
  "is_pinned"       BOOLEAN NOT NULL DEFAULT FALSE,
  "is_hidden"       BOOLEAN NOT NULL DEFAULT FALSE,
  "is_edited"       BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at"      TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"      TIMESTAMP WITH TIME ZONE,
  "deleted_at"      TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS "discussion_comment_thread_id_idx"
  ON "discussion_comment" ("thread_id", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

CREATE INDEX IF NOT EXISTS "discussion_comment_parent_id_idx"
  ON "discussion_comment" ("parent_id", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

CREATE INDEX IF NOT EXISTS "discussion_comment_author_idx"
  ON "discussion_comment" ("author_user_id")
  WHERE "deleted_at" IS NULL;

-- Hot sort index: pinned first, then by vote_score for root comments.
CREATE INDEX IF NOT EXISTS "discussion_comment_hot_idx"
  ON "discussion_comment" ("thread_id", "is_pinned" DESC, "vote_score" DESC, "created_at" DESC)
  WHERE "deleted_at" IS NULL AND "parent_id" IS NULL;

-- 4. Discussion votes: supports upvote (+1) and downvote (-1).
CREATE TABLE IF NOT EXISTS "discussion_comment_vote" (
  "id"         CHAR(26) NOT NULL PRIMARY KEY,
  "comment_id" CHAR(26) NOT NULL CONSTRAINT "discussion_vote_comment_id_fk" REFERENCES "discussion_comment" ("id"),
  "user_id"    CHAR(26) NOT NULL CONSTRAINT "discussion_vote_user_id_fk" REFERENCES "user" ("id"),
  "direction"  SMALLINT NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS "discussion_vote_comment_user_uniq"
  ON "discussion_comment_vote" ("comment_id", "user_id");

-- +goose Down

DROP INDEX IF EXISTS "discussion_vote_comment_user_uniq";
DROP TABLE IF EXISTS "discussion_comment_vote";
DROP INDEX IF EXISTS "discussion_comment_hot_idx";
DROP INDEX IF EXISTS "discussion_comment_author_idx";
DROP INDEX IF EXISTS "discussion_comment_parent_id_idx";
DROP INDEX IF EXISTS "discussion_comment_thread_id_idx";
DROP TABLE IF EXISTS "discussion_comment";
DROP INDEX IF EXISTS "discussion_thread_profile_id_uniq";
DROP INDEX IF EXISTS "discussion_thread_story_id_uniq";
DROP TABLE IF EXISTS "discussion_thread";
ALTER TABLE "profile" DROP COLUMN IF EXISTS "feature_discussions";
