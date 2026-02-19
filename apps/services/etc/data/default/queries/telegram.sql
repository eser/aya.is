-- name: CreateTelegramVerificationCode :exec
INSERT INTO "telegram_verification_code" (id, code, telegram_user_id, telegram_username, created_at, expires_at)
VALUES (sqlc.arg(id), sqlc.arg(code), sqlc.arg(telegram_user_id), sqlc.arg(telegram_username), NOW(), sqlc.arg(expires_at));

-- name: GetTelegramVerificationCodeByCode :one
SELECT *
FROM "telegram_verification_code"
WHERE code = sqlc.arg(code)
  AND consumed_at IS NULL
  AND expires_at > NOW()
LIMIT 1;

-- name: ConsumeTelegramVerificationCode :execrows
UPDATE "telegram_verification_code"
SET consumed_at = NOW()
WHERE code = sqlc.arg(code)
  AND consumed_at IS NULL;

-- name: CleanupExpiredTelegramVerificationCodes :execrows
DELETE FROM "telegram_verification_code"
WHERE expires_at < NOW() - INTERVAL '1 hour';

-- name: GetProfileLinkByTelegramRemoteID :one
SELECT pl.id, pl.profile_id, pl.remote_id, pl.public_id
FROM "profile_link" pl
WHERE pl.kind = 'telegram'
  AND pl.remote_id = sqlc.arg(remote_id)
  AND pl.deleted_at IS NULL
LIMIT 1;

-- name: GetProfileLinkByProfileIDAndTelegram :one
SELECT pl.id, pl.profile_id, pl.remote_id, pl.public_id
FROM "profile_link" pl
WHERE pl.profile_id = sqlc.arg(profile_id)
  AND pl.kind = 'telegram'
  AND pl.is_managed = TRUE
  AND pl.deleted_at IS NULL
LIMIT 1;

-- name: SoftDeleteTelegramProfileLink :execrows
UPDATE "profile_link"
SET deleted_at = NOW()
WHERE kind = 'telegram'
  AND remote_id = sqlc.arg(remote_id)
  AND deleted_at IS NULL;

-- name: ListManagedTelegramLinks :many
SELECT
  pl.id,
  pl.profile_id,
  pl.remote_id,
  pl.public_id
FROM "profile_link" pl
  INNER JOIN "profile" p ON p.id = pl.profile_id
    AND p.deleted_at IS NULL
WHERE pl.kind = 'telegram'
  AND pl.is_managed = TRUE
  AND pl.deleted_at IS NULL
ORDER BY pl.updated_at ASC NULLS FIRST
LIMIT sqlc.arg(limit_count);

-- name: GetProfileSlugByIDForTelegram :one
SELECT slug
FROM "profile"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
LIMIT 1;
