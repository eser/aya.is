-- +goose Up
ALTER TABLE "profile_custom_domain"
  ADD COLUMN "www_prefix" BOOLEAN NOT NULL DEFAULT TRUE;

-- +goose Down
ALTER TABLE "profile_custom_domain"
  DROP COLUMN "www_prefix";
