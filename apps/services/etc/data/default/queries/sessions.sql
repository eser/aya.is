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
