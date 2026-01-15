-- name: GetUserByID :one
SELECT *
FROM "user"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM "user"
WHERE email = sqlc.arg(email)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByGitHubRemoteID :one
SELECT *
FROM "user"
WHERE github_remote_id = sqlc.arg(github_remote_id)
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM "user"
WHERE (sqlc.narg(filter_kind)::TEXT IS NULL OR kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
  AND deleted_at IS NULL;

-- name: CreateUser :exec
INSERT INTO "user" (
    id,
    kind,
    NAME,
    email,
    phone,
    github_handle,
    github_remote_id,
    bsky_handle,
    bsky_remote_id,
    x_handle,
    x_remote_id,
    individual_profile_id,
    created_at,
    updated_at,
    deleted_at
  )
VALUES (
    sqlc.arg(id),
    sqlc.arg(kind),
    sqlc.arg(NAME),
    sqlc.arg(email),
    sqlc.arg(phone),
    sqlc.arg(github_handle),
    sqlc.arg(github_remote_id),
    sqlc.arg(bsky_handle),
    sqlc.arg(bsky_remote_id),
    sqlc.arg(x_handle),
    sqlc.arg(x_remote_id),
    sqlc.arg(individual_profile_id),
    sqlc.arg(created_at),
    sqlc.arg(updated_at),
    sqlc.arg(deleted_at)
  );

-- name: UpdateUser :execrows
UPDATE "user"
SET kind = sqlc.arg(kind),
  NAME = sqlc.arg(NAME),
  email = sqlc.arg(email),
  phone = sqlc.arg(phone),
  github_handle = sqlc.arg(github_handle),
  github_remote_id = sqlc.arg(github_remote_id),
  bsky_handle = sqlc.arg(bsky_handle),
  bsky_remote_id = sqlc.arg(bsky_remote_id),
  x_handle = sqlc.arg(x_handle),
  x_remote_id = sqlc.arg(x_remote_id),
  individual_profile_id = sqlc.arg(individual_profile_id)
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: SetUserIndividualProfileID :execrows
UPDATE "user"
SET individual_profile_id = sqlc.arg(individual_profile_id)
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: RemoveUser :execrows
UPDATE "user"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;
