-- name: EnqueueEvent :exec
INSERT INTO "event_queue" (
  id, type, payload, status, max_retries,
  visibility_timeout_secs, visible_at, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(type),
  sqlc.arg(payload),
  'pending',
  sqlc.arg(max_retries),
  sqlc.arg(visibility_timeout_secs),
  sqlc.arg(visible_at),
  NOW()
);

-- name: ClaimNextEvent :one
-- CTE-based claim: atomically selects + locks + updates.
-- Picks up both pending events that are due AND stale processing events
-- past their visibility timeout (crash recovery built into the claim).
-- Increments retry_count at claim time for crash safety.
WITH claimable AS (
  SELECT id FROM "event_queue"
  WHERE (
    (status = 'pending' AND visible_at <= NOW())
    OR (status = 'processing' AND visible_at <= NOW())
  )
  AND retry_count <= max_retries
  ORDER BY visible_at ASC
  LIMIT 1
  FOR UPDATE SKIP LOCKED
)
UPDATE "event_queue"
SET
  status = 'processing',
  started_at = NOW(),
  visible_at = NOW() + visibility_timeout_secs * INTERVAL '1 second',
  retry_count = retry_count + 1,
  worker_id = sqlc.arg(worker_id),
  updated_at = NOW()
FROM claimable
WHERE "event_queue".id = claimable.id
RETURNING "event_queue".*;

-- name: CompleteEvent :execrows
-- Worker ID check prevents a timed-out worker from completing
-- a job that was already re-claimed by another worker.
UPDATE "event_queue"
SET
  status = 'completed',
  completed_at = NOW(),
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND status = 'processing'
  AND worker_id = sqlc.arg(worker_id);

-- name: FailEvent :execrows
-- On failure: if retries exhausted -> dead, otherwise -> pending with backoff.
-- Worker ID check prevents stale workers from interfering.
UPDATE "event_queue"
SET
  status = CASE
    WHEN retry_count >= max_retries THEN 'dead'
    ELSE 'pending'
  END,
  error_message = sqlc.arg(error_message),
  failed_at = NOW(),
  visible_at = CASE
    WHEN retry_count >= max_retries THEN visible_at
    ELSE NOW() + sqlc.arg(backoff_seconds)::INTEGER * INTERVAL '1 second'
  END,
  worker_id = NULL,
  updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND status = 'processing'
  AND worker_id = sqlc.arg(worker_id);

-- name: ListEventsByType :many
SELECT *
FROM "event_queue"
WHERE type = sqlc.arg(type)
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count);
