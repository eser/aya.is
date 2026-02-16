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
  ap.slug AS author_profile_slug,
  apt.title AS author_profile_title,
  bp.slug AS answered_by_profile_slug,
  bpt.title AS answered_by_profile_title,
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
  LEFT JOIN "profile" ap ON ap.id = u.individual_profile_id
  LEFT JOIN "profile_tx" apt ON apt.profile_id = ap.id
    AND apt.locale_code = (
      SELECT aptf.locale_code FROM "profile_tx" aptf
      WHERE aptf.profile_id = ap.id
        AND (aptf.locale_code = sqlc.arg(locale_code) OR aptf.locale_code = ap.default_locale)
      ORDER BY CASE WHEN aptf.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
      LIMIT 1
    )
  LEFT JOIN "profile" bp ON bp.id = pq.answered_by
  LEFT JOIN "profile_tx" bpt ON bpt.profile_id = bp.id
    AND bpt.locale_code = (
      SELECT bptf.locale_code FROM "profile_tx" bptf
      WHERE bptf.profile_id = bp.id
        AND (bptf.locale_code = sqlc.arg(locale_code) OR bptf.locale_code = bp.default_locale)
      ORDER BY CASE WHEN bptf.locale_code = sqlc.arg(locale_code) THEN 0 ELSE 1 END
      LIMIT 1
    )
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

-- name: EditProfileQuestionAnswer :exec
UPDATE "profile_question"
SET
  answer_content = sqlc.arg(answer_content),
  answer_uri = sqlc.narg(answer_uri),
  answer_kind = sqlc.narg(answer_kind),
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
