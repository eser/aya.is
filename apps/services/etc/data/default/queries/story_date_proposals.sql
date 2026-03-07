-- name: InsertStoryDateProposal :one
INSERT INTO "story_date_proposal" (
  id, story_id, proposer_profile_id, datetime_start, datetime_end, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(story_id),
  sqlc.arg(proposer_profile_id),
  sqlc.arg(datetime_start),
  sqlc.narg(datetime_end),
  NOW()
) RETURNING *;

-- name: GetStoryDateProposal :one
SELECT * FROM "story_date_proposal"
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: ListStoryDateProposals :many
-- Lists proposals for a story with proposer profile info and viewer's vote direction.
SELECT
  sdp.id,
  sdp.story_id,
  sdp.proposer_profile_id,
  sdp.datetime_start,
  sdp.datetime_end,
  sdp.is_finalized,
  sdp.vote_score,
  sdp.upvote_count,
  sdp.downvote_count,
  sdp.created_at,
  sdp.updated_at,
  p.slug AS proposer_profile_slug,
  pt.title AS proposer_profile_title,
  p.profile_picture_uri AS proposer_profile_picture_uri,
  p.kind AS proposer_profile_kind,
  CASE
    WHEN sqlc.narg(viewer_profile_id)::TEXT IS NOT NULL THEN
      COALESCE(
        (SELECT sdpv.direction FROM "story_date_proposal_vote" sdpv
         WHERE sdpv.proposal_id = sdp.id
           AND sdpv.voter_profile_id = sqlc.narg(viewer_profile_id)::TEXT),
        0
      )::SMALLINT
    ELSE 0::SMALLINT
  END AS viewer_vote_direction
FROM "story_date_proposal" sdp
  INNER JOIN "profile" p ON p.id = sdp.proposer_profile_id AND p.deleted_at IS NULL
  INNER JOIN "profile_tx" pt ON pt.profile_id = p.id
    AND pt.locale_code = (
      SELECT ptx.locale_code FROM "profile_tx" ptx
      WHERE ptx.profile_id = p.id
      ORDER BY CASE
        WHEN ptx.locale_code = sqlc.arg(locale_code) THEN 0
        WHEN ptx.locale_code = p.default_locale THEN 1
        ELSE 2
      END
      LIMIT 1
    )
WHERE sdp.story_id = sqlc.arg(story_id)
  AND sdp.deleted_at IS NULL
ORDER BY sdp.vote_score DESC, sdp.created_at ASC;

-- name: SoftDeleteStoryDateProposal :exec
UPDATE "story_date_proposal"
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: MarkStoryDateProposalFinalized :exec
UPDATE "story_date_proposal"
SET is_finalized = TRUE, updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: GetStoryDateProposalVote :one
SELECT * FROM "story_date_proposal_vote"
WHERE proposal_id = sqlc.arg(proposal_id)
  AND voter_profile_id = sqlc.arg(voter_profile_id);

-- name: InsertStoryDateProposalVote :one
INSERT INTO "story_date_proposal_vote" (
  id, proposal_id, voter_profile_id, direction, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(proposal_id),
  sqlc.arg(voter_profile_id),
  sqlc.arg(direction),
  NOW()
) RETURNING *;

-- name: UpdateStoryDateProposalVoteDirection :exec
UPDATE "story_date_proposal_vote"
SET direction = sqlc.arg(direction), updated_at = NOW()
WHERE proposal_id = sqlc.arg(proposal_id)
  AND voter_profile_id = sqlc.arg(voter_profile_id);

-- name: DeleteStoryDateProposalVote :exec
DELETE FROM "story_date_proposal_vote"
WHERE proposal_id = sqlc.arg(proposal_id)
  AND voter_profile_id = sqlc.arg(voter_profile_id);

-- name: DeleteAllVotesForProposal :exec
DELETE FROM "story_date_proposal_vote"
WHERE proposal_id = sqlc.arg(proposal_id);

-- name: AdjustStoryDateProposalVoteScore :exec
UPDATE "story_date_proposal"
SET
  vote_score = vote_score + sqlc.arg(score_delta)::INTEGER,
  upvote_count = GREATEST(upvote_count + sqlc.arg(upvote_delta)::INTEGER, 0),
  downvote_count = GREATEST(downvote_count + sqlc.arg(downvote_delta)::INTEGER, 0),
  updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: GetStoryKindAndProperties :one
SELECT kind, properties, author_profile_id FROM "story"
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: UpdateStoryPropertiesForDateFinalization :exec
UPDATE "story"
SET properties = sqlc.arg(properties)::JSONB, updated_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;
