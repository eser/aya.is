-- ============================================================
-- Conversations
-- ============================================================

-- name: CreateMailboxConversation :exec
INSERT INTO "mailbox_conversation" (
  id, kind, title, created_by_profile_id, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(kind),
  sqlc.narg(title),
  sqlc.narg(created_by_profile_id),
  NOW()
);

-- name: GetMailboxConversationByID :one
SELECT *
FROM "mailbox_conversation"
WHERE id = sqlc.arg(id);

-- name: FindDirectConversation :one
SELECT mc.*
FROM "mailbox_conversation" mc
WHERE mc.kind = 'direct'
  AND EXISTS (
    SELECT 1 FROM "mailbox_participant" mp
    WHERE mp.conversation_id = mc.id
      AND mp.profile_id = sqlc.arg(profile_a)
      AND mp.left_at IS NULL
  )
  AND EXISTS (
    SELECT 1 FROM "mailbox_participant" mp
    WHERE mp.conversation_id = mc.id
      AND mp.profile_id = sqlc.arg(profile_b)
      AND mp.left_at IS NULL
  )
  AND (
    SELECT COUNT(*) FROM "mailbox_participant" mp
    WHERE mp.conversation_id = mc.id AND mp.left_at IS NULL
  ) = 2
LIMIT 1;

-- name: ListConversationsForProfile :many
SELECT
  mc.id,
  mc.kind,
  mc.title,
  mc.created_by_profile_id,
  mc.created_at,
  mc.updated_at,
  mp.last_read_at,
  mp.is_archived,
  (
    SELECT me.message FROM "mailbox_envelope" me
    WHERE me.conversation_id = mc.id AND me.deleted_at IS NULL
    ORDER BY me.created_at DESC LIMIT 1
  ) AS last_envelope_message,
  (
    SELECT me.kind FROM "mailbox_envelope" me
    WHERE me.conversation_id = mc.id AND me.deleted_at IS NULL
    ORDER BY me.created_at DESC LIMIT 1
  ) AS last_envelope_kind,
  (
    SELECT me.created_at FROM "mailbox_envelope" me
    WHERE me.conversation_id = mc.id AND me.deleted_at IS NULL
    ORDER BY me.created_at DESC LIMIT 1
  ) AS last_envelope_at,
  (
    SELECT me.sender_profile_id FROM "mailbox_envelope" me
    WHERE me.conversation_id = mc.id AND me.deleted_at IS NULL
    ORDER BY me.created_at DESC LIMIT 1
  ) AS last_envelope_sender_profile_id,
  (
    SELECT COUNT(*)::INT FROM "mailbox_envelope" me
    WHERE me.conversation_id = mc.id
      AND me.deleted_at IS NULL
      AND (mp.last_read_at IS NULL OR me.created_at > mp.last_read_at)
  ) AS unread_count
FROM "mailbox_conversation" mc
  INNER JOIN "mailbox_participant" mp
    ON mp.conversation_id = mc.id
    AND mp.profile_id = sqlc.arg(profile_id)
    AND mp.left_at IS NULL
WHERE mp.is_archived = sqlc.arg(include_archived)::BOOLEAN
ORDER BY last_envelope_at DESC NULLS LAST
LIMIT sqlc.arg(limit_count);

-- name: UpdateConversationTimestamp :exec
UPDATE "mailbox_conversation"
SET updated_at = NOW()
WHERE id = sqlc.arg(id);

-- ============================================================
-- Participants
-- ============================================================

-- name: AddMailboxParticipant :exec
INSERT INTO "mailbox_participant" (
  id, conversation_id, profile_id, joined_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(conversation_id),
  sqlc.arg(profile_id),
  NOW()
);

-- name: GetMailboxParticipant :one
SELECT
  mp.*,
  p.slug AS profile_slug,
  COALESCE(pt.title, '') AS profile_title,
  p.profile_picture_uri AS profile_picture_uri,
  p.kind AS profile_kind
FROM "mailbox_participant" mp
  INNER JOIN "profile" p ON p.id = mp.profile_id AND p.deleted_at IS NULL
  LEFT JOIN "profile_tx" pt ON pt.profile_id = p.id AND pt.locale_code = p.default_locale
