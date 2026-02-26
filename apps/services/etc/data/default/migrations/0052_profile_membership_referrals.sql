-- +goose Up

-- Membership referral records
CREATE TABLE IF NOT EXISTS "profile_membership_referral" (
  "id"                     CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id"             CHAR(26) NOT NULL
    CONSTRAINT "profile_membership_referral_profile_id_fk" REFERENCES "profile" ("id"),
  "referred_profile_id"    CHAR(26) NOT NULL
    CONSTRAINT "profile_membership_referral_referred_profile_id_fk" REFERENCES "profile" ("id"),
  "referrer_membership_id" CHAR(26) NOT NULL
    CONSTRAINT "profile_membership_referral_referrer_membership_id_fk" REFERENCES "profile_membership" ("id"),
  "status"                 TEXT NOT NULL DEFAULT 'voting',
  "vote_count"             INTEGER NOT NULL DEFAULT 0,
  "created_at"             TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"             TIMESTAMP WITH TIME ZONE,
  "deleted_at"             TIMESTAMP WITH TIME ZONE
);

-- One active referral per referred profile per organization
CREATE UNIQUE INDEX IF NOT EXISTS "profile_membership_referral_profile_referred_uniq"
  ON "profile_membership_referral" ("profile_id", "referred_profile_id")
  WHERE "deleted_at" IS NULL;

CREATE INDEX IF NOT EXISTS "profile_membership_referral_profile_id_idx"
  ON "profile_membership_referral" ("profile_id", "created_at" DESC)
  WHERE "deleted_at" IS NULL;

-- Teams suggested by referrer for the referred person
CREATE TABLE IF NOT EXISTS "profile_membership_referral_team" (
  "id"                              CHAR(26) NOT NULL PRIMARY KEY,
  "profile_membership_referral_id"  CHAR(26) NOT NULL
    CONSTRAINT "profile_membership_referral_team_referral_id_fk" REFERENCES "profile_membership_referral" ("id"),
  "profile_team_id"                 CHAR(26) NOT NULL
    CONSTRAINT "profile_membership_referral_team_team_id_fk" REFERENCES "profile_team" ("id"),
  "created_at"                      TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "deleted_at"                      TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX IF NOT EXISTS "profile_membership_referral_team_uniq"
  ON "profile_membership_referral_team" ("profile_membership_referral_id", "profile_team_id")
  WHERE "deleted_at" IS NULL;

-- Votes on referrals (5-level scale with optional comment)
CREATE TABLE IF NOT EXISTS "profile_membership_referral_vote" (
  "id"                              CHAR(26) NOT NULL PRIMARY KEY,
  "profile_membership_referral_id"  CHAR(26) NOT NULL
    CONSTRAINT "profile_membership_referral_vote_referral_id_fk" REFERENCES "profile_membership_referral" ("id"),
  "voter_membership_id"             CHAR(26) NOT NULL
    CONSTRAINT "profile_membership_referral_vote_voter_membership_id_fk" REFERENCES "profile_membership" ("id"),
  "score"                           SMALLINT NOT NULL,
  "comment"                         TEXT,
  "created_at"                      TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at"                      TIMESTAMP WITH TIME ZONE
);

-- One vote per voter per referral
CREATE UNIQUE INDEX IF NOT EXISTS "profile_membership_referral_vote_referral_voter_uniq"
  ON "profile_membership_referral_vote" ("profile_membership_referral_id", "voter_membership_id");

-- +goose Down
DROP INDEX IF EXISTS "profile_membership_referral_vote_referral_voter_uniq";
DROP TABLE IF EXISTS "profile_membership_referral_vote";
DROP INDEX IF EXISTS "profile_membership_referral_team_uniq";
DROP TABLE IF EXISTS "profile_membership_referral_team";
DROP INDEX IF EXISTS "profile_membership_referral_profile_id_idx";
DROP INDEX IF EXISTS "profile_membership_referral_profile_referred_uniq";
DROP TABLE IF EXISTS "profile_membership_referral";
