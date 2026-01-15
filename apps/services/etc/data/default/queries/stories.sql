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
