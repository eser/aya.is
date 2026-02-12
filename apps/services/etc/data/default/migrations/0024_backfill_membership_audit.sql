-- +goose Up
-- Backfill event_audit records from existing profile_membership data

-- Create a helper function that generates a ULID-like ID from a given timestamp
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_ulid_at(ts TIMESTAMP WITH TIME ZONE) RETURNS TEXT AS $$
DECLARE
  timestamp_part TEXT;
  random_part TEXT;
  chars TEXT := '0123456789ABCDEFGHJKMNPQRSTVWXYZ';
  i INT;
  ms BIGINT;
BEGIN
  ms := (EXTRACT(EPOCH FROM ts) * 1000)::BIGINT;

  timestamp_part := '';
  FOR i IN 1..10 LOOP
    timestamp_part := substring(chars FROM ((ms % 32) + 1)::INT FOR 1) || timestamp_part;
    ms := ms / 32;
  END LOOP;

  random_part := '';
  FOR i IN 1..16 LOOP
    random_part := random_part || substring(chars FROM (floor(random() * 32)::INT + 1) FOR 1);
  END LOOP;

  RETURN timestamp_part || random_part;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- For each profile_membership, insert a profile_membership_created audit event
INSERT INTO "event_audit" ("id", "event_type", "entity_type", "entity_id", "actor_kind", "payload", "created_at")
SELECT
  generate_ulid_at(pm.created_at),
  'profile_membership_created',
  'membership',
  pm.id,
  'system',
  jsonb_build_object(
    'profile_id', pm.profile_id,
    'member_profile_id', pm.member_profile_id,
    'kind', pm.kind,
    'backfilled', true
  ),
  pm.created_at
FROM "profile_membership" pm
WHERE pm.deleted_at IS NULL;

-- For each profile_membership with an updated_at, insert a profile_membership_updated audit event
INSERT INTO "event_audit" ("id", "event_type", "entity_type", "entity_id", "actor_kind", "payload", "created_at")
SELECT
  generate_ulid_at(pm.updated_at),
  'profile_membership_updated',
  'membership',
  pm.id,
  'system',
  jsonb_build_object(
    'profile_id', pm.profile_id,
    'member_profile_id', pm.member_profile_id,
    'kind', pm.kind,
    'backfilled', true
  ),
  pm.updated_at
FROM "profile_membership" pm
WHERE pm.deleted_at IS NULL
  AND pm.updated_at IS NOT NULL;

-- Drop the helper function
DROP FUNCTION IF EXISTS generate_ulid_at(TIMESTAMP WITH TIME ZONE);

-- +goose Down
DELETE FROM "event_audit"
WHERE event_type IN ('profile_membership_created', 'profile_membership_updated')
  AND payload->>'backfilled' = 'true';
