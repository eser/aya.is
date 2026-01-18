-- PoW Challenges

-- name: CreatePOWChallenge :exec
INSERT INTO
  protection_pow_challenge (
    id,
    prefix,
    difficulty,
    ip_hash,
    used,
    expires_at,
    created_at
  )
VALUES
  (
    sqlc.arg(id),
    sqlc.arg(prefix),
    sqlc.arg(difficulty),
    sqlc.arg(ip_hash),
    sqlc.arg(used),
    sqlc.arg(expires_at),
    sqlc.arg(created_at)
  );

-- name: GetPOWChallengeByID :one
SELECT
  id,
  prefix,
  difficulty,
  ip_hash,
  used,
  expires_at,
  created_at
FROM
  protection_pow_challenge
WHERE
  id = sqlc.arg(id);

-- name: MarkPOWChallengeUsed :exec
UPDATE protection_pow_challenge
SET
  used = TRUE
WHERE
  id = sqlc.arg(id);

-- name: DeleteExpiredPOWChallenges :exec
DELETE FROM protection_pow_challenge
WHERE
  expires_at < NOW();

-- name: CountPOWChallengesByIPHash :one
SELECT
  COUNT(*) AS count
FROM
  protection_pow_challenge
WHERE
  ip_hash = sqlc.arg(ip_hash)
  AND created_at > NOW() - INTERVAL '1 hour'
  AND used = FALSE;
