-- +goose Up
-- Backfill missing profile_membership records for profiles

-- Create a function to generate ULID-like IDs (26 char, Crockford base32)
-- Note: This is a simplified version - real ULIDs have more precise timestamp encoding
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_ulid() RETURNS TEXT AS $$
DECLARE
  timestamp_part TEXT;
  random_part TEXT;
  chars TEXT := '0123456789ABCDEFGHJKMNPQRSTVWXYZ';
  i INT;
  ms BIGINT;
BEGIN
  -- Get current timestamp in milliseconds
  ms := (EXTRACT(EPOCH FROM clock_timestamp()) * 1000)::BIGINT;

  -- Encode timestamp (10 chars)
  timestamp_part := '';
  FOR i IN 1..10 LOOP
    timestamp_part := substring(chars FROM ((ms % 32) + 1) FOR 1) || timestamp_part;
    ms := ms / 32;
  END LOOP;

  -- Generate random part (16 chars)
  random_part := '';
  FOR i IN 1..16 LOOP
    random_part := random_part || substring(chars FROM (floor(random() * 32)::INT + 1) FOR 1);
  END LOOP;

  RETURN timestamp_part || random_part;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- For individual profiles: create self-membership where the profile is its own owner
INSERT INTO "profile_membership" ("id", "profile_id", "member_profile_id", "kind", "started_at", "created_at")
SELECT
  generate_ulid(),
  p.id,
  p.id,
  'owner',
  p.created_at,
  NOW()
FROM "profile" p
WHERE p.kind = 'individual'
  AND p.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM "profile_membership" pm
    WHERE pm.profile_id = p.id
      AND pm.member_profile_id = p.id
      AND pm.deleted_at IS NULL
  );

-- For org/product profiles: create owner membership linking to the creator's individual profile
-- Links profiles to the user whose individual_profile_id matches stories they authored for that profile
INSERT INTO "profile_membership" ("id", "profile_id", "member_profile_id", "kind", "started_at", "created_at")
SELECT DISTINCT ON (p.id)
  generate_ulid(),
  p.id,
  u.individual_profile_id,
  'owner',
  p.created_at,
  NOW()
FROM "profile" p
INNER JOIN "story" s ON s.profile_id = p.id AND s.deleted_at IS NULL
INNER JOIN "user" u ON u.individual_profile_id = s.author_profile_id AND u.deleted_at IS NULL
WHERE p.kind IN ('organization', 'product')
  AND p.deleted_at IS NULL
  AND u.individual_profile_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM "profile_membership" pm
    WHERE pm.profile_id = p.id
      AND pm.kind = 'owner'
      AND pm.deleted_at IS NULL
  )
ORDER BY p.id, s.created_at ASC;

-- Drop the helper function
DROP FUNCTION IF EXISTS generate_ulid();

-- +goose Down
-- Remove memberships created by this migration
-- Note: This only removes 'owner' memberships that have matching self-references for individual profiles
DELETE FROM "profile_membership" pm
USING "profile" p
WHERE pm.profile_id = p.id
  AND pm.member_profile_id = p.id
  AND pm.kind = 'owner'
  AND p.kind = 'individual';
