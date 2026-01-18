-- name: GetSessionByID :one
SELECT
  id,
  status,
  oauth_request_state,
  oauth_request_code_verifier,
  oauth_redirect_uri,
  logged_in_user_id,
  logged_in_at,
  expires_at,
  created_at,
  updated_at
FROM
  session
WHERE
  id = sqlc.arg(id);

-- name: CreateSession :exec
INSERT INTO
  session (
    id,
    status,
    oauth_request_state,
    oauth_request_code_verifier,
    oauth_redirect_uri,
    logged_in_user_id,
    logged_in_at,
    expires_at,
    created_at,
    updated_at
  )
VALUES
  (
    sqlc.arg(id),
    sqlc.arg(status),
    sqlc.arg(oauth_request_state),
    sqlc.arg(oauth_request_code_verifier),
    sqlc.arg(oauth_redirect_uri),
    sqlc.arg(logged_in_user_id),
    sqlc.arg(logged_in_at),
    sqlc.arg(expires_at),
    sqlc.arg(created_at),
    sqlc.arg(updated_at)
  );

-- name: UpdateSessionLoggedInAt :exec
UPDATE
  session
SET
  logged_in_at = sqlc.arg(logged_in_at),
  updated_at = NOW()
WHERE
  id = sqlc.arg(id);

-- Session Preferences

-- name: GetSessionPreferences :many
SELECT
  session_id,
  key,
  value,
  updated_at
FROM
  session_preference
WHERE
  session_id = sqlc.arg(session_id);

-- name: GetSessionPreference :one
SELECT
  session_id,
  key,
  value,
  updated_at
FROM
  session_preference
WHERE
  session_id = sqlc.arg(session_id)
  AND key = sqlc.arg(key);

-- name: SetSessionPreference :exec
INSERT INTO
  session_preference (session_id, key, value, updated_at)
VALUES
  (
    sqlc.arg(session_id),
    sqlc.arg(key),
    sqlc.arg(value),
    NOW()
  )
ON CONFLICT (session_id, key) DO UPDATE
SET
  value = sqlc.arg(value),
  updated_at = NOW();

-- name: DeleteSessionPreference :exec
DELETE FROM session_preference
WHERE
  session_id = sqlc.arg(session_id)
  AND key = sqlc.arg(key);

-- name: DeleteAllSessionPreferences :exec
DELETE FROM session_preference
WHERE
  session_id = sqlc.arg(session_id);

-- Session Rate Limiting

-- name: GetSessionRateLimit :one
SELECT
  ip_hash,
  count,
  window_start
FROM
  session_rate_limit
WHERE
  ip_hash = sqlc.arg(ip_hash);

-- name: UpsertSessionRateLimit :exec
INSERT INTO
  session_rate_limit (ip_hash, count, window_start)
VALUES
  (sqlc.arg(ip_hash), 1, NOW())
ON CONFLICT (ip_hash) DO UPDATE
SET
  count = CASE
    WHEN session_rate_limit.window_start < NOW() - INTERVAL '1 hour'
    THEN 1
    ELSE session_rate_limit.count + 1
  END,
  window_start = CASE
    WHEN session_rate_limit.window_start < NOW() - INTERVAL '1 hour'
    THEN NOW()
    ELSE session_rate_limit.window_start
  END;

-- name: CleanupExpiredSessionRateLimits :exec
DELETE FROM session_rate_limit
WHERE
  window_start < NOW() - INTERVAL '2 hours';
