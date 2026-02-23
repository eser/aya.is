-- name: GetDiscussionVisibility :one
SELECT feature_discussions
FROM "profile"
WHERE id = sqlc.arg(id);

-- name: GetDiscussionThread :one
SELECT *
FROM "discussion_thread"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: GetDiscussionThreadByStoryID :one
SELECT *
FROM "discussion_thread"
WHERE story_id = sqlc.arg(story_id)
  AND deleted_at IS NULL;

-- name: GetDiscussionThreadByProfileID :one
SELECT *
FROM "discussion_thread"
WHERE profile_id = sqlc.arg(profile_id)
  AND deleted_at IS NULL;

-- name: InsertDiscussionThread :one
INSERT INTO "discussion_thread" (
  id,
  story_id,
  profile_id,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.narg(story_id),
  sqlc.narg(profile_id),
  NOW()
) RETURNING *;

-- name: UpdateDiscussionThreadLocked :exec
UPDATE "discussion_thread"
SET
  is_locked = sqlc.arg(is_locked),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: IncrementDiscussionThreadCommentCount :exec
UPDATE "discussion_thread"
SET
  comment_count = comment_count + 1,
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DecrementDiscussionThreadCommentCount :exec
UPDATE "discussion_thread"
SET
  comment_count = GREATEST(comment_count - 1, 0),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: GetDiscussionComment :one
SELECT
  dc.*,
  u.individual_profile_id AS author_profile_id,
  ap.slug AS author_profile_slug,
  apt.title AS author_profile_title,
  ap.profile_picture_uri AS author_profile_picture_uri
FROM "discussion_comment" dc
  LEFT JOIN "user" u ON u.id = dc.author_user_id
  LEFT JOIN "profile" ap ON ap.id = u.individual_profile_id
  LEFT JOIN "profile_tx" apt ON apt.profile_id = ap.id
    AND apt.locale_code = (
      SELECT aptf.locale_code FROM "profile_tx" aptf
      WHERE aptf.profile_id = ap.id
        AND (aptf.locale_code = sqlc.arg(locale_code) OR aptf.locale_code = ap.default_locale)
      ORDER BY CASE WHEN aptf.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
      LIMIT 1
    )
WHERE dc.id = sqlc.arg(id)
  AND dc.deleted_at IS NULL;

-- name: GetDiscussionCommentRaw :one
SELECT *
FROM "discussion_comment"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListTopLevelDiscussionComments :many
SELECT
  dc.*,
  u.individual_profile_id AS author_profile_id,
  ap.slug AS author_profile_slug,
  apt.title AS author_profile_title,
  ap.profile_picture_uri AS author_profile_picture_uri,
  CASE
    WHEN sqlc.narg(viewer_user_id)::TEXT IS NOT NULL THEN
      COALESCE(
        (SELECT dcv.direction FROM "discussion_comment_vote" dcv
         WHERE dcv.comment_id = dc.id AND dcv.user_id = sqlc.narg(viewer_user_id)::TEXT),
        0
      )
    ELSE 0
  END AS viewer_vote_direction
FROM "discussion_comment" dc
  LEFT JOIN "user" u ON u.id = dc.author_user_id
  LEFT JOIN "profile" ap ON ap.id = u.individual_profile_id
  LEFT JOIN "profile_tx" apt ON apt.profile_id = ap.id
    AND apt.locale_code = (
      SELECT aptf.locale_code FROM "profile_tx" aptf
      WHERE aptf.profile_id = ap.id
        AND (aptf.locale_code = sqlc.arg(locale_code) OR aptf.locale_code = ap.default_locale)
      ORDER BY CASE WHEN aptf.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
      LIMIT 1
    )
WHERE dc.thread_id = sqlc.arg(thread_id)
  AND dc.parent_id IS NULL
  AND dc.deleted_at IS NULL
  AND (sqlc.arg(include_hidden)::BOOLEAN = TRUE OR dc.is_hidden = FALSE)
ORDER BY
  dc.is_pinned DESC,
  CASE WHEN sqlc.arg(sort_mode)::TEXT = 'hot' THEN dc.vote_score ELSE 0 END DESC,
  CASE WHEN sqlc.arg(sort_mode)::TEXT = 'oldest' THEN EXTRACT(EPOCH FROM dc.created_at) ELSE 0 END ASC,
  dc.created_at DESC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);

