-- +goose Up

-- Form template: defines the set of questions for an org's application
CREATE TABLE IF NOT EXISTS "profile_application_form" (
  "id"                    CHAR(26) NOT NULL PRIMARY KEY,
  "profile_id"            CHAR(26) NOT NULL
    CONSTRAINT "profile_application_form_profile_id_fk" REFERENCES "profile" ("id"),
  "preset_key"            TEXT,
  "is_active"             BOOLEAN NOT NULL DEFAULT TRUE,
  "responses_visibility"  TEXT NOT NULL DEFAULT 'members',
  "created_at"            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "updated_at"            TIMESTAMP WITH TIME ZONE
);

ALTER TABLE "profile_application_form"
  ADD CONSTRAINT "profile_application_form_responses_visibility_check"
    CHECK (responses_visibility IN ('members', 'leads'));

-- One active form per profile
CREATE UNIQUE INDEX IF NOT EXISTS "profile_application_form_active_uniq"
  ON "profile_application_form" ("profile_id")
  WHERE "is_active" = TRUE;

-- Individual fields/questions in the form
CREATE TABLE IF NOT EXISTS "profile_application_form_field" (
  "id"          CHAR(26) NOT NULL PRIMARY KEY,
  "form_id"     CHAR(26) NOT NULL
    CONSTRAINT "profile_application_form_field_form_id_fk"
    REFERENCES "profile_application_form" ("id") ON DELETE CASCADE,
  "label"       TEXT NOT NULL,
  "field_type"  TEXT NOT NULL DEFAULT 'long_text',
  "is_required" BOOLEAN NOT NULL DEFAULT FALSE,
  "sort_order"  INTEGER NOT NULL DEFAULT 0,
  "placeholder" TEXT,
  "created_at"  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

ALTER TABLE "profile_application_form_field"
  ADD CONSTRAINT "profile_application_form_field_type_check"
    CHECK (field_type IN ('short_text', 'long_text', 'url'));

CREATE INDEX IF NOT EXISTS "profile_application_form_field_form_id_idx"
  ON "profile_application_form_field" ("form_id", "sort_order");

-- Applicant's responses to form fields
CREATE TABLE IF NOT EXISTS "profile_candidate_response" (
  "id"            CHAR(26) NOT NULL PRIMARY KEY,
  "candidate_id"  CHAR(26) NOT NULL
    CONSTRAINT "profile_candidate_response_candidate_id_fk"
    REFERENCES "profile_membership_candidate" ("id") ON DELETE CASCADE,
  "form_field_id" CHAR(26) NOT NULL
    CONSTRAINT "profile_candidate_response_field_id_fk"
    REFERENCES "profile_application_form_field" ("id"),
  "value"         TEXT NOT NULL DEFAULT '',
  "created_at"    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "profile_candidate_response_candidate_field_uniq"
  ON "profile_candidate_response" ("candidate_id", "form_field_id");

-- +goose Down
DROP INDEX IF EXISTS "profile_candidate_response_candidate_field_uniq";
DROP TABLE IF EXISTS "profile_candidate_response";
DROP INDEX IF EXISTS "profile_application_form_field_form_id_idx";
ALTER TABLE "profile_application_form_field"
  DROP CONSTRAINT IF EXISTS "profile_application_form_field_type_check";
DROP TABLE IF EXISTS "profile_application_form_field";
DROP INDEX IF EXISTS "profile_application_form_active_uniq";
ALTER TABLE "profile_application_form"
  DROP CONSTRAINT IF EXISTS "profile_application_form_responses_visibility_check";
DROP TABLE IF EXISTS "profile_application_form";
