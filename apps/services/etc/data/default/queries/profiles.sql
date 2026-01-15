-- name: GetProfileIDBySlug :one
SELECT id
FROM "profile"
WHERE slug = sqlc.arg(slug)
  AND deleted_at IS NULL
LIMIT 1;

-- name: CheckProfileSlugExists :one
SELECT EXISTS(
  SELECT 1 FROM "profile"
  WHERE slug = sqlc.arg(slug)
    AND deleted_at IS NULL
) AS exists;

-- name: GetProfileIDByCustomDomain :one
SELECT id
FROM "profile"
WHERE custom_domain = sqlc.arg(custom_domain)
  AND deleted_at IS NULL
LIMIT 1;

-- name: GetProfileByID :one
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.id = sqlc.arg(id)
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: GetProfileTxByID :many
SELECT sqlc.embed(pt)
FROM "profile_tx" pt
WHERE pt.profile_id = sqlc.arg(id);

-- name: ListProfiles :many
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE (sqlc.narg(filter_kind)::TEXT IS NULL OR p.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
  AND p.approved_at IS NOT NULL
  AND p.deleted_at IS NULL;

-- name: CreateProfile :exec
INSERT INTO "profile" (id, slug, kind, custom_domain, profile_picture_uri, pronouns, properties)
VALUES (sqlc.arg(id), sqlc.arg(slug), sqlc.arg(kind), sqlc.narg(custom_domain), sqlc.narg(profile_picture_uri), sqlc.narg(pronouns), sqlc.narg(properties));

-- name: CreateProfileTx :exec
INSERT INTO "profile_tx" (profile_id, locale_code, title, description, properties)
VALUES (sqlc.arg(profile_id), sqlc.arg(locale_code), sqlc.arg(title), sqlc.arg(description), sqlc.narg(properties));

-- name: UpdateProfile :execrows
UPDATE "profile"
SET
  profile_picture_uri = sqlc.narg(profile_picture_uri),
  pronouns = sqlc.narg(pronouns),
  properties = sqlc.narg(properties),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: UpdateProfileTx :execrows
UPDATE "profile_tx"
SET
  title = sqlc.arg(title),
  description = sqlc.arg(description),
  properties = sqlc.narg(properties)
WHERE profile_id = sqlc.arg(profile_id)
  AND locale_code = sqlc.arg(locale_code);

-- name: UpsertProfileTx :exec
INSERT INTO "profile_tx" (profile_id, locale_code, title, description, properties)
VALUES (sqlc.arg(profile_id), sqlc.arg(locale_code), sqlc.arg(title), sqlc.arg(description), sqlc.narg(properties))
ON CONFLICT (profile_id, locale_code)
DO UPDATE SET
  title = EXCLUDED.title,
  description = EXCLUDED.description,
  properties = EXCLUDED.properties;

-- name: GetUserProfilePermissions :many
SELECT
  p.id,
  p.slug,
  p.kind as profile_kind,
  COALESCE(pm.kind, '') as membership_kind,
  u.kind as user_kind
FROM "profile" p
LEFT JOIN "profile_membership" pm ON p.id = pm.profile_id AND pm.deleted_at IS NULL
LEFT JOIN "profile" up ON pm.member_profile_id = up.id AND up.deleted_at IS NULL
LEFT JOIN "user" u ON up.id = u.individual_profile_id
WHERE u.id = sqlc.arg(user_id)
  AND p.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW());

-- name: GetProfileOwnershipForUser :one
SELECT
  p.id,
  p.slug,
  p.kind as profile_kind,
  u.kind as user_kind,
  CASE
    WHEN u.kind = 'admin' THEN true
    WHEN p.kind = 'individual' AND u.individual_profile_id = p.id THEN true
    WHEN pm.kind IN ('owner', 'lead') THEN true
    ELSE false
  END as can_edit
