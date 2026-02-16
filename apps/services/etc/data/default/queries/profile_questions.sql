-- name: IsProfileQAHidden :one
SELECT hide_qa
FROM "profile"
WHERE id = sqlc.arg(id);

-- name: GetProfileQuestion :one
SELECT *
FROM "profile_question"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListProfileQuestionsByProfileID :many
SELECT
  pq.*,
  u.individual_profile_id AS author_profile_id,
  CASE
    WHEN sqlc.narg(viewer_user_id)::TEXT IS NOT NULL THEN
      EXISTS(
        SELECT 1 FROM "profile_question_vote" pqv
        WHERE pqv.question_id = pq.id
          AND pqv.user_id = sqlc.narg(viewer_user_id)::TEXT
      )
    ELSE FALSE
  END AS has_viewer_vote
FROM "profile_question" pq
  LEFT JOIN "user" u ON u.id = pq.author_user_id
WHERE pq.profile_id = sqlc.arg(profile_id)
  AND pq.deleted_at IS NULL
  AND (sqlc.arg(include_hidden)::BOOLEAN = TRUE OR pq.is_hidden = FALSE)
ORDER BY pq.vote_count DESC, pq.created_at DESC;

-- name: InsertProfileQuestion :one
INSERT INTO "profile_question" (
  id,
  profile_id,
  author_user_id,
  content,
  is_anonymous,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(profile_id),
  sqlc.arg(author_user_id),
  sqlc.arg(content),
  sqlc.arg(is_anonymous),
  NOW()
) RETURNING *;

-- name: UpdateProfileQuestionAnswer :exec
UPDATE "profile_question"
SET
  answer_content = sqlc.arg(answer_content),
  answer_uri = sqlc.narg(answer_uri),
  answer_kind = sqlc.narg(answer_kind),
  answered_at = NOW(),
  answered_by = sqlc.arg(answered_by),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: UpdateProfileQuestionHidden :exec
UPDATE "profile_question"
SET
  is_hidden = sqlc.arg(is_hidden),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: InsertProfileQuestionVote :one
INSERT INTO "profile_question_vote" (
  id,
  question_id,
  user_id,
  score,
  created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(question_id),
  sqlc.arg(user_id),
  sqlc.arg(score),
  NOW()
) RETURNING *;

-- name: IncrementProfileQuestionVoteCount :exec
UPDATE "profile_question"
SET
  vote_count = vote_count + 1,
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DecrementProfileQuestionVoteCount :exec
UPDATE "profile_question"
SET
  vote_count = GREATEST(vote_count - 1, 0),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: DeleteProfileQuestionVote :exec
DELETE FROM "profile_question_vote"
WHERE question_id = sqlc.arg(question_id)
  AND user_id = sqlc.arg(user_id);

-- name: GetProfileQuestionVote :one
SELECT *
FROM "profile_question_vote"
WHERE question_id = sqlc.arg(question_id)
  AND user_id = sqlc.arg(user_id);

-- name: CountProfileQuestionsByProfileID :one
SELECT COUNT(*) AS count
FROM "profile_question"
WHERE profile_id = sqlc.arg(profile_id)
  AND deleted_at IS NULL;
