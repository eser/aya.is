-- +goose Up
-- Change default values for feature_links and feature_qa to 'disabled'
-- so newly created profiles start with these modules off.
-- feature_relations keeps its default of 'public'.

ALTER TABLE "profile"
  ALTER COLUMN "feature_links" SET DEFAULT 'disabled';

ALTER TABLE "profile"
  ALTER COLUMN "feature_qa" SET DEFAULT 'disabled';

-- +goose Down
ALTER TABLE "profile"
  ALTER COLUMN "feature_links" SET DEFAULT 'public';

ALTER TABLE "profile"
  ALTER COLUMN "feature_qa" SET DEFAULT 'public';
