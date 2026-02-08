-- name: InsertEventAudit :exec
INSERT INTO "event_audit" (
  id, event_type, entity_type, entity_id,
  actor_id, actor_kind, session_id, payload, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(event_type),
  sqlc.arg(entity_type),
  sqlc.arg(entity_id),
  sqlc.arg(actor_id),
  sqlc.arg(actor_kind),
  sqlc.arg(session_id),
  sqlc.arg(payload),
  NOW()
);

-- name: ListEventAuditByEntity :many
SELECT *
FROM "event_audit"
WHERE entity_type = sqlc.arg(entity_type)
  AND entity_id = sqlc.arg(entity_id)
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_count);
