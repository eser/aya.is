-- +goose Up
CREATE TABLE IF NOT EXISTS "profile_custom_domain" (
  "id" CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id" CHAR(26) NOT NULL REFERENCES "profile"("id"),
  "domain" TEXT NOT NULL,
  "default_locale" TEXT,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX "profile_custom_domain_domain_unique" ON "profile_custom_domain" ("domain");
CREATE INDEX "profile_custom_domain_profile_id_idx" ON "profile_custom_domain" ("profile_id");

-- Migrate existing data from profile.custom_domain
INSERT INTO "profile_custom_domain" ("id", "profile_id", "domain", "default_locale", "created_at")
SELECT
  -- Generate a ULID-like ID from existing data (26 char, uppercase alphanumeric)
  UPPER(SUBSTRING(REPLACE(gen_random_uuid()::TEXT, '-', ''), 1, 26)),
  "id",
  "custom_domain",
  NULL,
  COALESCE("created_at", NOW())
FROM "profile"
WHERE "custom_domain" IS NOT NULL
  AND "deleted_at" IS NULL;

ALTER TABLE "profile" DROP COLUMN "custom_domain";

-- +goose Down
ALTER TABLE "profile" ADD COLUMN "custom_domain" TEXT;

-- Migrate data back to profile.custom_domain (take the first domain per profile)
UPDATE "profile" p
SET "custom_domain" = pcd."domain"
FROM (
  SELECT DISTINCT ON ("profile_id") "profile_id", "domain"
  FROM "profile_custom_domain"
  ORDER BY "profile_id", "created_at" ASC
) pcd
WHERE p."id" = pcd."profile_id";

DROP TABLE IF EXISTS "profile_custom_domain";
