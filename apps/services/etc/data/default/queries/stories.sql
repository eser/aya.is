-- name: GetStoryIDBySlug :one
SELECT id
FROM "story"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetStoryByID :one
SELECT
  sqlc.embed(s),
  sqlc.embed(st),
  sqlc.embed(p),
  sqlc.embed(pt),
  pb.publications
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile" p ON p.id = s.author_profile_id
  AND p.approved_at IS NOT NULL
  AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
  LEFT JOIN LATERAL (
    SELECT JSONB_AGG(
      JSONB_BUILD_OBJECT('profile', row_to_json(p2), 'profile_tx', row_to_json(p2t))
    ) AS "publications"
    FROM story_publication sp
      INNER JOIN "profile" p2 ON p2.id = sp.profile_id
      AND p2.deleted_at IS NULL
      INNER JOIN "profile_tx" p2t ON p2t.profile_id = p2.id
      AND p2t.locale_code = sqlc.arg(locale_code)
    WHERE sp.story_id = s.id
      AND (sqlc.narg(filter_publication_profile_id)::CHAR(26) IS NULL OR sp.profile_id = sqlc.narg(filter_publication_profile_id)::CHAR(26))
      AND sp.deleted_at IS NULL
  ) pb ON TRUE
WHERE s.id = sqlc.arg(id)
  AND (sqlc.narg(filter_author_profile_id)::CHAR(26) IS NULL OR s.author_profile_id = sqlc.narg(filter_author_profile_id)::CHAR(26))
  AND s.deleted_at IS NULL
LIMIT 1;

-- -- name: ListStories :many
-- SELECT sqlc.embed(s), sqlc.embed(st), sqlc.embed(p), sqlc.embed(pt)
-- FROM "story" s
--   INNER JOIN "story_tx" st ON st.story_id = s.id
--   AND (sqlc.narg(filter_kind)::TEXT IS NULL OR s.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
--   AND (sqlc.narg(filter_author_profile_id)::CHAR(26) IS NULL OR s.author_profile_id = sqlc.narg(filter_author_profile_id)::CHAR(26))
--   AND st.locale_code = sqlc.arg(locale_code)
--   LEFT JOIN "profile" p ON p.id = s.author_profile_id AND p.deleted_at IS NULL
--   INNER JOIN "profile_tx" pt ON pt.profile_id = p.id AND pt.locale_code = sqlc.arg(locale_code)
-- WHERE s.deleted_at IS NULL
-- ORDER BY s.created_at DESC;

-- name: InsertStory :one
INSERT INTO "story" (
  id,
  author_profile_id,
  slug,
  kind,
  status,
  is_featured,
  story_picture_uri,
  properties,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(author_profile_id),
  sqlc.arg(slug),
  sqlc.arg(kind),
  sqlc.arg(status),
  sqlc.arg(is_featured),
  sqlc.narg(story_picture_uri),
  sqlc.narg(properties),
  NOW()
) RETURNING *;

-- name: InsertStoryTx :exec
INSERT INTO "story_tx" (
  story_id,
  locale_code,
  title,
  summary,
  content
) VALUES (
  sqlc.arg(story_id),
  sqlc.arg(locale_code),
  sqlc.arg(title),
  sqlc.arg(summary),
  sqlc.arg(content)
);

-- name: InsertStoryPublication :one
INSERT INTO "story_publication" (
  id,
  story_id,
  profile_id,
  kind,
  properties,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(story_id),
  sqlc.arg(profile_id),
  sqlc.arg(kind),
  sqlc.narg(properties),
  NOW()
) RETURNING *;

-- name: UpdateStory :execrows
UPDATE "story"
SET
  slug = sqlc.arg(slug),
  status = sqlc.arg(status),
  is_featured = sqlc.arg(is_featured),
  story_picture_uri = sqlc.narg(story_picture_uri),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: UpdateStoryTx :execrows
UPDATE "story_tx"
SET
  title = sqlc.arg(title),
  summary = sqlc.arg(summary),
  content = sqlc.arg(content)
WHERE story_id = sqlc.arg(story_id)
  AND locale_code = sqlc.arg(locale_code);

-- name: UpsertStoryTx :exec
INSERT INTO "story_tx" (
  story_id,
  locale_code,
  title,
  summary,
  content
) VALUES (
  sqlc.arg(story_id),
  sqlc.arg(locale_code),
  sqlc.arg(title),
  sqlc.arg(summary),
  sqlc.arg(content)
) ON CONFLICT (story_id, locale_code) DO UPDATE SET
  title = EXCLUDED.title,
  summary = EXCLUDED.summary,
  content = EXCLUDED.content;

-- name: RemoveStory :execrows
UPDATE "story"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: GetStoryForEdit :one
SELECT
  s.*,
  st.locale_code,
  st.title,
  st.summary,
  st.content
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = sqlc.arg(locale_code)
WHERE s.id = sqlc.arg(id)
  AND s.deleted_at IS NULL
LIMIT 1;

-- name: GetStoryOwnershipForUser :one
SELECT
  s.id,
  s.slug,
  s.author_profile_id,
  u.kind as user_kind,
  CASE
    WHEN u.kind = 'admin' THEN true
    WHEN s.author_profile_id = u.individual_profile_id THEN true
    ELSE false
  END as can_edit
FROM "story" s
LEFT JOIN "user" u ON u.id = sqlc.arg(user_id)
WHERE s.id = sqlc.arg(story_id)
  AND s.deleted_at IS NULL
LIMIT 1;

-- name: ListStoriesOfPublication :many
SELECT
  sqlc.embed(s),
  sqlc.embed(st),
  sqlc.embed(p1),
  sqlc.embed(p1t),
  pb.publications
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile" p1 ON p1.id = s.author_profile_id
  AND p1.approved_at IS NOT NULL
  AND p1.deleted_at IS NULL
  INNER JOIN "profile_tx" p1t ON p1t.profile_id = p1.id
  AND p1t.locale_code = sqlc.arg(locale_code)
  LEFT JOIN LATERAL (
    SELECT JSONB_AGG(
      JSONB_BUILD_OBJECT('profile', row_to_json(p2), 'profile_tx', row_to_json(p2t))
    ) AS "publications"
    FROM story_publication sp
      INNER JOIN "profile" p2 ON p2.id = sp.profile_id
      AND p2.approved_at IS NOT NULL
      AND p2.deleted_at IS NULL
      INNER JOIN "profile_tx" p2t ON p2t.profile_id = p2.id
      AND p2t.locale_code = sqlc.arg(locale_code)
    WHERE sp.story_id = s.id
      AND (sqlc.narg(filter_publication_profile_id)::CHAR(26) IS NULL OR sp.profile_id = sqlc.narg(filter_publication_profile_id)::CHAR(26))
      AND sp.deleted_at IS NULL
  ) pb ON TRUE
WHERE
  pb.publications IS NOT NULL
  AND (sqlc.narg(filter_kind)::TEXT IS NULL OR s.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
  AND (sqlc.narg(filter_author_profile_id)::CHAR(26) IS NULL OR s.author_profile_id = sqlc.narg(filter_author_profile_id)::CHAR(26))
  AND s.deleted_at IS NULL
ORDER BY s.created_at DESC;
