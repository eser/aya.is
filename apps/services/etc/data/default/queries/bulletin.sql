-- name: GetActiveSubscriptionsForWindow :many
-- Returns active bulletin subscriptions whose preferred_time matches the given UTC hour
-- and whose last_bulletin_at respects the frequency-based cooldown.
SELECT
  bs.id,
  bs.profile_id,
  bs.channel,
  bs.frequency,
  bs.preferred_time,
  bs.last_bulletin_at,
  bs.created_at,
  bs.updated_at,
  p.default_locale,
  p.slug AS profile_slug
FROM "bulletin_subscription" bs
  INNER JOIN "profile" p ON p.id = bs.profile_id
    AND p.deleted_at IS NULL
    AND p.approved_at IS NOT NULL
WHERE bs.deleted_at IS NULL
  AND bs.preferred_time = sqlc.arg(utc_hour)::SMALLINT
  AND (
    bs.last_bulletin_at IS NULL
    OR (bs.frequency = 'daily'   AND bs.last_bulletin_at < NOW() - INTERVAL '20 hours')
    OR (bs.frequency = 'bidaily' AND bs.last_bulletin_at < NOW() - INTERVAL '44 hours')
    OR (bs.frequency = 'weekly'  AND bs.last_bulletin_at < NOW() - INTERVAL '164 hours')
  );

-- name: GetFollowedProfileStoriesSince :many
-- Returns published stories from profiles that the given subscriber follows,
-- published since the given timestamp, with 3-tier locale fallback on both
-- story_tx and author profile_tx.
SELECT
  s.id AS story_id,
  s.slug AS story_slug,
  s.kind AS story_kind,
  s.story_picture_uri,
  st.locale_code AS story_locale_code,
  st.title AS story_title,
  st.summary AS story_summary,
  st.summary_ai AS story_summary_ai,
  ap.id AS author_profile_id,
  ap.slug AS author_profile_slug,
  ap.profile_picture_uri AS author_profile_picture_uri,
  apt.title AS author_profile_title,
  (SELECT MIN(sp2.published_at) FROM story_publication sp2 WHERE sp2.story_id = s.id AND sp2.deleted_at IS NULL) AS published_at
FROM "profile_membership" pm
  INNER JOIN "story" s ON s.author_profile_id = pm.profile_id
    AND s.deleted_at IS NULL
    AND s.kind IN ('article', 'news', 'content')
  INNER JOIN "story_tx" st ON st.story_id = s.id
    AND st.locale_code = (
      SELECT stx.locale_code FROM "story_tx" stx
      WHERE stx.story_id = s.id
      ORDER BY CASE
        WHEN stx.locale_code = sqlc.arg(locale_code) THEN 0
        WHEN stx.locale_code = (SELECT p_loc.default_locale FROM "profile" p_loc WHERE p_loc.id = s.author_profile_id) THEN 1
        ELSE 2
      END
      LIMIT 1
    )
  INNER JOIN "profile" ap ON ap.id = s.author_profile_id
    AND ap.deleted_at IS NULL
    AND ap.approved_at IS NOT NULL
  INNER JOIN "profile_tx" apt ON apt.profile_id = ap.id
    AND apt.locale_code = (
      SELECT aptf.locale_code FROM "profile_tx" aptf
      WHERE aptf.profile_id = ap.id
      ORDER BY CASE
        WHEN aptf.locale_code = sqlc.arg(locale_code) THEN 0
        WHEN aptf.locale_code = ap.default_locale THEN 1
        ELSE 2
      END
      LIMIT 1
    )
WHERE pm.member_profile_id = sqlc.arg(subscriber_profile_id)
  AND pm.kind = 'follower'
  AND pm.deleted_at IS NULL
  AND (pm.finished_at IS NULL OR pm.finished_at > NOW())
  AND EXISTS (
    SELECT 1 FROM story_publication sp
    WHERE sp.story_id = s.id
      AND sp.published_at IS NOT NULL
      AND sp.published_at > sqlc.arg(since)
      AND sp.deleted_at IS NULL
  )
  AND s.visibility = 'public'
ORDER BY published_at DESC
LIMIT sqlc.arg(max_stories);

-- name: GetBulletinSubscriptionsByProfileID :many
-- Returns all active subscriptions for a given profile.
SELECT
  bs.id,
  bs.profile_id,
  bs.channel,
  bs.frequency,
  bs.preferred_time,
  bs.last_bulletin_at,
  bs.created_at,
  bs.updated_at
FROM "bulletin_subscription" bs
WHERE bs.profile_id = sqlc.arg(profile_id)
  AND bs.deleted_at IS NULL
ORDER BY bs.created_at;

-- name: GetBulletinSubscription :one
-- Returns a single subscription by ID.
SELECT
  bs.id,
  bs.profile_id,
  bs.channel,
  bs.frequency,
  bs.preferred_time,
  bs.last_bulletin_at,
  bs.created_at,
  bs.updated_at
FROM "bulletin_subscription" bs
WHERE bs.id = sqlc.arg(id)
  AND bs.deleted_at IS NULL;

-- name: UpsertBulletinSubscription :one
-- Creates or reactivates a subscription for a profile+channel combination.
INSERT INTO "bulletin_subscription" (
  "id", "profile_id", "channel", "frequency", "preferred_time", "created_at"
) VALUES (
  sqlc.arg(id),
  sqlc.arg(profile_id),
  sqlc.arg(channel),
  sqlc.arg(frequency),
  sqlc.arg(preferred_time),
  NOW()
)
ON CONFLICT ("profile_id", "channel") WHERE "deleted_at" IS NULL
DO UPDATE SET
  frequency = EXCLUDED.frequency,
  preferred_time = EXCLUDED.preferred_time,
  deleted_at = NULL,
  updated_at = NOW()
RETURNING *;

-- name: UpdateBulletinSubscriptionPreferences :exec
-- Updates the frequency and preferred time for a subscription.
UPDATE "bulletin_subscription"
SET frequency = sqlc.arg(frequency),
    preferred_time = sqlc.arg(preferred_time),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DeleteBulletinSubscription :exec
-- Soft-deletes a subscription.
UPDATE "bulletin_subscription"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: UpdateBulletinSubscriptionLastSentAt :exec
-- Updates the last_bulletin_at timestamp after successfully sending a bulletin.
UPDATE "bulletin_subscription"
SET last_bulletin_at = NOW(),
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: GetUserEmailByIndividualProfileID :one
-- Returns the user's email for a given individual profile ID.
SELECT u.email
FROM "user" u
WHERE u.individual_profile_id = sqlc.arg(profile_id)
  AND u.deleted_at IS NULL
LIMIT 1;

-- name: DeleteBulletinSubscriptionsByProfileID :exec
-- Soft-deletes all subscriptions for a profile (used for "Don't send").
UPDATE "bulletin_subscription"
SET deleted_at = NOW()
WHERE profile_id = sqlc.arg(profile_id)
  AND deleted_at IS NULL;

-- name: CreateBulletinLog :exec
-- Records a bulletin send attempt.
INSERT INTO "bulletin_log" (
  "id", "subscription_id", "story_count", "status", "error_message", "created_at"
) VALUES (
  sqlc.arg(id),
  sqlc.arg(subscription_id),
  sqlc.arg(story_count),
  sqlc.arg(status),
  sqlc.narg(error_message),
  NOW()
);
