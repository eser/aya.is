-- +goose Up
ALTER TABLE "profile_custom_domain"
  ADD COLUMN "verification_status" TEXT NOT NULL DEFAULT 'pending',
  ADD COLUMN "dns_verified_at" TIMESTAMP WITH TIME ZONE,
  ADD COLUMN "last_dns_check_at" TIMESTAMP WITH TIME ZONE,
  ADD COLUMN "expired_at" TIMESTAMP WITH TIME ZONE,
  ADD COLUMN "webserver_synced" BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX "profile_custom_domain_verification_status_idx"
  ON "profile_custom_domain" ("verification_status");

-- +goose Down
DROP INDEX IF EXISTS "profile_custom_domain_verification_status_idx";

ALTER TABLE "profile_custom_domain"
  DROP COLUMN "webserver_synced",
  DROP COLUMN "expired_at",
  DROP COLUMN "last_dns_check_at",
  DROP COLUMN "dns_verified_at",
  DROP COLUMN "verification_status";
