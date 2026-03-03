-- +goose Up
ALTER TABLE "profile_link" ADD COLUMN "is_online" BOOLEAN DEFAULT FALSE NOT NULL;

-- +goose Down
ALTER TABLE "profile_link" DROP COLUMN "is_online";
