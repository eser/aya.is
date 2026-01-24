-- name: GetRuntimeState :one
SELECT key, value, updated_at
FROM "runtime_state"
WHERE key = sqlc.arg(key)
LIMIT 1;

-- name: SetRuntimeState :exec
INSERT INTO "runtime_state" (key, value, updated_at)
VALUES (sqlc.arg(key), sqlc.arg(value), NOW())
ON CONFLICT ("key") DO UPDATE SET value = sqlc.arg(value), updated_at = NOW();

-- name: RemoveRuntimeState :execrows
DELETE FROM "runtime_state"
WHERE key = sqlc.arg(key);

-- name: TryAdvisoryLock :one
SELECT pg_try_advisory_lock(sqlc.arg(lock_id)::BIGINT) AS acquired;

-- name: ReleaseAdvisoryLock :one
SELECT pg_advisory_unlock(sqlc.arg(lock_id)::BIGINT) AS released;