WHERE mp.conversation_id = sqlc.arg(conversation_id)
  AND mp.profile_id = sqlc.arg(profile_id)
  AND mp.left_at IS NULL;

-- name: ListMailboxParticipants :many
SELECT
  mp.*,
  p.slug AS profile_slug,
  COALESCE(pt.title, '') AS profile_title,
  p.profile_picture_uri AS profile_picture_uri,
  p.kind AS profile_kind
FROM "mailbox_participant" mp
  INNER JOIN "profile" p ON p.id = mp.profile_id AND p.deleted_at IS NULL
  LEFT JOIN "profile_tx" pt ON pt.profile_id = p.id AND pt.locale_code = p.default_locale
WHERE mp.conversation_id = sqlc.arg(conversation_id)
  AND mp.left_at IS NULL
ORDER BY mp.joined_at, mp.id;

-- name: UpdateParticipantReadCursor :exec
UPDATE "mailbox_participant"
SET last_read_at = NOW()
WHERE conversation_id = sqlc.arg(conversation_id)
  AND profile_id = sqlc.arg(profile_id)
  AND left_at IS NULL;

-- name: SetParticipantArchived :exec
UPDATE "mailbox_participant"
SET is_archived = sqlc.arg(is_archived)
WHERE conversation_id = sqlc.arg(conversation_id)
  AND profile_id = sqlc.arg(profile_id)
  AND left_at IS NULL;

-- ============================================================
-- Envelopes (messages)
-- ============================================================

-- name: CreateMailboxEnvelope :exec
INSERT INTO "mailbox_envelope" (
  id, conversation_id, target_profile_id, sender_profile_id, sender_user_id,
  kind, status, message, properties, reply_to_id, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(conversation_id),
  sqlc.arg(target_profile_id),
  sqlc.narg(sender_profile_id),
  sqlc.narg(sender_user_id),
  sqlc.arg(kind),
  'pending',
  sqlc.narg(message),
  sqlc.narg(properties),
  sqlc.narg(reply_to_id),
  NOW()
);

-- name: GetMailboxEnvelopeByID :one
SELECT *
FROM "mailbox_envelope"
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListEnvelopesByConversation :many
SELECT
  me.*,
  sp.slug AS sender_profile_slug,
  COALESCE(spt.title, '') AS sender_profile_title,
  sp.profile_picture_uri AS sender_profile_picture_uri,
  sp.kind AS sender_profile_kind
FROM "mailbox_envelope" me
  LEFT JOIN "profile" sp ON sp.id = me.sender_profile_id AND sp.deleted_at IS NULL
  LEFT JOIN "profile_tx" spt ON spt.profile_id = sp.id AND spt.locale_code = sp.default_locale
WHERE me.conversation_id = sqlc.arg(conversation_id)
  AND me.deleted_at IS NULL
ORDER BY me.created_at ASC
LIMIT sqlc.arg(limit_count);

-- name: ListMailboxEnvelopesByTargetProfileID :many
SELECT
  me.*,
  sp.slug AS sender_profile_slug,
  COALESCE(spt.title, '') AS sender_profile_title,
  sp.profile_picture_uri AS sender_profile_picture_uri,
  sp.kind AS sender_profile_kind
FROM "mailbox_envelope" me
  LEFT JOIN "profile" sp ON sp.id = me.sender_profile_id AND sp.deleted_at IS NULL
  LEFT JOIN "profile_tx" spt ON spt.profile_id = sp.id AND spt.locale_code = sp.default_locale
WHERE me.target_profile_id = sqlc.arg(target_profile_id)
  AND me.deleted_at IS NULL
  AND (sqlc.narg(status_filter)::TEXT IS NULL OR me.status = sqlc.narg(status_filter))
ORDER BY me.created_at DESC
LIMIT sqlc.arg(limit_count);

-- name: UpdateMailboxEnvelopeStatusToAccepted :execrows
UPDATE "mailbox_envelope"
SET status = 'accepted',
    accepted_at = sqlc.arg(now),
    updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(id)
  AND status = 'pending'
  AND deleted_at IS NULL;

