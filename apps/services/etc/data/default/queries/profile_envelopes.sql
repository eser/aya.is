-- name: CreateProfileEnvelope :exec
INSERT INTO "profile_envelope" (
  id, target_profile_id, sender_profile_id, sender_user_id,
  kind, status, title, description, properties, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(target_profile_id),
  sqlc.narg(sender_profile_id),
  sqlc.narg(sender_user_id),
  sqlc.arg(kind),
  'pending',
  sqlc.arg(title),
  sqlc.narg(description),
  sqlc.narg(properties),
  NOW()
);

-- name: GetProfileEnvelopeByID :one
SELECT *
FROM "profile_envelope"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListProfileEnvelopesByTargetProfileID :many
SELECT
  pe.*,
  sp.slug as sender_profile_slug,
  COALESCE(spt.title, '') as sender_profile_title,
  sp.profile_picture_uri as sender_profile_picture_uri,
  sp.kind as sender_profile_kind
FROM "profile_envelope" pe
  LEFT JOIN "profile" sp ON sp.id = pe.sender_profile_id AND sp.deleted_at IS NULL
  LEFT JOIN "profile_tx" spt ON spt.profile_id = sp.id AND spt.locale_code = sp.default_locale
WHERE pe.target_profile_id = sqlc.arg(target_profile_id)
  AND pe.deleted_at IS NULL
  AND (sqlc.narg(status_filter)::TEXT IS NULL OR pe.status = sqlc.narg(status_filter))
ORDER BY pe.created_at DESC
LIMIT sqlc.arg(limit_count);

-- name: UpdateProfileEnvelopeStatusToAccepted :execrows
UPDATE "profile_envelope"
SET status = 'accepted',
    accepted_at = sqlc.arg(now),
    updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(id)
  AND status = 'pending'
  AND deleted_at IS NULL;

-- name: UpdateProfileEnvelopeStatusToRejected :execrows
UPDATE "profile_envelope"
SET status = 'rejected',
    rejected_at = sqlc.arg(now),
    updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(id)
  AND status = 'pending'
  AND deleted_at IS NULL;

-- name: UpdateProfileEnvelopeStatusToRevoked :execrows
UPDATE "profile_envelope"
SET status = 'revoked',
    revoked_at = sqlc.arg(now),
    updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(id)
  AND status = 'pending'
  AND deleted_at IS NULL;

-- name: UpdateProfileEnvelopeStatusToRedeemed :execrows
UPDATE "profile_envelope"
SET status = 'redeemed',
    redeemed_at = sqlc.arg(now),
    updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(id)
  AND status = 'accepted'
  AND deleted_at IS NULL;

-- name: UpdateProfileEnvelopeProperties :exec
UPDATE "profile_envelope"
SET properties = sqlc.arg(properties),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListAcceptedInvitations :many
SELECT
  pe.*,
  sp.slug as sender_profile_slug,
  COALESCE(spt.title, '') as sender_profile_title,
  sp.profile_picture_uri as sender_profile_picture_uri,
  sp.kind as sender_profile_kind
FROM "profile_envelope" pe
  LEFT JOIN "profile" sp ON sp.id = pe.sender_profile_id AND sp.deleted_at IS NULL
  LEFT JOIN "profile_tx" spt ON spt.profile_id = sp.id AND spt.locale_code = sp.default_locale
WHERE pe.target_profile_id = sqlc.arg(target_profile_id)
  AND pe.kind = 'invitation'
  AND pe.status = 'accepted'
  AND pe.deleted_at IS NULL
  AND (sqlc.narg(invitation_kind)::TEXT IS NULL
       OR pe.properties->>'invitation_kind' = sqlc.narg(invitation_kind))
ORDER BY pe.created_at DESC;

-- name: CountPendingProfileEnvelopes :one
SELECT COUNT(*)::INT as count
FROM "profile_envelope"
WHERE target_profile_id = sqlc.arg(target_profile_id)
  AND status = 'pending'
  AND deleted_at IS NULL;
