-- name: ListManagedLinksForKind :many
SELECT
  pl.id,
  pl.profile_id,
  pl.kind,
  pl.remote_id,
  pl.auth_access_token,
  pl.auth_access_token_expires_at,
  pl.auth_refresh_token
FROM "profile_link" pl
  INNER JOIN "profile" p ON p.id = pl.profile_id
    AND p.deleted_at IS NULL
WHERE pl.kind = sqlc.arg(kind)
  AND pl.is_managed = TRUE
  AND pl.auth_access_token IS NOT NULL
  AND pl.deleted_at IS NULL
ORDER BY pl.updated_at ASC NULLS FIRST
LIMIT sqlc.arg(limit_count);

-- name: GetLatestImportByLinkID :one
SELECT *
FROM "profile_link_import"
WHERE profile_link_id = sqlc.arg(profile_link_id)
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: GetLinkImportByRemoteID :one
SELECT *
FROM "profile_link_import"
WHERE profile_link_id = sqlc.arg(profile_link_id)
  AND remote_id = sqlc.arg(remote_id)
LIMIT 1;

-- name: CreateLinkImport :exec
INSERT INTO "profile_link_import" (id, profile_link_id, remote_id, properties, created_at)
VALUES (sqlc.arg(id), sqlc.arg(profile_link_id), sqlc.arg(remote_id), sqlc.arg(properties), NOW());

-- name: UpdateLinkImport :execrows
UPDATE "profile_link_import"
SET
  properties = sqlc.arg(properties),
  updated_at = NOW(),
  deleted_at = NULL
WHERE id = sqlc.arg(id);

-- name: MarkLinkImportsDeletedExcept :execrows
UPDATE "profile_link_import"
SET deleted_at = NOW()
WHERE profile_link_id = sqlc.arg(profile_link_id)
  AND deleted_at IS NULL
  AND remote_id IS NOT NULL
  AND remote_id != ALL(sqlc.arg(active_remote_ids)::TEXT[]);

-- name: ListLinkImportsForStoryCreation :many
SELECT
  pli.id,
  pli.profile_link_id,
  pli.remote_id,
  pli.properties,
  pli.created_at,
  pl.profile_id,
  COALESCE(NULLIF(TRIM(p.default_locale), ''), 'en')::TEXT AS profile_default_locale
FROM "profile_link_import" pli
  INNER JOIN "profile_link" pl ON pl.id = pli.profile_link_id
  INNER JOIN "profile" p ON p.id = pl.profile_id
WHERE pl.kind = sqlc.arg(kind)
  AND pl.is_managed = TRUE
  AND pli.deleted_at IS NULL
  AND pli.remote_id IS NOT NULL
  AND NOT EXISTS (
    SELECT 1 FROM "story" s
    WHERE s.author_profile_id = pl.profile_id
      AND s.is_managed = TRUE
      AND s.properties->>'remote_id' = pli.remote_id
      AND s.deleted_at IS NULL
  )
ORDER BY pli.created_at ASC
LIMIT sqlc.arg(limit_count);

-- name: ListImportsWithExistingStories :many
-- Returns imports that have matching managed stories (for reconciliation during full sync).
SELECT
  pli.id,
  pli.profile_link_id,
  pli.remote_id,
  pli.properties,
  pli.created_at,
  pl.profile_id,
  COALESCE(NULLIF(TRIM(p.default_locale), ''), 'en')::TEXT AS profile_default_locale,
  s.id AS story_id,
  sp.id AS publication_id
FROM "profile_link_import" pli
  INNER JOIN "profile_link" pl ON pl.id = pli.profile_link_id
  INNER JOIN "profile" p ON p.id = pl.profile_id
  INNER JOIN "story" s ON s.author_profile_id = pl.profile_id
    AND s.is_managed = TRUE
    AND s.properties->>'remote_id' = pli.remote_id
    AND s.deleted_at IS NULL
  LEFT JOIN "story_publication" sp ON sp.story_id = s.id
    AND sp.profile_id = pl.profile_id
    AND sp.deleted_at IS NULL
WHERE pl.kind = sqlc.arg(kind)
  AND pl.is_managed = TRUE
  AND pli.deleted_at IS NULL
  AND pli.remote_id IS NOT NULL
ORDER BY pli.created_at ASC
LIMIT sqlc.arg(limit_count);

-- name: UpdateProfileLinkTokens :execrows
UPDATE "profile_link"
SET
  auth_access_token = sqlc.arg(auth_access_token),
  auth_access_token_expires_at = sqlc.narg(auth_access_token_expires_at),
  auth_refresh_token = sqlc.narg(auth_refresh_token),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;