-- name: UpdateMailboxEnvelopeStatusToRejected :execrows
UPDATE "mailbox_envelope"
SET status = 'rejected',
    rejected_at = sqlc.arg(now),
    updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(id)
  AND status = 'pending'
  AND deleted_at IS NULL;

-- name: UpdateMailboxEnvelopeStatusToRevoked :execrows
UPDATE "mailbox_envelope"
SET status = 'revoked',
    revoked_at = sqlc.arg(now),
    updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(id)
  AND status = 'pending'
  AND deleted_at IS NULL;

-- name: UpdateMailboxEnvelopeStatusToRedeemed :execrows
UPDATE "mailbox_envelope"
SET status = 'redeemed',
    redeemed_at = sqlc.arg(now),
    updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(id)
  AND status = 'accepted'
  AND deleted_at IS NULL;

-- name: UpdateMailboxEnvelopeProperties :exec
UPDATE "mailbox_envelope"
SET properties = sqlc.arg(properties),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
  AND deleted_at IS NULL;

-- name: ListAcceptedMailboxInvitations :many
SELECT
  me.*,
  sp.slug AS sender_profile_slug,
  COALESCE(spt.title, '') AS sender_profile_title,
  sp.profile_picture_uri AS sender_profile_picture_uri,
  sp.kind AS sender_profile_kind
FROM "mailbox_envelope" me
  LEFT JOIN "profile" sp ON sp.id = me.sender_profile_id AND sp.deleted_at IS NULL
  LEFT JOIN "profile_tx" spt ON spt.profile_id = sp.id AND spt.locale_code = sp.default_locale
WHERE me.target_profile_id = sqlc.arg(target_profile_id)
  AND me.kind = 'invitation'
  AND me.status = 'accepted'
  AND me.deleted_at IS NULL
  AND (sqlc.narg(invitation_kind)::TEXT IS NULL
       OR me.properties->>'invitation_kind' = sqlc.narg(invitation_kind))
ORDER BY me.created_at DESC;

-- name: CountPendingMailboxEnvelopes :one
SELECT COUNT(*)::INT AS count
FROM "mailbox_envelope"
WHERE target_profile_id = sqlc.arg(target_profile_id)
  AND status = 'pending'
  AND deleted_at IS NULL;

-- ============================================================
-- Reactions
-- ============================================================

-- name: AddMailboxReaction :exec
INSERT INTO "mailbox_reaction" (id, envelope_id, profile_id, emoji, created_at)
VALUES (sqlc.arg(id), sqlc.arg(envelope_id), sqlc.arg(profile_id), sqlc.arg(emoji), NOW())
ON CONFLICT (envelope_id, profile_id, emoji) DO NOTHING;

-- name: RemoveMailboxReaction :execrows
DELETE FROM "mailbox_reaction"
WHERE envelope_id = sqlc.arg(envelope_id)
  AND profile_id = sqlc.arg(profile_id)
  AND emoji = sqlc.arg(emoji);

-- name: ListReactionsByEnvelope :many
SELECT
  mr.*,
  p.slug AS profile_slug,
  COALESCE(pt.title, '') AS profile_title
FROM "mailbox_reaction" mr
  INNER JOIN "profile" p ON p.id = mr.profile_id AND p.deleted_at IS NULL
  LEFT JOIN "profile_tx" pt ON pt.profile_id = p.id AND pt.locale_code = p.default_locale
WHERE mr.envelope_id = sqlc.arg(envelope_id)
ORDER BY mr.created_at;

-- ============================================================
-- Aggregate counts
-- ============================================================

-- name: CountUnreadConversations :one
SELECT COUNT(DISTINCT mc.id)::INT AS count
FROM "mailbox_conversation" mc
  INNER JOIN "mailbox_participant" mp
    ON mp.conversation_id = mc.id
    AND mp.profile_id = sqlc.arg(profile_id)
    AND mp.left_at IS NULL
    AND mp.is_archived = FALSE
WHERE EXISTS (
  SELECT 1 FROM "mailbox_envelope" me
  WHERE me.conversation_id = mc.id
    AND me.deleted_at IS NULL
    AND (mp.last_read_at IS NULL OR me.created_at > mp.last_read_at)
);
