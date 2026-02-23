-- +goose Up

ALTER TABLE "story"
ADD COLUMN "feat_discussions" BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE "profile"
ADD COLUMN "option_story_discussions_by_default" BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down

ALTER TABLE "profile" DROP COLUMN IF EXISTS "option_story_discussions_by_default";
ALTER TABLE "story" DROP COLUMN IF EXISTS "feat_discussions";
