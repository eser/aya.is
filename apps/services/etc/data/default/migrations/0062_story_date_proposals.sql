-- +goose Up

-- Date proposals: participants suggest date/time options for undecided activities.
CREATE TABLE IF NOT EXISTS "story_date_proposal" (
  "id"                  CHAR(26) NOT NULL PRIMARY KEY,
  "story_id"            CHAR(26) NOT NULL
    CONSTRAINT "story_date_proposal_story_id_fk" REFERENCES "story" ("id"),
  "proposer_profile_id" CHAR(26) NOT NULL
    CONSTRAINT "story_date_proposal_proposer_profile_id_fk" REFERENCES "profile" ("id"),
  "datetime_start"      TIMESTAMP WITH TIME ZONE NOT NULL,
  "datetime_end"        TIMESTAMP WITH TIME ZONE,
  "is_finalized"        BOOLEAN NOT NULL DEFAULT FALSE,
  "vote_score"          INTEGER NOT NULL DEFAULT 0,
  "upvote_count"        INTEGER NOT NULL DEFAULT 0,
  "downvote_count"      INTEGER NOT NULL DEFAULT 0,
  "created_at"          TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"          TIMESTAMP WITH TIME ZONE,
  "deleted_at"          TIMESTAMP WITH TIME ZONE
);

CREATE INDEX "story_date_proposal_story_id_idx"
  ON "story_date_proposal" ("story_id", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

-- Votes on date proposals: agree (+1) or disagree (-1).
CREATE TABLE IF NOT EXISTS "story_date_proposal_vote" (
  "id"               CHAR(26) NOT NULL PRIMARY KEY,
  "proposal_id"      CHAR(26) NOT NULL
    CONSTRAINT "story_date_proposal_vote_proposal_id_fk" REFERENCES "story_date_proposal" ("id") ON DELETE CASCADE,
  "voter_profile_id" CHAR(26) NOT NULL
    CONSTRAINT "story_date_proposal_vote_voter_profile_id_fk" REFERENCES "profile" ("id"),
  "direction"        SMALLINT NOT NULL,
  "created_at"       TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"       TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX "story_date_proposal_vote_proposal_voter_uniq"
  ON "story_date_proposal_vote" ("proposal_id", "voter_profile_id");

-- +goose Down

DROP INDEX IF EXISTS "story_date_proposal_vote_proposal_voter_uniq";
DROP TABLE IF EXISTS "story_date_proposal_vote";
DROP INDEX IF EXISTS "story_date_proposal_story_id_idx";
DROP TABLE IF EXISTS "story_date_proposal";
