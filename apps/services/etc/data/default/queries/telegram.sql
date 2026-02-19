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

-- name: GetMemberProfileTelegramLinks :many
-- For a given member profile, find all non-individual profiles they belong to
-- and return the telegram links on those profiles (visibility filtering happens in Go).
SELECT
  p.id as profile_id,
  p.slug as profile_slug,
  pt.title as profile_title,
  pm.kind as membership_kind,
  pl.id as link_id,
  pl.uri,
  pl.public_id as link_public_id,
  pl.visibility as link_visibility,
  COALESCE(plt.title, '') as link_title
FROM "profile_membership" pm
  INNER JOIN "profile" p ON p.id = pm.profile_id
    AND p.kind != 'individual'
    AND p.approved_at IS NOT NULL
    AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
    AND pt.locale_code = (
      SELECT ptf.locale_code FROM "profile_tx" ptf
      WHERE ptf.profile_id = p.id
      ORDER BY CASE WHEN ptf.locale_code = p.default_locale THEN 0 ELSE 1 END
      LIMIT 1
    )
  INNER JOIN "profile_link" pl ON pl.profile_id = p.id
    AND pl.kind = 'telegram'
    AND pl.deleted_at IS NULL
  LEFT JOIN "profile_link_tx" plt ON plt.profile_link_id = pl.id
    AND plt.locale_code = pt.locale_code
WHERE pm.member_profile_id = sqlc.arg(member_profile_id)
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
ORDER BY p.slug, pl."order";

-- name: CreateTelegramGroupInviteCode :exec
INSERT INTO "telegram_group_invite_code" (id, code, telegram_chat_id, telegram_chat_title, created_by_telegram_user_id, created_at, expires_at)
VALUES (sqlc.arg(id), sqlc.arg(code), sqlc.arg(telegram_chat_id), sqlc.arg(telegram_chat_title), sqlc.arg(created_by_telegram_user_id), NOW(), sqlc.arg(expires_at));

-- name: GetTelegramGroupInviteCodeByCode :one
SELECT *
FROM "telegram_group_invite_code"
WHERE code = sqlc.arg(code)
  AND consumed_at IS NULL
  AND expires_at > NOW()
LIMIT 1;

-- name: ConsumeTelegramGroupInviteCode :execrows
UPDATE "telegram_group_invite_code"
SET consumed_at = NOW()
WHERE code = sqlc.arg(code)
  AND consumed_at IS NULL;
