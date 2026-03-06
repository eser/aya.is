-- name: GetStorySeriesByID :one
SELECT
  ss.id, ss.slug, ss.series_picture_uri, ss.created_at, ss.updated_at,
  sst.locale_code, sst.title, sst.description
FROM "story_series" ss
  INNER JOIN "story_series_tx" sst ON sst.story_series_id = ss.id
  AND sst.locale_code = (
    SELECT sstx.locale_code FROM "story_series_tx" sstx
    WHERE sstx.story_series_id = ss.id
    ORDER BY CASE
      WHEN sstx.locale_code = sqlc.arg(locale_code) THEN 0
      WHEN RTRIM(sstx.locale_code) = 'en' THEN 1
      ELSE 2
    END
    LIMIT 1
  )
WHERE ss.id = sqlc.arg(id)
  AND ss.deleted_at IS NULL
LIMIT 1;

-- name: GetStorySeriesBySlug :one
SELECT
  ss.id, ss.slug, ss.series_picture_uri, ss.created_at, ss.updated_at,
  sst.locale_code, sst.title, sst.description
FROM "story_series" ss
  INNER JOIN "story_series_tx" sst ON sst.story_series_id = ss.id
  AND sst.locale_code = (
    SELECT sstx.locale_code FROM "story_series_tx" sstx
    WHERE sstx.story_series_id = ss.id
    ORDER BY CASE
      WHEN sstx.locale_code = sqlc.arg(locale_code) THEN 0
      WHEN RTRIM(sstx.locale_code) = 'en' THEN 1
      ELSE 2
    END
    LIMIT 1
  )
WHERE ss.slug = sqlc.arg(slug)
  AND ss.deleted_at IS NULL
LIMIT 1;

-- name: ListStorySeries :many
SELECT
  ss.id, ss.slug, ss.series_picture_uri, ss.created_at, ss.updated_at,
  sst.locale_code, sst.title, sst.description
FROM "story_series" ss
  INNER JOIN "story_series_tx" sst ON sst.story_series_id = ss.id
  AND sst.locale_code = (
    SELECT sstx.locale_code FROM "story_series_tx" sstx
    WHERE sstx.story_series_id = ss.id
    ORDER BY CASE
      WHEN sstx.locale_code = sqlc.arg(locale_code) THEN 0
      WHEN RTRIM(sstx.locale_code) = 'en' THEN 1
      ELSE 2
    END
    LIMIT 1
  )
WHERE ss.deleted_at IS NULL
ORDER BY ss.created_at DESC;

-- name: InsertStorySeries :one
INSERT INTO "story_series" (
  id, slug, series_picture_uri, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(slug),
  sqlc.narg(series_picture_uri),
  NOW()
) RETURNING *;

-- name: UpdateStorySeries :execrows
UPDATE "story_series"
SET
  slug = sqlc.arg(slug),
  series_picture_uri = sqlc.narg(series_picture_uri),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: RemoveStorySeries :execrows
UPDATE "story_series"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: UpsertStorySeriesTx :exec
INSERT INTO "story_series_tx" (story_series_id, locale_code, title, description)
VALUES (
  sqlc.arg(story_series_id),
  sqlc.arg(locale_code),
  sqlc.arg(title),
  sqlc.arg(description)
)
ON CONFLICT (story_series_id, locale_code) DO UPDATE SET
  title = EXCLUDED.title,
  description = EXCLUDED.description;
