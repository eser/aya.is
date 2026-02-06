-- name: GetStoryIDBySlug :one
SELECT s.id
FROM "story" s
WHERE s.slug = sqlc.arg(slug)
  AND s.deleted_at IS NULL
  AND EXISTS (
    SELECT 1 FROM story_publication sp
    WHERE sp.story_id = s.id AND sp.deleted_at IS NULL
  )
LIMIT 1;

-- name: GetStoryIDBySlugForViewer :one
-- Returns story ID if:
--   1. Story has at least one active publication, OR
--   2. Viewer is admin, OR
--   3. Viewer is the author (individual profile owner)
--   4. Viewer is owner/lead/maintainer of the author profile
SELECT s.id
FROM "story" s
LEFT JOIN "user" u ON u.id = sqlc.narg(viewer_user_id)::CHAR(26)
LEFT JOIN "profile_membership" pm ON s.author_profile_id = pm.profile_id
  AND pm.member_profile_id = u.individual_profile_id
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
WHERE s.slug = sqlc.arg(slug)
  AND s.deleted_at IS NULL
  AND (
    EXISTS (SELECT 1 FROM story_publication sp WHERE sp.story_id = s.id AND sp.deleted_at IS NULL)
    OR u.kind = 'admin'
    OR s.author_profile_id = u.individual_profile_id
    OR pm.kind IN ('owner', 'lead', 'maintainer')
  )
LIMIT 1;

-- name: GetStoryIDBySlugIncludingDeleted :one
SELECT id
FROM "story"
WHERE slug = sqlc.arg(slug)
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

-- name: InsertStory :one
INSERT INTO "story" (
  id,
  author_profile_id,
  slug,
  kind,
  story_picture_uri,
  properties,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(author_profile_id),
  sqlc.arg(slug),
  sqlc.arg(kind),
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
  is_featured,
  published_at,
  properties,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(story_id),
  sqlc.arg(profile_id),
  sqlc.arg(kind),
  sqlc.arg(is_featured),
  sqlc.narg(published_at),
  sqlc.narg(properties),
  NOW()
) RETURNING *;

-- name: UpdateStory :execrows
UPDATE "story"
SET
  slug = sqlc.arg(slug),
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
-- Uses locale fallback: prefers the requested locale, falls back to any translation.
-- The returned locale_code indicates which translation was actually found.
SELECT
  s.*,
  st.locale_code,
  st.title,
  st.summary,
  st.content,
  p.slug as author_profile_slug
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = (
    SELECT stx.locale_code FROM "story_tx" stx
    WHERE stx.story_id = s.id
    ORDER BY CASE WHEN stx.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
    LIMIT 1
  )
  LEFT JOIN "profile" p ON p.id = s.author_profile_id AND p.deleted_at IS NULL
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
    WHEN pm.kind IN ('owner', 'lead', 'maintainer') THEN true
    ELSE false
  END as can_edit
FROM "story" s
LEFT JOIN "user" u ON u.id = sqlc.arg(user_id)
LEFT JOIN "profile_membership" pm ON s.author_profile_id = pm.profile_id
  AND pm.member_profile_id = u.individual_profile_id
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
WHERE s.id = sqlc.arg(story_id)
  AND s.deleted_at IS NULL
LIMIT 1;

