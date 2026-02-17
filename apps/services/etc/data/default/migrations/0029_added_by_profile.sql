-- +goose Up
-- Add added_by_profile_id tracking to profile_link and profile_page tables
-- Nullable because existing rows have no adder data

ALTER TABLE "profile_link"
  ADD COLUMN "added_by_profile_id" CHAR(26)
  CONSTRAINT "profile_link_added_by_profile_id_fk" REFERENCES "profile" ("id");

ALTER TABLE "profile_page"
  ADD COLUMN "added_by_profile_id" CHAR(26)
  CONSTRAINT "profile_page_added_by_profile_id_fk" REFERENCES "profile" ("id");

-- +goose Down
ALTER TABLE "profile_page" DROP COLUMN IF EXISTS "added_by_profile_id";
ALTER TABLE "profile_link" DROP COLUMN IF EXISTS "added_by_profile_id";
