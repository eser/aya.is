-- name: GetProfilePoints :one
SELECT points
FROM "profile"
WHERE id = sqlc.arg(profile_id)
  AND deleted_at IS NULL;

-- name: RecordProfilePointTransaction :one
INSERT INTO "profile_point_transaction" (
  id,
  target_profile_id,
  origin_profile_id,
  transaction_type,
  triggering_event,
  description,
  amount,
  balance_after,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(target_profile_id),
  sqlc.narg(origin_profile_id),
  sqlc.arg(transaction_type),
  sqlc.narg(triggering_event),
  sqlc.arg(description),
  sqlc.arg(amount),
  sqlc.arg(balance_after),
  NOW()
) RETURNING *;

-- name: AddPointsToProfile :execrows
UPDATE "profile"
SET
  points = points + sqlc.arg(amount),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DeductPointsFromProfile :execrows
UPDATE "profile"
SET
  points = points - sqlc.arg(amount),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
  AND points >= sqlc.arg(amount);

-- name: ListProfilePointTransactionsByProfileID :many
SELECT *
FROM "profile_point_transaction"
WHERE target_profile_id = sqlc.arg(profile_id)
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count);

-- name: GetProfilePointTransactionByID :one
SELECT *
FROM "profile_point_transaction"
WHERE id = sqlc.arg(id);

-- Pending Award Queries

-- name: CreatePendingAward :one
INSERT INTO "profile_point_pending_award" (
  id,
  target_profile_id,
  triggering_event,
  description,
  amount,
  status,
  metadata,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(target_profile_id),
  sqlc.arg(triggering_event),
  sqlc.arg(description),
  sqlc.arg(amount),
  'pending',
  sqlc.narg(metadata),
  NOW()
) RETURNING *;

-- name: GetPendingAwardByID :one
SELECT *
FROM "profile_point_pending_award"
WHERE id = sqlc.arg(id);

-- name: ListPendingAwards :many
SELECT *
FROM "profile_point_pending_award"
WHERE (sqlc.narg(status)::text IS NULL OR status = sqlc.narg(status))
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count);

-- name: ListPendingAwardsByStatus :many
SELECT *
FROM "profile_point_pending_award"
WHERE status = sqlc.arg(status)
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count);

-- name: ApprovePendingAward :exec
UPDATE "profile_point_pending_award"
SET
  status = 'approved',
  reviewed_by = sqlc.arg(reviewed_by),
  reviewed_at = NOW()
WHERE id = sqlc.arg(id)
  AND status = 'pending';

-- name: RejectPendingAward :exec
UPDATE "profile_point_pending_award"
SET
  status = 'rejected',
  reviewed_by = sqlc.arg(reviewed_by),
  reviewed_at = NOW(),
  rejection_reason = sqlc.arg(rejection_reason)
WHERE id = sqlc.arg(id)
  AND status = 'pending';

-- name: GetPendingAwardsStats :one
SELECT
  COUNT(*) FILTER (WHERE status = 'pending') AS total_pending,
  COUNT(*) FILTER (WHERE status = 'approved') AS total_approved,
  COUNT(*) FILTER (WHERE status = 'rejected') AS total_rejected,
  COALESCE(SUM(amount) FILTER (WHERE status = 'approved'), 0) AS points_awarded
FROM "profile_point_pending_award";

-- name: GetPendingAwardsStatsByEventType :many
SELECT
  triggering_event,
  COUNT(*) AS count
FROM "profile_point_pending_award"
WHERE status = 'pending'
GROUP BY triggering_event;