FROM "profile" p
LEFT JOIN "profile_membership" pm ON p.id = pm.profile_id
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
LEFT JOIN "profile" up ON pm.member_profile_id = up.id AND up.deleted_at IS NULL
LEFT JOIN "user" u ON up.id = u.individual_profile_id
WHERE u.id = sqlc.arg(user_id)
  AND p.slug = sqlc.arg(profile_slug)
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: RemoveProfile :execrows
UPDATE "profile"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListProfileLinksForKind :many
SELECT pl.*
FROM "profile_link" pl
  INNER JOIN "profile" p ON p.id = pl.profile_id
  AND p.deleted_at IS NULL
WHERE pl.kind = sqlc.arg(kind)
  AND pl.deleted_at IS NULL
ORDER BY pl."order";

-- name: ListProfilePagesByProfileID :many
SELECT pp.*, ppt.*
FROM "profile_page" pp
  INNER JOIN "profile_page_tx" ppt ON ppt.profile_page_id = pp.id
  AND ppt.locale_code = sqlc.arg(locale_code)
WHERE pp.profile_id = sqlc.arg(profile_id)
  AND pp.deleted_at IS NULL
ORDER BY pp."order";

-- name: GetProfilePageByProfileIDAndSlug :one
SELECT pp.*, ppt.*
FROM "profile_page" pp
  INNER JOIN "profile_page_tx" ppt ON ppt.profile_page_id = pp.id
  AND ppt.locale_code = sqlc.arg(locale_code)
WHERE pp.profile_id = sqlc.arg(profile_id) AND pp.slug = sqlc.arg(page_slug) AND pp.deleted_at IS NULL
ORDER BY pp."order";

-- name: GetProfilePage :one
SELECT *
FROM "profile_page"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: CreateProfilePage :one
INSERT INTO "profile_page" (
  id,
  slug,
  profile_id,
  "order",
  cover_picture_uri,
  published_at,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(slug),
  sqlc.arg(profile_id),
  sqlc.arg(page_order),
  sqlc.narg(cover_picture_uri),
  sqlc.narg(published_at),
  NOW()
) RETURNING *;

-- name: CreateProfilePageTx :exec
INSERT INTO "profile_page_tx" (
  profile_page_id,
  locale_code,
  title,
  summary,
  content
) VALUES (
  sqlc.arg(profile_page_id),
  sqlc.arg(locale_code),
  sqlc.arg(title),
  sqlc.arg(summary),
  sqlc.arg(content)
);

-- name: UpdateProfilePage :execrows
UPDATE "profile_page"
SET
  slug = sqlc.arg(slug),
  "order" = sqlc.arg(page_order),
  cover_picture_uri = sqlc.narg(cover_picture_uri),
  published_at = sqlc.narg(published_at),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: UpdateProfilePageTx :execrows
UPDATE "profile_page_tx"
SET
  title = sqlc.arg(title),
  summary = sqlc.arg(summary),
  content = sqlc.arg(content)
WHERE profile_page_id = sqlc.arg(profile_page_id)
  AND locale_code = sqlc.arg(locale_code);

-- name: UpsertProfilePageTx :exec
INSERT INTO "profile_page_tx" (
  profile_page_id,
  locale_code,
  title,
  summary,
  content
) VALUES (
  sqlc.arg(profile_page_id),
  sqlc.arg(locale_code),
  sqlc.arg(title),
  sqlc.arg(summary),
  sqlc.arg(content)
) ON CONFLICT (profile_page_id, locale_code) DO UPDATE SET
  title = EXCLUDED.title,
  summary = EXCLUDED.summary,
  content = EXCLUDED.content;

-- name: DeleteProfilePage :execrows
UPDATE "profile_page"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListProfileLinksByProfileID :many
SELECT *
FROM "profile_link"
WHERE profile_id = sqlc.arg(profile_id)
  AND is_hidden = FALSE
  AND deleted_at IS NULL
ORDER BY "order";

