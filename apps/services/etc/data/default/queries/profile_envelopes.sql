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
SELECT *
FROM "profile_envelope"
WHERE target_profile_id = sqlc.arg(target_profile_id)
  AND deleted_at IS NULL
  AND (sqlc.narg(status_filter)::TEXT IS NULL OR status = sqlc.narg(status_filter))
ORDER BY created_at DESC
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
SELECT *
FROM "profile_envelope"
WHERE target_profile_id = sqlc.arg(target_profile_id)
  AND kind = 'invitation'
  AND status = 'accepted'
  AND deleted_at IS NULL
  AND (sqlc.narg(invitation_kind)::TEXT IS NULL
       OR properties->>'invitation_kind' = sqlc.narg(invitation_kind))
ORDER BY created_at DESC;

-- name: CountPendingProfileEnvelopes :one
SELECT COUNT(*)::INT as count
FROM "profile_envelope"
WHERE target_profile_id = sqlc.arg(target_profile_id)
  AND status = 'pending'
  AND deleted_at IS NULL;