-- name: ListStoriesOfPublication :many
-- Uses locale fallback: prefers the requested locale, falls back to any translation.
SELECT
  sqlc.embed(s),
  sqlc.embed(st),
  sqlc.embed(p1),
  sqlc.embed(p1t),
  pb.publications
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = (
    SELECT stx.locale_code FROM "story_tx" stx
    WHERE stx.story_id = s.id
    ORDER BY CASE WHEN stx.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
    LIMIT 1
  )
  LEFT JOIN "profile" p1 ON p1.id = s.author_profile_id
  AND p1.approved_at IS NOT NULL
  AND p1.deleted_at IS NULL
  INNER JOIN "profile_tx" p1t ON p1t.profile_id = p1.id
  AND p1t.locale_code = (
    SELECT ptx.locale_code FROM "profile_tx" ptx
    WHERE ptx.profile_id = p1.id
    ORDER BY CASE WHEN ptx.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
    LIMIT 1
  )
  LEFT JOIN LATERAL (
    SELECT JSONB_AGG(
      JSONB_BUILD_OBJECT('profile', row_to_json(p2), 'profile_tx', row_to_json(p2t))
    ) AS "publications"
    FROM story_publication sp
      INNER JOIN "profile" p2 ON p2.id = sp.profile_id
      AND p2.approved_at IS NOT NULL
      AND p2.deleted_at IS NULL
      INNER JOIN "profile_tx" p2t ON p2t.profile_id = p2.id
      AND p2t.locale_code = (
        SELECT ptx2.locale_code FROM "profile_tx" ptx2
        WHERE ptx2.profile_id = p2.id
        ORDER BY CASE WHEN ptx2.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
        LIMIT 1
      )
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

-- name: ListStoriesByAuthorProfileID :many
-- Lists all stories authored by a profile, including unpublished ones.
-- Uses locale fallback: prefers the requested locale, falls back to any translation.
-- Publications are included as optional data (LEFT JOIN).
SELECT
  sqlc.embed(s),
  sqlc.embed(st),
  sqlc.embed(p1),
  sqlc.embed(p1t),
  pb.publications
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
  AND st.locale_code = (
    SELECT stx.locale_code FROM "story_tx" stx
    WHERE stx.story_id = s.id
    ORDER BY CASE WHEN stx.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
    LIMIT 1
  )
  LEFT JOIN "profile" p1 ON p1.id = s.author_profile_id
  AND p1.approved_at IS NOT NULL
  AND p1.deleted_at IS NULL
  INNER JOIN "profile_tx" p1t ON p1t.profile_id = p1.id
  AND p1t.locale_code = (
    SELECT ptx.locale_code FROM "profile_tx" ptx
    WHERE ptx.profile_id = p1.id
    ORDER BY CASE WHEN ptx.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
    LIMIT 1
  )
  LEFT JOIN LATERAL (
    SELECT JSONB_AGG(
      JSONB_BUILD_OBJECT('profile', row_to_json(p2), 'profile_tx', row_to_json(p2t))
    ) AS "publications"
    FROM story_publication sp
      INNER JOIN "profile" p2 ON p2.id = sp.profile_id
      AND p2.approved_at IS NOT NULL
      AND p2.deleted_at IS NULL
      INNER JOIN "profile_tx" p2t ON p2t.profile_id = p2.id
      AND p2t.locale_code = (
        SELECT ptx2.locale_code FROM "profile_tx" ptx2
        WHERE ptx2.profile_id = p2.id
        ORDER BY CASE WHEN ptx2.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
        LIMIT 1
      )
    WHERE sp.story_id = s.id
      AND sp.deleted_at IS NULL
  ) pb ON TRUE
WHERE s.author_profile_id = sqlc.arg(author_profile_id)::CHAR(26)
  AND (sqlc.narg(filter_kind)::TEXT IS NULL OR s.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
  AND s.deleted_at IS NULL
ORDER BY s.created_at DESC;

-- name: SearchStories :many
SELECT
  s.id,
  s.slug,
  s.kind,
  s.story_picture_uri,
  s.author_profile_id,
  st.title,
  st.summary,
  p.slug as author_slug,
  pt.title as author_title,
  ts_rank(st.search_vector, plainto_tsquery(locale_to_regconfig(sqlc.arg(locale_code)), sqlc.arg(query))) as rank
FROM "story" s
  INNER JOIN "story_tx" st ON st.story_id = s.id
    AND st.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile" p ON p.id = s.author_profile_id AND p.deleted_at IS NULL
  LEFT JOIN "profile_tx" pt ON pt.profile_id = p.id
    AND pt.locale_code = sqlc.arg(locale_code)
WHERE st.search_vector @@ plainto_tsquery(locale_to_regconfig(sqlc.arg(locale_code)), sqlc.arg(query))
  AND s.deleted_at IS NULL
  AND EXISTS (
    SELECT 1 FROM story_publication sp
    WHERE sp.story_id = s.id AND sp.deleted_at IS NULL
  )
  AND (sqlc.narg(filter_profile_slug)::TEXT IS NULL OR p.slug = sqlc.narg(filter_profile_slug)::TEXT)
ORDER BY rank DESC
LIMIT sqlc.arg(limit_count);

-- name: ListStoryPublications :many
-- Lists all publications for a story with profile info (for publish popup)
SELECT
  sp.id,
  sp.story_id,
  sp.profile_id,
  sp.kind,
  sp.is_featured,
  sp.published_at,
  sp.created_at,
  p.slug as profile_slug,
  pt.title as profile_title,
  p.profile_picture_uri,
  p.kind as profile_kind
FROM story_publication sp
  INNER JOIN "profile" p ON p.id = sp.profile_id AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id AND pt.locale_code = sqlc.arg(locale_code)
WHERE sp.story_id = sqlc.arg(story_id)
  AND sp.deleted_at IS NULL
ORDER BY sp.created_at;

-- name: GetStoryPublicationProfileID :one
-- Returns the profile_id for a specific publication (used for auth checks).
SELECT profile_id
FROM story_publication
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateStoryPublication :execrows
UPDATE story_publication
SET
  is_featured = sqlc.arg(is_featured),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: RemoveStoryPublication :execrows
UPDATE story_publication
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: CountStoryPublications :one
SELECT COUNT(*) as count
FROM story_publication
WHERE story_id = sqlc.arg(story_id)
  AND deleted_at IS NULL;

-- name: GetStoryFirstPublishedAt :one
SELECT MIN(published_at) as first_published_at
FROM story_publication
WHERE story_id = sqlc.arg(story_id)
  AND deleted_at IS NULL
  AND published_at IS NOT NULL;

-- name: GetUserMembershipForProfile :one
-- Returns the membership kind a user has for a specific profile.
-- Used to verify a user has access to publish to a target profile.
-- Returns:
--   'admin' if the user is an admin
--   'owner' if the target profile is the user's individual profile
--   the membership kind (owner/lead/maintainer/contributor) if they have membership
--   '' if no membership
SELECT
  CAST(CASE
    WHEN u.kind = 'admin' THEN 'admin'
    WHEN u.individual_profile_id = sqlc.arg(profile_id)::CHAR(26) THEN 'owner'
    ELSE COALESCE(pm.kind, '')
  END AS TEXT) as membership_kind
FROM "user" u
LEFT JOIN "profile_membership" pm ON pm.profile_id = sqlc.arg(profile_id)::CHAR(26)
  AND pm.member_profile_id = u.individual_profile_id
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
WHERE u.id = sqlc.arg(user_id)
  AND u.deleted_at IS NULL
LIMIT 1;
