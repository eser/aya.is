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

-- name: CheckProfileSlugExistsIncludingDeleted :one
SELECT EXISTS(
  SELECT 1 FROM "profile"
  WHERE slug = sqlc.arg(slug)
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

-- name: GetProfilesByIDs :many
SELECT sqlc.embed(p), sqlc.embed(pt)
FROM "profile" p
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.id = ANY(sqlc.arg(ids)::TEXT[])
  AND p.deleted_at IS NULL;

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
  profile_picture_uri = COALESCE(sqlc.narg(profile_picture_uri), profile_picture_uri),
  pronouns = COALESCE(sqlc.narg(pronouns), pronouns),
  properties = COALESCE(sqlc.narg(properties), properties),
  hide_relations = COALESCE(sqlc.narg(hide_relations), hide_relations),
  hide_links = COALESCE(sqlc.narg(hide_links), hide_links),
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
CROSS JOIN "user" u
LEFT JOIN "profile_membership" pm ON p.id = pm.profile_id
  AND pm.member_profile_id = u.individual_profile_id
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
WHERE u.id = sqlc.arg(user_id)
  AND p.slug = sqlc.arg(profile_slug)
  AND p.deleted_at IS NULL
  AND u.deleted_at IS NULL
LIMIT 1;

-- name: RemoveProfile :execrows
UPDATE "profile"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListProfileLinksForKind :many
SELECT
  pl.*,
  COALESCE(plt.profile_link_id, plt_en.profile_link_id, pl.id) as profile_link_id,
  COALESCE(plt.locale_code, plt_en.locale_code, 'en') as locale_code,
  COALESCE(plt.title, plt_en.title, pl.kind) as title,
  plt."group" as "group",
  plt.description as description
FROM "profile_link" pl
  INNER JOIN "profile" p ON p.id = pl.profile_id
    AND p.deleted_at IS NULL
  LEFT JOIN "profile_link_tx" plt ON plt.profile_link_id = pl.id
    AND plt.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile_link_tx" plt_en ON plt_en.profile_link_id = pl.id
    AND plt_en.locale_code = 'en'
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

-- name: CheckPageSlugExistsIncludingDeleted :one
SELECT EXISTS(
  SELECT 1 FROM "profile_page"
  WHERE profile_id = sqlc.arg(profile_id)
    AND slug = sqlc.arg(page_slug)
) AS exists;

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
SELECT
  pl.*,
  COALESCE(plt.profile_link_id, plt_en.profile_link_id, pl.id) as profile_link_id,
  COALESCE(plt.locale_code, plt_en.locale_code, 'en') as locale_code,
  COALESCE(plt.title, plt_en.title, pl.kind) as title,
  COALESCE(plt.icon, plt_en.icon, '') as icon,
  plt."group" as "group",
  plt.description as description
FROM "profile_link" pl
  LEFT JOIN "profile_link_tx" plt ON plt.profile_link_id = pl.id
    AND plt.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile_link_tx" plt_en ON plt_en.profile_link_id = pl.id
    AND plt_en.locale_code = 'en'
WHERE pl.profile_id = sqlc.arg(profile_id)
  AND pl.deleted_at IS NULL
ORDER BY pl."order";

-- name: GetProfileLink :one
SELECT
  pl.*,
  COALESCE(plt.profile_link_id, plt_en.profile_link_id, pl.id) as profile_link_id,
  COALESCE(plt.locale_code, plt_en.locale_code, 'en') as locale_code,
  COALESCE(plt.title, plt_en.title, pl.kind) as title,
  COALESCE(plt.icon, plt_en.icon, '') as icon,
  plt."group" as "group",
  plt.description as description
FROM "profile_link" pl
  LEFT JOIN "profile_link_tx" plt ON plt.profile_link_id = pl.id
    AND plt.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile_link_tx" plt_en ON plt_en.profile_link_id = pl.id
    AND plt_en.locale_code = 'en'
WHERE pl.id = sqlc.arg(id)
  AND pl.deleted_at IS NULL;

-- name: CreateProfileLink :one
INSERT INTO "profile_link" (
  id,
  kind,
  profile_id,
  "order",
  is_managed,
  is_verified,
  is_featured,
  visibility,
  remote_id,
  public_id,
  uri,
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
  sqlc.arg(is_featured),
  sqlc.arg(visibility),
  sqlc.narg(remote_id),
  sqlc.narg(public_id),
  sqlc.narg(uri),
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
  is_featured = sqlc.arg(is_featured),
  visibility = sqlc.arg(visibility),
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

-- name: SearchProfiles :many
SELECT
  p.id,
  p.slug,
  p.kind,
  p.profile_picture_uri,
  pt.title,
  pt.description,
  ts_rank(pt.search_vector, plainto_tsquery(locale_to_regconfig(sqlc.arg(locale_code)), sqlc.arg(query))) as rank
FROM "profile" p
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
    AND pt.locale_code = sqlc.arg(locale_code)
WHERE pt.search_vector @@ plainto_tsquery(locale_to_regconfig(sqlc.arg(locale_code)), sqlc.arg(query))
  AND p.approved_at IS NOT NULL
  AND p.deleted_at IS NULL
  AND (sqlc.narg(filter_profile_slug)::TEXT IS NULL OR p.slug = sqlc.narg(filter_profile_slug)::TEXT)
ORDER BY rank DESC
LIMIT sqlc.arg(limit_count);

-- name: SearchProfilePages :many
SELECT
  pp.id,
  pp.slug,
  pp.profile_id,
  pp.cover_picture_uri,
  ppt.title,
  ppt.summary,
  p.slug as profile_slug,
  pt.title as profile_title,
  ts_rank(ppt.search_vector, plainto_tsquery(locale_to_regconfig(sqlc.arg(locale_code)), sqlc.arg(query))) as rank
FROM "profile_page" pp
  INNER JOIN "profile_page_tx" ppt ON ppt.profile_page_id = pp.id
    AND ppt.locale_code = sqlc.arg(locale_code)
  INNER JOIN "profile" p ON p.id = pp.profile_id AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
    AND pt.locale_code = sqlc.arg(locale_code)
WHERE ppt.search_vector @@ plainto_tsquery(locale_to_regconfig(sqlc.arg(locale_code)), sqlc.arg(query))
  AND pp.deleted_at IS NULL
  AND p.approved_at IS NOT NULL
  AND (sqlc.narg(filter_profile_slug)::TEXT IS NULL OR p.slug = sqlc.narg(filter_profile_slug)::TEXT)
ORDER BY rank DESC
LIMIT sqlc.arg(limit_count);

-- name: GetProfileLinkByRemoteID :one
SELECT *
FROM "profile_link"
WHERE profile_id = sqlc.arg(profile_id)
  AND kind = sqlc.arg(kind)
  AND remote_id = sqlc.arg(remote_id)
  AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateProfileLinkOAuthTokens :execrows
UPDATE "profile_link"
SET
  public_id = sqlc.narg(public_id),
  uri = sqlc.narg(uri),
  auth_access_token = sqlc.narg(auth_access_token),
  auth_access_token_expires_at = sqlc.narg(auth_access_token_expires_at),
  auth_refresh_token = sqlc.narg(auth_refresh_token),
  auth_access_token_scope = sqlc.narg(auth_access_token_scope),
  is_verified = TRUE,
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: GetMaxProfileLinkOrder :one
SELECT COALESCE(MAX("order"), 0) as max_order
FROM "profile_link"
WHERE profile_id = sqlc.arg(profile_id)
  AND deleted_at IS NULL;

-- name: GetMembershipBetweenProfiles :one
SELECT pm.kind
FROM "profile_membership" pm
WHERE pm.profile_id = sqlc.arg(profile_id)
  AND pm.member_profile_id = sqlc.arg(member_profile_id)
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
LIMIT 1;

-- name: CreateProfileMembership :exec
INSERT INTO "profile_membership" (
  "id",
  "profile_id",
  "member_profile_id",
  "kind",
  "properties",
  "started_at",
  "created_at"
) VALUES (
  sqlc.arg(id),
  sqlc.arg(profile_id),
  sqlc.narg(member_profile_id),
  sqlc.arg(kind),
  sqlc.narg(properties),
  NOW(),
  NOW()
);

-- name: ListProfileMembershipsForSettings :many
SELECT
  pm.id,
  pm.profile_id,
  pm.member_profile_id,
  pm.kind,
  pm.properties,
  pm.started_at,
  pm.finished_at,
  pm.created_at,
  pm.updated_at,
  sqlc.embed(mp),
  sqlc.embed(mpt)
FROM "profile_membership" pm
INNER JOIN "profile" mp ON mp.id = pm.member_profile_id
  AND mp.deleted_at IS NULL
INNER JOIN "profile_tx" mpt ON mpt.profile_id = mp.id
  AND mpt.locale_code = sqlc.arg(locale_code)
WHERE pm.profile_id = sqlc.arg(profile_id)
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
ORDER BY
  CASE pm.kind
    WHEN 'owner' THEN 1
    WHEN 'lead' THEN 2
    WHEN 'maintainer' THEN 3
    WHEN 'contributor' THEN 4
    WHEN 'sponsor' THEN 5
    WHEN 'follower' THEN 6
    ELSE 7
  END,
  pm.created_at ASC;

-- name: GetProfileMembershipByID :one
SELECT
  pm.id,
  pm.profile_id,
  pm.member_profile_id,
  pm.kind,
  pm.properties,
  pm.started_at,
  pm.finished_at,
  pm.created_at,
  pm.updated_at
FROM "profile_membership" pm
WHERE pm.id = sqlc.arg(id)
  AND pm.deleted_at IS NULL;

-- name: UpdateProfileMembership :execrows
UPDATE "profile_membership"
SET
  kind = sqlc.arg(kind),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DeleteProfileMembership :execrows
UPDATE "profile_membership"
SET
  deleted_at = NOW(),
  finished_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: CountProfileOwners :one
SELECT COUNT(*) as owner_count
FROM "profile_membership" pm
WHERE pm.profile_id = sqlc.arg(profile_id)
  AND pm.kind = 'owner'
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW());

-- name: SearchUsersForMembership :many
SELECT
  u.id as user_id,
  u.email,
  u.name,
  u.individual_profile_id,
  sqlc.embed(p),
  sqlc.embed(pt)
FROM "user" u
INNER JOIN "profile" p ON p.id = u.individual_profile_id
  AND p.deleted_at IS NULL
INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE u.deleted_at IS NULL
  AND u.individual_profile_id IS NOT NULL
  AND (
    u.email ILIKE '%' || sqlc.arg(query) || '%'
    OR u.name ILIKE '%' || sqlc.arg(query) || '%'
    OR p.slug ILIKE '%' || sqlc.arg(query) || '%'
    OR pt.title ILIKE '%' || sqlc.arg(query) || '%'
  )
  -- Exclude users already members of this profile
  AND NOT EXISTS (
    SELECT 1 FROM "profile_membership" pm
    WHERE pm.profile_id = sqlc.arg(profile_id)
      AND pm.member_profile_id = u.individual_profile_id
      AND pm.deleted_at IS NULL
      AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
  )
ORDER BY u.name ASC
LIMIT 10;

-- name: ListAllProfilesForAdmin :many
SELECT
  p.id,
  p.slug,
  p.kind,
  p.custom_domain,
  p.profile_picture_uri,
  p.pronouns,
  p.properties,
  p.created_at,
  p.updated_at,
  p.points,
  COALESCE(pt.title, '') as title,
  COALESCE(pt.description, '') as description,
  pt.profile_id IS NOT NULL as has_translation
FROM "profile" p
  LEFT JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.deleted_at IS NULL
  AND (sqlc.narg(filter_kind)::TEXT IS NULL OR p.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')))
ORDER BY p.created_at DESC
LIMIT sqlc.arg(limit_count)
OFFSET sqlc.arg(offset_count);

-- name: CountAllProfilesForAdmin :one
SELECT COUNT(*) as count
FROM "profile" p
WHERE p.deleted_at IS NULL
  AND (sqlc.narg(filter_kind)::TEXT IS NULL OR p.kind = ANY(string_to_array(sqlc.narg(filter_kind)::TEXT, ',')));

-- name: GetAdminProfileBySlug :one
SELECT
  p.id,
  p.slug,
  p.kind,
  p.custom_domain,
  p.profile_picture_uri,
  p.pronouns,
  p.properties,
  p.created_at,
  p.updated_at,
  p.points,
  COALESCE(pt.title, '') as title,
  COALESCE(pt.description, '') as description,
  pt.profile_id IS NOT NULL as has_translation
FROM "profile" p
  LEFT JOIN "profile_tx" pt ON pt.profile_id = p.id
  AND pt.locale_code = sqlc.arg(locale_code)
WHERE p.slug = sqlc.arg(slug)
  AND p.deleted_at IS NULL
LIMIT 1;

-- name: ListFeaturedProfileLinksByProfileID :many
SELECT
  pl.id,
  pl.kind,
  pl.public_id,
  pl.uri,
  pl.is_verified,
  pl.is_managed,
  pl.is_featured,
  pl.visibility,
  COALESCE(plt.title, plt_en.title, pl.kind) as title,
  COALESCE(plt.icon, plt_en.icon, '') as icon,
  COALESCE(plt."group", plt_en."group", '') as "group",
  COALESCE(plt.description, plt_en.description, '') as description
FROM "profile_link" pl
  LEFT JOIN "profile_link_tx" plt ON plt.profile_link_id = pl.id
    AND plt.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile_link_tx" plt_en ON plt_en.profile_link_id = pl.id
    AND plt_en.locale_code = 'en'
WHERE pl.profile_id = sqlc.arg(profile_id)
  AND pl.is_featured = TRUE
  AND pl.deleted_at IS NULL
ORDER BY pl."order";

-- name: ListAllProfileLinksByProfileID :many
SELECT
  pl.id,
  pl.kind,
  pl.public_id,
  pl.uri,
  pl.is_verified,
  pl.is_managed,
  pl.is_featured,
  pl.visibility,
  COALESCE(plt.title, plt_en.title, pl.kind) as title,
  COALESCE(plt.icon, plt_en.icon, '') as icon,
  COALESCE(plt."group", plt_en."group", '') as "group",
  COALESCE(plt.description, plt_en.description, '') as description
FROM "profile_link" pl
  LEFT JOIN "profile_link_tx" plt ON plt.profile_link_id = pl.id
    AND plt.locale_code = sqlc.arg(locale_code)
  LEFT JOIN "profile_link_tx" plt_en ON plt_en.profile_link_id = pl.id
    AND plt_en.locale_code = 'en'
WHERE pl.profile_id = sqlc.arg(profile_id)
  AND pl.deleted_at IS NULL
ORDER BY pl."order";

-- name: GetProfileLinkTx :one
SELECT *
FROM "profile_link_tx"
WHERE profile_link_id = sqlc.arg(profile_link_id)
  AND locale_code = sqlc.arg(locale_code)
LIMIT 1;

-- name: CreateProfileLinkTx :exec
INSERT INTO "profile_link_tx" (
  profile_link_id,
  locale_code,
  title,
  icon,
  "group",
  description
) VALUES (
  sqlc.arg(profile_link_id),
  sqlc.arg(locale_code),
  sqlc.arg(title),
  sqlc.narg(icon),
  sqlc.narg(link_group),
  sqlc.narg(description)
);

-- name: UpdateProfileLinkTx :execrows
UPDATE "profile_link_tx"
SET
  title = sqlc.arg(title),
  icon = sqlc.narg(icon),
  "group" = sqlc.narg(link_group),
  description = sqlc.narg(description)
WHERE profile_link_id = sqlc.arg(profile_link_id)
  AND locale_code = sqlc.arg(locale_code);

-- name: UpsertProfileLinkTx :exec
INSERT INTO "profile_link_tx" (
  profile_link_id,
  locale_code,
  title,
  icon,
  "group",
  description
) VALUES (
  sqlc.arg(profile_link_id),
  sqlc.arg(locale_code),
  sqlc.arg(title),
  sqlc.narg(icon),
  sqlc.narg(link_group),
  sqlc.narg(description)
) ON CONFLICT (profile_link_id, locale_code) DO UPDATE SET
  title = EXCLUDED.title,
  icon = EXCLUDED.icon,
  "group" = EXCLUDED."group",
  description = EXCLUDED.description;
