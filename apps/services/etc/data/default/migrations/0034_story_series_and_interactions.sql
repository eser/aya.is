-- +goose Up

-- story_series: groups any stories (article series, activity series, etc.)
CREATE TABLE IF NOT EXISTS "story_series" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "slug" TEXT NOT NULL CONSTRAINT "story_series_slug_unique" UNIQUE,
  "series_picture_uri" TEXT,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE
);

-- Add series_id to story table
ALTER TABLE "story"
  ADD COLUMN "series_id" CHAR(26)
    CONSTRAINT "story_series_id_fk" REFERENCES "story_series";

-- story_interaction: profile interactions with stories (RSVP, likes, bookmarks, etc.)
CREATE TABLE IF NOT EXISTS "story_interaction" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "story_id" CHAR(26) NOT NULL
    CONSTRAINT "story_interaction_story_id_fk" REFERENCES "story",
  "profile_id" CHAR(26) NOT NULL
    CONSTRAINT "story_interaction_profile_id_fk" REFERENCES "profile",
  "kind" TEXT NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE,
  "deleted_at" TIMESTAMP WITH TIME ZONE,
  CONSTRAINT "story_interaction_story_profile_kind_unique"
    UNIQUE ("story_id", "profile_id", "kind")
);

-- Index for activity listings sorted by time_start from JSONB properties
-- Uses text ordering (ISO 8601 strings sort correctly as text)
CREATE INDEX "story_activity_time_start_idx"
  ON "story" ((properties->>'activity_time_start') DESC NULLS LAST)
  WHERE kind = 'activity' AND deleted_at IS NULL;

-- Indexes for interaction lookups
CREATE INDEX "story_interaction_story_id_idx"
  ON "story_interaction" ("story_id")
  WHERE "deleted_at" IS NULL;

CREATE INDEX "story_interaction_profile_id_idx"
  ON "story_interaction" ("profile_id")
  WHERE "deleted_at" IS NULL;

-- +goose Down

DROP INDEX IF EXISTS "story_interaction_profile_id_idx";
DROP INDEX IF EXISTS "story_interaction_story_id_idx";
DROP INDEX IF EXISTS "story_activity_time_start_idx";
DROP TABLE IF EXISTS "story_interaction";
ALTER TABLE "story" DROP COLUMN IF EXISTS "series_id";
DROP TABLE IF EXISTS "story_series";
