-- +goose Up
ALTER TABLE "profile" ADD COLUMN "default_locale" TEXT NOT NULL DEFAULT 'en';

-- +goose Down
ALTER TABLE "profile" DROP COLUMN "default_locale";
