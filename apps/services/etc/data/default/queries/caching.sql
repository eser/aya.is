-- name: GetFromCache :one
SELECT value, updated_at
FROM "cache"
WHERE key = sqlc.arg(key)
LIMIT 1;

-- name: GetFromCacheSince :one
SELECT value, updated_at
FROM "cache"
WHERE key = sqlc.arg(key)
  AND updated_at > sqlc.arg(since)
LIMIT 1;

-- name: SetInCache :execrows
INSERT INTO "cache" (key, value, updated_at)
VALUES (sqlc.arg(key), sqlc.arg(value), NOW())
ON CONFLICT ("key") DO UPDATE SET value = sqlc.arg(value), updated_at = NOW();

-- name: RemoveFromCache :execrows
DELETE FROM "cache"
WHERE key = sqlc.arg(key);

-- name: RemoveAllFromCache :execrows
DELETE FROM "cache";

-- name: RemoveExpiredFromCache :execrows
DELETE FROM "cache"
WHERE updated_at < sqlc.arg(before);
