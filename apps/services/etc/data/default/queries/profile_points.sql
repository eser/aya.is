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
