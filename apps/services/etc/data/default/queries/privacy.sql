-- name: GetPrivacyStateByUserID :one
SELECT
  u.id AS user_id,
  p.id AS profile_id,
  p.slug AS profile_slug,
  p.hidden_at,
  p.deleted_at AS profile_deleted_at,
  p.redacted_at AS profile_redacted_at,
  u.deleted_at AS user_deleted_at,
  u.redacted_at AS user_redacted_at
FROM "user" u
  INNER JOIN "profile" p ON p.id = u.individual_profile_id
WHERE u.id = sqlc.arg(user_id);

-- name: HideProfileByID :execrows
UPDATE "profile"
SET deleted_at = NOW(),
  hidden_at = NOW(),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
  AND redacted_at IS NULL;

-- name: UnhideProfileByUserID :execrows
UPDATE "profile"
SET deleted_at = NULL,
  hidden_at = NULL,
  updated_at = NOW()
WHERE id = (
    SELECT u.individual_profile_id
    FROM "user" u
    WHERE u.id = sqlc.arg(user_id)
  )
  AND hidden_at IS NOT NULL
  AND redacted_at IS NULL;

-- name: TerminateActiveSessionsByUserID :execrows
UPDATE "session"
SET status = 'terminated',
  updated_at = NOW()
WHERE logged_in_user_id = sqlc.arg(user_id)
  AND status = 'active';

-- name: ListSoleOwnedProfiles :many
-- Profiles (other than the user's own) that would be left without any owner if this profile left them.
SELECT
  p.id,
  p.slug,
  p.kind,
  (
    SELECT ptx.title
    FROM "profile_tx" ptx
    WHERE ptx.profile_id = p.id
    ORDER BY CASE WHEN ptx.locale_code = p.default_locale THEN 0 ELSE 1 END
    LIMIT 1
  ) AS title
FROM "profile_membership" pm
  INNER JOIN "profile" p ON p.id = pm.profile_id
WHERE pm.member_profile_id = sqlc.arg(profile_id)
  AND pm.profile_id <> sqlc.arg(profile_id)
  AND pm.kind = 'owner'
  AND pm.finished_at IS NULL
  AND pm.deleted_at IS NULL
  AND p.deleted_at IS NULL
  AND p.kind <> 'individual'
  AND NOT EXISTS (
    SELECT 1
    FROM "profile_membership" other
    WHERE other.profile_id = pm.profile_id
      AND other.member_profile_id <> pm.member_profile_id
      AND other.kind = 'owner'
      AND other.finished_at IS NULL
      AND other.deleted_at IS NULL
  )
ORDER BY p.slug;

-- name: CreateDataExport :exec
INSERT INTO "profile_data_export" (id, user_id, profile_id, status, requested_at)
VALUES (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(profile_id),
    'pending',
    NOW()
  );

-- name: GetLatestDataExportByUserID :one
SELECT *
FROM "profile_data_export"
WHERE user_id = sqlc.arg(user_id)
ORDER BY requested_at DESC
LIMIT 1;

-- name: GetDataExportByID :one
SELECT *
FROM "profile_data_export"
WHERE id = sqlc.arg(id);

-- name: ClaimDataExport :execrows
-- Moves a pending (or stuck processing) export to processing; returns 0 rows if someone else already finished it.
UPDATE "profile_data_export"
SET status = 'processing',
  started_at = NOW(),
  error_message = NULL,
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND status IN ('pending', 'processing');

-- name: MarkDataExportReady :exec
UPDATE "profile_data_export"
SET status = 'ready',
  completed_at = NOW(),
  expires_at = sqlc.arg(expires_at),
  storage_key = sqlc.arg(storage_key),
  size_bytes = sqlc.arg(size_bytes),
  error_message = NULL,
  updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: MarkDataExportFailed :exec
UPDATE "profile_data_export"
SET status = 'failed',
  error_message = sqlc.arg(error_message),
  updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: MarkDataExportNotified :exec
UPDATE "profile_data_export"
SET notified_at = NOW(),
  updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: ListExpiredDataExports :many
SELECT *
FROM "profile_data_export"
WHERE status = 'ready'
  AND expires_at < NOW()
ORDER BY expires_at
LIMIT sqlc.arg(batch_size);

-- name: MarkDataExportExpired :exec
UPDATE "profile_data_export"
SET status = 'expired',
  storage_key = NULL,
  updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: ListDataExportsByUserID :many
SELECT *
FROM "profile_data_export"
WHERE user_id = sqlc.arg(user_id);

-- name: DeleteDataExportsByUserID :exec
DELETE FROM "profile_data_export"
WHERE user_id = sqlc.arg(user_id);

-- name: GetPrivacyContactByUserID :one
SELECT u.email,
  p.slug,
  p.default_locale
FROM "user" u
  INNER JOIN "profile" p ON p.id = u.individual_profile_id
WHERE u.id = sqlc.arg(user_id);
