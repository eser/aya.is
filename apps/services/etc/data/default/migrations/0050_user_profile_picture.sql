-- +goose Up
ALTER TABLE "user" ADD COLUMN profile_picture_uri TEXT;

-- +goose Down
ALTER TABLE "user" DROP COLUMN IF EXISTS profile_picture_uri;