-- name: ListChildDiscussionComments :many
SELECT
  dc.*,
  u.individual_profile_id AS author_profile_id,
  ap.slug AS author_profile_slug,
  apt.title AS author_profile_title,
  ap.profile_picture_uri AS author_profile_picture_uri,
  CASE
    WHEN sqlc.narg(viewer_user_id)::TEXT IS NOT NULL THEN
      COALESCE(
        (SELECT dcv.direction FROM "discussion_comment_vote" dcv
         WHERE dcv.comment_id = dc.id AND dcv.user_id = sqlc.narg(viewer_user_id)::TEXT),
        0
      )
    ELSE 0
  END AS viewer_vote_direction
FROM "discussion_comment" dc
  LEFT JOIN "user" u ON u.id = dc.author_user_id
  LEFT JOIN "profile" ap ON ap.id = u.individual_profile_id
  LEFT JOIN "profile_tx" apt ON apt.profile_id = ap.id
    AND apt.locale_code = (
      SELECT aptf.locale_code FROM "profile_tx" aptf
      WHERE aptf.profile_id = ap.id
        AND (aptf.locale_code = sqlc.arg(locale_code) OR aptf.locale_code = ap.default_locale)
      ORDER BY CASE WHEN aptf.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
      LIMIT 1
    )
WHERE dc.thread_id = sqlc.arg(thread_id)
  AND dc.parent_id = sqlc.arg(parent_id)
  AND dc.deleted_at IS NULL
  AND (sqlc.arg(include_hidden)::BOOLEAN = TRUE OR dc.is_hidden = FALSE)
ORDER BY dc.created_at ASC
LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);

-- name: InsertDiscussionComment :one
INSERT INTO "discussion_comment" (
  id,
  thread_id,
  parent_id,
  author_user_id,
  content,
  depth,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(thread_id),
  sqlc.narg(parent_id),
  sqlc.arg(author_user_id),
  sqlc.arg(content),
  sqlc.arg(depth),
  NOW()
) RETURNING *;

-- name: UpdateDiscussionCommentContent :exec
UPDATE "discussion_comment"
SET
  content = sqlc.arg(content),
  is_edited = TRUE,
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: SoftDeleteDiscussionComment :exec
UPDATE "discussion_comment"
SET
  deleted_at = NOW(),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: UpdateDiscussionCommentHidden :exec
UPDATE "discussion_comment"
SET
  is_hidden = sqlc.arg(is_hidden),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: UpdateDiscussionCommentPinned :exec
UPDATE "discussion_comment"
SET
  is_pinned = sqlc.arg(is_pinned),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: IncrementDiscussionCommentReplyCount :exec
UPDATE "discussion_comment"
SET
  reply_count = reply_count + 1,
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DecrementDiscussionCommentReplyCount :exec
UPDATE "discussion_comment"
SET
  reply_count = GREATEST(reply_count - 1, 0),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: GetDiscussionCommentVote :one
SELECT *
FROM "discussion_comment_vote"
WHERE comment_id = sqlc.arg(comment_id)
  AND user_id = sqlc.arg(user_id);

-- name: InsertDiscussionCommentVote :one
INSERT INTO "discussion_comment_vote" (
  id,
  comment_id,
  user_id,
  direction,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(comment_id),
  sqlc.arg(user_id),
  sqlc.arg(direction),
  NOW()
) RETURNING *;

-- name: UpdateDiscussionCommentVoteDirection :exec
UPDATE "discussion_comment_vote"
SET direction = sqlc.arg(direction)
WHERE comment_id = sqlc.arg(comment_id)
  AND user_id = sqlc.arg(user_id);

-- name: DeleteDiscussionCommentVote :exec
DELETE FROM "discussion_comment_vote"
WHERE comment_id = sqlc.arg(comment_id)
  AND user_id = sqlc.arg(user_id);

-- name: AdjustDiscussionCommentVoteScore :exec
UPDATE "discussion_comment"
SET
  vote_score = vote_score + sqlc.arg(score_delta)::INTEGER,
  upvote_count = GREATEST(upvote_count + sqlc.arg(upvote_delta)::INTEGER, 0),
  downvote_count = GREATEST(downvote_count + sqlc.arg(downvote_delta)::INTEGER, 0),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: GetStoryAuthorProfileID :one
SELECT author_profile_id
FROM "story"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;
