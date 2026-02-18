-- name: GetStorySeriesByID :one
SELECT * FROM "story_series"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetStorySeriesBySlug :one
SELECT * FROM "story_series"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListStorySeries :many
SELECT * FROM "story_series"
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: InsertStorySeries :one
INSERT INTO "story_series" (
  id, slug, series_picture_uri, title, description, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(slug),
  sqlc.narg(series_picture_uri),
  sqlc.arg(title),
  sqlc.arg(description),
  NOW()
) RETURNING *;

-- name: UpdateStorySeries :execrows
UPDATE "story_series"
SET
  slug = sqlc.arg(slug),
  series_picture_uri = sqlc.narg(series_picture_uri),
  title = sqlc.arg(title),
  description = sqlc.arg(description),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: RemoveStorySeries :execrows
UPDATE "story_series"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;