-- name: GetProfileLink :one
SELECT *
FROM "profile_link"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: CreateProfileLink :one
INSERT INTO "profile_link" (
  id,
  kind,
  profile_id,
  "order",
  is_managed,
  is_verified,
  is_hidden,
  remote_id,
  public_id,
  uri,
  title,
  auth_provider,
  auth_access_token_scope,
  auth_access_token,
  auth_access_token_expires_at,
  auth_refresh_token,
  auth_refresh_token_expires_at,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(kind),
  sqlc.arg(profile_id),
  sqlc.arg(link_order),
  sqlc.arg(is_managed),
  sqlc.arg(is_verified),
  sqlc.arg(is_hidden),
  sqlc.narg(remote_id),
  sqlc.narg(public_id),
  sqlc.narg(uri),
  sqlc.arg(title),
  sqlc.narg(auth_provider),
  sqlc.narg(auth_access_token_scope),
  sqlc.narg(auth_access_token),
  sqlc.narg(auth_access_token_expires_at),
  sqlc.narg(auth_refresh_token),
  sqlc.narg(auth_refresh_token_expires_at),
  NOW()
) RETURNING *;

-- name: UpdateProfileLink :execrows
UPDATE "profile_link"
SET
  kind = sqlc.arg(kind),
  "order" = sqlc.arg(link_order),
  uri = sqlc.narg(uri),
  title = sqlc.arg(title),
  is_hidden = sqlc.arg(is_hidden),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DeleteProfileLink :execrows
UPDATE "profile_link"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListProfileMemberships :many
SELECT
  sqlc.embed(pm),
  sqlc.embed(p1),
  sqlc.embed(p1t),
  sqlc.embed(p2),
  sqlc.embed(p2t)
FROM
	"profile_membership" pm
  INNER JOIN "profile" p1 ON p1.id = pm.profile_id
    AND (sqlc.narg(filter_profile_kind)::TEXT IS NULL OR p1.kind = ANY(string_to_array(sqlc.narg(filter_profile_kind)::TEXT, ',')))
    AND p1.approved_at IS NOT NULL
    AND p1.deleted_at IS NULL
  INNER JOIN "profile_tx" p1t ON p1t.profile_id = p1.id
	  AND p1t.locale_code = sqlc.arg(locale_code)
  INNER JOIN "profile" p2 ON p2.id = pm.member_profile_id
    AND (sqlc.narg(filter_member_profile_kind)::TEXT IS NULL OR p2.kind = ANY(string_to_array(sqlc.narg(filter_member_profile_kind)::TEXT, ',')))
    AND p2.approved_at IS NOT NULL
    AND p2.deleted_at IS NULL
  INNER JOIN "profile_tx" p2t ON p2t.profile_id = p2.id
	  AND p2t.locale_code = sqlc.arg(locale_code)
WHERE pm.deleted_at IS NULL
    AND (sqlc.narg(filter_profile_id)::TEXT IS NULL OR pm.profile_id = sqlc.narg(filter_profile_id)::TEXT)
    AND (sqlc.narg(filter_member_profile_id)::TEXT IS NULL OR pm.member_profile_id = sqlc.narg(filter_member_profile_id)::TEXT);

-- name: GetProfileMembershipsByMemberProfileID :many
SELECT
  pm.id as membership_id,
  pm.kind as membership_kind,
  pm.started_at,
  pm.finished_at,
  pm.properties as membership_properties,
  pm.created_at as membership_created_at,
  sqlc.embed(p),
  sqlc.embed(pt)
FROM
  "profile_membership" pm
  INNER JOIN "profile" p ON p.id = pm.profile_id
    AND p.approved_at IS NOT NULL
    AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
    AND pt.locale_code = sqlc.arg(locale_code)
WHERE
  pm.deleted_at IS NULL
  AND pm.member_profile_id = sqlc.arg(member_profile_id)
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
ORDER BY pm.created_at DESC;
