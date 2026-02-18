-- name: UpsertStoryInteraction :one
INSERT INTO "story_interaction" (
  id, story_id, profile_id, kind, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(story_id),
  sqlc.arg(profile_id),
  sqlc.arg(kind),
  NOW()
) ON CONFLICT (story_id, profile_id, kind)
DO UPDATE SET updated_at = NOW()
RETURNING *;

-- name: RemoveStoryInteraction :execrows
UPDATE "story_interaction"
SET deleted_at = NOW()
WHERE story_id = sqlc.arg(story_id)
  AND profile_id = sqlc.arg(profile_id)
  AND kind = sqlc.arg(kind)
  AND deleted_at IS NULL;

-- name: RemoveStoryInteractionsByKinds :execrows
-- Removes all interactions matching any of the given kinds for a user on a story.
-- Used to enforce mutual exclusivity for RSVP kinds.
UPDATE "story_interaction"
SET deleted_at = NOW()
WHERE story_id = sqlc.arg(story_id)
  AND profile_id = sqlc.arg(profile_id)
  AND kind = ANY(string_to_array(sqlc.arg(kinds), ','))
  AND deleted_at IS NULL;

-- name: GetStoryInteraction :one
SELECT * FROM "story_interaction"
WHERE story_id = sqlc.arg(story_id)
  AND profile_id = sqlc.arg(profile_id)
  AND kind = sqlc.arg(kind)
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListStoryInteractionsForProfile :many
-- Returns all active interactions a profile has on a specific story.
SELECT * FROM "story_interaction"
WHERE story_id = sqlc.arg(story_id)
  AND profile_id = sqlc.arg(profile_id)
  AND deleted_at IS NULL
ORDER BY created_at;

-- name: ListStoryInteractions :many
-- Lists interactions on a story with profile info, optionally filtered by kind.
SELECT
  si.id,
  si.story_id,
  si.profile_id,
  si.kind,
  si.created_at,
  p.slug as profile_slug,
  pt.title as profile_title,
  p.profile_picture_uri,
  p.kind as profile_kind
FROM "story_interaction" si
  INNER JOIN "profile" p ON p.id = si.profile_id
    AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
    AND pt.locale_code = (
      SELECT ptx.locale_code FROM "profile_tx" ptx
      WHERE ptx.profile_id = p.id
      ORDER BY CASE WHEN ptx.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
      LIMIT 1
    )
WHERE si.story_id = sqlc.arg(story_id)
  AND si.deleted_at IS NULL
  AND (sqlc.narg(filter_kind)::TEXT IS NULL OR si.kind = sqlc.narg(filter_kind)::TEXT)
ORDER BY si.created_at;

-- name: CountStoryInteractionsByKind :many
-- Returns interaction counts grouped by kind for a story.
SELECT kind, COUNT(*) as count
FROM "story_interaction"
WHERE story_id = sqlc.arg(story_id)
  AND deleted_at IS NULL
GROUP BY kind;
