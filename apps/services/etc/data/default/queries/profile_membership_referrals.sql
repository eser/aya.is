-- name: CreateProfileMembershipReferral :one
INSERT INTO "profile_membership_referral" (
  id, profile_id, referred_profile_id, referrer_membership_id, status, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(profile_id),
  sqlc.arg(referred_profile_id),
  sqlc.arg(referrer_membership_id),
  'voting',
  NOW()
) RETURNING *;

-- name: GetProfileMembershipReferralByID :one
SELECT * FROM "profile_membership_referral"
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: GetProfileMembershipReferralByProfileAndReferred :one
SELECT * FROM "profile_membership_referral"
WHERE profile_id = sqlc.arg(profile_id)
  AND referred_profile_id = sqlc.arg(referred_profile_id)
  AND deleted_at IS NULL;

-- name: ListProfileMembershipReferralsByProfileID :many
SELECT
  pmr.*,
  ref_p.slug AS referrer_profile_slug,
  ref_p.kind AS referrer_profile_kind,
  ref_p.profile_picture_uri AS referrer_profile_picture_uri,
  ref_pt.title AS referrer_profile_title,
  tgt_p.slug AS referred_profile_slug,
  tgt_p.kind AS referred_profile_kind,
  tgt_p.profile_picture_uri AS referred_profile_picture_uri,
  tgt_pt.title AS referred_profile_title,
  (SELECT COUNT(*) FROM "profile_membership_referral_vote" v
   WHERE v.profile_membership_referral_id = pmr.id)::BIGINT AS total_votes,
  COALESCE((SELECT AVG(v.score)::NUMERIC(3,2) FROM "profile_membership_referral_vote" v
   WHERE v.profile_membership_referral_id = pmr.id), 0)::TEXT AS average_score,
  COALESCE((SELECT v.score FROM "profile_membership_referral_vote" v
   WHERE v.profile_membership_referral_id = pmr.id
     AND v.voter_membership_id = sqlc.narg(viewer_membership_id)
  ), 0)::SMALLINT AS viewer_vote_score,
  (SELECT v.comment FROM "profile_membership_referral_vote" v
   WHERE v.profile_membership_referral_id = pmr.id
     AND v.voter_membership_id = sqlc.narg(viewer_membership_id)
  ) AS viewer_vote_comment
FROM "profile_membership_referral" pmr
  INNER JOIN "profile_membership" ref_pm ON ref_pm.id = pmr.referrer_membership_id
  INNER JOIN "profile" ref_p ON ref_p.id = ref_pm.member_profile_id AND ref_p.deleted_at IS NULL
  INNER JOIN "profile_tx" ref_pt ON ref_pt.profile_id = ref_p.id
    AND ref_pt.locale_code = (
      SELECT rptf.locale_code FROM "profile_tx" rptf
      WHERE rptf.profile_id = ref_p.id
      ORDER BY CASE
        WHEN rptf.locale_code = sqlc.arg(locale_code) THEN 0
        WHEN rptf.locale_code = ref_p.default_locale THEN 1
        ELSE 2
      END
      LIMIT 1
    )
  INNER JOIN "profile" tgt_p ON tgt_p.id = pmr.referred_profile_id AND tgt_p.deleted_at IS NULL
  INNER JOIN "profile_tx" tgt_pt ON tgt_pt.profile_id = tgt_p.id
    AND tgt_pt.locale_code = (
      SELECT tptf.locale_code FROM "profile_tx" tptf
      WHERE tptf.profile_id = tgt_p.id
      ORDER BY CASE
        WHEN tptf.locale_code = sqlc.arg(locale_code) THEN 0
        WHEN tptf.locale_code = tgt_p.default_locale THEN 1
        ELSE 2
      END
      LIMIT 1
    )
WHERE pmr.profile_id = sqlc.arg(profile_id)
  AND pmr.deleted_at IS NULL
ORDER BY pmr.created_at DESC;

-- name: UpsertReferralVote :one
INSERT INTO "profile_membership_referral_vote" (
  id, profile_membership_referral_id, voter_membership_id, score, comment, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(profile_membership_referral_id),
  sqlc.arg(voter_membership_id),
  sqlc.arg(score),
  sqlc.narg(comment),
  NOW()
) ON CONFLICT (profile_membership_referral_id, voter_membership_id)
DO UPDATE SET
  score = EXCLUDED.score,
  comment = EXCLUDED.comment,
  updated_at = NOW()
RETURNING *;

-- name: ListReferralVotes :many
SELECT
  v.*,
  vp.slug AS voter_profile_slug,
  vp.kind AS voter_profile_kind,
  vp.profile_picture_uri AS voter_profile_picture_uri,
  vpt.title AS voter_profile_title
FROM "profile_membership_referral_vote" v
  INNER JOIN "profile_membership" vm ON vm.id = v.voter_membership_id
  INNER JOIN "profile" vp ON vp.id = vm.member_profile_id AND vp.deleted_at IS NULL
  INNER JOIN "profile_tx" vpt ON vpt.profile_id = vp.id
    AND vpt.locale_code = (
      SELECT vptf.locale_code FROM "profile_tx" vptf
      WHERE vptf.profile_id = vp.id
      ORDER BY CASE
        WHEN vptf.locale_code = sqlc.arg(locale_code) THEN 0
        WHEN vptf.locale_code = vp.default_locale THEN 1
        ELSE 2
      END
      LIMIT 1
    )
WHERE v.profile_membership_referral_id = sqlc.arg(referral_id)
ORDER BY v.created_at DESC;

-- name: GetReferralVoteBreakdown :many
SELECT score, COUNT(*)::BIGINT AS count
FROM "profile_membership_referral_vote"
WHERE profile_membership_referral_id = sqlc.arg(referral_id)
GROUP BY score
ORDER BY score;

-- name: UpdateReferralVoteCount :exec
UPDATE "profile_membership_referral"
SET vote_count = (
  SELECT COUNT(*) FROM "profile_membership_referral_vote"
  WHERE profile_membership_referral_id = sqlc.arg(id)
), updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: InsertReferralTeam :one
INSERT INTO "profile_membership_referral_team" (
  id, profile_membership_referral_id, profile_team_id, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(profile_membership_referral_id),
  sqlc.arg(profile_team_id),
  NOW()
) RETURNING *;

-- name: ListReferralTeams :many
SELECT pt.* FROM "profile_team" pt
JOIN "profile_membership_referral_team" pmrt
  ON pmrt.profile_team_id = pt.id AND pmrt.deleted_at IS NULL
WHERE pmrt.profile_membership_referral_id = sqlc.arg(referral_id)
  AND pt.deleted_at IS NULL
ORDER BY pt.name ASC;

-- name: SoftDeleteReferral :execrows
UPDATE "profile_membership_referral"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;
